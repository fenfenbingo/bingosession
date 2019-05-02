package util

import (
	"sync"
	"os"
	"log"
	"strings"
	"time"
	"bytes"
	"math/rand"
	"encoding/hex"
)

var uuid *UUID

var signKey = []byte("bingo")

func init() {
	addr, err := GetMacAddr()
	if err != nil {
		log.Fatal("get mac address failed")
	}
	addr = strings.Replace(addr, ":", "", -1)
	bytes, err := hex.DecodeString(addr)
	if err != nil {
		log.Fatal("get mac address format error")
	}
	uuid = &UUID{
		macAddr: bytes,
		pid:     os.Getpid(),
	}
}

//全局唯一ID生成器
type UUID struct {
	//自增id
	id int

	//MAC地址转二进制格式
	macAddr []byte

	locker sync.Mutex

	pid int
}

func (self *UUID) generateUUID() string {
	self.locker.Lock()
	ts := time.Now().UnixNano()
	self.id++
	id := self.id
	//先解锁再生成UUID，避免资源占用影响并发
	self.locker.Unlock()

	b1 := Int64ToBytes(ts)
	b2 := IntToBytes(id)
	b3 := IntToBytes(self.pid)
	b4 := Int64ToBytes(rand.Int63())
	var buffer bytes.Buffer
	buffer.Write(b1)
	buffer.Write(self.macAddr)
	buffer.Write(b2)
	buffer.Write(b3)
	buffer.Write(b4)
	buffer.Write(signKey)
	return MD5(buffer.Bytes())
}

//生成32位全局ID
func GenerateUUID() string {
	return uuid.generateUUID()
}
