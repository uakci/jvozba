package jvozba

import (
	"fmt"
	"testing"
)

func TestGeneration(t *testing.T) {
	type example struct {
		tanru string
		selci [][]string
	}
	examples := []example{
		example{"gerku zdani", [][]string{{"ge'u", "ger", "gerk"}, {"zda", "zdani"}}},
		example{"bloti klesi", [][]string{{"blo", "lo'i", "lot", "blot"}, {"kle", "lei", "klesi"}}},
		example{"logji bangu girzu", [][]string{{"loj", "logj"}, {"ban", "bau", "bang"}, {"gri", "girzu"}}},
		example{"nakni ke cinse ctuca", [][]string{{"nak", "nakn"}, {"kem"}, {"cin", "cins"}, {"ctu", "ctuca"}}}}
	for i, e := range examples {
		s, err := selci(e.tanru, Rafsi, false)
		if err != nil {
			t.Errorf("(example #%d): error %v", i, err)
		} else if fmt.Sprintf("%v", s) != fmt.Sprintf("%v", e.selci) {
			t.Errorf("(example #%d): got %v, expected %v", i, s, e.selci)
		}
	}
}
