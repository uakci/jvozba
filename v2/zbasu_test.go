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
		example{"mi prenu", [][]string{{"mib"}, {"pre", "pren", "prenu"}}},
		example{"gerku zdani", [][]string{{"ge'u", "ger", "gerk"}, {"zda", "zdan", "zdani"}}},
		example{"bloti klesi", [][]string{{"blo", "lo'i", "lot", "blot"}, {"kle", "lei", "kles", "klesi"}}},
		example{"logji bangu girzu", [][]string{{"loj", "logj"}, {"ban", "bau", "bang"}, {"gir", "gri", "girz", "girzu"}}},
		example{"nakni ke cinse ctuca", [][]string{{"nak", "nakn"}, {"kem"}, {"cin", "cins"}, {"ctu", "ctuc", "ctuca"}}}}
	for i, e := range examples {
		s, err := selci(e.tanru, Rafsi, Brivla|Cmevla)
		if err != nil {
			t.Errorf("(example #%d): error %v", i, err)
		} else if fmt.Sprintf("%v", s) != fmt.Sprintf("%v", e.selci) {
			t.Errorf("(example #%d): got %v, expected %v", i, s, e.selci)
		}
	}
}
