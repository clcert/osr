	package models

import (
	"github.com/go-pg/pg"
	"net"
	"time"
)

// BlacklistedSubnetModel contains the metainformation related to the respective model.
var BlacklistedSubnetModel = Model{
	Name:                "Blacklisted Sub Networks",
	Description:         "Subnetworks that we are not allowed to check",
	StructType:          &BlacklistedSubnet{},
	AfterCreateFunction: defaultBlacklistSubnets,
	AfterCreateStmts: []string{
		"CREATE INDEX IF NOT EXISTS blacklist_index ON ?TableName USING gist (subnet inet_ops)",
	},
}

// BlacklistedSubnet represents a subnet tthat we must not scan.
type BlacklistedSubnet struct {
	// Task structure
	DateAdded time.Time  `sql:",notnull,default:now()"` // Date of the blacklisted element submission
	Subnet    *net.IPNet `sql:",pk"`                    // IP range blacklisted
	Reason    string     `sql:",type:varchar(255)"`     // Reason for blacklisting
	ContactID string     `sql:",type:varchar(32)"`   // Contact ID
	Contact   *Contact
}

// DefaultBlacklistSubnets inserts the reserved ranges of subnets to the database when it is created.
func defaultBlacklistSubnets(db *pg.DB) error {
	reservedSubnets := [][3]string{
		{"2013-05-22", "0.0.0.0/8", "RFC1122: This host on this network"},
		{"2013-05-22", "10.0.0.0/8", "RFC1918: Private-Use"},
		{"2013-05-22", "100.64.0.0/10", "RFC6598: Shared Access Space"},
		{"2013-05-22", "127.0.0.0/8", "RFC1122: Loopback"},
		{"2013-05-22", "169.254.0.0/16", "RFC3927: Link Local"},
		{"2013-05-22", "172.16.0.0/12", "RFC1918: Private-Use"},
		{"2013-05-22", "192.0.0.0/24", "RFC6890: IETF Protocol Assignments"},
		{"2013-05-22", "192.0.2.0/24", "RFC5737: Documentation (TEST-NET-1)"},
		{"2013-05-22", "192.88.99.0/24", "RFC3068: 6to4 Relay Anycast"},
		{"2013-05-22", "192.168.0.0/16", "RFC1918: Private-Use"},
		{"2013-05-22", "198.18.0.0/15", "RFC2544: Benchmarking"},
		{"2013-05-22", "198.51.100.0/24", "RFC5737: Documentation (TEST-NET-2)"},
		{"2013-05-22", "203.0.113.0/24", "RFC5737: Documentation (TEST-NET-3)"},
		{"2013-05-22", "240.0.0.0/4", "RFC1122: Reserved"},
		{"2013-05-22", "255.255.255.255/32", "RFC919: Limited Broadcast"},
		{"2013-06-25", "224.0.0.0/4", "RFC5771: Multicast/Reserved"},
	}

	subnets := make([]*BlacklistedSubnet, 0)
	for _, subnet := range reservedSubnets {
		subnetDate, err := time.Parse("2006-01-02", subnet[0])
		if err != nil {
			return err
		}
		_, subnetCIDR, err := net.ParseCIDR(subnet[1])
		if err != nil {
			return err
		}
		subnets = append(subnets, &BlacklistedSubnet{
			DateAdded: subnetDate,
			Subnet:    subnetCIDR,
			Reason:    subnet[2],
		})
	}
	_, err := db.Model(&subnets).OnConflict("DO NOTHING").Insert()
	if err != nil {
		return err
	}
	return nil
}
