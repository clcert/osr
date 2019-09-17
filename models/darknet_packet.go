package models

import (
	"fmt"
	"net"
	"time"
)

func init() {
	DefaultModels.Append(DarknetPacketModel)
}

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
	Hash       string // Hash of the package, useful for repetition detection
	TaskID     int `sql:",notnull,type:bigint"` // ID of related task
	Task       *Task // Task oObject
	SourceID   DataSourceID `sql:",pk,type:bigint"` // ID of source
	Source     *Source // Source object
	Count      uint32    `sql:",notnull,type:bigint"` // Number of received packets
	Time       time.Time `sql:",pk"` // Time of package reception
	Ihl        uint32    `sql:",notnull"`
	Tos        uint32    `sql:",notnull"`
	Length     uint32    `sql:",notnull"`
	Ipid       uint32    `sql:",notnull"`
	Flags      string    `sql:",notnull"`
	FragOffset uint32    `sql:",notnull"`
	TTLMax     uint32    `sql:",notnull"`
	TTLMin     uint32    `sql:",notnull"`
	Protocol   string    `sql:",notnull"`
	IPChecksum uint32    `sql:",notnull"`
	SrcIP      net.IP    `sql:",pk"`
	SrcPort    uint16    `sql:",pk,type:integer"`
	DstIP      net.IP    `sql:",pk"`
	DstPort    uint16    `sql:",pk,type:integer"`
	Seq        uint64    `sql:",notnull"`
	Ack        uint64    `sql:",notnull"`
	DataOffset uint64    `sql:",notnull"`
	Window     uint32    `sql:",notnull"`
	Checksum   uint32    `sql:",pk,type:bigint,notnull"`
	Urgent     uint32    `sql:",notnull"`
	Fin        bool      `sql:",notnull"`
	Syn        bool      `sql:",notnull"`
	Rst        bool      `sql:",notnull"`
	Psh        bool      `sql:",notnull"`
	AckFlag    bool      `sql:",notnull"`
	Urg        bool      `sql:",notnull"`
	Ece        bool      `sql:",notnull"`
	Cwr        bool      `sql:",notnull"`
}

// String shows a text representation of the package.
func (darknetPacket DarknetPacket) String() string {
	return fmt.Sprintf("Address packet<%+v %+v %d %+v %d>\n", darknetPacket.Time, darknetPacket.SrcIP, darknetPacket.SrcPort, darknetPacket.DstIP, darknetPacket.DstPort)
}
