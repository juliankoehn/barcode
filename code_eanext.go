package barcode

import (
	"fmt"
	"strconv"
)

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
