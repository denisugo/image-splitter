package main

import (
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"math"
	"os"
	"regexp"
	"strings"
)

func main() {
	var inputFile, outputDir, w, h = parseFlags()
	r, _ := regexp.Compile(`\/(\w*)\.`)
	var prefix, _, _ = strings.Cut(r.FindString(inputFile), ".")
	var img = load(inputFile)
	var width, height int = img.Bounds().Dx(), img.Bounds().Dy()
	var frame, frameW, frameH image.Rectangle
	var remW, remH = math.Mod(float64(width), float64(w)), math.Mod(float64(height), float64(h))
	if remW != 0 {
		if math.Mod(remW, 2.0) == 1 {
			frameW = image.Rect(int(remW/2)+1, 0, width-int(remW/2), height)
		} else {
			frameW = image.Rect(int(remW/2), 0, width-int(remW/2), height)
		}
	}
	if remH != 0 {
		if math.Mod(remH, 2.0) == 1 {
			frameH = image.Rect(0, int(remH/2)+1, width, height-int(remH/2))
		} else {
			frameH = image.Rect(0, int(remH/2), width, height-int(remH/2))
		}
	}

	if frameW.Max != image.Pt(0, 0) && frameH.Max != image.Pt(0, 0) {
		img = img.SubImage(frameW.Intersect(frameH)).(*image.RGBA)
		width, height = img.Bounds().Dx(), img.Bounds().Dy()
	} else if frameW.Max != image.Pt(0, 0) {
		img = img.SubImage(frameW).(*image.RGBA)
		width = img.Bounds().Dx()
	} else if frameH.Max != image.Pt(0, 0) {
		img = img.SubImage(frameH).(*image.RGBA)
		height = img.Bounds().Dy()
	}

	var count int
	if (height/h) >= 1 && (width/w) >= 1 {
		for x := 0; x < width; x += w {
			for y := 0; y < height; y += h {
				//fmt.Println(x, y, x+w, y+h, count)
				frame = image.Rect(x, y, x+w, y+h)
				save(fmt.Sprintf("%s%s_%d.jpg", outputDir, prefix, count), img.SubImage(frame)) // prefix already has /
				count++
			}
		}
	} else {
		log.Println("Image is too small")
	}
}

func parseFlags() (inputFile, outputDir string, width, height int) {
	flag.StringVar(&inputFile, "i", "", "input file path")
	flag.StringVar(&outputDir, "o", "./", "output directory")
	flag.IntVar(&width, "w", 54, "width of subimage")
	flag.IntVar(&height, "h", 54, "height of subimage")

	flag.Parse()
	if inputFile == "" {
		log.Fatalf("input file must be specified")
	}
	return
}

func load(filePath string) *image.RGBA {
	imgFile, err := os.Open(filePath)
	if err != nil {
		log.Println("Cannot read file:", err)
	}
	defer imgFile.Close()

	var img image.Image
	var _, ext, _ = strings.Cut(filePath, ".")
	if ext == "png" {
		img, err = png.Decode(imgFile)
	} else {
		img, err = jpeg.Decode(imgFile)
	}
	if err != nil {
		log.Println("Cannot decode file:", err)
	}
	return img.(*image.RGBA)
}

func save(filePath string, img image.Image) {
	imgFile, err := os.Create(filePath)
	if err != nil {
		log.Println("Cannot create file:", err)
	}
	defer imgFile.Close()
	jpeg.Encode(imgFile, img, nil)
}
