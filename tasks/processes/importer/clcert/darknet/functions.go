package darknet

import (
	"fmt"
	"github.com/clcert/osr/savers"
	"github.com/clcert/osr/sources"
	"github.com/clcert/osr/tasks"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/pcapgo"
	"github.com/sirupsen/logrus"
	"sync"
)

// worker this assigns a file to each worker jobs is used as queue containing all the files to be processed.
// All the files processed are moved to "scanned" folder.
func worker(id int, wg *sync.WaitGroup, jobs chan sources.Entry, saver savers.Saver, args *tasks.Context) error {
	defer wg.Done()
	for entry := range jobs {
		packetsSeen := NewPacketDictionary(saver)
		packetsSeen.SetArgs(args)
		msg := fmt.Sprintf("Reading file: %s\n", entry.Name())
		packetsSeen.Args.Log.Info(msg)
		if err := readFromFile(entry, packetsSeen, "tcp"); err != nil { // TODO make it configurable
			args.Log.WithFields(logrus.Fields{
				"file":     entry,
				"workerID": id,
			}).Errorf("could not read file: %v", err)
			continue
		}
		args.Log.Infof("Done reading file %s\n", entry.Name())
	}
	args.Log.Info("Done importing all files\n")
	return nil
}

// readFromFile read a entry containing a pcap file (compressed or uncompressed)
// It process the packets in the file using the PacketDict
func readFromFile(f sources.Entry, packetsSeen *PacketDict, filter string) error {
	pcapReader, err := f.Open()
	if err != nil {
		return err
	}
	defer f.Close()
	reader, err := pcapgo.NewReader(pcapReader)
	if err != nil {
		msg := fmt.Sprintf("Couldn't open reader for file %s\n", f.Path())
		packetsSeen.Args.Log.Error(msg)
		return err
	}

	data, ci, err := reader.ReadPacketData()
	if err != nil {
		msg := fmt.Sprintf("Couldn't read first packet from file %s\n", f.Path())
		packetsSeen.Args.Log.Error(msg)
		return err
	}
	packetSource := gopacket.NewPacketSource(reader, reader.LinkType())

	packet := gopacket.NewPacket(data, reader.LinkType(), gopacket.Default)
	m := packet.Metadata()
	m.CaptureInfo = ci
	m.Truncated = m.Truncated || ci.CaptureLength < ci.Length

	err = packetsSeen.ProcessPacket(packet)
	if err != nil {
		msg := fmt.Sprintf("Couldn't process first packet from file %s\n", f.Name())
		packetsSeen.Args.Log.Error(msg)
		return err
	}
	bpfi, err := pcap.NewBPF(reader.LinkType(), ci.CaptureLength, filter)
	if err != nil {
		msg := fmt.Sprintf("Couldn't make BPFFilter for file %s\n", f.Path())
		packetsSeen.Args.Log.Error(msg)
		return err
	}
	for packet := range packetSource.Packets() {
		ci = packet.Metadata().CaptureInfo
		data = packet.Data()
		if len(data) == 0 {
			msg := fmt.Sprintf("Packet Data from file %s is empty\n", f.Name())
			packetsSeen.Args.Log.Error(msg)
			return err
		}
		if bpfi.Matches(ci, data) {
			err = packetsSeen.ProcessPacket(packet)
			if err != nil {
				msg := fmt.Sprintf("Couldn't process packet from file %s\n", f.Name())
				packetsSeen.Args.Log.Error(msg)
				return err
			}
		}
	}
	return nil
}
