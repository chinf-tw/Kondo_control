package eeprom

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/pkg/errors"
)

func printHex(bs []byte) string {
	sum := ""
	for _, b := range bs {
		sum += fmt.Sprintf("%X", b) + " "
	}
	return sum
}
func Test(t *testing.T) {
	testingFilePath := "./data"
	dat, err := ioutil.ReadFile(testingFilePath)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(printHex(dat)))
	eeprom, err := Parse(dat)
	if err != nil {
		t.Fatalf("%+v\n", err)
	}
	j, err := json.Marshal(eeprom.Address)
	if err != nil {
		t.Fatalf("json.Marshal: %v", err)
	}
	t.Log(string(j))
}

func TestEEPROM(t *testing.T) {
	testingFilePath := "./data"
	dat, err := ioutil.ReadFile(testingFilePath)
	if err != nil {
		t.Error(err)
	}
	eeprom, err := Parse(dat)
	if err != nil {
		t.Fatalf("%+v\n", err)
	}
	t.Logf("%+v\n", eeprom)
}
func TestCompose(t *testing.T) {
	var testingFilePath string = "./data"
	dat, err := ioutil.ReadFile(testingFilePath)
	if err != nil {
		t.Fatal(err)
	}
	// checking data is legal data and
	ee, err := Parse(dat)
	if err != nil {
		t.Fatalf("%+v\n", err)
	}

	checkFunc := func(source []byte, targetEEPROM EEPROM) error {
		// get source
		sourceData := make([]byte, len(source))
		copy(sourceData, source)

		// log out before data
		// t.Logf("*** Before Data ***\n%+v\n", printHex(sourceData))
		// Compose data
		composeData, err := Compose(sourceData, targetEEPROM)
		if err != nil {
			return err
		}
		// log out after data
		// t.Logf("*** After Data ***\n%+v\n", printHex(composeData))
		// checking compose data equal targetEEPROM
		composeEEPROM, err := Parse(composeData)
		if err != nil {
			return err
		}

		// t.Logf("%+v\n", composeEEPROM)
		if composeEEPROM != targetEEPROM {
			return errors.Errorf("not equal\ncomposeEEPROM:%+v\ntargetEEPROM:%+v", composeEEPROM, targetEEPROM)
		}
		if len(sourceData) != len(composeData) {
			return errors.New("length of compose data doesn't equal to length of source data")
		}

		// checking "Must not be changed" data is not changed
		if !bytes.Equal(sourceData[24:26], composeData[24:26]) {
			return errors.Errorf(`"Must not be changed" data [24:26] is changed, source: %v, compose: %v`, sourceData[24:26], composeData[24:26])
		}
		if !bytes.Equal(sourceData[32:50], composeData[32:50]) {
			return errors.Errorf(`"Must not be changed" data [32:50] is changed, source: %v, compose: %v`, sourceData[32:50], composeData[32:50])
		}
		return nil
	}
	t.Run("StretchGain", func(t *testing.T) {
		targetEEPROM := ee
		for i := uint8(0); i <= 100; i++ {
			targetEEPROM.StretchGain = i
			if i%2 == 1 {
				if err := checkFunc(dat, targetEEPROM); err == nil {
					t.Error("This should be fail, becouse StretchGain should be even number")
				}
			} else {
				if err := checkFunc(dat, targetEEPROM); err != nil {
					t.Error(err)
				}
			}
		}
		targetEEPROM.StretchGain = 101
		if err := checkFunc(dat, targetEEPROM); err == nil {
			t.Error("This should be fail, becouse StretchGain should be less than or equal to 100")
		}
	})
	t.Run("Speed", func(t *testing.T) {
		targetEEPROM := ee
		for i := uint8(0); i <= 127; i++ {
			targetEEPROM.Speed = i
			if err := checkFunc(dat, targetEEPROM); err != nil {
				t.Fatalf("%+v\n", err)
			}
		}
		targetEEPROM.Speed = 128
		if err := checkFunc(dat, targetEEPROM); err == nil {
			t.Error("This should be fail, becouse Speed should be less than or equal to 127")
		}
	})
	t.Run("Punch", func(t *testing.T) {
		targetEEPROM := ee
		for i := uint8(0); i <= 10; i++ {
			targetEEPROM.Punch = i
			if err := checkFunc(dat, targetEEPROM); err != nil {
				t.Fatalf("%+v\n", err)
			}
		}
		targetEEPROM.Punch = 11
		if err := checkFunc(dat, targetEEPROM); err == nil {
			t.Error("This should be fail, becouse Punch should be less than or equal to 10")
		}
	})
	t.Run("DeadBand", func(t *testing.T) {
		targetEEPROM := ee
		for i := uint8(0); i <= 10; i++ {
			targetEEPROM.DeadBand = i
			if err := checkFunc(dat, targetEEPROM); err != nil {
				t.Fatalf("%+v\n", err)

			}
		}
		targetEEPROM.DeadBand = 11
		if err := checkFunc(dat, targetEEPROM); err == nil {
			t.Error("This should be fail, becouse DeadBand should be less than or equal to 10")
		}
	})
	t.Run("Damping", func(t *testing.T) {
		targetEEPROM := ee
		for i := uint8(1); i <= 254; i++ {
			targetEEPROM.Damping = i
			if err := checkFunc(dat, targetEEPROM); err != nil {
				t.Fatalf("%+v\n", err)
			}
		}
		targetEEPROM.Damping = 255
		if err := checkFunc(dat, targetEEPROM); err != nil {
			t.Fatalf("%+v\n", err)
		}
		targetEEPROM.Damping = 0
		if err := checkFunc(dat, targetEEPROM); err == nil {
			t.Fatal("This should be fail, becouse DeadBand should be not equal to 0")
		}
	})
	t.Run("SafeTimer", func(t *testing.T) {
		targetEEPROM := ee
		for i := uint8(1); i <= 254; i++ {
			targetEEPROM.SafeTimer = i
			if err := checkFunc(dat, targetEEPROM); err != nil {
				t.Fatalf("%+v\n", err)
			}
		}
		targetEEPROM.SafeTimer = 255
		if err := checkFunc(dat, targetEEPROM); err != nil {
			t.Fatalf("%+v\n", err)
		}
		targetEEPROM.SafeTimer = 0
		if err := checkFunc(dat, targetEEPROM); err == nil {
			t.Fatal("This should be fail, becouse SafeTimer should be not equal to 0")
		}
	})
	t.Run("Flag", func(t *testing.T) {
		targetEEPROM := ee
		for i := 0; i <= 0x0011111; i++ {
			targetEEPROM.Flag = Flag{
				i&0x1<<4 != 0,
				i&0x1<<3 != 0,
				i&0x1<<2 != 0,
				i&0x1<<1 != 0,
				i&0x1<<0 != 0,
			}
			if err := checkFunc(dat, targetEEPROM); err != nil {
				t.Fatalf("%+v\n", err)
			}
		}
	})
	t.Run("MaximumPulseLimit", func(t *testing.T) {
		targetEEPROM := ee
		for i := uint16(3500); i <= 11500; i++ {
			targetEEPROM.MaximumPulseLimit = i
			if err := checkFunc(dat, targetEEPROM); err != nil {
				t.Fatalf("%+v\n", err)
			}
		}
		for i := uint16(0); i < 3500; i++ {
			targetEEPROM.MaximumPulseLimit = i
			if err := checkFunc(dat, targetEEPROM); err == nil {
				t.Fail()
			}
		}
		// this will be block
		// for i := uint16(11501); i < 1<<16-1; i++ {
		// 	targetEEPROM.MaximumPulseLimit = i
		// 	if err := checkFunc(dat, targetEEPROM); err == nil {
		// 		t.Fail()
		// 	}
		// }
	})
	t.Run("MinimumPulseLimit", func(t *testing.T) {
		targetEEPROM := ee
		for i := uint16(3500); i <= 11500; i++ {
			targetEEPROM.MinimumPulseLimit = i
			if err := checkFunc(dat, targetEEPROM); err != nil {
				t.Fatalf("%+v\n", err)
			}
		}
		for i := uint16(0); i < 3500; i++ {
			targetEEPROM.MinimumPulseLimit = i
			if err := checkFunc(dat, targetEEPROM); err == nil {
				t.Fail()
			}
		}
		// this will be block
		// for i := uint16(11501); i < 1<<16-1; i++ {
		// 	targetEEPROM.MaximumPulseLimit = i
		// 	if err := checkFunc(dat, targetEEPROM); err == nil {
		// 		t.Fail()
		// 	}
		// }
	})
	t.Run("SignalSpeed", func(t *testing.T) {
		targetEEPROM := ee
		targetEEPROM.SignalSpeed = High
		if err := checkFunc(dat, targetEEPROM); err != nil {
			t.Fatalf("%+v\n", err)
		}
		targetEEPROM.SignalSpeed = Mid
		if err := checkFunc(dat, targetEEPROM); err != nil {
			t.Fatalf("%+v\n", err)
		}
		targetEEPROM.SignalSpeed = Low
		if err := checkFunc(dat, targetEEPROM); err != nil {
			t.Fatalf("%+v\n", err)
		}
	})
	t.Run("TemperatureLimit", func(t *testing.T) {
		targetEEPROM := ee
		for i := uint8(1); i <= 127; i++ {
			targetEEPROM.TemperatureLimit = i
			if err := checkFunc(dat, targetEEPROM); err != nil {
				t.Fatalf("%+v\n", err)
			}
		}
		targetEEPROM.TemperatureLimit = 128
		if err := checkFunc(dat, targetEEPROM); err == nil {
			t.Error("This should be fail, becouse Speed should be less than or equal to 127")
		}
		targetEEPROM.TemperatureLimit = 0
		if err := checkFunc(dat, targetEEPROM); err == nil {
			t.Error("This should be fail, becouse Speed should be equal to 0")
		}
	})
	t.Run("CurrentLimit", func(t *testing.T) {
		targetEEPROM := ee
		for i := uint8(1); i <= 63; i++ {
			targetEEPROM.CurrentLimit = i
			if err := checkFunc(dat, targetEEPROM); err != nil {
				t.Fatalf("%+v\n", err)
			}
		}
		targetEEPROM.CurrentLimit = 128
		if err := checkFunc(dat, targetEEPROM); err == nil {
			t.Error("This should be fail, becouse Speed should be less than or equal to 63")
		}
		targetEEPROM.CurrentLimit = 0
		if err := checkFunc(dat, targetEEPROM); err == nil {
			t.Error("This should be fail, becouse Speed should be equal to 0")
		}
	})
	t.Run("Response", func(t *testing.T) {
		targetEEPROM := ee
		for i := uint8(1); i <= 5; i++ {
			targetEEPROM.Response = i
			if err := checkFunc(dat, targetEEPROM); err != nil {
				t.Fatalf("%+v\n", err)
			}
		}
		targetEEPROM.Response = 6
		if err := checkFunc(dat, targetEEPROM); err == nil {
			t.Error("This should be fail, becouse Speed should be less than or equal to 6")
		}
		targetEEPROM.Response = 0
		if err := checkFunc(dat, targetEEPROM); err == nil {
			t.Error("This should be fail, becouse Speed should be equal to 0")
		}
	})
	t.Run("UserOffset", func(t *testing.T) {
		targetEEPROM := ee
		for i := int8(-127); i < 127; i++ {
			targetEEPROM.UserOffset = i
			if err := checkFunc(dat, targetEEPROM); err != nil {
				t.Fatalf("%+v\n", err)
			}
		}
		targetEEPROM.UserOffset = 127
		if err := checkFunc(dat, targetEEPROM); err != nil {
			t.Fatalf("%+v\n", err)
		}
	})
	t.Run("ID", func(t *testing.T) {
		targetEEPROM := ee
		for i := uint8(0); i <= 31; i++ {
			targetEEPROM.ID = i
			if err := checkFunc(dat, targetEEPROM); err != nil {
				t.Fatalf("%+v\n", err)
			}
		}
		targetEEPROM.ID = 32
		if err := checkFunc(dat, targetEEPROM); err == nil {
			t.Fatal("This should be fail, becouse Speed should be equal to 31")
		}
	})
	t.Run("CharacteristicChangeStretch1", func(t *testing.T) {
		targetEEPROM := ee
		for i := uint8(2); i <= 254; i++ {
			targetEEPROM.CharacteristicChangeStretch1 = i
			if i%2 == 1 {
				if err := checkFunc(dat, targetEEPROM); err == nil {
					t.Fatal("This should be fail, becouse StretchGain should be even number")
				}
			} else {
				if err := checkFunc(dat, targetEEPROM); err != nil {
					t.Fatal(err)
				}
			}
		}
		targetEEPROM.CharacteristicChangeStretch1 = 255
		if err := checkFunc(dat, targetEEPROM); err == nil {
			t.Fatal("This should be fail, becouse StretchGain should be even number")
		}
		targetEEPROM.CharacteristicChangeStretch1 = 0
		if err := checkFunc(dat, targetEEPROM); err == nil {
			t.Error("This should be fail, becouse Speed should be equal to 0")
		}
	})
	t.Run("CharacteristicChangeStretch2", func(t *testing.T) {
		targetEEPROM := ee
		for i := uint8(2); i <= 254; i++ {
			targetEEPROM.CharacteristicChangeStretch2 = i
			if i%2 == 1 {
				if err := checkFunc(dat, targetEEPROM); err == nil {
					t.Fatal("This should be fail, becouse StretchGain should be even number")
				}
			} else {
				if err := checkFunc(dat, targetEEPROM); err != nil {
					t.Fatal(err)
				}
			}
		}
		targetEEPROM.CharacteristicChangeStretch2 = 255
		if err := checkFunc(dat, targetEEPROM); err == nil {
			t.Fatal("This should be fail, becouse StretchGain should be even number")
		}
		targetEEPROM.CharacteristicChangeStretch2 = 0
		if err := checkFunc(dat, targetEEPROM); err == nil {
			t.Error("This should be fail, becouse Speed should be equal to 0")
		}
	})
	t.Run("CharacteristicChangeStretch3", func(t *testing.T) {
		targetEEPROM := ee
		for i := uint8(2); i <= 254; i++ {
			targetEEPROM.CharacteristicChangeStretch3 = i
			if i%2 == 1 {
				if err := checkFunc(dat, targetEEPROM); err == nil {
					t.Fatal("This should be fail, becouse StretchGain should be even number")
				}
			} else {
				if err := checkFunc(dat, targetEEPROM); err != nil {
					t.Fatal(err)
				}
			}
		}
		targetEEPROM.CharacteristicChangeStretch3 = 255
		if err := checkFunc(dat, targetEEPROM); err == nil {
			t.Fatal("This should be fail, becouse StretchGain should be even number")
		}
		targetEEPROM.CharacteristicChangeStretch3 = 0
		if err := checkFunc(dat, targetEEPROM); err == nil {
			t.Error("This should be fail, becouse Speed should be equal to 0")
		}
	})

	// t.Run("TemperatureLimit", func(t *testing.T) {
	// 	targetEEPROM := ee
	// 	for i := uint8(1); i <= 127; i++ {
	// 		targetEEPROM.Damping = i
	// 		if err := checkFunc(dat, targetEEPROM); err != nil {
	// 			t.Fatalf("%+v\n", err)
	// 		}
	// 	}
	// 	targetEEPROM.Damping = 128
	// 	if err := checkFunc(dat, targetEEPROM); err == nil {
	// 		t.Fail()
	// 	}
	// 	targetEEPROM.Damping = 0
	// 	if err := checkFunc(dat, targetEEPROM); err == nil {
	// 		t.Fatal("This should be fail, becouse DeadBand should be not equal to 0")
	// 	}
	// })
	// t.Logf("%+v\n", ee)
}
