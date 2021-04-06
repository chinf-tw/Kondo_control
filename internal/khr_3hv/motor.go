package khr_3hv

import (
	_ "embed"
	"errors"
	"io"
	"kondocontrol/internal/eeprom"
	"kondocontrol/internal/serial"
	"reflect"
	"strconv"

	"gopkg.in/yaml.v2"
)

// Robot
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

// RobotNum
type RobotNum [22]Motor

// Kind
type Kind uint8

const (
	Head Kind = iota
	Waist
	LeftShoulderRoll
	LeftShoulderPitch
	LeftElbowRoll
	LeftElbowPitch
	LeftHipRoll
	LeftHipPitch
	LeftHipYaw
	LeftKnee
	LeftAnkleRoll
	LeftAnklePitch
	RightShoulderRoll
	RightShoulderPitch
	RightElbowRoll
	RightElbowPitch
	RightHipRoll
	RightHipPitch
	RightHipYaw
	RightKnee
	RightAnkleRoll
	RightAnklePitch
)

//go:embed id.yaml
var defaultYaml []byte

// DefaultRobot
func DefaultRobot(leftPort, rightPort io.ReadWriteCloser) (Robot, error) {
	r := Robot{}
	tt := make(map[string]string)
	if err := yaml.Unmarshal(defaultYaml, &tt); err != nil {
		return r, err
	}
	r.LoadIDWithYaml(tt)
	return r, nil
}

// DefaultRobotNum
func DefaultRobotNum(leftPort, rightPort io.ReadWriteCloser) (RobotNum, error) {
	r := RobotNum{}
	// setting all ID
	settingRobotNumID(&r)
	// setting all port
	r[Head].port = leftPort
	r[Waist].port = rightPort
	for i := LeftShoulderRoll; i <= LeftAnklePitch; i++ {
		r[i].port = leftPort
	}
	for i := RightShoulderRoll; i <= RightAnklePitch; i++ {
		r[i].port = rightPort
	}
	return r, nil
}

func settingRobotNumID(r *RobotNum) {
	r[Head].EEPROM.ID = 0
	r[Waist].EEPROM.ID = 0
	for i := LeftShoulderRoll; i <= LeftAnklePitch; i++ {
		r[i].EEPROM.ID = uint8(i-LeftShoulderRoll) + 1
	}
	for i := RightShoulderRoll; i <= RightAnklePitch; i++ {
		r[i].EEPROM.ID = uint8(i-RightShoulderRoll) + 1
	}
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

// UpdateCurrentPositionWithFree
func (m *Motor) UpdateCurrentPositionWithFree() error {
	position, err := serial.SetFree(m.EEPROM.ID, m.port)
	if err != nil {
		return err
	}
	m.Position = position
	return nil
}

// SetPosition
func (m *Motor) SetPosition(target uint) error {
	currentPos, err := serial.SetPosition(m.GetID(), target, m.port)
	if err != nil {
		return err
	}
	m.Position = currentPos
	return err
}

// GetID
func (m Motor) GetID() uint8 {
	return m.EEPROM.ID
}

// SetID
func (m *Motor) SetID(id uint8) {
	m.EEPROM.ID = id
}

// LoadIDWithYaml
func (r *Robot) LoadIDWithYaml(y map[string]string) error {
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
	return nil
}
