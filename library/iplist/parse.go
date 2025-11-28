package iplist

import (
	"net"
	"net/netip"
	"strings"
)

type Lister interface {
	Contains(ip netip.Addr) bool
}

func Parse(ips []string) (Lister, error) {
	var ls listIPs
	for _, ip := range ips {
		if l, err := parse(ip); err != nil {
			return nil, err
		} else {
			ls = append(ls, l)
		}
	}

	return ls, nil
}

// parse 处理 IP/CIDR/IPv4-Range.
//
//   - IP:                  172.31.61.168, fe80::47:79a6:be7e:3a2d
//   - CIDR:                172.31.61.0/24, 2001:db8::/32
//   - IPv4 Range:          172.31.61.100-172.31.61.200
//   - [IPv4-Mapped IPv6]:  ::ffff:172.31.61.168
//
// [IPv4-Mapped IPv6]: https://datatracker.ietf.org/doc/html/rfc4291#section-2.5.5.2
func parse(strIP string) (Lister, error) {
	if strings.Contains(strIP, "/") {
		prefix, err := netip.ParsePrefix(strIP)
		if err != nil {
			return nil, &net.ParseError{Type: "CIDR address", Text: strIP}
		}
		return prefix, nil
	} else if strings.Contains(strIP, "-") { // IP 区间
		before, after, _ := strings.Cut(strIP, "-")

		minIP, err := netip.ParseAddr(before)
		if err != nil {
			return nil, &net.ParseError{Type: "IP range", Text: strIP}
		}
		maxIP, err := netip.ParseAddr(after)
		if err != nil {
			return nil, &net.ParseError{Type: "IP range", Text: strIP}
		}
		if !minIP.Is4() || !maxIP.Is4() {
			return nil, &net.ParseError{Type: "IP range", Text: strIP}
		}
		if minIP.Compare(maxIP) > 0 {
			return nil, &net.ParseError{Type: "IP range", Text: strIP}
		}
		rng := &ipv4Range{minIP: minIP, maxIP: maxIP}

		return rng, nil
	} else {
		addr, err := netip.ParseAddr(strIP)
		if err != nil {
			return nil, &net.ParseError{Type: "IP address", Text: strIP}
		}
		// IPv4: 172.31.61.168
		// IPv6: fe80::47:79a6:be7e:3a2d
		// IPv4-Mapped IPv6: ::ffff:172.31.61.168
		if addr.Is4In6() { // 归一化到 IPv4 表示形式。
			addr = netip.AddrFrom4(addr.As4())
		}

		return &constIP{ip: addr}, nil
	}
}

type listIPs []Lister

func (lps listIPs) Contains(ip netip.Addr) bool {
	for _, l := range lps {
		if l.Contains(ip) {
			return true
		}
	}

	return false
}
