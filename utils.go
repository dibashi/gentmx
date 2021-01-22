package main

import (
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"strings"
)

func filenameToImage(filename string) image.Image {
	srcFile, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer srcFile.Close()
	var decode func(io.Reader) (image.Image, error)
	if strings.HasSuffix(filename, ".png") {
		decode = png.Decode
	} else if strings.HasSuffix(filename, ".jpg") {
		decode = jpeg.Decode
	} else {
		panic(filename + "不是png或jpg")
	}
	srcImage, err := decode(srcFile)
	if err != nil {
		panic(err)
	}
	return srcImage
}
