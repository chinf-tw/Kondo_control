package khr_3hv

import (
	"errors"
	"fmt"
	"io"
	"kondocontrol/internal/eeprom"
	"kondocontrol/internal/serial"
	"reflect"
	"strconv"
)

type Robot struct {
	Head               Motor
	Waist              Motor
	LeftShoulderRoll   Motor
	LeftShoulderPitch  Motor
	LeftElbowRoll      Motor
	LeftElbowPitch     Motor
	LeftHipRoll        Motor
	LeftHipPitch       Motor
	LeftHipYaw         Motor
	LeftKnee           Motor
	LeftAnkleRoll      Motor
	LeftAnklePitch     Motor
	RightShoulderRoll  Motor
	RightShoulderPitch Motor
	RightElbowRoll     Motor
	RightElbowPitch    Motor
	RightHipRoll       Motor
	RightHipPitch      Motor
	RightHipYaw        Motor
	RightKnee          Motor
	RightAnkleRoll     Motor
	RightAnklePitch    Motor
}

func NewRobot(leftPort, rightPort io.ReadWriteCloser) {

}

type Motor struct {
	EEPROM      eeprom.EEPROM
	Position    uint
	Stretch     uint8
	Speed       uint8
	Current     uint8
	Temperature uint8
	port        io.ReadWriteCloser
}

func (m *Motor) UpdateCurrentPositionWithFree() error {
	position, err := serial.SetFree(m.EEPROM.ID, m.port)
	if err != nil {
		return err
	}
	m.Position = position
	return nil
}

// GetID
func (m Motor) GetID() uint8 {
	return m.EEPROM.ID
}

// SetID
func (m *Motor) SetID(id uint8) {
	m.EEPROM.ID = id
}

func (r *Robot) LoadYaml(y map[string]string) error {
	val := reflect.Indirect(reflect.ValueOf(r))
	for key, v := range y {
		f := val.FieldByName(key)
		if f.IsValid() {
			// A Value can be changed only if it is
			// addressable and was not obtained by
			// the use of unexported struct fields.
			if f.CanSet() {
				// change value of N
				if f.Kind() == reflect.Struct {
					if f.Type().Name() != "Motor" {
						return errors.New("LoadYaml have Something wrong")
					}
					if !f.CanInterface() {
						return errors.New("LoadYaml have Something wrong with Interface")
					}
					motor, ok := f.Interface().(Motor)
					if !ok {
						return errors.New("LoadYaml have Something wrong with Transformation")
					}
					i, err := strconv.Atoi(v)
					if err != nil {
						return err
					}
					motor.SetID(uint8(i))
					f.Set(reflect.ValueOf(motor))
				}
			}
		}
	}

	fmt.Println(val.Type().Field(0).Name)
	fmt.Println(y)
	return nil
}
