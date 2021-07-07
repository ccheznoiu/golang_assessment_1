package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
)

// USD is intended to convert dollars to cents for int calculation
// as opposed to float64 calculation as goes the prevailing wisdom on the matter.
// Go seems to fully lack support for fractional cents.
type USD int

// UnmarshalJSON wraps the json number marshal after which converting from dollars to cents
func (d *USD) UnmarshalJSON(data []byte) error {
	var f float64

	if i := bytes.IndexByte(data, '.'); i > 0 {
		data = bytes.Replace(data, []byte(`.`), nil, 1)
		l := len(data[i:])
		switch {
		case l == 1:
			data = append(data, []byte(`0`)...)
		case l > 2:
			data = append(data, 0)
			copy(data[i+2:], data[i+1:])
			data[i+2] = '.'
		}
		isNeg := bytes.HasPrefix(data, []byte(`-`))

		if isNeg {
			data = bytes.TrimPrefix(data, []byte(`-`))
		}

		for bytes.HasPrefix(data, []byte(`0`)) && !bytes.HasPrefix(data, []byte(`0.`)) {
			data = bytes.TrimPrefix(data, []byte(`0`))
		}

		if isNeg {
			data = append([]byte(`-`), data...)
		}

		if err := json.Unmarshal(data, &f); err != nil {
			return err
		}

		*d = USD(math.Round(f))

		return nil
	}

	if err := json.Unmarshal(data, &f); err != nil {
		return err
	}

	*d = USD(100 * f)

	return nil
}

// MarshalJSON converts cents to dollars and prints two decimal places
func (d USD) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%.2f", float64(d)/100)), nil
}
