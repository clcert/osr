package geoasn

import (
	"bytes"
	"github.com/google/btree"
	"net"
)

// IPNetNode represents a subnet node with extra values related to it
type IPNetNode struct {
	*net.IPNet
	Value interface{}
}

// Less allows us to implement btree.Item interface
func (node *IPNetNode) Less(than btree.Item) bool {
	ipCompare := bytes.Compare(node.IP, than.(*IPNetNode).IP)
	switch ipCompare {
	case 0:
		maskCompare := bytes.Compare(node.Mask, than.(*IPNetNode).Mask)
		switch maskCompare {
		case 0:
			return false
		default:
			return maskCompare == -1 // 10.0.0.1/32 > 10.0.0.1/24 => 32 > 24
		}
	default:
		return ipCompare == -1
	}
}

// IPNetTree represents a BTree with IPNets and extra value
type IPNetTree struct {
	*btree.BTree
}

// NewIPNetTree returns a new IPNetTree with degree 5.
func NewIPNetTree() *IPNetTree {
	return &IPNetTree{
		btree.New(5),
	}
}

// AddNode adds a node to the tree. It returns nil if the element is new, or the previous
// element if this operation replaces another element.
func (tree *IPNetTree) AddNode(node *IPNetNode) btree.Item {
	return tree.BTree.ReplaceOrInsert(node)
}

// GetIPData returns the associated value to the smallest subnet the IP is part of.
func (tree *IPNetTree) GetIPData(ip net.IP) (interface{}, bool) {
	item := &IPNetNode{IPNet: &net.IPNet{ip, net.IPv4Mask(255, 255, 255, 255)}}
	var found *IPNetNode
	tree.BTree.DescendLessOrEqual(item, func(i btree.Item) bool {
		if i.(*IPNetNode).IPNet.Contains(item.IP) {
			found = i.(*IPNetNode)
		}
		return false
	})
	return found.Value, found.Value != nil
}
