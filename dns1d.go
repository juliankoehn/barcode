package barcode

import (
	"image"
	clr "image/color"
	"image/draw"
	"os"
	"strconv"
	"strings"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

var (
	chars = map[rune]int{
		'0': 0, '1': 1, '2': 2, '3': 3, '4': 4, '5': 5, '6': 6, '7': 7, '8': 8, '9': 9,
		'A': 10, 'B': 11, 'C': 12, 'D': 13, 'E': 14, 'F': 15, 'G': 16, 'H': 17, 'I': 18, 'J': 19, 'K': 20,
		'L': 21, 'M': 22, 'N': 23, 'O': 24, 'P': 25, 'Q': 26, 'R': 27, 'S': 28, 'T': 29, 'U': 30, 'V': 31,
		'W': 32, 'X': 33, 'Y': 34, 'Z': 35, '-': 36, '.': 37, ' ': 38, '$': 39, '/': 40, '+': 41, '%': 42,
	}

	// 131 = ƒ

	dict93 = map[int]string{
		0: "(%)U", 1: "€A", 2: "€B", 3: "€C", 4: "€D", 5: "€E", 6: "€F",
		7: "€G", 8: "€H", 9: "€I", 10: "€J", 11: "£K", 12: "€L", 13: "€M",
		14: "€N", 15: "€O", 16: "€P", 17: "€Q", 18: "€R", 19: "€S",
		20: "€T", 21: "€U", 22: "€V", 23: "€W", 24: "€X", 25: "€Y", 26: "€Z",
		27: "ƒE", 28: "ƒB", 29: "ƒC", 30: "ƒD", 31: "ƒE", 32: " ", 33: "A",
		34: "B", 35: "C",
	}

	encodeDictionary = map[int]string{
		0: "%u", 1: "$A", 2: "$B", 3: "$C",
		4: "$D", 5: "$E", 6: "$F", 7: "$G",
		8: "$H", 9: "$I", 10: "$J", 11: "£K",
		12: "$L", 13: "$M", 14: "$N", 15: "$0",
		16: "$P", 17: "$Q", 18: "$R", 19: "$S",
		20: "$T", 21: "$U", 22: "$V", 23: "$W",
		24: "$X", 25: "$Y", 26: "$Z", 27: "%A",
		28: "%B", 29: "%C", 30: "%D", 31: "%E",
		32: " ", 33: "/A", 34: "/B", 35: "/C",
		36: "/D", 37: "/E", 38: "/F", 39: "/G",
		40: "/H", 41: "/I", 42: "/J", 43: "/K",
		44: "/L", 45: "-", 46: ".", 47: "/O",
		48: "0", 49: "1", 50: "2", 51: "3",
		52: "4", 53: "5", 54: "6", 55: "7",
		56: "8", 57: "9", 58: "/Z", 59: "%F",
		60: "%G", 61: "%H", 62: "%I", 63: "%J",
		64: "%V", 65: "A", 66: "B", 67: "C",
		68: "D", 69: "E", 70: "F", 71: "G",
		72: "H", 73: "I", 74: "J", 75: "K",
		76: "L", 77: "M", 78: "N", 79: "O",
		80: "P", 81: "Q", 82: "R", 83: "S",
		84: "T", 85: "U", 86: "V", 87: "W",
		88: "X", 89: "Y", 90: "Z", 91: "%K",
		92: "%L", 93: "%M", 94: "%N", 95: "%O",
		96: "%W", 97: "+A", 98: "+B", 99: "+C",
		100: "+D", 101: "+E", 102: "+F", 103: "+G",
		104: "+H", 105: "+I", 106: "+J", 107: "+K",
		108: "+L", 109: "+M", 110: "+N", 111: "+O",
		112: "+P", 113: "+Q", 114: "+R", 115: "+S",
		116: "+T", 117: "+U", 118: "+V", 119: "+W",
		120: "+X", 121: "+Y", 122: "+Z", 123: "%P",
		124: "%Q", 125: "%R", 126: "%S", 127: "%T",
	}

	chr93 = map[byte]string{
		48:  "131112", // 0
		49:  "111213", // 1
		50:  "111312", // 2
		51:  "111411", // 3
		52:  "121113", // 4
		53:  "121212", // 5
		54:  "121311", // 6
		55:  "111114", // 7
		56:  "131211", // 8
		57:  "141111", // 9
		65:  "211113", // A
		66:  "211212", // B
		67:  "211311", // C
		68:  "221112", // D
		69:  "221211", // E
		70:  "231111", // F
		71:  "112113", // G
		72:  "112212", // H
		73:  "112311", // I
		74:  "122112", // J
		75:  "132111", // K
		76:  "111123", // L
		77:  "111222", // M
		78:  "111321", // N
		79:  "121122", // O
		80:  "131121", // P
		81:  "212112", // Q
		82:  "212211", // R
		83:  "211122", // S
		84:  "211221", // T
		85:  "221121", // U
		86:  "222111", // V
		87:  "112122", // W
		88:  "112221", // X
		89:  "122121", // Y
		90:  "123111", // Z
		45:  "121131", // -
		46:  "311112", // .
		32:  "311211", // " "
		36:  "321111", // $
		47:  "112131", // /
		43:  "113121", // +
		37:  "211131", // %
		128: "121221", // ($)
		129: "311121", // (/)
		130: "122211", // (+)
		131: "312111", // (%),
		42:  "111141", // start-stop
	}

	chr = map[string]string{
		"0": "111331311",
		"1": "311311113",
		"2": "113311113",
		"3": "313311111",
		"4": "111331113",
		"5": "311331111",
		"6": "113331111",
		"7": "111311313",
		"8": "311311311",
		"9": "113311311",
		"A": "311113113",
		"B": "113113113",
		"C": "313113111",
		"D": "111133113",
		"E": "311133111",
		"F": "113133111",
		"G": "111113313",
		"H": "311113311",
		"I": "113113311",
		"J": "111133311",
		"K": "311111133",
		"L": "113111133",
		"M": "313111131",
		"N": "111131133",
		"O": "311131131",
		"P": "113131131",
		"Q": "111111333",
		"R": "311111331",
		"S": "113111331",
		"T": "111131331",
		"U": "331111113",
		"V": "133111113",
		"W": "333111111",
		"X": "131131113",
		"Y": "331131111",
		"Z": "133131111",
		"-": "131111313",
		".": "331111311",
		" ": "133111311",
		"$": "131313111",
		"/": "131311131",
		"+": "131113131",
		"%": "111313131",
		"*": "131131311",
	}
)

type barArray struct {
	Code  string
	MaxW  int
	MaxH  int
	BCode []bCode
}

type bCode struct {
	T bool
	W int
	H int
	P int
}

// GetBarcodeSVG generates a SVG xml for given code
func GetBarcodeSVG(code, variant string, w, h int, color string, showCode bool, inline bool) string {
	var svg string
	barcodeArray := setBarcode(code, variant)

	if barcodeArray == nil {
		return ""
	}

	if !inline {
		svg = "<?xml version=\"1.0\" standalone=\"no\" ?>\n"
		svg = svg + "<!DOCTYPE svg PUBLIC \"-//W3C//DTD SVG 1.1//EN\" \"http://www.w3.org/Graphics/SVG/1.1/DTD/svg11.dtd\">\n"
	}

	svg = svg + "<svg width=\"" + strconv.Itoa(barcodeArray.MaxW*w) + "\" height=\"" + strconv.Itoa(h) + "\" version=\"1.1\" xmlns=\"http://www.w3.org/2000/svg\">\n"
	svg = svg + "\t<g id=\"bars\" fill=\"" + color + "\" stroke=\"none\">\n"
	x := 0
	bw := 0
	bh := 0
	for _, value := range barcodeArray.BCode {
		bw = value.W * w
		bh = value.H * h / barcodeArray.MaxH

		if showCode {
			bh = bh - 12
		}

		if value.T {
			y := value.P * h / barcodeArray.MaxH
			svg = svg + "\t\t<rect x=\"" + strconv.Itoa(x) + "\" y=\"" + strconv.Itoa(y) + "\" width=\"" + strconv.Itoa(bw) + "\" height=\"" + strconv.Itoa(bh) + "\" />\n"
		}
		x = (x + bw)
	}
	if showCode {
		xCode := (barcodeArray.MaxW * w) / 2
		codeX := strconv.FormatInt(int64(xCode), 10)
		svg = svg + "\t <text x=\"" + codeX + "\" text-anchor=\"middle\" y=\"" + strconv.FormatInt(int64((bh+12)), 10) + "\" id=\"code\" fill=\"" + color + "\" font-size=\"12px\">" + barcodeArray.Code + "</text>\n"
	}
	svg = svg + "\t</g>\n</svg>\n"
	return svg
}

// GetBarcodeFile returns a Barcode as PNG representation
func GetBarcodeFile(code, variant string, w, h int, color string, showCode bool, inline bool, transparant bool) (*os.File, *image.RGBA) {
	barcodeArray := setBarcode(code, variant)
	// calculate image size
	width := barcodeArray.MaxW * w
	height := h

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	// background white
	if transparant {
		draw.Draw(img, img.Bounds(), &image.Uniform{clr.Transparent}, image.ZP, draw.Src)
	} else {
		draw.Draw(img, img.Bounds(), &image.Uniform{clr.White}, image.ZP, draw.Src)
	}

	// print bars

	x := 0
	bw := 0
	bh := 0

	for _, value := range barcodeArray.BCode {
		bw = value.W * w
		bh = value.H * h / barcodeArray.MaxH
		if showCode {
			bh = bh - 14
			// reduce by font height <.<
		}
		if value.T {
			y := value.P * h / barcodeArray.MaxH

			draw.Draw(img, image.Rect(x, y, (x+bw), (y+bh)), &image.Uniform{clr.Black}, image.ZP, draw.Src)
		}
		x = (x + bw)
	}

	if showCode {
		addLabel(img, ((width / 2) - (len(barcodeArray.Code) * 3)), (h - 2), barcodeArray.Code)
		// add code to image
	}

	f, _ := os.Create("barcode.png")
	return f, img
}

func addLabel(img *image.RGBA, x, y int, label string) {
	col := clr.Black
	//x = x - 10
	point := fixed.Point26_6{X: fixed.Int26_6(x * 64), Y: fixed.Int26_6(y * 64)}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(label)
}

func setBarcode(code, variant string) *barArray {
	variant = strings.ToUpper(variant)
	var arrcode *barArray

	switch variant {
	case "C39":
		// CODE 39 - ANSI MH10.8M-1983 - USD-3 - 3 of 9
		arrcode = barcodeCode39(code, false, false)
	case "C39+":
		// CODE 39 with checksum
		return barcodeCode39(code, false, true)
	case "C39E":
		// Code 39 extended
		return barcodeCode39(code, true, false)
	case "C39E+":
		return barcodeCode39(code, true, true)
	case "C93":
		return barcodeCode93(code)
	}

	return arrcode
}

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
