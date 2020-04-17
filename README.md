# barcode
Golang Barcode Generation Package

This package generates barcodes as SVG / or File. You can Encode the File Outputs with `image/* *.Encode(w, f)` See `example` folder.

Supported Barcodes:

* C39
* C39+
* C39E
* C39E+
* C93
* S25
* S25+
* I25
* I25+

## Call:

* `code`: {string} Your Code
* `variant`: {string} one of Supported Barcodes
* `w`: {int} barcode with * w multiplier
* `h`: {int} height of the barcode in px
* `color`: {string} color as CSS compatible string value
* `showCode`: {bool} display code under BARCODE
* `inline`: {bool} removes XML/SVG headers from output

### Returns
SVG as `string`

```go
GetBarcodeSVG(code, variant, w, h, color, showCode, inline)
```