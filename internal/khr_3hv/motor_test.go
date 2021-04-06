package khr_3hv

import (
	"io/ioutil"
	"testing"

	"gopkg.in/yaml.v2"
)

func TestYaml(t *testing.T) {
	tt := make(map[string]string)
	data, err := ioutil.ReadFile("./id.yaml")
	if err != nil {
		t.Error(err)
	}
	if err := yaml.Unmarshal([]byte(data), &tt); err != nil {
		t.Error(err)
		t.Log(tt)
	}
	r := Robot{}
	if err := r.LoadYaml(tt); err != nil {
		t.Fatal(err)
	}
	var id uint8
	id = 0
	if r.Head.GetID() != id {
		t.Errorf("ID Head is wrong, should be %d, but actual %d", id, r.Head.GetID())
	}
	if r.Waist.GetID() != id {
		t.Errorf("ID Waist is wrong, should be %d, but actual %d", id, r.Waist.GetID())
	}
	id = 1
	if r.LeftShoulderRoll.GetID() != id {
		t.Errorf("ID LeftShoulderRoll is wrong, should be %d, but actual %d", id, r.LeftShoulderRoll.GetID())
	}
	if r.RightShoulderRoll.GetID() != id {
		t.Errorf("ID RightShoulderRoll is wrong, should be %d, but actual %d", id, r.RightShoulderRoll.GetID())
	}
	id = 2
	if r.LeftShoulderPitch.GetID() != id {
		t.Errorf("ID LeftShoulderPitch is wrong, should be %d, but actual %d", id, r.LeftShoulderPitch.GetID())
	}
	if r.RightShoulderPitch.GetID() != id {
		t.Errorf("ID RightShoulderPitch is wrong, should be %d, but actual %d", id, r.RightShoulderPitch.GetID())
	}
	id = 3
	if r.LeftElbowRoll.GetID() != id {
		t.Errorf("ID LeftElbowRoll is wrong, should be %d, but actual %d", id, r.LeftElbowRoll.GetID())
	}
	if r.RightElbowRoll.GetID() != id {
		t.Errorf("ID RightElbowRoll is wrong, should be %d, but actual %d", id, r.RightElbowRoll.GetID())
	}
	id = 4
	if r.LeftElbowPitch.GetID() != id {
		t.Errorf("ID RightElbowRoll is wrong, should be %d, but actual %d", id, r.RightElbowRoll.GetID())
	}
	if r.RightElbowPitch.GetID() != id {
		t.Errorf("ID RightElbowPitch is wrong, should be %d, but actual %d", id, r.RightElbowPitch.GetID())
	}
	id = 5
	if r.LeftHipRoll.GetID() != id {
		t.Errorf("ID LeftHipRoll is wrong, should be %d, but actual %d", id, r.LeftHipRoll.GetID())
	}
	if r.RightHipRoll.GetID() != id {
		t.Errorf("ID RightHipRoll is wrong, should be %d, but actual %d", id, r.RightHipRoll.GetID())
	}
	id = 6
	if r.LeftHipPitch.GetID() != id {
		t.Errorf("ID LeftHipPitch is wrong, should be %d, but actual %d", id, r.LeftHipPitch.GetID())
	}
	if r.RightHipPitch.GetID() != id {
		t.Errorf("ID RightHipPitch is wrong, should be %d, but actual %d", id, r.RightHipPitch.GetID())
	}
	id = 7
	if r.LeftHipYaw.GetID() != id {
		t.Errorf("ID LeftHipYaw is wrong, should be %d, but actual %d", id, r.LeftHipYaw.GetID())
	}
	if r.RightHipYaw.GetID() != id {
		t.Errorf("ID RightHipYaw is wrong, should be %d, but actual %d", id, r.RightHipYaw.GetID())
	}
	id = 8
	if r.LeftKnee.GetID() != id {
		t.Errorf("ID LeftKnee is wrong, should be %d, but actual %d", id, r.LeftKnee.GetID())
	}
	if r.RightKnee.GetID() != id {
		t.Errorf("ID RightKnee is wrong, should be %d, but actual %d", id, r.RightKnee.GetID())
	}
	id = 9
	if r.LeftAnkleRoll.GetID() != id {
		t.Errorf("ID LeftAnkleRoll is wrong, should be %d, but actual %d", id, r.LeftAnkleRoll.GetID())
	}
	if r.RightAnkleRoll.GetID() != id {
		t.Errorf("ID RightAnkleRoll is wrong, should be %d, but actual %d", id, r.RightAnkleRoll.GetID())
	}
	id = 10
	if r.LeftAnklePitch.GetID() != id {
		t.Errorf("ID LeftAnklePitch is wrong, should be %d, but actual %d", id, r.LeftAnklePitch.GetID())
	}
	if r.RightAnklePitch.GetID() != id {
		t.Errorf("ID RightAnklePitch is wrong, should be %d, but actual %d", id, r.RightAnklePitch.GetID())
	}

}
