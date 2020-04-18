package barcode

import (
	"strings"
)

// 1 = pos 1, length 2
// 2 = pos 2, length 3
// 3 = pos 2, length 1
// 4 = pos 2, length 2
var barmode = map[string][]int{
	"0": {3, 3, 2, 2},
	"1": {3, 4, 1, 2},
	"2": {3, 4, 2, 1},
	"3": {4, 3, 1, 2},
	"4": {4, 3, 2, 1},
	"5": {4, 4, 1, 1},
	"6": {3, 1, 4, 2},
	"7": {3, 2, 3, 2},
	"8": {3, 2, 4, 1},
	"9": {4, 1, 3, 2},
	"A": {4, 1, 4, 1},
	"B": {4, 2, 3, 1},
	"C": {3, 1, 2, 4},
	"D": {3, 2, 1, 4},
	"E": {3, 2, 2, 3},
	"F": {4, 1, 1, 4},
	"G": {4, 1, 2, 3},
	"H": {4, 2, 1, 3},
	"I": {1, 3, 4, 2},
	"J": {1, 4, 3, 2},
	"K": {1, 4, 4, 1},
	"L": {2, 3, 3, 2},
	"M": {2, 3, 4, 1},
	"N": {2, 4, 3, 1},
	"O": {1, 3, 2, 4},
	"P": {1, 4, 1, 4},
	"Q": {1, 4, 2, 3},
	"R": {2, 3, 1, 4},
	"S": {2, 3, 2, 3},
	"T": {2, 4, 1, 3},
	"U": {1, 1, 4, 4},
	"V": {1, 2, 3, 4},
	"W": {1, 2, 4, 3},
	"X": {2, 1, 3, 4},
	"Y": {2, 1, 4, 3},
	"Z": {2, 2, 3, 3},
}

// table for checksum calculation for CBC
var checktable = map[string][2]int{
	"0": {1, 1},
	"1": {1, 2},
	"2": {1, 3},
	"3": {1, 4},
	"4": {1, 5},
	"5": {1, 0},
	"6": {2, 1},
	"7": {2, 2},
	"8": {2, 3},
	"9": {2, 4},
	"A": {2, 5},
	"B": {2, 0},
	"C": {3, 1},
	"D": {3, 2},
	"E": {3, 3},
	"F": {3, 4},
	"G": {3, 5},
	"H": {3, 0},
	"I": {4, 1},
	"J": {4, 2},
	"K": {4, 3},
	"L": {4, 4},
	"M": {4, 5},
	"N": {4, 0},
	"O": {5, 1},
	"P": {5, 2},
	"Q": {5, 3},
	"R": {5, 4},
	"S": {5, 5},
	"T": {5, 0},
	"U": {0, 1},
	"V": {0, 2},
	"W": {0, 3},
	"X": {0, 4},
	"Y": {0, 5},
	"Z": {0, 0},
}

/**	barcodeCBCKIX - CBC - KIX
 * RMS4CC (Royal Mail 4-state Customer Code) - CBC (Customer Bar Code) - KIX (Klant index - Customer index)
 * RM4SCC is the name of the barcode symbology used by the Royal Mail for its Cleanmail service.
 */
func barcodeCBCKIX(code string, kix bool) *barArray {
	code = strings.ToUpper(code)
	bararray := &barArray{
		Code:  code,
		MaxW:  0,
		MaxH:  3,
		BCode: []bCode{},
	}

	// it not kix
	if !kix {
		row := 0
		col := 0

		for i := 0; i < len(code); i++ {
			row += checktable[string(code[i])][0]
			col += checktable[string(code[i])][1]
		}

		row %= 6
		col %= 6
		// row, col fits
		searchSlice := [2]int{row, col}

		for key, value := range checktable {
			if value == searchSlice {
				code = code + key
				break
			}
		}
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
	}

	for i := 0; i < len(code); i++ {
		for j := 0; j < 4; j++ {
			var p int
			var h int
			switch barmode[string(code[i])][j] {
			case 1:
				p = 0
				h = 2
			case 2:
				p = 0
				h = 3
			case 3:
				p = 1
				h = 1
			case 4:
				p = 1
				h = 2
			}
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

	if !kix {
		// stop bar
		bararray.BCode = append(bararray.BCode, bCode{
			T: true,
			W: 1,
			H: 3,
			P: 0,
		})
		bararray.MaxW++
	}

	return bararray
}
