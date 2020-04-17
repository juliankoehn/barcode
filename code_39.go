package barcode

import (
	"strconv"
	"strings"
)

func barcodeCode39(code string, extended, checksum bool) *barArray {
	if code == "" {
		return nil
	}

	if extended {
		// code = encode_code39_ext(code)
		code = encodeCode39Ext(code)
	}

	if checksum {
		code = code + checksumCode39(code)
	}
	code = strings.ToUpper(code)

	// add start and stop codes if they does not exists on code
	if code[len(code)-1:] != "*" {
		code = code + "*"
	}
	if code[:1] != "*" {
		code = "*" + code
	}

	bararray := barArray{
		Code:  code,
		MaxW:  0,
		MaxH:  1,
		BCode: []bCode{},
	}

	// avg 7 iterations
	for i := 0; i < len(code); i++ {
		char := string([]rune(code)[i])

		chrs := chr[char]
		if chrs == "" {
			return nil
		}
		for j := 0; j < 9; j++ {
			var t bool
			if j%2 == 0 {
				t = true
			} else {
				t = false
			}
			w := string([]rune(chr[char])[j])
			wValue, _ := strconv.Atoi(w)
			x := bCode{
				T: t,
				W: wValue,
				H: 1,
				P: 0,
			}

			bararray.BCode = append(bararray.BCode, x)
			bararray.MaxW = bararray.MaxW + wValue
		}
		// gaps
		bararray.BCode = append(bararray.BCode, bCode{
			T: false,
			W: 1,
			H: 1,
			P: 0,
		})
	}
	// 128
	bararray.MaxW += len(code)

	return &bararray
}

/**
 * Calculate Code 39 checksum (modulo 43)
 */
func checksumCode39(code string) string {
	sum := 0

	for _, r := range code {
		v := chars[r]
		sum += v
	}

	sum = sum % 43

	for r, v := range chars {
		if v == sum {
			return string(r)
		}
	}
	return "#"
}

// encodeCode39Ext encode a string to be used for code 39 extended mode
func encodeCode39Ext(code string) string {
	var codeExt string

	for _, r := range code {
		if int(r) > len(encodeDictionary) {
			return ""
		}
		codeExt = codeExt + encodeDictionary[int(r)]
	}

	return codeExt
}
