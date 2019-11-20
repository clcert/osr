package compare_historic

import (
	"fmt"
	"github.com/clcert/osr/tasks/processes/transform/ips"
	"github.com/clcert/osr/utils"
	"strconv"
	"time"
)

func GetPortAndDate(m map[string]string) (port uint16, date time.Time, err error) {
	dateStr, ok := m["date"]
	if !ok {
		err = fmt.Errorf("date key not found on row")
		return
	}
	portStr, ok := m["port_number"]
	if !ok {
		err = fmt.Errorf("date key not found on row")
		return
	}
	date, err = time.Parse(dateFormat, dateStr)
	if err != nil {
		return
	}
	portInt, err := strconv.Atoi(portStr)
	if err != nil {
		return
	}
	port = uint16(portInt)
	return
}

func IPCompareUntilPortAndDate(port uint16, date time.Time, dateFormat string) utils.RowChanCompareFunc {
	return func(map1, map2 map[string]string) (cmp int8, err error) {
		port1, date1, err := GetPortAndDate(map1)
		if err != nil {
			return
		}
		port2, date2, err := GetPortAndDate(map2)
		if err != nil {
			return
		}
		if !(port1 == port || port2 == port) { // None of them are of the port asked
			err = fmt.Errorf("both ports are distinct to current port")
		} else if !(date1.Equal(date) || date2.Equal(date)) { // None of them are of the date asked
			err = fmt.Errorf("both dates are distinct to current port")
		} else if port1 == port && date1.Equal(date) && port2 == port && date2.Equal(date) { // Both are of the port asked
			cmp, err = ips.Compare(map1, map2)
		} else if port1 == port && date1.Equal(date) { // The first one is from the protocol and port asked
			cmp = -1
		} else if port2 == port && date2.Equal(date) { // The seond one is from the port and protocol asked
			cmp = 1
		} else {
			err = fmt.Errorf("unknown error")
		}
		return
	}
}

func getMinDate(chan1, chan2 *utils.RowChan) (date time.Time, err error) {
	var date1, date2 time.Time
	if !chan1.IsOpen() && !chan2.IsOpen() {
		err = fmt.Errorf("both channels closed")
		return
	}
	if chan1.IsOpen() {
		elem1 := chan1.Peek()
		date1, err = time.Parse(dateFormat, elem1["date"])
		if err != nil {
			return
		}
		if date1.IsZero() {
			err = fmt.Errorf("date 1 cannot be zero")
			return
		}
	}
	if chan2.IsOpen() {
		elem2 := chan2.Peek()
		date2, err = time.Parse(dateFormat, elem2["date"])
		if err != nil {
			return
		}
		if date2.IsZero() {
			err = fmt.Errorf("date 2 cannot be zero")
			return
		}
	}
	if !date1.IsZero() && (date1.Before(date2) || date1.Equal(date2)) {
		date = date1
		return
	} else if !date2.IsZero() && date2.Before(date1) {
		date = date2
		return
	} else { // both dates are zero
		err = fmt.Errorf("both dates are zero")
		return
	}
}

func getMinPort(chan1, chan2 *utils.RowChan) (port uint16, err error) {
	var port1, port2 int
	if !chan1.IsOpen() && !chan2.IsOpen() {
		err = fmt.Errorf("both channels closed")
		return
	}
	if chan1.IsOpen() {
		elem1 := chan1.Peek()
		port1, err = strconv.Atoi(elem1["port_number"])
		if err != nil {
			return
		}
		if port1 == 0 {
			err = fmt.Errorf("port 1 cannot be zero")
			return
		}
	}
	if chan2.IsOpen() {
		elem2 := chan2.Peek()
		port2, err = strconv.Atoi(elem2["port_number"])
		if err != nil {
			return
		}
		if port2 == 0 {
			err = fmt.Errorf("port 2 cannot be zero")
			return
		}

	}
	if port1 != 0 && port1 <= port2 {
		port = uint16(port1)
		return
	} else if port2 != 0 && port2 < port1 {
		port2 = port2
		return
	} else { // both ports are zero
		err = fmt.Errorf("both ports are zero")
		return
	}
}
