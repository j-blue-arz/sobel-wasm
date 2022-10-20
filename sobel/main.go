//go:build wasm

package main

import "syscall/js"

func sobelOperator(this js.Value, args []js.Value) interface{} {
	inputBuffer := make([]byte, args[0].Get("length").Int())
	js.CopyBytesToGo(inputBuffer, args[0])
	width := args[1].Int()
	height := args[2].Int()

	image := canvasImage[byte]{inputBuffer, 4, width, height}

	resultImage := sobelRGBA(image)

	size := len(resultImage.buffer)
	result := js.Global().Get("Uint8ClampedArray").New(size)
	js.CopyBytesToJS(result, resultImage.buffer)

	return result
}

func main() {
	js.Global().Set("sobelOperator", js.FuncOf(sobelOperator))

	<-make(chan bool)
}
