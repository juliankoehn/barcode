package barcode

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

func barcodeC128(code string, variant string) *barArray {
	codeRunes := strToRunes(code)
	if len(codeRunes) <= 0 || len(codeRunes) > 80 {
		fmt.Printf("code length should be between 1 and 80 runes but got %d", len(codeRunes))
		return nil
	}
	chrs128 := []string{
		"212222", /* 00 */
		"222122", /* 01 */
		"222221", /* 02 */
		"121223", /* 03 */
		"121322", /* 04 */
		"131222", /* 05 */
		"122213", /* 06 */
		"122312", /* 07 */
		"132212", /* 08 */
		"221213", /* 09 */
		"221312", /* 10 */
		"231212", /* 11 */
		"112232", /* 12 */
		"122132", /* 13 */
		"122231", /* 14 */
		"113222", /* 15 */
		"123122", /* 16 */
		"123221", /* 17 */
		"223211", /* 18 */
		"221132", /* 19 */
		"221231", /* 20 */
		"213212", /* 21 */
		"223112", /* 22 */
		"312131", /* 23 */
		"311222", /* 24 */
		"321122", /* 25 */
		"321221", /* 26 */
		"312212", /* 27 */
		"322112", /* 28 */
		"322211", /* 29 */
		"212123", /* 30 */
		"212321", /* 31 */
		"232121", /* 32 */
		"111323", /* 33 */
		"131123", /* 34 */
		"131321", /* 35 */
		"112313", /* 36 */
		"132113", /* 37 */
		"132311", /* 38 */
		"211313", /* 39 */
		"231113", /* 40 */
		"231311", /* 41 */
		"112133", /* 42 */
		"112331", /* 43 */
		"132131", /* 44 */
		"113123", /* 45 */
		"113321", /* 46 */
		"133121", /* 47 */
		"313121", /* 48 */
		"211331", /* 49 */
		"231131", /* 50 */
		"213113", /* 51 */
		"213311", /* 52 */
		"213131", /* 53 */
		"311123", /* 54 */
		"311321", /* 55 */
		"331121", /* 56 */
		"312113", /* 57 */
		"312311", /* 58 */
		"332111", /* 59 */
		"314111", /* 60 */
		"221411", /* 61 */
		"431111", /* 62 */
		"111224", /* 63 */
		"111422", /* 64 */
		"121124", /* 65 */
		"121421", /* 66 */
		"141122", /* 67 */
		"141221", /* 68 */
		"112214", /* 69 */
		"112412", /* 70 */
		"122114", /* 71 */
		"122411", /* 72 */
		"142112", /* 73 */
		"142211", /* 74 */
		"241211", /* 75 */
		"221114", /* 76 */
		"413111", /* 77 */
		"241112", /* 78 */
		"134111", /* 79 */
		"111242", /* 80 */
		"121142", /* 81 */
		"121241", /* 82 */
		"114212", /* 83 */
		"124112", /* 84 */
		"124211", /* 85 */
		"411212", /* 86 */
		"421112", /* 87 */
		"421211", /* 88 */
		"212141", /* 89 */
		"214121", /* 90 */
		"412121", /* 91 */
		"111143", /* 92 */
		"111341", /* 93 */
		"131141", /* 94 */
		"114113", /* 95 */
		"114311", /* 96 */
		"411113", /* 97 */
		"411311", /* 98 */
		"113141", /* 99 */
		"114131", /* 100 */
		"311141", /* 101 */
		"411131", /* 102 */
		"211412", /* 103 START A */
		"211214", /* 104 START B */
		"211232", /* 105 START C */
		"233111", /* STOP */
		"200000", /* END */
	}
	// ASCII chars for code A (ASCII 00 - 95)
	keysAB := " !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_"
	keysB := keysAB + "`abcdefghijklmnopqrstuvwxyz{|}~\u007F"
	keysAOnly := "\u0000\u0001\u0002\u0003\u0004" + // NUL, SOH, STX, ETX, EOT
		"\u0005\u0006\u0007\u0008\u0009" + // ENQ, ACK, BEL, BS,  HT
		"\u000A\u000B\u000C\u000D\u000E" + // LF,  VT,  FF,  CR,  SO
		"\u000F\u0010\u0011\u0012\u0013" + // SI,  DLE, DC1, DC2, DC3
		"\u0014\u0015\u0016\u0017\u0018" + // DC4, NAK, SYN, ETB, CAN
		"\u0019\u001A\u001B\u001C\u001D" + // EM,  SUB, ESC, FS,  GS
		"\u001E\u001F"
	keysA := keysAB + keysAOnly
	fncA := map[rune]int{
		241: 102,
		242: 97,
		243: 96,
		244: 101,
	}
	fncB := map[rune]int{
		241: 102,
		242: 97,
		243: 96,
		244: 100,
	}

	codeData := []int{}
	var startID int

	switch strings.ToUpper(variant) {
	case "A":
		startID = 103
		for _, value := range codeRunes {
			if value >= 241 && value <= 244 {
				codeData = append(codeData, fncA[value])
			} else if value >= 0 && value <= 95 {
				codeData = append(codeData, strings.IndexRune(keysA, value))
			} else {
				// out of range
				return nil
			}

		}
	case "B":
		startID = 104
		for _, value := range codeRunes {
			if value >= 241 && value <= 244 {
				codeData = append(codeData, fncB[value])
			} else if value >= 32 && value <= 127 {
				codeData = append(codeData, strings.IndexRune(keysB, value))
			} else {
				// out of range
				return nil
			}
		}
	case "C":
		startID = 105
		if codeRunes[0] == 241 {
			codeData = append(codeData, 102)
			code = string(code[:1])
			codeRunes = strToRunes(code)
		}
		if len(codeRunes)%2 != 0 {
			fmt.Print("code length must be even")
			return nil
		}
		for i := 0; i < len(codeRunes); i += 2 {
			chrnum := string(codeRunes[i]) + string(codeRunes[i+1])
			matched, err := regexp.MatchString("([0-9]{2})", chrnum)
			if err != nil {
				fmt.Print("Regex error with given C128C")
				return nil
			}
			if matched {
				val, _ := strconv.Atoi(chrnum)
				codeData = append(codeData, val)
			} else {
				fmt.Print("No matches in given Code")
				return nil
			}
		}
	default:
		re := regexp.MustCompile("([0-9]{4,})")
		matches := re.FindAllString(string(codeRunes), -1)
		var seq []sequence
		if len(matches) > 0 {
			endOffset := 0
			for _, value := range matches {
				offset := strings.Index(string(codeRunes), value)
				if offset > endOffset {
					seq = append(seq, get128ABsequence(string(codeRunes[endOffset:offset])))
				}
				slen := len(value)
				if slen%2 != 0 {
					slen--
				}
				// sequence{"B", code, len(code)}
				seq = append(seq, sequence{"C", value, slen, false})
				endOffset = offset + slen
			}

			if endOffset < len(codeRunes) {
				seq = append(seq, get128ABsequence(string(codeRunes[endOffset:])))
			}
		} else {
			// (10)123456(8020)198798787
			// text code (non C mode)
			seq = append(seq, get128ABsequence(string(codeRunes)))
		}
		// process the sequence
		for key, value := range seq {
			switch value.Variant {
			case "A":
				if key == 0 {
					startID = 103
				} else if seq[key-1].Variant != "A" {
					if value.Length == 1 && key > 0 && seq[key-1].Variant == "B" && !seq[key-1].B {
						// single character shift
						codeData = append(codeData, 98)
						value.B = true
					} else if !seq[key-1].B {
						codeData = append(codeData, 101)
					}
				}
				for _, value := range value.Code {
					// value rune
					if value >= 241 && value <= 244 {
						codeData = append(codeData, fncA[value])
					} else {
						codeData = append(codeData, strings.IndexRune(keysA, value))
					}
				}
			case "B":
				if key == 0 {
					tempchr := []rune(value.Code[:1])
					if value.Length == 1 &&
						tempchr[0] >= 241 &&
						tempchr[0] <= 244 &&
						len(seq)-1 > key &&
						seq[key+1].Variant != "B" {
						switch seq[key+1].Variant {
						case "A":
							startID = 103
							value.Variant = "A"
							codeData = append(codeData, fncA[tempchr[0]])
						case "C":
							startID = 105
							value.Variant = "C"
							codeData = append(codeData, fncA[tempchr[0]])
						}
					} else {
						startID = 104
					}
				} else if seq[key-1].Variant != "B" {
					if value.Length == 1 && key > 0 && seq[key-1].Variant == "A" && seq[key-1].B == false {
						codeData = append(codeData, 98)
						value.B = true
					} else if seq[key-1].B == false {
						codeData = append(codeData, 100)
					}
				}
				for _, value := range value.Code {
					// value rune
					if value >= 241 && value <= 244 {
						codeData = append(codeData, fncB[value])
					} else {
						codeData = append(codeData, strings.IndexRune(keysB, value))
					}
				}
			case "C":
				if key == 0 {
					startID = 105
				} else if seq[key-1].Variant != "C" {
					codeData = append(codeData, 99)
				}
				for i := 0; i < len(value.Code); i += 2 {
					chrnum := string(byte(value.Code[i])) + string(byte(value.Code[i+1]))
					intval, _ := strconv.Atoi(chrnum)
					codeData = append(codeData, intval)
				}
			}
		}
	}
	sum := startID

	for key, value := range codeData {
		sum += (value * (key + 1))
	}
	// add check character
	codeData = append(codeData, sum%103)
	codeData = append(codeData, 106)
	codeData = append(codeData, 107)
	// add start code at the beginning
	codeData = append([]int{startID}, codeData...)

	// build barcode array
	arr := &barArray{
		Code:  string(codeRunes),
		MaxW:  0,
		MaxH:  1,
		BCode: []bCode{},
	}
	for _, value := range codeData {
		seq := chrs128[value]
		for j := 0; j < 6; j++ {
			var t bool
			if j%2 == 0 {
				t = true
			}
			w, _ := strconv.Atoi(string(seq[j]))
			arr.BCode = append(arr.BCode, bCode{
				T: t,
				W: w,
				H: 1,
				P: 0,
			})
			arr.MaxW += w
		}
	}
	return arr
}

type sequence struct {
	Variant string
	Code    string
	Length  int
	B       bool
}

// TODO: finish....
func get128ABsequence(code string) sequence {
	fmt.Println(code)
	re := regexp.MustCompile("([0-31])")
	matches := re.FindAllString(code, -1)
	if len(matches) > 0 {
		// it's A or B
		return sequence{"B", code, len(code), false}
	} else {
		// it's B - for sure
		return sequence{"B", code, len(code), false}
	}
}

func strToRunes(str string) []rune {
	result := make([]rune, utf8.RuneCountInString(str))
	i := 0
	for _, r := range str {
		result[i] = r
		i++
	}
	return result
}
