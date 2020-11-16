package jvozba

import (
	"fmt"
	"strings"
	"testing"
)

func byter(s string) [][]byte {
	b, e := Normalize(s)
	if e != nil {
		panic(e)
	}
	return b
}

func TestGeneration(t *testing.T) {
	type example struct {
		config Config
		tanru  string
		selci  [][]string
	}
	examples := []example{
		{Brivla | Cmevla, "co’e cohe", [][]string{{"co'e", "com"}, {"co'e", "com"}}},
		{Brivla | Cmevla, "mi prenu", [][]string{{"mib"}, {"pre", "pren", "prenu"}}},
		{Brivla | Cmevla, "gerku zdani", [][]string{{"ge'u", "ger", "gerk"}, {"zda", "zdan", "zdani"}}},
		{Brivla | Cmevla, "bloti klesi", [][]string{{"blo", "lo'i", "lot", "blot"}, {"kle", "lei", "kles", "klesi"}}},
		{Brivla | Cmevla, "logji bangu girzu", [][]string{{"loj", "logj"}, {"ban", "bau", "bang"}, {"gir", "gri", "girz", "girzu"}}},
		{Brivla | Cmevla, "nakni ke cinse ctuca", [][]string{{"nak", "nakn"}, {"kem"}, {"cin", "cins"}, {"ctu", "ctuc", "ctuca"}}},
		{Brivla, "bakni rectu", [][]string{{"bak", "bakn"}, {"re'u", "rectu"}}},
		{Brivla, "kukte zukte", [][]string{{"kuk", "kukt"}, {"zu'e", "zukte"}}},
		{Cmevla, "kukte zukte", [][]string{{"kuk", "kukt"}, {"zuk", "zukt"}}},
		{Brivla, "alma nelci", [][]string{{"almy"}, {"nei", "nelci"}}},
		{Brivla | LongFuhivla, "alma nelci", [][]string{{"alma'y"}, {"nei", "nelci"}}},
		{Brivla, "barduku nelci", [][]string{{"barduku'y"}, {"nei", "nelci"}}},
		{Brivla | LongFuhivla, "barduku nelci", [][]string{{"barduku'y"}, {"nei", "nelci"}}},
	}
	for i, e := range examples {
		s, err := Selci(byter(e.tanru), Rafsi, e.config)
		s_ := make([][]string, len(s))
		for i, sss := range s {
			ss := make([]string, len(sss))
			for j, ssss := range sss {
				ss[j] = string(ssss)
			}
			s_[i] = ss
		}
		if err != nil {
			t.Errorf("(example #%d): error %v", i, err)
		} else if fmt.Sprintf("%v", s_) != fmt.Sprintf("%v", e.selci) {
			t.Errorf("(example #%d): got %v, expected %v (config=%v)", i, s_, e.selci, e.config)
		}
	}
}

func BenchmarkSelciWithBloti1000(b *testing.B) {
	bloti := strings.Repeat("bloti ", 1000)
	bl, _ := Normalize(bloti)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Selci(bl, Rafsi, Brivla)
	}
}

func BenchmarkSelciWithBlaci1000(b *testing.B) {
	blaci := strings.Repeat("blaci ", 1000)
	bl, _ := Normalize(blaci)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Selci(bl, Rafsi, Brivla)
	}
}
