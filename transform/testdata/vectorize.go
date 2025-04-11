//go:build tinygo.wasm
// +build tinygo.wasm

package main

//go:noinline
func vectorize(n int, fn func(int, int, *[16]bool)) {
	// This dummy body will be ignored and replaced by LLVM IR pass.
	for i := 0; i < n; i++ {
		fn(i, i%16, nil)
	}
}

func main() {
	dst := [4]int32{}
	src := [4]int32{1, 2, 3, 4}

	vectorize(4, func(i, lane int, mask *[16]bool) {
		dst[i] = src[i] + 1
	})
}
