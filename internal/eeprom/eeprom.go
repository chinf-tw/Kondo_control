package eeprom

import (
	"encoding/json"
	"io/ioutil"
	"reflect"
	"strconv"

	"github.com/pkg/errors"
)

// EEPROM is KONDO servo motor eeprom
type EEPROM struct {
	StretchGain                  uint8       `start:"2"  end:"4"`
	Speed                        uint8       `start:"4"  end:"6"`
	Punch                        uint8       `start:"6"  end:"8"`
	DeadBand                     uint8       `start:"8"  end:"10"`
	Damping                      uint8       `start:"10" end:"12"`
	SafeTimer                    uint8       `start:"12" end:"14"`
	Flag                         Flag        `start:"14" end:"16"`
	MaximumPulseLimit            uint16      `start:"16" end:"20"`
	MinimumPulseLimit            uint16      `start:"20" end:"24"`
	SignalSpeed                  SignalSpeed `start:"26" end:"28"`
	TemperatureLimit             uint8       `start:"28" end:"30"`
	CurrentLimit                 uint8       `start:"30" end:"32"`
	Response                     uint8       `start:"50" end:"52"`
	UserOffset                   int8        `start:"52" end:"54"`
	ID                           uint8       `start:"56" end:"58"`
	CharacteristicChangeStretch1 uint8       `start:"58" end:"60"`
	CharacteristicChangeStretch2 uint8       `start:"60" end:"62"`
	CharacteristicChangeStretch3 uint8       `start:"62" end:"64"`
	Address                      Address
}

// Flag is KONDO servo motor flag
type Flag struct {
	Reverse      bool
	Free         bool
	PWMINH       bool
	RotationMode bool
	SlaveMode    bool
}
type SignalSpeed uint8

const (
	High SignalSpeed = 0
	Mid  SignalSpeed = 1
	Low  SignalSpeed = 10
)

var (
	// ErrDataLength is when the data length is not 64
	ErrDataLength = errors.New("The data length is not 64")
	// ErrDataMismatch is when The data is mismatch
	ErrDataMismatch = errors.New("The data is mismatch")
)

// Parse resolve bytes to EEPROM
func Parse(bs []byte) (EEPROM, error) {
	if len(bs) != 64 {
		return EEPROM{}, ErrDataLength
	}
	var (
		result = EEPROM{}
		mark   = uint8(0)
	)
	var recordFunc func(key *Interval, add uint8)
	if buildAddress {
		recordFunc = func(key *Interval, add uint8) {
			*key = NewInterval(mark, mark+add)
			mark += add
		}
	} else {
		recordFunc = func(_ *Interval, add uint8) {
			mark += add
		}
		{
			dat, err := ioutil.ReadFile(testingFilePath)
			if err != nil {
				return EEPROM{}, errors.WithStack(err)
			}
			err = json.Unmarshal(dat, &result.Address)
			if err != nil {
				return EEPROM{}, errors.WithStack(err)
			}
			if result.Address == (Address{}) {
				return EEPROM{}, errors.New("json Data unmarshal failed")
			}
		}
	}

	// *** DON'T DO IT ***
	{ // Fixed as 0x5A
		beginning, err := sliceByteToUint8(bs[mark : mark+2])
		if err != nil {
			return EEPROM{}, errors.WithStack(err)
		}
		if beginning != 0x5A {
			return EEPROM{},
				errors.Wrapf(
					ErrDataMismatch,
					"Target is %X, but actual %X, origin: %v",
					0x5A, beginning, bs)
		}
		recordFunc(&result.Address.Fixed, 2)
	}
	{ // Stretch gain
		// 2,4…254 2-step sequence
		// consisting only of even
		stretchGain, err := sliceByteToUint8(bs[mark : mark+2])
		if err != nil {
			return EEPROM{}, errors.WithStack(err)
		}
		if stretchGain%2 != 0 {
			return EEPROM{}, errors.WithStack(ErrDataMismatch)
		}
		result.StretchGain = stretchGain
		recordFunc(&result.Address.StretchGain, 2)
	}
	{ // Speed
		// 1,2,3…127
		speed, err := sliceByteToUint8(bs[mark : mark+2])
		if err != nil {
			return EEPROM{}, errors.WithStack(err)
		}
		if speed > 127 {
			return EEPROM{}, errors.WithStack(ErrDataMismatch)
		}
		result.Speed = speed
		recordFunc(&result.Address.Speed, 2)
	}
	{ // Punch
		// 0,1,2,3…10
		punch, err := sliceByteToUint8(bs[mark : mark+2])
		if err != nil {
			return EEPROM{}, errors.WithStack(err)
		}
		if punch > 10 {
			return EEPROM{}, errors.WithStack(ErrDataMismatch)
		}
		result.Punch = punch
		recordFunc(&result.Address.Punch, 2)
	}
	{ // Dead band
		// 0,1,2,3,4,5,...10
		deadBand, err := sliceByteToUint8(bs[mark : mark+2])
		if err != nil {
			return EEPROM{}, errors.WithStack(err)
		}
		if deadBand > 10 { //FIXME: that need to check
			return EEPROM{},
				errors.Wrapf(
					ErrDataMismatch,
					"Dead band cann't bigger then 5, but actual is %d, data is %v, origin: %v, result: %v",
					bs[mark:mark+2], deadBand, bs, result)
		}
		result.DeadBand = deadBand
		recordFunc(&result.Address.DeadBand, 2)
	}
	{ // Damping
		// 1,2…255
		damping, err := sliceByteToUint8(bs[mark : mark+2])
		if err != nil {
			return EEPROM{}, errors.WithStack(err)
		}
		if damping == 0 {
			return EEPROM{}, errors.WithStack(ErrDataMismatch)
		}
		result.Damping = damping
		recordFunc(&result.Address.Damping, 2)
	}
	{ // Safe timer
		// 10,11…255
		// (0x01-0xFF)
		safeTimer, err := sliceByteToUint8(bs[mark : mark+2])
		if err != nil {
			return EEPROM{}, errors.WithStack(err)
		}
		if safeTimer == 0 {
			return EEPROM{}, errors.WithStack(ErrDataMismatch)
		}
		result.SafeTimer = safeTimer
		recordFunc(&result.Address.SafeTimer, 2)
	}
	{ // Flag
		flagDetail := bs[mark : mark+2]
		if flagDetail[0]&0xF0 != 0 || flagDetail[1]&0xF0 != 0 {
			return EEPROM{}, errors.WithStack(ErrDataMismatch)
		}
		if flagDetail[0]&0b00000110 != 0 {
			return EEPROM{}, errors.Wrapf(ErrDataMismatch, "flagDetail: %b\n", flagDetail)
		}
		if (flagDetail[1] & 0b00000100 >> 2) != 1 {
			return EEPROM{}, errors.Wrapf(ErrDataMismatch, "flagDetail[1]&0b00000100 != 1,actual: %v, data: %v", flagDetail[1]&0b00000100, flagDetail)
		}
		flag := Flag{
			SlaveMode:    flagDetail[0]&0b00001000>>3 == 1,
			RotationMode: flagDetail[0]&0b00000001 == 1,
			PWMINH:       flagDetail[1]&0b00001000>>3 == 1,
			Free:         flagDetail[1]&0b00000010>>1 == 1,
			Reverse:      flagDetail[1]&0b00000001 == 1,
		}
		result.Flag = flag
		recordFunc(&result.Address.Flag, 2)
	}
	{ // Maximum pulse limit
		maximumPulseLimit, err := sliceByteToUint16(bs[mark : mark+4])
		if err != nil {
			return EEPROM{}, errors.WithStack(err)
		}
		if maximumPulseLimit < 3500 || maximumPulseLimit > 11500 {
			return EEPROM{}, errors.WithStack(ErrDataMismatch)
		}
		result.MaximumPulseLimit = maximumPulseLimit
		recordFunc(&result.Address.MaximumPulseLimit, 4)
	}
	{ // Minimum pulse limit
		minimumPulseLimit, err := sliceByteToUint16(bs[mark : mark+4])
		if err != nil {
			return EEPROM{}, errors.WithStack(err)
		}
		if minimumPulseLimit < 3500 || minimumPulseLimit > 11500 {
			return EEPROM{}, errors.WithStack(ErrDataMismatch)
		}
		result.MinimumPulseLimit = minimumPulseLimit
		recordFunc(&result.Address.MinimumPulseLimit, 4)
	}
	mark += 2
	{ // Signal speed
		// switch ss {
		// case 0x00:
		// 	result.SignalSpeed = High
		// case 0x01:
		// 	result.SignalSpeed = Mid
		// case 0x10:
		// 	result.SignalSpeed = Low
		// }
		result.SignalSpeed = SignalSpeed(bs[mark])
		recordFunc(&result.Address.SignalSpeed, 2)
	}
	{ // Temperature limit
		temperatureLimit, err := sliceByteToUint8(bs[mark : mark+2])
		if err != nil {
			return EEPROM{}, errors.WithStack(err)
		}
		if temperatureLimit < 1 || temperatureLimit > 127 {
			return EEPROM{}, errors.WithStack(ErrDataMismatch)
		}
		result.TemperatureLimit = temperatureLimit
		recordFunc(&result.Address.TemperatureLimit, 2)
	}
	{ // Current limit
		currentLimit, err := sliceByteToUint8(bs[mark : mark+2])
		if err != nil {
			return EEPROM{}, errors.WithStack(err)
		}
		if currentLimit < 1 || currentLimit > 63 {
			return EEPROM{}, errors.WithStack(ErrDataMismatch)
		}
		result.CurrentLimit = currentLimit
		recordFunc(&result.Address.CurrentLimit, 2)
	}
	for i := 0; i < 9; i++ {
		mark += 2
	}
	{ // Response
		response, err := sliceByteToUint8(bs[mark : mark+2])
		if err != nil {
			return EEPROM{}, errors.WithStack(err)
		}
		if response < 1 || response > 5 {
			return EEPROM{}, errors.WithStack(ErrDataMismatch)
		}
		result.Response = response
		recordFunc(&result.Address.Response, 2)
	}
	{ // User offset
		userOffset, err := sliceByteToInt8(bs[mark : mark+2])
		if err != nil {
			return EEPROM{}, errors.WithStack(err)
		}
		result.UserOffset = userOffset
		recordFunc(&result.Address.UserOffset, 2)
	}
	mark += 2
	{ // ID
		id, err := sliceByteToUint8(bs[mark : mark+2])
		if err != nil {
			return EEPROM{}, errors.WithStack(err)
		}
		if id > 31 {
			return EEPROM{}, errors.WithStack(ErrDataMismatch)
		}
		result.ID = id
		recordFunc(&result.Address.ID, 2)
	}
	{ // Characteristic change stretch 1
		stretch1, err := sliceByteToUint8(bs[mark : mark+2])
		if err != nil {
			return EEPROM{}, errors.WithStack(err)
		}
		if stretch1 == 0 || stretch1%2 == 1 {
			return EEPROM{}, errors.Wrapf(ErrDataMismatch, "actual: %d", stretch1)
		}
		result.CharacteristicChangeStretch1 = stretch1
		recordFunc(&result.Address.CharacteristicChangeStretch1, 2)
	}
	{ // Characteristic change stretch 2
		stretch2, err := sliceByteToUint8(bs[mark : mark+2])
		if err != nil {
			return EEPROM{}, errors.WithStack(err)
		}
		if stretch2 == 0 || stretch2%2 == 1 {
			return EEPROM{}, errors.WithStack(ErrDataMismatch)
		}
		result.CharacteristicChangeStretch2 = stretch2
		recordFunc(&result.Address.CharacteristicChangeStretch2, 2)
	}
	{ // Characteristic change stretch 3
		stretch3, err := sliceByteToUint8(bs[mark : mark+2])
		if err != nil {
			return EEPROM{}, errors.WithStack(err)
		}
		if stretch3 == 0 || stretch3%2 == 1 {
			return EEPROM{}, errors.WithStack(ErrDataMismatch)
		}
		result.CharacteristicChangeStretch3 = stretch3
		recordFunc(&result.Address.CharacteristicChangeStretch3, 2)
	}
	// *** ****** ***
	return result, nil
}

func Compose(origin []byte, target EEPROM) ([]byte, error) {
	result := origin[:]
	valueOfEEPROM := reflect.ValueOf(target)
	typeOfEEPROM := valueOfEEPROM.Type()
	checkTag := func(tag reflect.StructTag) (Interval, error) {
		var (
			start uint8
			end   uint8
			err   error
			// ok    bool
		)
		s := tag.Get("start")
		e := tag.Get("end")
		if s == "" || e == "" {
			return Interval{}, errors.Errorf("address tag is not getting")
		}
		ss, err := strconv.ParseUint(s, 10, 8)
		if err != nil {
			return Interval{}, errors.Errorf("初四啦阿伯，strconv.ParseUint，%+v\n", err)
		}
		// start, ok = interface{}(ss).(uint8)
		// if !ok {
		// 	return Interval{}, errors.Errorf("初四啦阿伯，start 轉型失敗。%+v %T, %v %T", ss, ss, start, start)
		// }
		start = uint8(ss)
		ee, err := strconv.ParseUint(e, 10, 8)
		if err != nil {
			return Interval{}, errors.Errorf("初四啦阿伯，strconv.ParseUint，%+v\n", err)
		}
		// end, ok = interface{}(ee).(uint8)
		// if !ok {
		// 	return Interval{}, errors.New("初四啦阿伯， end 轉型失敗。")
		// }
		end = uint8(ee)
		return Interval{start, end}, nil
	}
	for i := 0; i < typeOfEEPROM.NumField(); i++ {
		fieldValue := valueOfEEPROM.Field(i)
		fieldType := typeOfEEPROM.Field(i)
		s := fieldType.Tag.Get("start")
		e := fieldType.Tag.Get("end")
		if s == "" || e == "" {
			continue
		}
		interval, err := checkTag(fieldType.Tag)
		if err != nil {
			return []byte{}, err
		}
		switch fieldType.Type.Kind() {
		case reflect.Uint8:
			// fmt.Printf("%d. %s Uint8 %v\n", i, fieldType.Name, fieldValue.Uint())
			if fieldType.Name == "SignalSpeed" {
				b := signalSpeedToSliceByte(SignalSpeed(fieldValue.Uint()))
				if uint8(len(b)) != interval.End-interval.Start {
					return []byte{}, errors.Errorf("%d. %s type tag length doesn't equal to function length of bytes, %v\n", i, fieldType.Name, fieldValue.Interface())
				}
				for i := interval.Start; i < interval.End; i++ {
					result[i] = b[i-interval.Start]
				}
				continue
			}
			b := uint8ToSliceByte(uint8(fieldValue.Uint()))
			if uint8(len(b)) != interval.End-interval.Start {
				return []byte{}, errors.Errorf("%d. %s type tag length doesn't equal to function length of bytes, %v\n", i, fieldType.Name, fieldValue.Interface())
			}
			for i := interval.Start; i < interval.End; i++ {
				result[i] = b[i-interval.Start]
			}
		case reflect.Struct:
			// fmt.Printf("%d. %s %+v\n", i, fieldType.Name, fieldType.Type.Name())
			if fieldType.Name == "Flag" {
				flag, ok := (fieldValue.Interface()).(Flag)
				if !ok {
					return []byte{}, errors.Errorf("%d. %s type Flag is transformation problematic, %v\n", i, fieldType.Name, fieldValue.Interface())
				}
				b := flagToSliceByte(flag)
				if uint8(len(b)) != interval.End-interval.Start {
					return []byte{}, errors.Errorf("%d. %s type tag length doesn't equal to function length of bytes, %v\n", i, fieldType.Name, fieldValue.Interface())
				}
				for i := interval.Start; i < interval.End; i++ {
					result[i] = b[i-interval.Start]
				}
			}
		case reflect.Uint16:
			// fmt.Printf("%d. %s Uint16\n", i, fieldType.Name)
			b := uint16ToSliceByte(uint16(fieldValue.Uint()))
			if uint8(len(b)) != interval.End-interval.Start {
				return []byte{}, errors.Errorf("%d. %s type tag length doesn't equal to function length of bytes, %v\n", i, fieldType.Name, fieldValue.Interface())
			}
			for i := interval.Start; i < interval.End; i++ {
				result[i] = b[i-interval.Start]
			}
		// case reflect.Uint:
		// 	// fmt.Printf("%d. %s Uint32\n", i, fieldType.Name)

		case reflect.Int8:
			// fmt.Printf("%d. %s Int8\n", i, fieldType.Name)
			b := int8ToSliceByte(int8(fieldValue.Int()))
			if uint8(len(b)) != interval.End-interval.Start {
				return []byte{}, errors.Errorf("%d. %s type tag length doesn't equal to function length of bytes, %v\n", i, fieldType.Name, fieldValue.Interface())
			}
			for i := interval.Start; i < interval.End; i++ {
				result[i] = b[i-interval.Start]
			}
		default:
			return []byte{}, errors.Errorf("%d. %s type is not catched, %v\n", i, fieldType.Name, fieldType.Type.Kind())
		}
		// fmt.Printf("name: %v ,tag:'%v'\n", fieldType.Name, fieldType.Tag)

	}
	return result, nil
}
