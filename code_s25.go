package barcode

import (
	"strconv"
)

/**
 * Standard 2 of 5 barcodes
 * Used in airline ticket markin, photofinishing
 * Contains digits (0 to 9) and encodes the data only in the width of bars
 */
func barcodeS25(code string, checksum bool) *barArray {
	chrS25 := map[string]string{
		"0": "10101110111010",
		"1": "11101010101110",
		"2": "10111010101110",
		"3": "11101110101010",
		"4": "10101110101110",
		"5": "11101011101010",
		"6": "10111011101010",
		"7": "10101011101110",
		"8": "10101110111010",
		"9": "10111010111010",
	}
	if checksum {
		code = code + checksumS25(code)
	}
	if len(code)%2 != 0 {
		// add leading zero if code-length is odd
		code = "0" + code
	}
	seq := "11011010"

	for i := 0; i < len(code); i++ {
		digit := string([]rune(code)[i])
		if chrS25[digit] == "" {
			return nil
		}

		seq = seq + chrS25[digit]
	}
	seq = seq + "1101011"

	arr := &barArray{
		Code:  code,
		MaxW:  0,
		MaxH:  1,
		BCode: []bCode{},
	}

	return binaryToArr(seq, arr)
}

func binaryToArr(seq string, arr *barArray) *barArray {
	w := 0

	for i := 0; i < len(seq); i++ {
		w++
		v, _ := strconv.Atoi(string(seq[i]))
		var vX int
		if len(seq) != i+1 {
			vX, _ = strconv.Atoi(string(seq[i+1]))
		} else {
			vX = 0
		}

		if i == (len(seq)-1) || i < (len(seq)-1) && v != vX {
			var t bool
			if v == 1 {
				t = true
			} else {
				t = false
			}

			_ = t
			arr.BCode = append(arr.BCode, bCode{
				T: t,
				W: w,
				H: 1,
				P: 0,
			})
			arr.MaxW += w
			w = 0
		}
	}

	return arr
}

// checksumS25 calculates the digit checksum of S25 code
func checksumS25(code string) string {
	sum := 0
	for i := 0; i < len(code); i += 2 {
		v, _ := strconv.Atoi(string(code[i]))
		sum += v
	}
	sum *= 3

	for i := 1; i < len(code); i += 2 {
		v, _ := strconv.Atoi(string(code[i]))
		sum += v
	}

	r := sum % 10
	if r > 0 {
		r = (10 - r)
	}

	return strconv.Itoa(r)
}
