package utils

import (
	"fmt"
	"math/rand"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBitwiseRange(t *testing.T) {
	// random generate 100 pieces of data.
	for iii := 0; iii < 100; iii++ {
		bitwiseMap := map[int]int{}
		yes := 0
		no := 0

		maxRandNum := 100000
		a := rand.Intn(maxRandNum)
		b := rand.Intn(maxRandNum)
		if a > b {
			a, b = b, a
		}

		bitwiseRange := BitwiseRange(a, b)
		for _, bitwise := range bitwiseRange {
			bitwiseMap[bitwise.num] = bitwise.freeBits
		}

		for i := 0; i <= maxRandNum; i++ {
			curNum := i
			curFreeBits := 0
			for curNum > 0 {
				v, ok := bitwiseMap[curNum<<curFreeBits]
				if ok && v == curFreeBits {
					yes++
					break
				}
				curNum = curNum >> 1
				curFreeBits++
			}
			if _, ok := bitwiseMap[curNum<<curFreeBits]; !ok {
				no++
			}
		}
		assert.Equal(t, b-a+1, yes)
		assert.Equal(t, maxRandNum-b+a, no)
	}
}
