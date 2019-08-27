package remote

import (
	"bytes"
	"fmt"
	"text/tabwriter"
	"time"
)

// ServerInfo defines the health information associated to a server.
type ServerInfo struct {
	Name    string    // name of the server
	Address string    // Address of the server
	Date    time.Time // Date of info obtained
	Ping    int       // Ping average of the server
	Disks   DiskList  // List of disks on the server
}

// Disk defines the status of a disk in a remote server
type Disk struct {
	Name       string // Filesystem name of the disk
	Blocks     int    // The size in 512-blocks of the disk
	Used       int    // The amount of 512-blocks used of the disk
	Available  int    // Available blocks in disk
	Capacity   int    // Percentage of disk used, as shown by df.
	Mountpoint string // The remote folder where the disk is mounted
}

type DiskList []*Disk

// Returns a text representation of Server Information.
func (info ServerInfo) String() string {
	w := new(bytes.Buffer)
	_, _ = fmt.Fprintf(w, "\nDatos para %s\n", info.Name)
	_, _ = fmt.Fprintf(w, "%s", info.Disks)
	return w.String()
}

// Returns a text table representation of Disk Space Usage.
func (diskList DiskList) String() string {
	b := new(bytes.Buffer)
	w := new(tabwriter.Writer)
	w.Init(b, 5, 8, 1, '\t', 0)
	_, _ = fmt.Fprintf(w, "Sist. Archivos\tTotal (MiB)\tUsado (MiB)\tDisponible (MiB)\tCapacidad\tMontado en\t\n")
	for _, disk := range diskList {
		_, _ = fmt.Fprintf(w, "%s", disk)
	}
	_ = w.Flush()
	return b.String()
}

func (diskList DiskList) getMaxCapacity() int {
	maxCap := 0
	for _, disk := range diskList {
		if disk.Capacity > maxCap {
			maxCap = disk.Capacity
		}
	}
	return maxCap
}

// Returns a text table representation of Disk Space Usage.
func (disk Disk) String() string {
	return fmt.Sprintf("%s\t%d\t%d\t%d\t%d%%\t%s\t\n", disk.Name, disk.Blocks>>11, disk.Used>>11,
		disk.Available>>11, disk.Capacity, disk.Mountpoint)
}
