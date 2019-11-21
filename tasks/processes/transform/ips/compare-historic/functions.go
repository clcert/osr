package compare_historic

import (
	"fmt"
	"github.com/clcert/osr/utils"
	"strconv"
	"time"
)

func GetPortAndDate(m map[string]string, dateFormat string) (port uint16, date time.Time, err error) {
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

func IPCompareUntilPortAndDate(port uint16, date time.Time, dateFormat string) utils.RowChanCompareFunc {
	return func(map1, map2 map[string]string) (cmp int8, err error) {
		if map1 == nil {
			cmp = 1
			return
		} else if map2 == nil {
			cmp = -1
			return
		}
		port1, date1, err := GetPortAndDate(map1, dateFormat)
		if err != nil {
			return
		}
		port2, date2, err := GetPortAndDate(map2, dateFormat)
		if err != nil {
			return
		}
		if port1 != port && port2 != port { // None of them are of the port asked
			err = fmt.Errorf("both ports are distinct to current port")
		} else if !date1.Equal(date) && !date2.Equal(date) { // None of them are of the date asked
			err = fmt.Errorf("both dates are distinct to current date")
		} else if port1 == port && port2 == port && date1.Equal(date) && date2.Equal(date) { //Both are equal
			cmp = 0
		} else if port1 == port && date1.Equal(date) {
			cmp = -1
		} else if port2 == port && date2.Equal(date) {
			cmp = 1
		} else {
			err = fmt.Errorf("unknown situation")
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
		port1, date1, err = GetPortAndDate(elem1, dateFormat)
		if err != nil {
			return
		}
	}
	if chan2.IsOpen() {
		elem2 := chan2.Peek()
		port2, date2, err = GetPortAndDate(elem2, dateFormat)
		if err != nil {
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
