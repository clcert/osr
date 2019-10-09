package models

import (
	"net"
)

func init() {
	DefaultModels.Append(SubnetContactModel)
}

var SubnetContactModel = Model{
	Name:                "Subnet to Contact",
	Description:         "Groups the subnets we know some organization controls.",
	StructType:          &SubnetContact{},
}

type SubnetContact struct {
	ContactId string `sql:",pk"`
	Contact   Contact
	Subnet    net.IPNet `sql:",pk"`
}
