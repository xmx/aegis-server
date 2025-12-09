package firewalld

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"net/netip"

	"github.com/oschwald/geoip2-golang/v2"
	"github.com/xmx/aegis-server/library/iplist"
)

type Config struct {
	TrustHeaders []string
	TrustProxies iplist.Lister
	Blacklist    bool
	CountryMode  bool
	IPNets       iplist.Lister       // IP 模式时必填
	Countries    map[string]struct{} // 国家地区模式时必填
	MaxmindDB    MaxmindReader       // 国家地区模式时必填
}

type ConfigureFunc func(context.Context) (*Config, error)

func (cf ConfigureFunc) Configure(ctx context.Context) (*Config, error) {
	return cf(ctx)
}

type Configurer interface {
	// Configure 加载防火墙配置文件，如果未开启且配置未出错，请返回 nil, nil
	Configure(ctx context.Context) (*Config, error)
}

type MaxmindReader interface {
	ReadMaxmind(ctx context.Context) (*geoip2.Reader, error)
}

func New(cfg Configurer, log *slog.Logger) *Firewalld {
	return &Firewalld{
		cfg: cfg,
		log: log,
	}
}

type Firewalld struct {
	cfg Configurer
	log *slog.Logger
}

func (fw *Firewalld) Allowed(r *http.Request) (bool, error) {
	ctx, remoteAddr := r.Context(), r.RemoteAddr
	attrs := []any{"remote_addr", remoteAddr}
	cfg, err := fw.cfg.Configure(ctx)
	if err != nil {
		attrs = append(attrs, "error", err)
		fw.log.WarnContext(ctx, "加载防火墙配置出错", attrs...)
		return false, err
	}
	if cfg == nil {
		fw.log.Debug("未启用防火墙策略")
		return true, nil
	}

	remoteIP, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		attrs = append(attrs, "error", err)
		fw.log.WarnContext(ctx, "解析远程地址错误", attrs...)
		return false, err
	}
	directIP, err := netip.ParseAddr(remoteIP)
	if err != nil {
		attrs = append(attrs, "error", err)
		fw.log.WarnContext(ctx, "解析远程地址错误", attrs...)
		return false, err
	}
	attrs = append(attrs, "direct_ip", directIP)

	clientIP := directIP
	// TrustProxies 和 TrustHeaders 同时有值才有意义。
	if cfg.TrustProxies != nil && len(cfg.TrustHeaders) != 0 {
		if !cfg.TrustProxies.Contains(directIP) {
			fw.log.WarnContext(ctx, "来自不可信网关的请求", attrs...)
			return false, nil
		}

		var find bool
		for _, key := range cfg.TrustHeaders {
			if vals := r.Header.Values(key); len(vals) > 0 {
				if addr, err := netip.ParseAddr(vals[0]); err == nil {
					find = true
					clientIP = addr
				} else {
					attrs = append(attrs, "error", err)
					fw.log.WarnContext(ctx, "从可信 Header 中解析客户端 IP 出错", attrs...)
				}
			}
		}
		if !find {
			fw.log.WarnContext(ctx, "从可信 Header 中没有找到客户端 IP", attrs...)
			return false, nil
		}
	}

	var allowed bool
	if cfg.CountryMode {
		mdb, err := cfg.MaxmindDB.ReadMaxmind(ctx)
		if err != nil {
			attrs = append(attrs, "error", err)
			fw.log.WarnContext(ctx, "加载 IP 库出错", attrs...)
			return false, err
		}
		ret, err1 := mdb.Country(clientIP)
		if err1 != nil {
			attrs = append(attrs, "error", err1)
			fw.log.WarnContext(ctx, "查询 IP 库出错", attrs...)
			return false, err
		}
		isoCode := ret.Country.ISOCode
		if isoCode == "" {
			isoCode = ret.RegisteredCountry.ISOCode
		}
		if isoCode == "" {
			isoCode = ret.RepresentedCountry.ISOCode
		}
		attrs = append(attrs, "country_code", isoCode)
		fw.log.DebugContext(ctx, "已查询到所属国家", attrs...)
		_, allowed = cfg.Countries[isoCode]
	} else {
		allowed = cfg.IPNets.Contains(clientIP)
	}
	if cfg.Blacklist {
		return !allowed, nil
	}

	return allowed, nil
}
