package barcode

import (
	"strconv"
	"strings"
)

/* barcodeI25 Interleaved 2 of 5 barcodes.
 * Compact numeric code, widely used in industry, air cargo
 * Contains digits (0 to 9) and encodes the data in the width of both bars and spaces.
 */
func barcodeI25(code string, checksum bool) *barArray {
	chrsI25 := map[string]string{
		"0": "11221",
		"1": "21112",
		"2": "12112",
		"3": "22111",
		"4": "11212",
		"5": "21211",
		"6": "12211",
		"7": "11122",
		"8": "21121",
		"9": "12121",
		"A": "11",
		"Z": "21",
	}
	if checksum {
		code = code + checksumI25(code)
	}
	if len(code)%2 != 0 {
		// add leading zero if code-length is odd
		code = "0" + code
	}

	code = "AA" + strings.ToLower(code) + "ZA"
	arr := barArray{
		Code:  code,
		MaxW:  0,
		MaxH:  1,
		BCode: []bCode{},
	}

	for i := 0; i < len(code); i = i + 2 {
		var charBar string
		var charSpace string

		charBar = string([]rune(code)[i])

		if len(code) != i+1 {
			charSpace = string([]rune(code)[i])
		} else {
			charSpace = ""
		}

		if chrsI25[charBar] == "" || chrsI25[charSpace] == "" {
			return nil
		}

		seq := ""
		for s := 0; s < len(chrsI25[charBar]); s++ {
			seq = seq + string([]rune(chrsI25[charBar])[s])
			if charSpace != "" {
				seq = seq + string([]rune(chrsI25[charSpace])[s])
			}
		}

		for j := 0; j < len(seq); j++ {
			var t bool
			if j%2 == 0 {
				t = true
			} else {
				t = false
			}
			w := string(seq[j])
			wValue, _ := strconv.Atoi(w)
			arr.BCode = append(arr.BCode, bCode{
				T: t,
				W: wValue,
				H: 1,
				P: 0,
			})
			arr.MaxW += wValue
		}

	}

	return &arr
}

func checksumI25(code string) string {
	sum := 0
	for i := 0; i < len(code); i += 2 {
		v, _ := strconv.Atoi(string(code[i]))
		sum += v
	}
	sum *= 3

	for i := 0; i < len(code); i += 2 {
		v, _ := strconv.Atoi(string(code[i]))
		sum += v
	}
	r := sum % 10
	if r > 0 {
		r = (10 - r)
	}
	return strconv.Itoa(r)
}
