package cluster

import (
	"math"
	"net"
)

type Bridges []Bridge

type Bridge struct {
	ClusterId    string
	Name    string
	Network string
	Gateway string
}

func (b *Bridge) IPRequest(subnet *net.IPNet, pos uint) net.IP {

	retAddr := make([]byte, len(subnet.IP))
	copy(retAddr, subnet.IP)

	mask, _ := subnet.Mask.Size()
	var tb, byteCount, bitCount int
	if subnet.IP.To4() != nil {
		tb = 4
		byteCount = (32 - mask) / 8
		bitCount = (32 - mask) % 8
	} else {
		tb = 16
		byteCount = (128 - mask) / 8
		bitCount = (128 - mask) % 8
	}

	for i := 0; i <= byteCount; i++ {
		maskLen := 0xFF
		if i == byteCount {
			if bitCount != 0 {
				maskLen = int(math.Pow(2, float64(bitCount))) - 1
			} else {
				maskLen = 0
			}
		}
		masked := pos & uint((0xFF&maskLen)<<uint(8*i))
		retAddr[tb-i-1] |= byte(masked >> uint(8*i))
	}
	return net.IP(retAddr)
}
