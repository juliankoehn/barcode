package main

import (
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/juliankoehn/barcode"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.GET("/:type/:value", handler)
	e.Logger.Fatal(e.Start(":1323"))
}

func handler(c echo.Context) error {
	width := 2
	height := 60
	color := "black"
	showCode := true

	which := c.Param("type")
	value := c.Param("value")

	var extension = filepath.Ext(value)
	if extension != "" {
		value = value[0 : len(value)-len(extension)]
	}

	widthParam := c.QueryParam("width")
	if widthParam != "" {
		x, err := strconv.Atoi(widthParam)
		if err == nil {
			width = x
		}
	}
	heightParam := c.QueryParam("height")
	if heightParam != "" {
		x, err := strconv.Atoi(heightParam)
		if err == nil {
			height = x
		}
	}

	colorParam := c.QueryParam("color")
	if colorParam != "" {
		color = colorParam
	}

	showCodeParam := c.QueryParam("showCode")
	if showCodeParam != "" {
		y, err := strconv.ParseBool(showCodeParam)
		if err == nil {
			showCode = y
		}
	}

	if extension == ".png" || extension == ".jpg" || extension == ".jpeg" || extension == ".gif" {
		var Openfile *os.File
		var img *image.RGBA
		if extension == ".png" {
			Openfile, img = barcode.GetBarcodeFile(value, which, width, height, color, showCode, false, true)
		} else {
			Openfile, img = barcode.GetBarcodeFile(value, which, width, height, color, showCode, false, false)
		}

		var contentType string

		switch extension {
		case ".png":
			png.Encode(Openfile, img)
			contentType = "image/png"
		case ".jpg", ".jpeg":
			contentType = "image/jpeg"
			jpeg.Encode(Openfile, img, &jpeg.Options{Quality: 100})
		case ".gif":
			gif.Encode(Openfile, img, nil)
		}

		FileHeader := make([]byte, 512)
		Openfile.Read(FileHeader)

		if contentType == "" {
			contentType = http.DetectContentType(FileHeader)
		}

		FileStat, _ := Openfile.Stat()                     // Get info from file
		FileSize := strconv.FormatInt(FileStat.Size(), 10) // Get file size as a string

		//c.Response().Header().Set("Content-Disposition", "attachment; filename=\""+Openfile.Name()+"\"")
		c.Response().Header().Set("Content-Type", contentType)
		c.Response().Header().Set("Content-Length", FileSize)
		Openfile.Seek(0, 0)

		io.Copy(c.Response().Writer, Openfile)
		return nil
	}
	// by default return svg
	return c.HTML(200, barcode.GetBarcodeSVG(value, which, width, height, color, showCode, false))
}
