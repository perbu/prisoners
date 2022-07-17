package main

import (
	"testing"
)

func Test_newPool(t *testing.T) {
	m := make(map[int]int)
	p := newPool()
	for b := range p {
		m[b]++
	}
	for i := 0; i < prisoners; i++ {
		v := m[i]
		if v != 1 {
			t.Errorf("%d: %d", i, v)
		}
	}
}
