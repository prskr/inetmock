package rules

import "net"

type CIDR struct {
	*net.IPNet
}

func (c *CIDR) UnmarshalText(text []byte) (err error) {
	_, c.IPNet, err = net.ParseCIDR(string(text))
	return
}

func ParseCIDR(cidr string) (*CIDR, error) {
	var (
		c   = new(CIDR)
		err error
	)

	if _, c.IPNet, err = net.ParseCIDR(cidr); err != nil {
		return nil, err
	} else {
		return c, nil
	}
}

func MustParseCIDR(cidr string) *CIDR {
	c, err := ParseCIDR(cidr)
	if err != nil {
		panic(err)
	}
	return c
}
