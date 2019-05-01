package util

import (
	"net"
	"crypto/md5"
	"encoding/hex"
	"encoding/binary"
)


// 获取本机的MAC地址
func GetMacAddr() (addr string, err error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	first := ""
	for _, inter := range interfaces {
		name, addr := inter.Name, inter.HardwareAddr.String()
		//以linux系统为标准，返回eth0的MAC地址，如果获取不到eth0则返回第一张网卡MAC地址（其它操作系统可自行调整）
		if name == "eth0" && addr != "" {
			return addr, nil
		}
		if first == "" {
			first = addr
		}
	}
	return first, nil
}



func MD5(data []byte) string  {
	h := md5.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

func Int64ToBytes(i int64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf
}


func Int32ToBytes(i int32) []byte {
	var buf = make([]byte,4)
	binary.BigEndian.PutUint32(buf, uint32(i))
	return buf
}

func IntToBytes(i int) []byte {
	var buf = make([]byte,4)
	binary.BigEndian.PutUint32(buf, uint32(i))
	return buf
}