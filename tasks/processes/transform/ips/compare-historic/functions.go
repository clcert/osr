package compare_historic

import (
	"fmt"
	"github.com/clcert/osr/tasks/processes/transform/ips"
	"github.com/clcert/osr/utils"
	"strconv"
	"time"
)

func GetPortDateAndIP(m map[string]string, dateFormat string) (port uint16, date time.Time, err error) {
	if m == nil {
		err = fmt.Errorf("map is nil")
		return
	}
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
	if portInt == 0 {
		err = fmt.Errorf("port is zero")
		return
	} else if date.IsZero() {
		err = fmt.Errorf("date is zero")
		return
	}
	port = uint16(portInt)
	return
}

//noinspection ALL
func IPCompareUntilPortAndDate(port uint16, date time.Time, dateFormat string) utils.RowChanCompareFunc {
	return func(map1, map2 map[string]string) (cmp int8, err error) {
		var port1, port2 uint16
		var date1, date2 time.Time
		port1, date1, err1 := GetPortDateAndIP(map1, dateFormat)
		port2, date2, err2 := GetPortDateAndIP(map2, dateFormat)
		if err1 != nil && err2 != nil {
			err = fmt.Errorf("cannot get port and date of neither of both maps")
		} else if err1 == nil && err2 != nil && date1.Equal(date) && port1 == port {
			cmp = -1
		} else if err1 != nil && err2 == nil && date2.Equal(date) && port2 == port {
			cmp = 1
		} else if port1 == port && port2 == port && date1.Equal(date) && date2.Equal(date) { //Both are equal
			cmp, err = ips.Compare(map1, map2)
		} else if port1 == port && date1.Equal(date) {
			cmp = -1
		} else if port2 == port && date2.Equal(date) {
			cmp = 1
		} else {
			err = fmt.Errorf("neither map has the port and the date desired")
		}
		return
	}
}

func GetMinPortAndDate(chan1, chan2 *utils.RowChan) (port uint16, date time.Time, err error) {
	var port1, port2 uint16
	var date1, date2 time.Time
	if !chan1.IsOpen() && !chan2.IsOpen() {
		err = fmt.Errorf("both channels closed")
		return
	}
	if chan1.IsOpen() {
		elem1 := chan1.Peek()
		port1, date1, err = GetPortDateAndIP(elem1, dateFormat)
		if err != nil {
			return
		}
		if !chan2.IsOpen() {
			port = port1
			date = date1
			return
		}
	}
	if chan2.IsOpen() {
		elem2 := chan2.Peek()
		port2, date2, err = GetPortDateAndIP(elem2, dateFormat)
		if err != nil {
			return
		}
		if !chan1.IsOpen() {
			port = port2
			date = date2
			return
		}
	}
	if port1 < port2 {
		port = port1
		date = date1
	} else if port1 > port2 {
		port = port2
		date = date2
	} else if date1.Before(date2) {
		port = port1
		date = date1
	} else if date2.Before(date1) || date2.Equal(date1) {
		port = port2
		date = date2
	} else { // both dates are zero
		err = fmt.Errorf("this should not happen")
	}
	return
}
