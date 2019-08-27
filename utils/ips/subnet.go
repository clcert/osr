package ips

import (
	"net"
)

// Represents a sorted list of Subnets
// If it is not sorted, everything will fail.
type SubnetList []*net.IPNet

func (list SubnetList) ToIPChan() IPChan {
	ipChan := make(IPChan, 0)
	go func() {
		for _, subnet := range list {
			SubnetToIPs(subnet, ipChan)
		}
		close(ipChan)
	}()
	return ipChan
}

// Matches returns a subnet in the SubnetList if the IP passed as argument matches it
// and nil if no ips match it.
func (list SubnetList) Matches(ip net.IP) *net.IPNet {
	for _, subnet := range list {
		if subnet.Contains(ip) {
			return subnet
		}
	}
	return nil
}

func IntersectSubnets(subnet1, subnet2 *net.IPNet) (intersected *net.IPNet, before, after SubnetList) {
	before = make(SubnetList, 0)
	after = make(SubnetList, 0)
	if subnet1 == nil && subnet2 == nil {
		return
	} else if subnet1 == nil || subnet1.IP == nil || subnet1.Mask == nil {
		return
	} else if subnet2 == nil || subnet1.IP == nil || subnet1.Mask == nil {
		return
	}
	mask1, _ := subnet1.Mask.Size()
	mask2, _ := subnet2.Mask.Size()
	compare := CompareBytes(subnet1.IP, subnet2.IP)
	if compare > 0 || (compare == 0 && mask1 > mask2) {
		subnet1, subnet2 = subnet2, subnet1 // symmetry
		mask1, mask2 = mask2, mask1
	}
	switch {
	// Subnet with lowest IP could contain subnet with greatest IP
	case subnet1.Contains(subnet2.IP):
		switch {
		case mask1 < mask2:
			// subnet2 is the smallest subnet (its mask is bigger (range is smaller) and its IP is contained)

			// Split the biggest subnet into two smaller ones.
			subnet1Left, subnet1Right := SplitSubnet(subnet1)
			switch {
			case subnet1Left.Contains(subnet2.IP):
				subInter, subBefore, subAfter := IntersectSubnets(subnet1Left, subnet2)
				intersected = subInter
				before = append(before, subBefore...)
				after = append(after, subAfter...)
				after = append(after, subnet1Right)
			case subnet1Right.Contains(subnet2.IP):
				subInter, subBefore, subAfter := IntersectSubnets(subnet1Right, subnet2)
				intersected = subInter
				before = append(before, subnet1Left)
				before = append(before, subBefore...)
				after = append(after, subAfter...)
			default:
				panic("possible bug: split subnet but previously contained subnet is not contained by any split parts")
			}
			return
		case mask1 == mask2:
			// they are the same nets (intersection is one of them, there is no before or after)
			intersected = CopySubnet(subnet1)
			return
		default:
			// mask1 cannot be less than mask2 at the same time that subnet1 contains subnet2 ip.
			panic("subnet2 is invalid")
		}

	default:
		// disjoint ips, return nil as intersected and empty sets as before/after
		return
	}
}

func SplitSubnet(subnet *net.IPNet) (left, right *net.IPNet) {
	// Xoring previous mask with new one
	newMask := NewIPMask(subnet.Mask, 1)
	deltaMask := OperateBytes(subnet.Mask, newMask, func(a, b byte) byte { return a ^ b })

	left = &net.IPNet{
		IP:   CopyIP(subnet.IP),
		Mask: NewIPMask(subnet.Mask, 1),
	}
	right = &net.IPNet{
		// IP has a 1 on the last position of its mask (so we need to "or" delta and original IP)
		IP:   OperateBytes(subnet.IP, deltaMask, func(a, b byte) byte { return a | b }),
		Mask: NewIPMask(subnet.Mask, 1),
	}
	return
}

func CopySubnet(subnet *net.IPNet) (subnetCopy *net.IPNet) {
	subnetCopy = &net.IPNet{
		IP:   make(net.IP, len(subnet.IP)),
		Mask: make(net.IPMask, len(subnet.Mask)),
	}
	copy(subnetCopy.IP, subnet.IP)
	copy(subnetCopy.Mask, subnet.Mask)
	return
}

func NewIPMask(oldMask net.IPMask, addBits int) (newMask net.IPMask) {
	ones, bits := oldMask.Size()
	return net.CIDRMask(ones+addBits, bits)
}

// SubnetToIPs transforms a subnet on a list of IPs
func SubnetToIPs(subnet *net.IPNet, ch IPChan) {
	startIP := Start(subnet)
	endIP := End(subnet)
	for curIP := startIP; CompareBytes(curIP, endIP) <= 0; curIP = AddIP(curIP, 1) {
		ch <- curIP
	}
}

// Start returns the first IP on the interval defined by the Subnet
func Start(subnet *net.IPNet) net.IP {
	if subnet == nil {
		return nil
	}
	return subnet.IP.Mask(subnet.Mask)
}

// End returns the last IP on the interval defined by the Subnet
func End(subnet *net.IPNet) net.IP {
	if subnet == nil {
		return nil
	}
	masked := subnet.IP.Mask(subnet.Mask)
	return OperateBytes(masked, subnet.Mask, func(b1, b2 byte) byte { return b1 ^ ^b2 })
}
