package main

import (
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
	"runtime"
	"testing"
)

func TestSobel(t *testing.T) {
	sourceImage, _ := getImageFromFilePath("skyline.jpg")

	resultRGBA := sobelRGBA(sourceImage)

	result := image.NewGray(resultRGBA.Bounds())
	draw.Draw(result, resultRGBA.Bounds(), resultRGBA, resultRGBA.Bounds().Min, draw.Src)

	expected, _ := getGrayImageFromFilePath("expected.png")

	if !expected.Bounds().Max.Eq(result.Bounds().Max) {
		t.Fatalf("Expected size %s, but was %s", expected.Bounds().Max, result.Bounds().Max)
	}

	for y := 0; y < expected.Bounds().Max.Y; y++ {
		for x := 0; x < expected.Bounds().Max.X; x++ {
			if delta(expected.GrayAt(x, y).Y, result.GrayAt(x, y).Y) > 1 {
				writeImageToFile("result.png", result)
				t.Fatalf("The result image differs from the expected image at (%d, %d). Check result.jpg if it is visually OK, then rename it to expected.jpg", x, y)
			}
		}
	}
}

func BenchmarkSobel(b *testing.B) {
	sourceImage, _ := getImageFromFilePath("skyline.jpg")
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		runtime.KeepAlive(sobelRGBA(sourceImage))
	}
}

func delta(a, b uint8) int {
	d := int(a) - int(b)
	if d < 0 {
		return -d
	}
	return d
}

func getImageFromFilePath(filePath string) (*image.RGBA, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, err := jpeg.Decode(f)
	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, img.Bounds(), img, img.Bounds().Min, draw.Src)
	return rgba, err
}

func getGrayImageFromFilePath(filePath string) (*image.Gray, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, err := png.Decode(f)
	gray := image.NewGray(img.Bounds())
	draw.Draw(gray, img.Bounds(), img, img.Bounds().Min, draw.Src)
	return gray, err
}

func writeImageToFile(filePath string, img image.Image) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	err = png.Encode(f, img)
	return err
}
