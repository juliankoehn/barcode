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
	switch variant {
	case "C39":
		return barcodeCode39(code, false, false)
	case "C39+":
		return barcodeCode39(code, false, true)
	case "C39E":
		return barcodeCode39(code, true, false)
	case "C39E+":
		return barcodeCode39(code, true, true)
	case "C93":
		return barcodeCode93(code)
	case "S25":
		return barcodeS25(code, false)
	case "S25+":
		return barcodeS25(code, true)
	case "I25":
		return barcodeI25(code, false)
	case "I25+":
		return barcodeI25(code, true)
	case "C128":
		return barcodeC128(code, "")
	case "C128A":
		return barcodeC128(code, "A")
	case "C128B":
		return barcodeC128(code, "B")
	case "C128C":
		return barcodeC128(code, "C")
	case "EAN2":
		return barcodeEANNEXT(code, 2)
	case "EAN5":
		return barcodeEANNEXT(code, 5)
	case "EAN8":
		return barcodeEANUPC(code, 8)
	case "EAN13":
		return barcodeEANUPC(code, 13)
	case "UPCA":
		return barcodeEANUPC(code, 12)
	case "UPCE":
		return barcodeEANUPC(code, 6)
	case "MSI":
		return barcodeMSI(code, false)
	case "MSI+":
		return barcodeMSI(code, true)
	case "POSTNET":
		return barcodePOSTNET(code, false)
	case "PLANET":
		return barcodePOSTNET(code, true)
	case "RMS4CC":
		return barcodeCBCKIX(code, false)
	case "KIX":
		return barcodeCBCKIX(code, true)
	case "IMB":
		return imd(code)
	default:
		return nil
	}
}
