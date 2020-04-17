package barcode

import (
	"strconv"
	"strings"
)

func barcodeCode93(code string) *barArray {
	code = strings.ToUpper(code)

	var codeExt string

	// rune to string
	for _, r := range code {
		codeExt = codeExt + string(r)
	}

	codeExt = codeExt + checksumCode93(code)

	code = "*" + codeExt + "*"

	bararray := barArray{
		Code:  code,
		MaxW:  0,
		MaxH:  1,
		BCode: []bCode{},
	}
	for i := 0; i < len(code); i++ {
		byt := []byte(code)[i]
		if chr93[byt] == "" {
			return nil
		}
		for j := 0; j < 6; j++ {
			var t bool
			if j%2 == 0 {
				t = true
			} else {
				t = false
			}
			w := string([]rune(chr93[byt])[j])
			wValue, _ := strconv.Atoi(w)
			bararray.BCode = append(bararray.BCode, bCode{
				T: t,
				W: wValue,
				H: 1,
				P: 0,
			})
			bararray.MaxW += wValue
		}

	}
	bararray.BCode = append(bararray.BCode, bCode{
		T: true,
		W: 1,
		H: 1,
		P: 0,
	})
	bararray.MaxW++

	return &bararray
}

// checksumCode93 calculates checksum of
// code digit C & K
func checksumCode93(code string) string {
	chars := map[string]int{
		"0": 0, "1": 1, "2": 2, "3": 3, "4": 4, "5": 5, "6": 6, "7": 7, "8": 8, "9": 9,
		"A": 10, "B": 11, "C": 12, "D": 13, "E": 14, "F": 15, "G": 16, "H": 17, "I": 18, "J": 19, "K": 20,
		"L": 21, "M": 22, "N": 23, "O": 24, "P": 25, "Q": 26, "R": 27, "S": 28, "T": 29, "U": 30, "V": 31,
		"W": 32, "X": 33, "Y": 34, "Z": 35, "-": 36, ".": 37, " ": 38, "$": 39, "/": 40, "+": 41, "%": 42,
		"<": 43, "=": 44, ">": 45, "?": 46,
	}

	// calc check digit C
	weight := 1
	check := 0
	for i := len(code) - 1; i >= 0; i-- {
		char := string([]rune(code)[i]) // 654321
		k := chars[char]
		check = check + (k * weight)
		// reset weight to 1
		if weight++; weight > 20 {
			weight = 1
		}
	}
	check = check % 47
	var c string
	for key, value := range chars {
		if value == check {
			c = key
		}
	}

	code = code + c

	// calc check digit K
	weight = 1
	check = 0
	for i := len(code) - 1; i >= 0; i-- {
		//fmt.Println(i)
		char := string([]rune(code)[i]) // 654321
		k := chars[char]
		check = check + (k * weight)
		// reset weight to 1
		if weight++; weight > 15 {
			weight = 1
		}
	}
	// should be 26 for QWERTZU
	check = check % 47
	var k string
	for key, value := range chars {
		if value == check {
			k = key
		}
	}
	checksum := c + k
	return checksum
}
