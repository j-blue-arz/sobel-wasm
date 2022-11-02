package main

import (
	"image"
	"image/color"
	"math"
)

// the buffer has size width*height
type float64Image struct {
	buffer []float64
	width  int
	height int
}

func makeGrayImage(width, height int) float64Image {
	return float64Image{make([]float64, width*height), width, height}
}

func (image float64Image) index(row, col int) int {
	return (row*image.width + col)
}

func (image float64Image) get(row, col int) float64 {
	return image.buffer[image.index(row, col)]
}

func (image float64Image) set(row, col int, value float64) {
	image.buffer[image.index(row, col)] = value
}

// The returned image has its size reduced by 2 in both dimensions.
func sobelRGBA(rgba image.RGBA) *image.RGBA {
	grayImage := toGrayImage(rgba)
	convolved, min, max := sobelGray(grayImage)
	return toRGBAImage(convolved, min, max)
}

func toGrayImage(rgba image.RGBA) float64Image {
	width, height := rgba.Bounds().Dx(), rgba.Bounds().Dy()
	grayImage := makeGrayImage(width, height)
	for x := rgba.Bounds().Min.X; x < rgba.Bounds().Max.X; x++ {
		for y := rgba.Bounds().Min.Y; y < rgba.Bounds().Max.Y; y++ {
			red, green, blue, _ := rgba.At(x, y).RGBA()
			grayImage.set(y, x, float64(toGray(red, green, blue)))
		}
	}
	return grayImage
}

func toGray(red, green, blue uint32) float64 {
	return 0.2989*float64(red) + 0.5870*float64(green) + 0.1140*float64(blue)
}

type kernel [9]float64

func (k kernel) get(row, col int) float64 {
	return k[row*3+col]
}

var kernel_x = kernel{
	1.0, 0.0, -1.0,
	2.0, 0.0, -2.0,
	1.0, 0.0, -1.0,
}

var kernel_y = kernel{
	1.0, 2.0, 1.0,
	0.0, 0.0, 0.0,
	-1.0, -2.0, -1.0,
}

func sobelGray(grayImage float64Image) (float64Image, float64, float64) {
	width := grayImage.width - 2
	height := grayImage.height - 2
	convolved := makeGrayImage(width, height)
	min, max := math.MaxFloat64, 0.0
	for row := 1; row < grayImage.height-1; row++ {
		for col := 1; col < grayImage.width-1; col++ {
			value_x := convolvePixel(grayImage, kernel_x, row, col)
			value_y := convolvePixel(grayImage, kernel_y, row, col)
			value := math.Sqrt(value_x*value_x + value_y*value_y)
			if min > value {
				min = value
			}
			if max < value {
				max = value
			}
			convolved.set(row-1, col-1, value)
		}
	}
	return convolved, min, max
}

func convolvePixel(img float64Image, kernel kernel, row, col int) float64 {
	var value float64
	for x, kx := col-1, 2; x <= col+1; x, kx = x+1, kx-1 {
		for y, ky := row-1, 2; y <= row+1; y, ky = y+1, ky-1 {
			value += float64(img.get(y, x)) * kernel.get(ky, kx)
		}
	}
	return value
}

func toRGBAImage(grayImage float64Image, min float64, max float64) *image.RGBA {
	result := image.NewRGBA(image.Rect(0, 0, grayImage.width, grayImage.height))
	for x := 0; x < grayImage.width; x++ {
		for y := 0; y < grayImage.height; y++ {
			value := grayImage.get(y, x)
			outValue := byte((value - min) / (max - min) * 255)
			result.Set(x, y, color.RGBA{outValue, outValue, outValue, 255})
		}
	}
	return result
}
