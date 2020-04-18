package barcode

import (
	"fmt"
	"math"
	"strconv"
)

/**
 * EAN13 and UPC-A barcodes.
 * EAN13: European Article Numbering international retail product code
 * UPC-A: Universal product code seen on almost all retail products in the USA and Canada
 * UPC-E: Short version of UPC Symbol
 */
func barcodeEANUPC(code string, digit int) *barArray {
	if !onlyDigits(code) {
		fmt.Printf("code may only contain digits but got %s\n", code)
		return nil
	}
	upce := false
	if digit == 6 {
		digit = 12  // UPC-A
		upce = true // UPC-E mode
	}
	dataLen := digit - 1
	// Padding
	if upce {
		code = upce2a(code)
	} else {
		code = fmt.Sprintf("%0*s", dataLen, code)
	}

	// calculate check digit
	sumA := 0
	for i := 1; i < dataLen; i += 2 {
		intVal, _ := strconv.Atoi(string(code[i]))
		sumA += intVal
	}
	if digit > 12 {
		sumA *= 3
	}
	sumB := 0
	for i := 0; i < dataLen; i += 2 {
		intVal, _ := strconv.Atoi(string(code[i]))
		sumB += (intVal)
	}
	if digit < 13 {
		sumB *= 3
	}
	r := (sumA + sumB) % 10
	if r > 0 {
		r = (10 - r)
	}
	if len(code) == dataLen {
		// add check digit
		code = code + strconv.Itoa(r)
	} else {
		// validate digit
		// last char of code must be r
		checkDigitInt, _ := strconv.Atoi(code[len(code)-1:])
		if checkDigitInt != r {
			fmt.Printf("[BARCODE] Check digit of given code %s is invalid must be %d", code, r)
			return nil
		}
	}
	if digit == 12 {
		// UPC-A
		code = "0" + code
	}
	// convert upc-a to upc-e
	var upcecode string
	if upce {
		tmp := code[4:7]
		if tmp == "000" || tmp == "100" || tmp == "200" {
			// manufacturer code ends in 000, 100, 200
			upcecode = code[2:4] + code[9:12] + code[4:5]
		} else {
			tmp = code[5:7]
			if tmp == "00" {
				upcecode = code[2:5] + code[10:12] + "3"
			} else {
				tmp = code[6:7]
				if tmp == "0" {
					upcecode = code[2:6] + code[11:12] + "4"
				} else {
					upcecode = code[2:7] + code[11:12]
				}
			}

		}
	}
	// Convert digits to bars
	seq := "101"
	var bararray barArray
	if upce {
		bararray = barArray{
			Code:  upcecode,
			MaxW:  0,
			MaxH:  1,
			BCode: []bCode{},
		}
		codeIntVal, _ := strconv.Atoi(code[:1])
		p := upcParities[codeIntVal][r]

		for i := 0; i < 6; i++ {
			seq = seq + codes[p[i]][string(upcecode[i])]
		}
		seq = seq + "010101" // right guard bar
	} else {
		bararray = barArray{
			Code:  code,
			MaxW:  0,
			MaxH:  1,
			BCode: []bCode{},
		}
		halfLen := int(math.Ceil(float64(digit) / float64(2)))
		if digit == 8 {
			for i := 0; i < halfLen; i++ {
				seq = seq + codes["A"][string(code[i])]
			}
		} else {
			codeIntVal, _ := strconv.Atoi(code[:1])
			p := parities[0][codeIntVal]
			for i := 1; i < halfLen; i++ {
				seq = seq + codes[p[i-1]][string(code[i])]
			}
			// all others except upc
		}
		// center guard bar
		seq = seq + "01010"
		for i := halfLen; i < digit; i++ {
			seq = seq + codes["C"][string(code[i])]
		}
		seq = seq + "101" // right guard bar
	}

	// calc width
	w := 0
	for i := 0; i < len(seq); i++ {
		w++
		if i == (len(seq)-1) || i < (len(seq)-1) && string(seq[i]) != string(seq[i+1]) {
			var t bool
			if string(seq[i]) == "1" {
				t = true
			}
			bararray.BCode = append(bararray.BCode, bCode{
				T: t,
				W: w,
				H: 1,
				P: 0,
			})
			bararray.MaxW += w
			w = 0
		}
	}

	return &bararray
}

// upce2a converts UPC-E to UPC-A
func upce2a(code string) string {
	var manufacturer string
	var itemNumber string

	if len(code) > 6 {
		code = code[len(code)-6:]
	} else {
		code = fmt.Sprintf("%06v", code)
	}

	switch code[5:6] {
	case "0":
		// substr($code, 0, 1) = 1
		manufacturer = code[:1] + code[1:2] + code[5:6] + "00"
		itemNumber = "00" + code[2:5]

	case "1":
		manufacturer = code[:2] + code[5:6] + "00"
		itemNumber = "00" + code[2:5]
	case "2":
		manufacturer = code[:2] + code[5:6] + "00"
		itemNumber = "00" + code[2:5]
	case "3":
		manufacturer = code[:2] + code[2:3] + "00"
		itemNumber = "000" + code[3:5]
	case "4":
		manufacturer = code[:4] + "0"
		itemNumber = "0000" + code[4:5]
	default:
		manufacturer = code[:5]
		itemNumber = "0000" + code[5:6]
	}
	return "0" + manufacturer + itemNumber
}
