package utils

type Combine struct {
	a          []int
	rawLen     int
	combineLen int
	f          func([]int)
}

func NewCombine(a []int, rawLen int, combineLen int) *Combine {
	return &Combine{a: a, rawLen: rawLen, combineLen: combineLen}
}

func (c Combine) Combine(f func([]int)) {
	if c.combineLen < 1 {
		panic("combineLen must > 0")
	}
	if c.combineLen == 1 {
		for _, v := range c.a {
			f([]int{v})
		}
		return
	}

	c.f = f

	arrLen := len(c.a) - c.combineLen
	for i := 0; i <= arrLen; i++ {
		result := make([]int, c.combineLen)
		result[0] = c.a[i]
		c.doProcess(result, i, 1)
	}
}

func (c Combine) doProcess(result []int, rawIndex int, curIndex int) {
	var choice = c.rawLen - rawIndex + curIndex - c.combineLen
	var tResult []int
	for i := 0; i < choice; i++ {
		if i != 0 {
			tResult := make([]int, c.combineLen)
			copyArr(result, tResult)
		} else {
			tResult = result
		}
		tResult[curIndex] = c.a[i+1+rawIndex]

		if curIndex+1 == c.combineLen {
			c.f(tResult)
			continue
		} else {
			c.doProcess(tResult, rawIndex+i+1, curIndex+1)
		}

	}
}

func copyArr(rawArr []int, target []int) {
	for i := 0; i < len(rawArr); i++ {
		target[i] = rawArr[i]
	}
}
