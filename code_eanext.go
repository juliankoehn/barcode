package barcode

import (
	"fmt"
	"strconv"
)

var parities = make(map[int]map[int][]string, 5)

var codes = map[string]map[string]string{
	"A": { // left odd parity
		"0": "0001101",
		"1": "0011001",
		"2": "0010011",
		"3": "0111101",
		"4": "0100011",
		"5": "0110001",
		"6": "0101111",
		"7": "0111011",
		"8": "0110111",
		"9": "0001011",
	},
	"B": {
		"0": "0100111",
		"1": "0110011",
		"2": "0011011",
		"3": "0100001",
		"4": "0011101",
		"5": "0111001",
		"6": "0000101",
		"7": "0010001",
		"8": "0001001",
		"9": "0010111",
	},
}

func init() {
	parities[2] = map[int][]string{
		0: {"A", "1"},
		1: {"A", "B"},
		2: {"B", "A"},
		3: {"B", "B"},
	}
	parities[5] = map[int][]string{
		0: {"B", "B", "A", "A", "A"},
		1: {"B", "A", "B", "A", "A"},
		2: {"B", "A", "A", "B", "A"},
		3: {"B", "A", "A", "A", "B"},
		4: {"A", "B", "B", "A", "A"},
		5: {"A", "A", "B", "B", "A"},
		6: {"A", "A", "A", "B", "B"},
		7: {"A", "B", "A", "B", "A"},
		8: {"A", "B", "A", "A", "B"},
		9: {"A", "A", "B", "A", "B"},
	}
}

// barcodeEANNEXT UPC-Based Extensions
// 2 Digit ext: used to indicate magazines and newspaper issue numbers
// 5 digit ext: used to mark suggested retail price of books
func barcodeEANNEXT(code string, digits int) *barArray {
	if !onlyDigits(code) {
		fmt.Printf("code may only contain digits but got %s\n", code)
		return nil
	}
	// padding
	switch digits {
	case 2:
		code = fmt.Sprintf("%02v", code)
	case 5:
		code = fmt.Sprintf("%05v", code)
	default:
		fmt.Printf("invalid extension: %s\n", code)
		return nil
	}
	intCode, _ := strconv.Atoi(code)
	var r int
	if digits == 2 {
		r = intCode % 4
	} else if digits == 5 {
		r2 := ((intCode / 1000) % 10) + ((intCode % 100) / 10)
		r2 = 9 * r2
		r = ((intCode / 10000) + ((intCode % 1000) / 100) + (intCode % 10)) * 3
		r = r + r2
		r %= 10
	} else {
		fmt.Printf("invalid extension: %s\n", code)
		return nil
	}
	// convert digits to bars
	p := parities[digits][r]
	seq := "1011"

	seq = seq + codes[p[0]][string(code[0])]

	for i := 1; i < digits; i++ {
		seq = seq + "01" // seperator
		seq = seq + codes[p[i]][string(code[i])]
	}
	bararray := &barArray{
		Code:  code,
		MaxW:  0,
		MaxH:  1,
		BCode: []bCode{},
	}

	return binaryToArr(seq, bararray)
}

func onlyDigits(code string) bool {
	b := true
	for _, c := range code {
		if c < '0' || c > '9' {
			b = false
			break
		}
	}
	return b
}
