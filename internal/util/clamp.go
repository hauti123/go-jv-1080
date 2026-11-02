// Package util provides small shared helpers used by several layers
// (widgets, device layer, commands).
package util

import "cmp"

// Clamp returns v limited to the range [lo, hi].
func Clamp[T cmp.Ordered](v, lo, hi T) T {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
