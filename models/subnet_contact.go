package models

import (
	"net"
)

var SubnetContactModel = Model{
	Name:        "Subnet to Contact",
	Description: "Groups the subnets we know some organization controls.",
	StructType:  &SubnetContact{},
}

type SubnetContact struct {
	ContactID string `sql:",pk,type:varchar(32)"` // Contact ID
	Contact   Contact
	Subnet    net.IPNet `sql:",pk"`
}
