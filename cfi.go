package cfi

import (
	"strconv"
)

// http://www.iotafinance.com/en/Classification-of-Financial-Instrument-codes-CFI-ISO-10962.html
// https://en.wikipedia.org/wiki/ISO_10962

type CFI struct {
	Type       string
	Subtype    string
	Attributes []Attribute
}

type Attribute struct {
	Position int
	Symbol   string
	Name     string
	Value    string
}

func Decode(code string) CFI {
	st := types[code[:2]]

	return CFI{
		Type:       types[code[:1]],
		Subtype:    st,
		Attributes: decodeAttributes(code[:2], code[2:]),
	}
}

func decodeAttributes(st string, code string) (r []Attribute) {
	for i, c := range code {
		k := st + strconv.Itoa(i+1)

		n, ok := attributes[k]
		if !ok {
			continue
		}

		v, ok := attributes[k+string(c)]
		if !ok {
			continue
		}

		a := Attribute{
			Position: i + 1,
			Symbol:   string(c),
			Name:     n,
			Value:    v,
		}

		r = append(r, a)
	}

	return
}
