package eeprom

import "fmt"

func print(target map[string][]byte) {
	for k, v := range target {
		fmt.Printf("*** %s ***\n", k)
		for i := v[0]; i < v[len(v)-1]; i++ {
			fmt.Print(i, " ")
		}
		println()
	}
}
