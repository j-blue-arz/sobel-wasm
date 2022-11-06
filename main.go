//go:build wasm

package main

import "bytes"
import "syscall/js"
import "image"
import _ "image/jpeg"
import "image/png"

func sobelOperator(this js.Value, args []js.Value) interface{} {
	inputBuffer := make([]byte, args[0].Get("byteLength").Int())
	js.CopyBytesToGo(inputBuffer, args[0])

	img, _, _ := image.Decode(bytes.NewReader(inputBuffer))

	resultImage := sobelRGBA(img)

	var outputBuffer bytes.Buffer
	png.Encode(&outputBuffer, resultImage) // todo: check error?

	outputBytes := outputBuffer.Bytes()
	size := len(outputBytes)
	result := js.Global().Get("Uint8Array").New(size)
	js.CopyBytesToJS(result, outputBytes)

	return result
}

func main() {
	js.Global().Set("sobelOperator", js.FuncOf(sobelOperator))

	<-make(chan bool)
}
