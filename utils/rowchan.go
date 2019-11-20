package utils

import (
	"fmt"
	"io"
)

// Represents a sorted list of rows
// If it is not sorted, everything will fail.
type RowChan struct {
	Tag     string
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

func NewRowChan(tag string) *RowChan {
	return &RowChan{
		Tag: tag,
		ch:  make(chan map[string]string),
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

func (ch1 *RowChan) Join(ch2 *RowChan, cmpFun RowChanCompareFunc) (chUnion *RowChan) {
	chUnion = NewRowChan(fmt.Sprintf("union_%s_%s", ch1.Tag, ch2.Tag))
	go func() {
		for {
			var row map[string]string
			var tag string
			if !ch1.IsOpen() && !ch2.IsOpen() {
				// Both channels are closed, finishing
				break
			} else if !ch2.IsOpen() {
				row = ch1.Get()
				tag = ch1.Tag
			} else if !ch1.IsOpen() {
				row = ch2.Get()
				tag = ch2.Tag
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
					row = ch1.Get()
					tag = ch1.Tag
				case 0:
					row = ch2.Get()
					_ = ch1.Get() // clean this IP
					tag = "both"
				default:
					row = ch2.Get()
					tag = ch2.Tag
				}
			}
			row["tag"] = tag
			chUnion.Put(row)
		}
		chUnion.Close()
	}()
	return
}

func CSVToRowChan(csv *HeadedCSV) *RowChan {
	ch := NewRowChan(csv.Name)
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
