package iplist

import "net/netip"

// NormalizeCIDR 将
func NormalizeCIDR(cidr string) (netip.Prefix, error) {
	pre, err := netip.ParsePrefix(cidr)
	if err != nil {
		return pre, err
	}

	return pre.Masked(), nil
}

type constIP struct {
	ip netip.Addr
}

func (cp *constIP) Contains(ip netip.Addr) bool {
	if ip.Is4In6() { // IPv4-Mapped IPv6: ::ffff:172.31.61.168
		ip = netip.AddrFrom4(ip.As4())
	}

	return cp.ip.Compare(ip) == 0
}

type ipv4Range struct {
	minIP netip.Addr
	maxIP netip.Addr
}

func (ir *ipv4Range) Contains(ip netip.Addr) bool {
	if !ip.Is4() && !ip.Is4In6() {
		return false
	}
	ipv4 := netip.AddrFrom4(ip.As4())

	return ir.minIP.Compare(ipv4) <= 0 &&
		ir.maxIP.Compare(ipv4) >= 0
}
