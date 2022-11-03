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
	return toRGBAImage(*convolved, min, max)
}

func toGrayImage(img image.RGBA) *image.Gray16 {
	grayImage := image.NewGray16(img.Rect)
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			rgba := img.RGBAAt(x, y)
			red, green, blue := rgba.R, rgba.G, rgba.B
			grayImage.SetGray16(x, y, color.Gray16{uint16(toGray(red, green, blue))})
		}
	}
	return grayImage
}

func toGray(red, green, blue uint8) float64 {
	return 0.2989*float64(red) + 0.5870*float64(green) + 0.1140*float64(blue)
}

type kernel [9]int

func (k kernel) get(row, col int) int {
	return k[row*3+col]
}

var kernel_x = kernel{
	1, 0, -1,
	2, 0, -2,
	1, 0, -1,
}

var kernel_y = kernel{
	1, 2, 1,
	0, 0, 0,
	-1, -2, -1,
}

func sobelGray(grayImage *image.Gray16) (*image.Gray16, uint16, uint16) {
	width := grayImage.Bounds().Dx() - 2
	height := grayImage.Bounds().Dy() - 2
	convolved := image.NewGray16(image.Rect(0, 0, width, height))
	min, max := uint16(math.MaxUint16), uint16(0)
	for y := 1; y < grayImage.Bounds().Max.Y-1; y++ {
		for x := 1; x < grayImage.Bounds().Max.X-1; x++ {
			value_x := convolvePixel(*grayImage, kernel_x, y, x)
			value_y := convolvePixel(*grayImage, kernel_y, y, x)
			value := uint16(math.Sqrt(float64(value_x*value_x + value_y*value_y)))
			if min > value {
				min = value
			}
			if max < value {
				max = value
			}
			convolved.SetGray16(x-1, y-1, color.Gray16{value})
		}
	}
	return convolved, min, max
}

func convolvePixel(img image.Gray16, kernel kernel, row, col int) int {
	var value int
	for x, kx := col-1, 2; x <= col+1; x, kx = x+1, kx-1 {
		for y, ky := row-1, 2; y <= row+1; y, ky = y+1, ky-1 {
			value += int(img.Gray16At(x, y).Y) * kernel.get(ky, kx)
		}
	}
	return value
}

func toRGBAImage(img image.Gray16, min uint16, max uint16) *image.RGBA {
	result := image.NewRGBA(img.Bounds())
	valueRange := float64(max - min)
	for y := 0; y < img.Bounds().Max.Y; y++ {
		for x := 0; x < img.Bounds().Max.X; x++ {
			value := img.Gray16At(x, y).Y
			outValue := byte(float64(value-min) / valueRange * 255)
			result.Set(x, y, color.RGBA{outValue, outValue, outValue, 255})
		}
	}
	return result
}
