package jvozba

import (
	"fmt"
	"testing"
)

func TestGeneration(t *testing.T) {
	type example struct {
		config Config
		tanru  string
		selci  [][]string
	}
	examples := []example{
		example{Brivla | Cmevla, "mi prenu", [][]string{{"mib"}, {"pre", "pren", "prenu"}}},
		example{Brivla | Cmevla, "gerku zdani", [][]string{{"ge'u", "ger", "gerk"}, {"zda", "zdan", "zdani"}}},
		example{Brivla | Cmevla, "bloti klesi", [][]string{{"blo", "lo'i", "lot", "blot"}, {"kle", "lei", "kles", "klesi"}}},
		example{Brivla | Cmevla, "logji bangu girzu", [][]string{{"loj", "logj"}, {"ban", "bau", "bang"}, {"gir", "gri", "girz", "girzu"}}},
		example{Brivla | Cmevla, "nakni ke cinse ctuca", [][]string{{"nak", "nakn"}, {"kem"}, {"cin", "cins"}, {"ctu", "ctuc", "ctuca"}}},
		example{Brivla, "alma nelci", [][]string{{"almy"}, {"nei", "nelci"}}},
		example{Brivla | LongFuhivla, "alma nelci", [][]string{{"alma'y"}, {"nei", "nelci"}}},
		example{Brivla, "barduku nelci", [][]string{{"barduku'y"}, {"nei", "nelci"}}},
		example{Brivla | LongFuhivla, "barduku nelci", [][]string{{"barduku'y"}, {"nei", "nelci"}}},
	}
	for i, e := range examples {
		s, err := selci(e.tanru, Rafsi, e.config)
		if err != nil {
			t.Errorf("(example #%d): error %v", i, err)
		} else if fmt.Sprintf("%v", s) != fmt.Sprintf("%v", e.selci) {
			t.Errorf("(example #%d): got %v, expected %v (config=%v)", i, s, e.selci, e.config)
		}
	}
}
