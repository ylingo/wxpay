package wxpay

import (
	"math/rand"
)

type RANDTYPE int

const (
	LOWER RANDTYPE = iota
	UPPER
	NUMBER
	SPECIALCHAR
)

type RandStr struct {
	randChan  chan string
	size      int
	randTypes []RANDTYPE
}

func NewRandStr(length int, rType []RANDTYPE) *RandStr {
	ret := &RandStr{
		randChan:  make(chan string, 10),
		size:      length,
		randTypes: rType,
	}
	go ret.creatRandStr()
	return ret
}

func (r RandStr) GetRandString() string {
	return <-r.randChan
}

func (r *RandStr) creatRandStr() {
	for {
		result := make([]byte, r.size)
		kinds := [][]int{[]int{97, 26}, []int{65, 26}, []int{48, 10}, []int{33, 15}, []int{58, 7}, []int{91, 6}, []int{123, 4}}

		//rand.Seed(<-r.timeStamp)
		rand.Read(result)
		for i := 0; i < r.size; i++ {
			kind := r.randTypes[rand.Intn(len(r.randTypes))]
			if kind == SPECIALCHAR {
				kind = (RANDTYPE)(3 + rand.Intn(4)) //取特殊字符
			}
			ikind := int(kind)
			base, scope := kinds[ikind][0], kinds[ikind][1]
			result[i] = uint8(base + rand.Intn(scope))
		}
		r.randChan <- string(result)
	}
}
