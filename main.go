//go:build wasm

package main

import "syscall/js"
import "image"

func sobelOperator(this js.Value, args []js.Value) interface{} {
	inputBuffer := make([]byte, args[0].Get("length").Int())
	js.CopyBytesToGo(inputBuffer, args[0])
	width := args[1].Int()
	height := args[2].Int()

	inputImage := image.NewRGBA(image.Rect(0, 0, width, height))
	inputImage.Pix = inputBuffer

	resultImage := sobelRGBA(*inputImage)

	size := len(resultImage.Pix)
	result := js.Global().Get("Uint8ClampedArray").New(size)
	js.CopyBytesToJS(result, resultImage.Pix)

	return result
}

func main() {
	js.Global().Set("sobelOperator", js.FuncOf(sobelOperator))

	<-make(chan bool)
}
