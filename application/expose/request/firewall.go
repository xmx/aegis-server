package request

import (
	"net/http"
	"net/netip"
	"strings"

	"github.com/xmx/aegis-server/library/iplist"
)

type FirewallUpsert struct {
	Name         string   `json:"name"          validate:"required,lte=50"`
	Enabled      bool     `json:"enabled"`
	Blacklist    bool     `json:"blacklist"`                                                              // 是否黑名单模式，反之白名单模式。
	TrustHeaders []string `json:"trust_headers" validate:"lte=100,unique,dive,required,lte=100"`          // 取 IP 的可信 Headers。
	TrustProxies []string `json:"trust_proxies" validate:"lte=100,unique,dive,required,ip_range"`         // 可信网关。
	Inets        []string `json:"inets"         validate:"lte=1000,unique,dive,required,ip_range"`        // IP 列表
	Countries    []string `json:"countries"     validate:"lte=300,unique,dive,required,iso3166_1_alpha2"` // https://www.iso.org/iso-3166-country-codes.html
}

func (fu FirewallUpsert) Format() (FirewallUpsert, error) {
	ret := FirewallUpsert{Name: fu.Name, Enabled: fu.Enabled, Blacklist: fu.Blacklist, Countries: fu.Countries}
	uniq := make(map[string]struct{}, 16)

	proxies := make([]string, 0, len(fu.TrustProxies))
	for _, inet := range fu.TrustProxies {
		if strings.Contains(inet, "/") { // CIDR 归一化处理
			pre, err := netip.ParsePrefix(inet)
			if err != nil {
				return ret, err
			}
			inet = pre.Masked().String()
		}
		if _, exists := uniq[inet]; !exists {
			uniq[inet] = struct{}{}
			proxies = append(proxies, inet)
		}
	}
	if _, err := iplist.Parse(proxies); err != nil {
		return ret, err
	}

	clear(uniq)
	inets := make([]string, 0, len(fu.Inets))
	for _, inet := range fu.Inets {
		if strings.Contains(inet, "/") { // CIDR 归一化处理
			pre, err := netip.ParsePrefix(inet)
			if err != nil {
				return ret, err
			}
			inet = pre.Masked().String()
		}
		if _, exists := uniq[inet]; !exists {
			uniq[inet] = struct{}{}
			inets = append(inets, inet)
		}
	}
	if _, err := iplist.Parse(inets); err != nil {
		return ret, err
	}

	clear(uniq)
	headers := make([]string, 0, len(fu.TrustHeaders))
	for _, header := range fu.TrustHeaders {
		header = http.CanonicalHeaderKey(header)
		if _, exists := uniq[header]; !exists {
			uniq[header] = struct{}{}
			headers = append(headers, header)
		}
	}
	if _, err := iplist.Parse(inets); err != nil {
		return ret, err
	}

	ret.TrustProxies = proxies
	ret.Inets = inets
	ret.TrustHeaders = headers

	return ret, nil
}
