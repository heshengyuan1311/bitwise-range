package utils

import (
"math"
"net"
"strconv"
"strings"
)

type res struct {
	num      int
	freeBits int
}

func BitwiseRange(a, b int) []res {
	// make sure that a < b.
	if a > b {
		a, b = b, a
	}

	// when a == b, only one result in the range list.
	if a == b {
		return []res{res{
			num:      a,
			freeBits: 0,
		}}
	}
	// init var
	var firstDiffIndex int
	var aTop, bBottom int64
	totalBits := len(strconv.FormatInt(int64(b), 2))
	result := []res{}

	var appendRes func(int, int)
	appendRes = func(nums, freeBits int) {
		result = append(result, res{num: nums, freeBits: freeBits})
	}

	// trans int into string binary.
	var strr func(int) string
	strr = func(n int) string {
		res := strconv.FormatInt(int64(n), 2)
		curBits := len(res)
		if len(res) < totalBits {
			for i := 0; i < totalBits-curBits; i++ {
				res = "0" + res
			}
		}
		return res
	}

	sa := strr(a)
	sb := strr(b)

	// Handle A -> a_top.
	// a_top: the bit on first_diff_index is 0, all bits after first_diff_index is 1.
	var part1 func()
	part1 = func() {
		aa := a
		ss := sa

		// special: A == 0
		if aa == 0 {
			appendRes(aa, totalBits-firstDiffIndex-1)
			return
		}

		// NOTE: can be handled in common loop, but this is OK.
		if strings.LastIndex(ss, "1") > 0 {
			result = append(result, res{num: aa, freeBits: 0})
			aa += 1
			ss = strr(aa)
		}

		// loop: look from low bit to high bit
		for {
			if aa > int(aTop) {
				break
			}

			// find tail '0's
			i := strings.LastIndex(ss, "1")
			if i < 0 || i < firstDiffIndex {
				i = firstDiffIndex
			}

			freeBits := totalBits - i - 1
			appendRes(aa, freeBits)
			aa = aa + int(math.Pow(2, float64(freeBits)))
			ss = strr(aa)
		}
	}

	// Handle b_bottom -> B.
	// b_bottom: the bit on first_diff_index is 1, all bits after first_diff_index is 0.
	var part2 func()
	part2 = func() {
		bb := int(bBottom)

		index := firstDiffIndex + 1
		for {
			if index >= totalBits {
				break
			}

			if !strings.Contains(sb[index:], "0") {
				appendRes(bb, totalBits-index)
				break
			}

			if !strings.Contains(sb[index:], "1") {
				appendRes(bb, 0)
				break
			}

			if sb[index] == byte('1') {
				appendRes(bb, totalBits-index-1)
				bb = bb + int(math.Pow(2, float64(totalBits-index-1)))
				index += 1
			} else {
				index += 1
			}
		}
	}

	var handle func()
	handle = func() {
		// get the first different index between a and b.
		if len(sa) < len(sb) {
			for i := 0; i < len(sb)-len(sa); i++ {
				sa = "0" + sa
			}
		}

		for i := 0; i < totalBits; i++ {
			if sa[i] != sb[i] {
				firstDiffIndex = i
				break
			}
		}

		// only the last (lowest) bit is different. 2 result: A, B
		if firstDiffIndex == totalBits-1 {
			appendRes(a, 0)
			appendRes(b, 0)
			return
		}

		// from firstDiffIndex, A is all 0 and B is all 1. a big result
		if !strings.Contains(sa[firstDiffIndex:], "1") && !strings.Contains(sb[firstDiffIndex:], "0") {
			appendRes(a, totalBits-firstDiffIndex)
			return
		}

		// split into 2 parts: A -> a_top, b_bottom -> B
		//   a_top: the bit on first_diff_index is 0, all bits after first_diff_index is 1.
		//   b_bottom: the bit on first_diff_index is 1, all bits after first_diff_index is 0.
		saTop := sa[:firstDiffIndex+1]
		for i := 0; i < (totalBits - firstDiffIndex - 1); i++ {
			saTop += "1"
		}
		aTop, _ = strconv.ParseInt(saTop, 2, 64)
		bBottom = aTop + 1
		part1()
		part2()
	}

	handle()
	return result
}


func BitwiseIPV4Range(ipA, ipB net.IP) []net.IPNet {
	var ipv4ToInt func(net.IP) int
	result := []net.IPNet{}
	ipv4ToInt = func(ipv4 net.IP) int {
		ipv4 = ipv4.To4()
		num := (int(ipv4[0]) << 24) + (int(ipv4[1]) << 16) + (int(ipv4[2]) << 8) + (int(ipv4[3]))
		return num
	}

	var IntToIpv4 func(int) net.IP
	IntToIpv4 = func(num int) net.IP {
		strconv.FormatInt(int64(num), 2)
		ip4 := byte(num & 0xff)
		ip3 := byte((num >> 8) & 0xff)
		ip2 := byte((num >> 16) & 0xff)
		ip1 := byte((num >> 24) & 0xff)
		return net.IPv4(ip1, ip2, ip3, ip4).To4()
	}

	ipIntA := ipv4ToInt(ipA)
	ipIntB := ipv4ToInt(ipB)

	bitWiseRange := BitwiseRange(ipIntA, ipIntB)
	for _, v := range bitWiseRange {
		ipv4 := IntToIpv4(v.num)
		ipv4Mask := net.CIDRMask(8*net.IPv4len-v.freeBits, 8*net.IPv4len)
		result = append(result, net.IPNet{IP: ipv4, Mask: ipv4Mask})
	}
	return result
}
