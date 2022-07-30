package tool

import (
	"math"
	"strconv"
)

func StringTOInt(number string) int {
	if number == "" {
		return -1
	}
	lenn := len(number)
	var trannum int = 0
	for i := 0; i < lenn; i++ {
		trannum = trannum*10 + int(number[i]-'0')
		//fmt.Println(trannum)
	}
	return trannum
}

func StringToFloat(number string) float32 {
	if number == "" {
		return -1
	}
	var atrannum float32 = 0
	var btrannum float32 = 0
	var i = 0
	for ; number[i] != '.' && i < len(number); i++ {
		if i == len(number) {
			return 0.00
		}
	}

	for j := 0; j < i; j++ {
		atrannum = atrannum*10 + float32(number[j]-'0')
	}
	for j := len(number) - 1; j > i; j-- {
		btrannum = btrannum*0.1 + float32(number[j]-'0')
	}
	return atrannum + btrannum*0.1
}

func Float64ToFloat2(f float64) float64 {
	f1 := math.Trunc(f*1e2+0.5) * 1e-2
	f1Str := strconv.FormatFloat(f1, 'f', 2, 64)
	value, _ := strconv.ParseFloat(f1Str, 64)
	return value
}

func GetInterfaceToInt(t1 interface{}) int {
	var t2 int
	switch t1.(type) {
	case uint:
		t2 = int(t1.(uint))
		break
	case int8:
		t2 = int(t1.(int8))
		break
	case uint8:
		t2 = int(t1.(uint8))
		break
	case int16:
		t2 = int(t1.(int16))
		break
	case uint16:
		t2 = int(t1.(uint16))
		break
	case int32:
		t2 = int(t1.(int32))
		break
	case uint32:
		t2 = int(t1.(uint32))
		break
	case int64:
		t2 = int(t1.(int64))
		break
	case uint64:
		t2 = int(t1.(uint64))
		break
	case float32:
		t2 = int(t1.(float32))
		break
	case float64:
		t2 = int(t1.(float64))
		break
	case string:
		t2, _ = strconv.Atoi(t1.(string))
		break
	default:
		t2 = t1.(int)
		break
	}
	return t2
}
