rm -rf build/
mkdir build
GOOS=js GOARCH=wasm go build -o build/convolve.wasm .
cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" build/wasm_exec.js
cp index.html build/
cp *.jpg build/