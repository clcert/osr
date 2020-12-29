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
	TaskID     int          `pg:",use_zero,type:bigint"` // ID of related task
	Task       *Task        `pg:"rel:has-one"`           // Task oObject
	SourceID   DataSourceID `pg:",pk,type:bigint"`       // ID of source
	Source     *Source      `pg:"rel:has-one"`           // Source object
	Count      uint32       `pg:",use_zero,type:bigint"` // Number of received packets
	Time       time.Time    `pg:",pk"`                   // Time of package reception
	Ihl        uint32       `pg:",use_zero"`
	Tos        uint32       `pg:",use_zero"`
	Length     uint32       `pg:",use_zero"`
	Ipid       uint32       `pg:",use_zero"`
	Flags      string       `pg:",use_zero"`
	FragOffset uint32       `pg:",use_zero"`
	TTLMax     uint32       `pg:",use_zero"`
	TTLMin     uint32       `pg:",use_zero"`
	Protocol   string       `pg:",use_zero"`
	IPChecksum uint32       `pg:",use_zero"`
	SrcIP      net.IP       `pg:",pk"`
	SrcPort    uint16       `pg:",pk,type:integer"`
	DstIP      net.IP       `pg:",pk"`
	DstPort    uint16       `pg:",pk,type:integer"`
	Seq        uint64       `pg:",use_zero"`
	Ack        uint64       `pg:",use_zero"`
	DataOffset uint64       `pg:",use_zero"`
	Window     uint32       `pg:",use_zero"`
	Checksum   uint32       `pg:",pk,type:bigint,use_zero"`
	Urgent     uint32       `pg:",use_zero"`
	Fin        bool         `pg:",use_zero"`
	Syn        bool         `pg:",use_zero"`
	Rst        bool         `pg:",use_zero"`
	Psh        bool         `pg:",use_zero"`
	AckFlag    bool         `pg:",use_zero"`
	Urg        bool         `pg:",use_zero"`
	Ece        bool         `pg:",use_zero"`
	Cwr        bool         `pg:",use_zero"`
}

// String shows a text representation of the package.
func (darknetPacket DarknetPacket) String() string {
	return fmt.Sprintf("Address packet<%+v %+v %d %+v %d>\n", darknetPacket.Time, darknetPacket.SrcIP, darknetPacket.SrcPort, darknetPacket.DstIP, darknetPacket.DstPort)
}
