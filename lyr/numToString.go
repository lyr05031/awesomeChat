package lyr

import (
	"math"
)

func NumToString(num int) string {
	var tmpNum = num
	var length int
	var finalString string
	var numToStringmap = map[int]string{
		0: "0",
		1: "1",
		2: "2",
		3: "3",
		4: "4",
		5: "5",
		6: "6",
		7: "7",
		8: "8",
		9: "9",
	}

	for {
		length += 1
		tmpNum /= 10
		if tmpNum == 0 {
			break
		}
	}

	for {
		length -= 1
		if length < 0 {
			break
		}

		finalString += numToStringmap[num/int(math.Pow(10.0, float64(length)))]
		num %= int(math.Pow(10.0, float64(length)))
	}
	return finalString
}
