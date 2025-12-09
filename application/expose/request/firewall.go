package request

import (
	"net/http"
	"net/netip"
	"strings"

	"github.com/xmx/aegis-server/application/errcode"
	"github.com/xmx/aegis-server/library/iplist"
)

type FirewallUpsert struct {
	Name         string   `json:"name"          validate:"required,gte=2,lte=20"`
	Enabled      bool     `json:"enabled"`
	Blacklist    bool     `json:"blacklist"`                                                              // 是否黑名单模式，反之白名单模式。
	TrustHeaders []string `json:"trust_headers" validate:"lte=100,unique,dive,required,lte=100"`          // 取 IP 的可信 Headers。
	TrustProxies []string `json:"trust_proxies" validate:"lte=100,unique,dive,required,ip_range"`         // 可信网关。
	CountryMode  bool     `json:"country_mode"`                                                           // 是否启用国家（地区）模式，否则为 IPNets 模式
	IPNets       []string `json:"ip_nets"       validate:"lte=1000,unique,dive,required,ip_range"`        // IP 列表
	Countries    []string `json:"countries"     validate:"lte=300,unique,dive,required,iso3166_1_alpha2"` // https://www.iso.org/iso-3166-country-codes.html
}

func (fu FirewallUpsert) Format() (FirewallUpsert, error) {
	ret := FirewallUpsert{
		Name:        fu.Name,
		Enabled:     fu.Enabled,
		Blacklist:   fu.Blacklist,
		CountryMode: fu.CountryMode,
		Countries:   fu.Countries,
	}
	// TrustHeaders TrustProxies 要么同时为空，要么同时有值。
	if (len(fu.TrustHeaders) == 0) != (len(fu.TrustProxies) == 0) {
		return ret, errcode.ErrTrustHeaderProxy
	}
	if fu.CountryMode { // 国家地区模式，那么
		if len(fu.Countries) == 0 {
			return ret, errcode.ErrISOCodeRequired
		}
	} else if len(fu.IPNets) == 0 {
		return ret, errcode.ErrIPNetsRequired
	}

	proxies, err := fu.deduplicationIP(fu.TrustProxies)
	if err != nil {
		return ret, err
	}
	ret.TrustProxies = proxies

	inets, err := fu.deduplicationIP(fu.IPNets)
	if err != nil {
		return ret, err
	}
	ret.IPNets = inets

	ret.TrustProxies = proxies
	ret.IPNets = inets
	ret.TrustHeaders = fu.deduplicationHeader(fu.TrustHeaders)

	return ret, nil
}

func (fu FirewallUpsert) deduplicationIP(inets []string) ([]string, error) {
	uniq := make(map[string]struct{}, len(inets))
	result := make([]string, 0, len(inets))
	for _, inet := range inets {
		if strings.Contains(inet, "/") { // CIDR 归一化处理
			pre, err := netip.ParsePrefix(inet)
			if err != nil {
				return nil, err
			}
			inet = pre.Masked().String()
		}
		if _, exists := uniq[inet]; !exists {
			uniq[inet] = struct{}{}
			result = append(result, inet)
		}
	}
	if _, err := iplist.Parse(result); err != nil {
		return nil, err
	}

	return result, nil
}

func (fu FirewallUpsert) deduplicationHeader(headers []string) []string {
	uniq := make(map[string]struct{}, len(headers))
	result := make([]string, 0, len(headers))
	for _, header := range headers {
		header = http.CanonicalHeaderKey(header)
		if _, exists := uniq[header]; !exists {
			uniq[header] = struct{}{}
			headers = append(headers, header)
		}
	}

	return result
}
