package io

import (
	"math"

	"github.com/patrickhuber/go-wasm/abi/values"
)

func CanonicalizeFloat32(f float32) float32 {
	if math.IsNaN(float64(f)) {
		return math.Float32frombits(values.Float32Nan)
	}
	return f
}

func CanonicalizeFloat64(f float64) float64 {
	if math.IsNaN(f) {
		return math.Float64frombits(values.Float64Nan)
	}
	return f
}
