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
	thanNet := than.(*IPNetNode)
	cmp := bytes.Compare(node.IP.To4(), thanNet.IP.To4())
	return cmp == -1 && !(node.Contains(thanNet.IP))
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
	found := tree.BTree.Get(item)
	if found == nil {
		return nil, false
	}
	return found.(*IPNetNode).Value, found.(*IPNetNode).Value != nil
}
