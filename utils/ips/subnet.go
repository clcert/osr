package ips

import (
	"github.com/clcert/osr/utils"
	"net"
)

// Represents a sorted list of Subnets
// If it is not sorted, everything will fail.
type SubnetList []*net.IPNet

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

func (list1 SubnetList) Compare(list2 SubnetList) (common, less, more SubnetList) {
	common, less, more = make(SubnetList, 0), make(SubnetList, 0), make(SubnetList, 0)
	stack1 := &SubnetStack{
		stack: make([]*net.IPNet, 0),
	}
	stack1.PushReverse(list1...)

	stack2 := &SubnetStack{
		stack: make([]*net.IPNet, 0),
	}
	stack2.PushReverse(list2...)

	for {
		subnet1 := stack1.Pop()
		subnet2 := stack2.Pop()
		if subnet1 == nil && subnet2 == nil {
			break
		} else if subnet2 == nil {
			// only subnet1 has still IPs, add it to less.
			less = append(less, CopySubnet(subnet1))
		} else if subnet1 == nil {
			// onlySubnet2 has still IPs, add it to more.
			more = append(more, CopySubnet(subnet2))
		} else {
			thisCommon, thisBefore, thisAfter := IntersectSubnets(subnet1, subnet2)
			if thisCommon != nil {
				// There was intersection, there could be before and after
				common = append(common, thisCommon)
				// Before should be pushed to more or less, depending if thisCommon is equal to subnet1 or subnet2
				// After should be pushed to the other stack
				if CompareBytes(thisCommon.IP, subnet1.IP) == 0 &&
					CompareBytes(thisCommon.Mask, subnet1.Mask) == 0 {
					// Extra is from subnet2 because subnet1 is the intersection
					more = append(more, thisBefore...)
					stack2.PushReverse(thisAfter...)
				} else if CompareBytes(thisCommon.IP, subnet2.IP) == 0 &&
					CompareBytes(thisCommon.Mask, subnet2.Mask) == 0 {
					// Extra is from subnet1 because subnet2 is the intersection
					less = append(less, thisBefore...)
					stack1.PushReverse(thisAfter...)
				}
			} else {
				if CompareBytes(subnet1.IP, subnet2.IP) < 0 {
					less = append(less, subnet1)
					stack2.Push(subnet2)
				} else if CompareBytes(subnet1.IP, subnet2.IP) > 0 {
					more = append(more, subnet2)
					stack1.Push(subnet1)
				}
			}
		}
	}
	return
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

func (list SubnetList) ToRowChan(tag string) *utils.RowChan {
	rowChan := utils.NewRowChan(tag)
	go func() {
		for _, subnet := range list {
			startIP := Start(subnet)
			endIP := End(subnet)
			for curIP := startIP; CompareBytes(curIP, endIP) <= 0; curIP = AddIP(curIP, 1) {
				rowChan.Put(map[string]string{"ip": curIP.String()})
			}
		}
		rowChan.Close()
	}()
	return rowChan
}
