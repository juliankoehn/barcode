package barcode

import (
	"fmt"
	"strconv"
	"strings"
)

var planet = map[int][]int{
	0: {1, 1, 2, 2, 2},
	1: {2, 2, 2, 1, 1},
	2: {2, 2, 1, 2, 1},
	3: {2, 2, 1, 1, 2},
	4: {2, 1, 2, 2, 1},
	5: {2, 1, 2, 1, 2},
	6: {2, 1, 1, 2, 2},
	7: {1, 2, 2, 2, 1},
	8: {1, 2, 2, 1, 2},
	9: {1, 2, 1, 2, 2},
}

var postnet = map[int][]int{
	0: {2, 2, 1, 1, 1},
	1: {1, 1, 1, 2, 2},
	2: {1, 1, 2, 1, 2},
	3: {1, 1, 2, 2, 1},
	4: {1, 2, 1, 1, 2},
	5: {1, 2, 1, 2, 1},
	6: {1, 2, 2, 1, 1},
	7: {2, 1, 1, 1, 2},
	8: {2, 1, 1, 2, 1},
	9: {2, 1, 2, 1, 1},
}

/** POSTNET and PLANET barcodes.
 * Used by U.S. Postal Service for automated mail sorting
 * @param code {string} zip code to represent.
 * Must be a string containing a zip code of the form DDDDD or DDDDD-DDDD
 * @param planet {bool} if true print the PLANET barcode, otherwise print POSTNET
 */
func barcodePOSTNET(code string, isPlanet bool) *barArray {
	var barlen map[int][]int
	if isPlanet {
		barlen = planet
	} else {
		barlen = postnet
	}
	bararray := &barArray{
		Code:  code,
		MaxW:  0,
		MaxH:  2,
		BCode: []bCode{},
	}
	// remove dashes
	code = strings.Replace(code, "-", "", -1)
	// remove whitespaces
	code = strings.Replace(code, " ", "", -1)

	if !onlyDigits(code) {
		if !onlyDigits(code) {
			fmt.Printf("[BARCODE] invalid POSTNET Code: %s must be Digit e.g.: 55555-1237\n", code)
			return nil
		}
	}

	// calculate checksum
	sum := 0
	for i := 0; i < len(code); i++ {
		intVal, _ := strconv.Atoi(string(code[i]))
		sum = sum + intVal
	}
	chkd := sum % 10
	if chkd > 0 {
		chkd = 10 - chkd
	}
	code = code + strconv.Itoa(chkd)
	bararray.BCode = append(bararray.BCode, bCode{
		T: true,
		W: 1,
		H: 2,
		P: 0,
	})
	bararray.BCode = append(bararray.BCode, bCode{
		T: false,
		W: 1,
		H: 2,
		P: 0,
	})
	bararray.MaxW += 2

	for i := 0; i < len(code); i++ {
		for j := 0; j < 5; j++ {
			intVal, _ := strconv.Atoi(string(code[i]))
			h := barlen[intVal][j]
			p := 1 / h
			bararray.BCode = append(bararray.BCode, bCode{
				T: true,
				W: 1,
				H: h,
				P: p,
			})
			bararray.BCode = append(bararray.BCode, bCode{
				T: false,
				W: 1,
				H: 2,
				P: 0,
			})
			bararray.MaxW += 2
		}
	}

	bararray.BCode = append(bararray.BCode, bCode{
		T: true,
		W: 1,
		H: 2,
		P: 0,
	})
	bararray.MaxW++
	return bararray
}
