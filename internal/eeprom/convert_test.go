package eeprom

import "testing"

func TestInt8ToSliceByte(t *testing.T) {
	const (
		max int8 = 127
		min int8 = -127
	)
	abs := func(i int8) (uint8, bool) {
		isNegative := i < 0
		if isNegative {
			i = -i
		}
		return uint8(i), isNegative
	}
	for i := min; i < max; i++ {
		value, isNegative := abs(i)
		var signBit uint8 = 0
		if isNegative {
			signBit = 1 << 7
		}
		ans := value + signBit
		answerHigh := (ans & 0xF0) >> 4
		answerLow := ans & 0x0F
		answer := []byte{answerHigh, answerLow}
		result := int8ToSliceByte(i)
		for index, r := range result {
			if r > 0xF {
				t.Errorf("every byte should be less than or equal to 15,index: %d, result %d", index, r)
				return
			}
		}
		if len(result) != 2 {
			t.Error("length result should be 2")
		}
		for j := 0; j < 2; j++ {
			if result[j] != answer[j] {
				t.Logf("signBit: %d, value: %d, i: %d, result: %v, answer: %v", signBit, value, i, result, answer)
				t.Errorf("result doesn't equal answer,result: %d != answer: %d", result[i], answer[i])
				return
			}
		}
	}
}
func TestUint16ToSliceByte(t *testing.T) {
	const (
		max uint16 = 11500
		min uint16 = 0
	)
	for i := min; i < max; i++ {
		result := uint16ToSliceByte(i)
		u, err := sliceByteToUint16(result)
		if err != nil {
			t.Error(err)
		}
		if u != i {
			t.Errorf("u: %d doesn't equal i: %d", u, i)
		}
	}
}
