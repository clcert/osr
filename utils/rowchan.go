package utils

import (
	"io"
)

// Represents a sorted list of rows
// If it is not sorted, everything will fail.
type RowChan struct {
	ch      chan map[string]string
	nextVal map[string]string
	isOpen  bool
}

type RowChanCompareFunc func(map[string]string, map[string]string) (int8, error)

func (rch *RowChan) get() {
	val, ok := <-rch.ch
	rch.nextVal = val
	rch.isOpen = ok
}

func NewRowChan() *RowChan {
	return &RowChan{
		ch: make(chan map[string]string),
	}
}

func (rch *RowChan) IsOpen() bool {
	if rch.nextVal == nil {
		rch.get()
	}
	return rch.isOpen
}

func (rch *RowChan) Peek() map[string]string {
	if rch.nextVal == nil {
		rch.get()
	}
	return rch.nextVal
}

func (rch *RowChan) Get() map[string]string {
	if rch.nextVal == nil {
		rch.get()
	}
	v := rch.nextVal
	rch.get()
	return v
}

func (rch *RowChan) Put(m map[string]string) {
	rch.ch <- m
}

func (rch *RowChan) Count() int {
	count := 0
	for rch.IsOpen() {
		_ = rch.Get()
		count++
	}
	return count
}

func (rch *RowChan) Close() {
	close(rch.ch)
}

func (ch1 *RowChan) Compare(ch2 *RowChan, cmpFun RowChanCompareFunc) (chBoth, ch1Uniq, ch2Uniq *RowChan) {
	chBoth = NewRowChan()
	ch2Uniq = NewRowChan()
	ch1Uniq = NewRowChan()
	go func() {
		for {
			if !ch1.IsOpen() && !ch2.IsOpen() {
				// Both channels are closed, finishing
				break
			} else if !ch2.IsOpen() {
				ch1Uniq.Put(ch1.Get())
			} else if !ch1.IsOpen() {
				ch2Uniq.Put(ch2.Get())
			} else {
				map1 := ch1.Peek()
				map2 := ch2.Peek()
				cmp, err := cmpFun(map1, map2)
				if err != nil {
					// In case of error, finish all
					break
				}
				switch cmp {
				case -1:
					ch1Uniq.Put(ch1.Get())
				case 0:
					chBoth.Put(ch2.Get())
				case 1:
					ch2Uniq.Put(ch2.Get())
				}
			}
		}
		chBoth.Close()
		ch2Uniq.Close()
		ch1Uniq.Close()
	}()
	return
}

func CSVToRowChan(csv *HeadedCSV) *RowChan {
	ch := NewRowChan()
	go func() {
		for {
			next, err := csv.NextRow()
			if err != nil {
				if err == io.EOF {
					break
				}
			}
			ch.Put(next)
		}
		ch.Close()
	}()
	return ch
}
