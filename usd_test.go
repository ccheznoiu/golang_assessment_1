package main

import (
	"encoding/json"
	"testing"
)

func TestUSD(t *testing.T) {
	for _, tc := range [][]string{
		{"1.0", "0.99", "1.99", "0.01", "0.99", "1.01"},
		{"1.0", "0.994", "1.99", "0.01", "0.99", "1.01"},
		{"1.0", "0.995", "2.00", "0.00", "1.00", "1.00"},
		{"1.01", "0.01", "1.02", "1.00", "0.01", "101.00"},
		{"1.014", "0.005", "1.02", "1.01", "0.00", "10.14"},
		{"1.015", "0.005", "1.03", "1.01", "0.01", "2.03"},
		{"1", "1", "2.00", "0.00", "1.00", "1.00"},
	} {
		var a, b USD
		if err := json.Unmarshal([]byte(tc[0]), &a); err != nil {
			t.Error(err)
		} else if err = json.Unmarshal([]byte(tc[1]), &b); err != nil {
			t.Error(err)
		} else if sum, err := json.Marshal(a + b); err != nil {
			t.Error(err)
		} else if act := string(sum); act != tc[2] {
			t.Errorf("expected %s + %s = %s, but got %s", tc[0], tc[1], tc[2], act)
		}
	}
}
