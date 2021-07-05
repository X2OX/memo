package model

import (
	"crypto/rand"
	"sync/atomic"
	"time"
)

var (
	keyValue atomic.Value
)

func loadKey() {
	UpdateKey()

	if Conf.Token.AutoUpdate != 0 {
		go func() {
			for {
				<-time.NewTimer(time.Duration(Conf.Token.AutoUpdate) * time.Minute).C
				UpdateKey()
			}
		}()
	}
}

func GetKey() [16]byte {
	if k, ok := keyValue.Load().(*[16]byte); ok {
		return *k
	}
	return UpdateKey()
}

func UpdateKey() [16]byte {
	var arr [16]byte
	if _, err := rand.Read(arr[:]); err != nil {
		arr = [16]byte{0x54, 0x68, 0x69, 0x73, 0x27, 0x73, 0x20, 0x70, 0x75, 0x72, 0x65, 0x20, 0x6d, 0x65, 0x6d, 0x6f}
	}
	keyValue.Store(&arr)
	return arr
}
