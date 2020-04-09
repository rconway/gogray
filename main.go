package main

import (
	"image"
	"image/color"
	_ "image/jpeg" // Import just for the side-effect of its init() method to perform format registration
	"image/png"
	"log"
	"os"
	"time"
)

func check(err error) {
	if err != nil {
		log.Panicln(err)
	}
}

// ConvertFunc zzz
type ConvertFunc func(c color.Color) color.Color

// Use build-in grayscale conversion
func stdGrayConvert(c color.Color) color.Color {
	return c
}

// Basic 'unweighted' average conversion
func avgGrayConvert(c color.Color) color.Color {
	r, g, b, _ := c.RGBA()
	cgray := ((r + g + b) / 3) >> 8
	return color.Gray{uint8(cgray)}
}

// red channel only
func redGrayConvert(c color.Color) color.Color {
	r, _, _, _ := c.RGBA()
	cgray := r >> 8
	return color.Gray{uint8(cgray)}
}

// green channel only
func greenGrayConvert(c color.Color) color.Color {
	_, g, _, _ := c.RGBA()
	cgray := g >> 8
	return color.Gray{uint8(cgray)}
}

// blue channel only
func blueGrayConvert(c color.Color) color.Color {
	_, _, b, _ := c.RGBA()
	cgray := b >> 8
	return color.Gray{uint8(cgray)}
}

// blue channel only
func ndviGrayConvert(c color.Color) color.Color {
	r, g, b, _ := c.RGBA()
	red := float64(r) / 0xffff
	nir := (float64(g+b) / 2) / 0xffff
	ndvi := 0.0
	if nir != 0 || red != 0 {
		ndvi = (nir - red) / (nir + red)
	}
	cgray := (((ndvi) * 255) + 255) / 2
	return color.Gray{uint8(cgray)}
}

func convert(img image.Image, convert ConvertFunc, outFileName string) {
	log.Println("START create of output: " + outFileName)
	start := time.Now()

	// Converting image to grayscale
	grayImg := image.NewGray(img.Bounds())

	// Calculate ewach pixel
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			grayImg.Set(x, y, convert(img.At(x, y)))
		}
	}

	// Working with grayscale image, e.g. convert to png
	fout, err := os.Create(outFileName)
	check(err)
	defer fout.Close()

	err = png.Encode(fout, grayImg)
	check(err)

	elapsed := time.Since(start)
	log.Printf("...DONE create of output: %v [%v]", outFileName, elapsed)
}

func main() {
	log.Println("...gogray...")

	if len(os.Args) < 2 {
		log.Fatalln("Image path is required")
	}
	imgPath := os.Args[1]

	f, err := os.Open(imgPath)
	check(err)
	defer f.Close()

	img, format, err := image.Decode(f)
	check(err)
	log.Printf("Image (%v) object-type: %T, bounds: %v / %v\n", format, img, img.Bounds().Min, img.Bounds().Max)

	convert(img, stdGrayConvert, "std.png")
	convert(img, avgGrayConvert, "avg.png")
	convert(img, redGrayConvert, "red.png")
	convert(img, greenGrayConvert, "green.png")
	convert(img, blueGrayConvert, "blue.png")
	convert(img, ndviGrayConvert, "ndvi.png")
}
