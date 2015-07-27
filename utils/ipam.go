package utils

import (
	"math"
	"net"
   log "code.google.com/p/log4go"
   "github.com/megamsys/megamd/global"
)

func IPRequest(subnet net.IPNet) (net.IP, uint, error) {
	bits := bitCount(subnet)
	bc := int(bits / 8)
	partial := int(math.Mod(bits, float64(8)))
	if partial != 0 {
		bc += 1
	}
	index := global.IPIndex{}
	res, err := index.Get(global.IPINDEXKEY)
	if err != nil {
		log.Error("Error: Riak didn't cooperate:\n%s.", err)
		return nil, 0, err
	}
	
	return getIP(subnet, res.Index+1), res.Index+1, nil
}

// Given Subnet of interest and free bit position, this method returns the corresponding ip address
// This method is functional and tested. Refer to ipam_test.go But can be improved

func getIP(subnet net.IPNet, pos uint) net.IP {
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

func bitCount(addr net.IPNet) float64 {
	mask, _ := addr.Mask.Size()
	if addr.IP.To4() != nil {
		return math.Pow(2, float64(32-mask))
	} else {
		return math.Pow(2, float64(128-mask))
	}
}

func testAndSetBit(a []byte) uint {
	var i uint
	for i = uint(0); i < uint(len(a)*8); i++ {
		if !testBit(a, i) {
			setBit(a, i)
			return i + 1
		}
	}
	return i
}

func testBit(a []byte, k uint) bool {
	return ((a[k/8] & (1 << (k % 8))) != 0)
}

func setBit(a []byte, k uint) {
	a[k/8] |= 1 << (k % 8)
}
