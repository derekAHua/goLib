package utils

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

var a = []int{0, 1, 2}
var rawLen int
var combineLen int

func TestCom(t *testing.T) {
	rawLen = len(a)
	combineLen = 2
	combine()
}

func combine() {
	fmt.Println("start combine")
	arrLen := len(a) - combineLen
	for i := 0; i <= arrLen; i++ {
		result := make([]int, combineLen)
		result[0] = a[i]
		doProcess(result, i, 1)
	}
}

func doProcess(result []int, rawIndex int, curIndex int) {
	var choice = rawLen - rawIndex + curIndex - combineLen

	var tResult []int
	for i := 0; i < choice; i++ {
		if i != 0 {
			tResult := make([]int, combineLen)
			copyArr(result, tResult)
		} else {
			tResult = result
		}

		tResult[curIndex] = a[i+1+rawIndex]

		if curIndex+1 == combineLen {
			PrintIntArr(tResult)
			continue
		} else {
			doProcess(tResult, rawIndex+i+1, curIndex+1)
		}

	}
}

func PrintIntArr(arr []int) {
	var valuesText []string
	for i := range arr {
		number := arr[i]
		text := strconv.Itoa(number)
		valuesText = append(valuesText, text)
	}

	fmt.Println(strings.Join(valuesText, ","))
}

func TestMyCom(t *testing.T) {
	contentTypeList := []int{0, 1, 2, 3, 4, 5}

	ret := make([]interface{}, 0)

	for i := 1; i <= len(contentTypeList); i++ {
		c := NewCombine(contentTypeList, len(contentTypeList), i)
		c.Combine(func(ints []int) {
			ret = append(ret, ints)
		})
	}

	t.Log(ret)
	t.Log(len(ret))
}
