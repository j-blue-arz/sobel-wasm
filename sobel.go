package main

import (
	"image"
	"math"
)

type numeric interface {
	byte | float64
}

// the buffer has size width*height*numChannels.
type canvasImage[pixelType numeric] struct {
	buffer      []pixelType
	numChannels int
	width       int
	height      int
}

func makeGrayImage(width, height int) canvasImage[float64] {
	return canvasImage[float64]{make([]float64, width*height), 1, width, height}
}

type channelIndex int

func (image canvasImage[_]) index(row, col int, c channelIndex) int {
	return (row*image.width+col)*image.numChannels + int(c)
}

func (image canvasImage[pixelType]) get(row, col int, c channelIndex) pixelType {
	return image.buffer[image.index(row, col, c)]
}

func (image canvasImage[pixelType]) set(row, col int, c channelIndex, value pixelType) {
	image.buffer[image.index(row, col, c)] = value
}

// The returned image has its size reduced by 2 in both dimensions.
func sobelRGBA(rgba image.RGBA) *image.RGBA {
	width, height := rgba.Bounds().Dx(), rgba.Bounds().Dy()

	sourceImage := canvasImage[byte]{rgba.Pix, 4, width, height}
	grayImage := toGrayImage(sourceImage)
	convolved, min, max := sobelGray(grayImage)
	return toRGBAImage(convolved, min, max)
}

func toGrayImage(img canvasImage[byte]) canvasImage[float64] {
	grayImage := makeGrayImage(img.width, img.height)
	for pixel := 0; pixel < len(grayImage.buffer); pixel++ {
		colors := img.buffer[pixel*4 : pixel*4+3]
		red, green, blue := colors[0], colors[1], colors[2]
		grayImage.buffer[pixel] = float64(toGray(red, green, blue))
	}
	return grayImage
}

func toGray(red, green, blue byte) float64 {
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

func sobelGray(grayImage canvasImage[float64]) (canvasImage[float64], float64, float64) {
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
			convolved.set(row-1, col-1, 0, value)
		}
	}
	return convolved, min, max
}

func convolvePixel[pixelType numeric](img canvasImage[pixelType], kernel kernel, row, col int) float64 {
	var value float64
	for x, kx := col-1, 2; x <= col+1; x, kx = x+1, kx-1 {
		for y, ky := row-1, 2; y <= row+1; y, ky = y+1, ky-1 {
			value += float64(img.get(y, x, 0)) * kernel.get(ky, kx)
		}
	}
	return value
}

func toRGBAImage(grayImage canvasImage[float64], min float64, max float64) *image.RGBA {
	result := image.NewRGBA(image.Rect(0, 0, grayImage.width, grayImage.height))
	for pixel, value := range grayImage.buffer {
		outValue := byte((value - min) / (max - min) * 255)
		result.Pix[pixel*4] = outValue
		result.Pix[pixel*4+1] = outValue
		result.Pix[pixel*4+2] = outValue
		result.Pix[pixel*4+3] = 255
	}
	return result
}
