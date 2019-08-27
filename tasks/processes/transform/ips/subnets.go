package ips

import (
	"github.com/clcert/osr/tasks"
	"github.com/clcert/osr/utils"
	"github.com/clcert/osr/utils/ips"
	"io"
	"net"
)

type SubnetStack struct {
	stack []*net.IPNet
}

func (stack *SubnetStack) Push(subnets ...*net.IPNet) {
	stack.stack = append(stack.stack, subnets...)
}
func (stack *SubnetStack) PushReverse(subnets ...*net.IPNet) {
	for i := len(subnets) - 1; i >= 0; i-- {
		stack.Push(subnets[i])
	}
}

func (stack *SubnetStack) Pop() (subnet *net.IPNet) {
	if stack.Len() == 0 {
		return nil
	}
	l := stack.stack[len(stack.stack)-1]
	stack.stack = stack.stack[:len(stack.stack)-1]
	return l
}

func (stack *SubnetStack) Peek() (subnet *net.IPNet) {
	if stack.Len() == 0 {
		return nil
	}
	return stack.stack[len(stack.stack)-1]
}

func (stack *SubnetStack) Len() int {
	return len(stack.stack)
}

func (stack *SubnetStack) Empty() bool {
	return stack.Len() > 0
}

func GetSubnets(csv *utils.HeadedCSV, args *tasks.Args) (ips.SubnetList, error) {
	list := make(ips.SubnetList, 0)
	for {
		line, err := csv.NextRow()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		_, subnet, err := net.ParseCIDR(line["subnet"]) // header should be subnet
		if err != nil {
			args.Log.Error("cannot save ip: %s", err)
			continue
		}
		list = append(list, subnet)
	}
	return list, nil
}


func CompareSubnets(list1, list2 ips.SubnetList) (common, less, more ips.SubnetList) {
	common = make(ips.SubnetList, 0)
	less = make(ips.SubnetList, 0)
	more = make(ips.SubnetList, 0)

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
			less = append(less, ips.CopySubnet(subnet1))
		} else if subnet1 == nil {
			// onlySubnet2 has still IPs, add it to more.
			more = append(more, ips.CopySubnet(subnet2))
		} else {
			thisCommon, thisBefore, thisAfter := ips.IntersectSubnets(subnet1, subnet2)
			if thisCommon != nil {
				// There was intersection, there could be before and after
				common = append(common, thisCommon)
				// Before should be pushed to more or less, depending if thisCommon is equal to subnet1 or subnet2
				// After should be pushed to the other stack
				if ips.CompareBytes(thisCommon.IP, subnet1.IP) == 0 &&
					ips.CompareBytes(thisCommon.Mask, subnet1.Mask) == 0 {
					// Extra is from subnet2 because subnet1 is the intersection
					more = append(more, thisBefore...)
					stack2.PushReverse(thisAfter...)
				} else if ips.CompareBytes(thisCommon.IP, subnet2.IP) == 0 &&
					ips.CompareBytes(thisCommon.Mask, subnet2.Mask) == 0 {
					// Extra is from subnet1 because subnet2 is the intersection
					less = append(less, thisBefore...)
					stack1.PushReverse(thisAfter...)
				}
			} else {
				if ips.CompareBytes(subnet1.IP, subnet2.IP) < 0 {
					less = append(less, subnet1)
					stack2.Push(subnet2)
				} else if ips.CompareBytes(subnet1.IP, subnet2.IP) > 0 {
					more = append(more, subnet2)
					stack1.Push(subnet1)
				}
			}
		}
	}
	return
}


