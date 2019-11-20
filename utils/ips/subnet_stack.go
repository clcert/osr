package ips

import (
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



