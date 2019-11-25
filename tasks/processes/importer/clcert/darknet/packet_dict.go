package darknet

import (
	"crypto"
	"fmt"
	"github.com/clcert/osr/models"
	"github.com/clcert/osr/savers"
	"github.com/clcert/osr/tasks"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/sirupsen/logrus"
	"net"
	"strconv"
	"strings"
	"time"
)

// PacketDict is a useful struct that automatically stores all the packets received
// within 1 second of arrival.
type PacketDict struct {
	count       map[string]int      // Number of packages
	ttlmin      map[string]int      // Minimum value of ttl field
	ttlmax      map[string]int      // Maximum value of ttl field
	packets     map[string][]string // List of packets
	arrivalTime int64               // Time of arrival of packets
	importer    savers.Saver        // Channel to DB
	Args        *tasks.Context      // Related Task struct
}

// Reset sets to default values the fields of the struct except for the csvWriter and the channel
// for writing to the DB. For setting the channel use SetChannelToDB() instead.
func (packetsSeen *PacketDict) Reset() {
	packetsSeen.count = make(map[string]int)
	packetsSeen.ttlmin = make(map[string]int)
	packetsSeen.ttlmax = make(map[string]int)
	packetsSeen.packets = make(map[string][]string)
	packetsSeen.arrivalTime = 0
}

// NewPacketDictionary creates a new PacketDict with a defined importer.
// and writer for the CSVWriter.
func NewPacketDictionary(dbImporter savers.Saver) *PacketDict {
	packetsSeen := new(PacketDict)
	packetsSeen.Reset()
	packetsSeen.importer = dbImporter
	return packetsSeen
}

//SetArgs sets the *tasks.Context
func (packetsSeen *PacketDict) SetArgs(args *tasks.Context) {
	packetsSeen.Args = args
}

// tomodels.DarknetPackets returns a slice of packets containing the data from the packets received within 1 second.
func (packetsSeen *PacketDict) toDBPacket() []*models.DarknetPacket {
	packs := make([]*models.DarknetPacket, 0)
	count := packetsSeen.count
	ttlmin := packetsSeen.ttlmin
	ttlmax := packetsSeen.ttlmax
	packets := packetsSeen.packets
	arrivalTime := packetsSeen.arrivalTime
	for k, v := range count {

		ipData := packets[k]

		ihl, _ := strconv.ParseInt(ipData[0], 10, 32)
		tos, _ := strconv.ParseInt(ipData[1], 10, 32)
		length, _ := strconv.ParseInt(ipData[2], 10, 32)
		ipid, _ := strconv.ParseInt(ipData[3], 10, 32)
		flags := fmt.Sprintf("%s", ipData[4])
		fragoffset, _ := strconv.ParseInt(ipData[5], 10, 32)

		protocol := fmt.Sprintf("%s", ipData[7])

		ipchecksum, _ := strconv.ParseInt(ipData[8], 10, 32)
		srcip := net.ParseIP(ipData[9])
		dstip := net.ParseIP(ipData[10])
		tcppacket := strings.Split(fmt.Sprintf("%s", k), "\t")
		if len(tcppacket) < 16 {
			packetsSeen.Args.Log.WithFields(logrus.Fields{
				"tcppacket": tcppacket,
			}).Error("Packet length is less than 16")
			continue
		}
		srcport, _ := strconv.ParseUint(tcppacket[0], 10, 32)
		dstport, _ := strconv.ParseUint(tcppacket[1], 10, 32)
		seq, _ := strconv.ParseUint(tcppacket[2], 10, 64)
		ack, _ := strconv.ParseUint(tcppacket[3], 10, 64)
		dataoffset, _ := strconv.ParseUint(tcppacket[4], 10, 64)
		window, _ := strconv.ParseUint(tcppacket[5], 10, 32)
		checksum, _ := strconv.ParseUint(tcppacket[6], 10, 32)
		urgent, _ := strconv.ParseUint(tcppacket[7], 10, 32)
		fin, _ := strconv.ParseBool(fmt.Sprintf("%s", tcppacket[8]))
		syn, _ := strconv.ParseBool(fmt.Sprintf("%s", tcppacket[9]))
		rst, _ := strconv.ParseBool(fmt.Sprintf("%s", tcppacket[10]))
		psh, _ := strconv.ParseBool(fmt.Sprintf("%s", tcppacket[11]))
		ackFlag, _ := strconv.ParseBool(fmt.Sprintf("%s", tcppacket[12]))
		urg, _ := strconv.ParseBool(fmt.Sprintf("%s", tcppacket[13]))
		ece, _ := strconv.ParseBool(fmt.Sprintf("%s", tcppacket[14]))
		cwr, _ := strconv.ParseBool(fmt.Sprintf("%s", tcppacket[15]))

		pack := &models.DarknetPacket{
			TaskID:     packetsSeen.Args.GetTaskID(),
			SourceID:   packetsSeen.Args.GetSourceID(),
			Count:      uint32(v),
			Time:       time.Unix(arrivalTime, 0),
			Ihl:        uint32(ihl),
			Tos:        uint32(tos),
			Length:     uint32(length),
			Ipid:       uint32(ipid),
			Flags:      flags,
			FragOffset: uint32(fragoffset),
			TTLMin:     uint32(ttlmin[k]),
			TTLMax:     uint32(ttlmax[k]),
			Protocol:   protocol,
			IPChecksum: uint32(ipchecksum),
			SrcIP:      srcip,
			SrcPort:    uint16(srcport),
			DstIP:      dstip,
			DstPort:    uint16(dstport),
			Seq:        seq,
			Ack:        ack,
			DataOffset: dataoffset,
			Window:     uint32(window),
			Checksum:   uint32(checksum),
			Urgent:     uint32(urgent),
			Fin:        fin,
			Syn:        syn,
			Rst:        rst,
			Psh:        psh,
			AckFlag:    ackFlag,
			Urg:        urg,
			Ece:        ece,
			Cwr:        cwr,
		}
		hash := crypto.SHA256.New()
		packStr := []byte(fmt.Sprintf("%+v", pack))
		_, _ = hash.Write(packStr)
		pack.Hash = fmt.Sprintf("%x", hash.Sum(nil))[:8]
		packs = append(packs, pack)
	}
	return packs
}

// writeToDB sends all the captured data to the db channel.
func (packetsSeen *PacketDict) writeToDB() {
	dbPackets := packetsSeen.toDBPacket()
	for _, packet := range dbPackets {
		err := packetsSeen.importer.Save(packet)
		if err != nil {
			// TODO: log this
		}
	}
}

//checkWrite checks if the packetDictionary should print the data contained within.
//if the arrivalTime is 0, it means it was resetted so the functions sets it to the arrivalTime.
//if the arrivalTime is correctly setup it checks if it should write the data to the db instead.
//also cleans the structure when there are no more packets to write.
func (packetsSeen *PacketDict) checkWrite(arrivalTime int64) {
	if packetsSeen.arrivalTime == 0 {
		packetsSeen.arrivalTime = arrivalTime
	} else if packetsSeen.arrivalTime < arrivalTime {
		packetsSeen.writeToDB()
		packetsSeen.Reset()
		packetsSeen.arrivalTime = arrivalTime
	}
}

// addPacket takes the header of the packet (tipically the transport layer), the ipData and the current ttl
// Adds the given packet to the dictionary and also updates the current min and max for the ttl of that packet.
func (packetsSeen *PacketDict) addPacket(header []string, ipData []string, ttl int) {
	joined := strings.Join(header, "\t")
	if _, ok := packetsSeen.packets[joined]; ok {
		if ttl < packetsSeen.ttlmin[joined] {
			packetsSeen.ttlmin[joined] = ttl
		}
		if ttl > packetsSeen.ttlmax[joined] {
			packetsSeen.ttlmax[joined] = ttl
		}
		packetsSeen.count[joined] += 1
	} else {
		packetsSeen.packets[joined] = ipData
		packetsSeen.ttlmin[joined] = ttl
		packetsSeen.ttlmax[joined] = ttl
		packetsSeen.count[joined] = 1
	}
}

// ProcessPacket takes a raw gopacket.
// It parses the gopacket and then adds it to the dictionary
// if an error occurs it prints an error
func (packetsSeen *PacketDict) ProcessPacket(packet gopacket.Packet) error {

	arrivalTime := packet.Metadata().Timestamp.Unix()
	packetsSeen.checkWrite(arrivalTime)

	ipData, ttl := processIPPacket(packet)

	header := processTCPPacket(packet)
	packetsSeen.addPacket(header, ipData, ttl)

	if err := packet.ErrorLayer(); err != nil {
		return err.Error()
	}
	return nil
}

// processIPPacket takes a raw gopacket and tries to decode it
// returns a slice of string containing the fields of the ip layer and the ttl.
func processIPPacket(packet gopacket.Packet) (ipData []string, ttl int) {
	ipLayer := packet.Layer(layers.LayerTypeIPv4)
	if ipLayer != nil {
		//fmt.Println("IPv4 layer detected.")
		ip, _ := ipLayer.(*layers.IPv4)

		// Address layer variables:
		// Version (Either 4 or 6)
		// IHL (Address Header Length in 32-bit words)
		// TOS, Length, Id, Flags, FragOffset, TTL, Protocol (TCP?),
		// Checksum, IP, DstIP
		ttl = int(ip.TTL)
		ipData = ParseIP(ip)

	}
	return
}

// processTCPPacket takes a raw gopacket and tries to decode it
// returns a slice of string containing the fields of the transport layer (this should work only for TCP).
func processTCPPacket(packet gopacket.Packet) (header []string) {
	// Let's see if the packet is TCP
	tcpLayer := packet.Layer(layers.LayerTypeTCP)

	if tcpLayer != nil {
		//fmt.Println("TCP layer detected.")
		tcp, _ := tcpLayer.(*layers.TCP)

		// TCP layer variables:
		// SrcPort, DstPort, Seq, Ack, DataOffset, Window, Checksum, Urgent
		// Bool flags: FIN, SYN, RST, PSH, ACK, URG, ECE, CWR, NS

		header = ParseTCP(tcp)
		return

	}
	return
}

//ParseIP takes an *layers.IPv4 and returns a formatted slice of string containing the fields of the Address packet.
func ParseIP(ip *layers.IPv4) []string {
	return []string{fmt.Sprintf("%d", ip.IHL), fmt.Sprintf("%d", ip.TOS), fmt.Sprintf("%d", ip.Length), fmt.Sprintf("%d", ip.Id),
		fmt.Sprintf("%v", ip.Flags), fmt.Sprintf("%v", ip.FragOffset), fmt.Sprintf("%v", ip.TTL),
		fmt.Sprintf("%v", ip.Protocol), fmt.Sprintf("%v", ip.Checksum), fmt.Sprintf("%v", ip.SrcIP), fmt.Sprintf("%v", ip.DstIP)}
}

//ParseTCP takes an *layers.TCP and returns a formatted slice of string containing the fields of the TCP packets.
func ParseTCP(tcp *layers.TCP) []string {

	return []string{fmt.Sprintf("%d", tcp.SrcPort), fmt.Sprintf("%d", tcp.DstPort),
		fmt.Sprintf("%v", tcp.Seq), fmt.Sprintf("%v", tcp.Ack),
		fmt.Sprintf("%v", tcp.DataOffset), fmt.Sprintf("%v", tcp.Window),
		fmt.Sprintf("%v", tcp.Checksum), fmt.Sprintf("%v", tcp.Urgent),
		fmt.Sprintf("%v", tcp.FIN), fmt.Sprintf("%v", tcp.SYN),
		fmt.Sprintf("%v", tcp.RST), fmt.Sprintf("%v", tcp.PSH),
		fmt.Sprintf("%v", tcp.ACK), fmt.Sprintf("%v", tcp.URG),
		fmt.Sprintf("%v", tcp.ECE), fmt.Sprintf("%v", tcp.CWR)}
}
