package models

import (
	"fmt"
	"net"
	"time"
)

// DarknetPacketModel contains the metainformation related to the respective model.
var DarknetPacketModel = Model{
	Name:        "Darknet Packet",
	Description: "Darknet Packet Model",
	StructType:  &DarknetPacket{},
	AfterCreateStmts: []string{
		"CREATE INDEX IF NOT EXISTS darknet_packet_index ON ?TableName USING gist (src_ip inet_ops)",
		"CREATE INDEX IF NOT EXISTS darknet_timestamp ON ?TableName USING btree (time)",
	},
}

// DarknetPacket represents the TCP/IP headers of a darknet packet.
type DarknetPacket struct {
	Hash       string       // Hash of the package, useful for repetition detection
	TaskID     int          `pg:",notnull,type:bigint"` // ID of related task
	Task       *Task        `pg:"rel:has-one"`          // Task oObject
	SourceID   DataSourceID `pg:",pk,type:bigint"`      // ID of source
	Source     *Source      `pg:"rel:has-one"`          // Source object
	Count      uint32       `pg:",notnull,type:bigint"` // Number of received packets
	Time       time.Time    `pg:",pk"`                  // Time of package reception
	Ihl        uint32       `pg:",notnull"`
	Tos        uint32       `pg:",notnull"`
	Length     uint32       `pg:",notnull"`
	Ipid       uint32       `pg:",notnull"`
	Flags      string       `pg:",notnull"`
	FragOffset uint32       `pg:",notnull"`
	TTLMax     uint32       `pg:",notnull"`
	TTLMin     uint32       `pg:",notnull"`
	Protocol   string       `pg:",notnull"`
	IPChecksum uint32       `pg:",notnull"`
	SrcIP      net.IP       `pg:",pk"`
	SrcPort    uint16       `pg:",pk,type:integer"`
	DstIP      net.IP       `pg:",pk"`
	DstPort    uint16       `pg:",pk,type:integer"`
	Seq        uint64       `pg:",notnull"`
	Ack        uint64       `pg:",notnull"`
	DataOffset uint64       `pg:",notnull"`
	Window     uint32       `pg:",notnull"`
	Checksum   uint32       `pg:",pk,type:bigint,notnull"`
	Urgent     uint32       `pg:",notnull"`
	Fin        bool         `pg:",notnull"`
	Syn        bool         `pg:",notnull"`
	Rst        bool         `pg:",notnull"`
	Psh        bool         `pg:",notnull"`
	AckFlag    bool         `pg:",notnull"`
	Urg        bool         `pg:",notnull"`
	Ece        bool         `pg:",notnull"`
	Cwr        bool         `pg:",notnull"`
}

// String shows a text representation of the package.
func (darknetPacket DarknetPacket) String() string {
	return fmt.Sprintf("Address packet<%+v %+v %d %+v %d>\n", darknetPacket.Time, darknetPacket.SrcIP, darknetPacket.SrcPort, darknetPacket.DstIP, darknetPacket.DstPort)
}
