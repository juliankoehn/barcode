package barcode

import (
	"fmt"
	"strconv"
)

var msiChars = map[string]string{
	"0": "100100100100",
	"1": "100100100110",
	"2": "100100110100",
	"3": "100100110110",
	"4": "100110100100",
	"5": "100110100110",
	"6": "100110110100",
	"7": "100110110110",
	"8": "110100100100",
	"9": "110100100110",
	"A": "110100110100",
	"B": "110100110110",
	"C": "110110100100",
	"D": "110110100110",
	"E": "110110110100",
	"F": "110110110110",
}

/** barcodeMSI
 * Variation of Plessey code, with similiar applications
 * Contains digits (0-9) and encodes the data only in the width of bars
 */
func barcodeMSI(code string, checksum bool) *barArray {
	if !onlyDigits(code) {
		fmt.Printf("[BARCODE] code may only contain digits but got %s\n", code)
		return nil
	}
	if checksum {
		// add checksum
		p := 2
		check := 0
		for i := (len(code) - 1); i >= 0; i-- {
			intVal, _ := strconv.Atoi(string(code[i]))
			check = check + (intVal * p)
			p++
			if p > 7 {
				p = 2
			}
		}
		check %= 11
		if check > 0 {
			check = 11 - check
		}
		code = code + strconv.Itoa(check)
	}

	seq := "110" // left guard

	for i := 0; i < len(code); i++ {
		digit := string(code[i])
		if _, ok := msiChars[digit]; !ok {
			return nil
		}
		seq = seq + msiChars[digit]
	}

	seq = seq + "1001" // right guard

	bararray := &barArray{
		Code:  code,
		MaxW:  0,
		MaxH:  1,
		BCode: []bCode{},
	}
	return binaryToArr(seq, bararray)
}
