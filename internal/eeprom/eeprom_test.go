package eeprom

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
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
	eeprom, err := Parsing(dat)
	if err != nil {
		t.Errorf("%+v\n", err)
	}
	j, err := json.Marshal(eeprom.Address)
	if err != nil {
		t.Errorf("json.Marshal: %v", err)
	}
	t.Log(string(j))
}
