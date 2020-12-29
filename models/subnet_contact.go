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
	ContactID string    `pg:",pk,type:varchar(32)"` // Contact ID
	Contact   *Contact  `pg:"rel:has-one"`
	Subnet    net.IPNet `pg:",pk"`
}
