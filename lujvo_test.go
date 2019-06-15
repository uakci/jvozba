package jvozba

import (
	"fmt"
	"strings"
	"testing"
)

func TestHyphenation(t *testing.T) {
	examples := []struct {
		whole string
		parts string
		types string
	}{
		{"zbasai", "zba-sai", "78"},
		{"nunynau", "nun-y-nau", "508"},
		{"sairzbata'u", "sai-r-zba-ta'u", "8076"},
		{"zbazbasysarji", "zba-zbas-y-sarji", "7401"}}
	for i, e := range examples {
		total := []string{}
		types := ""
		s := e.whole
		var sle string
		for len(s) > 0 {
			sle, s = katna(s)
			total = append(total, sle)
			types += fmt.Sprintf("%d", int(rafsiTarmi(sle)))
		}
		attempt := strings.Join(total, "-")
		if e.parts != attempt || e.types != types {
			t.Errorf("'%s' (example #%d): got %s %s, expected %s %s",
				e.whole, i, attempt, types, e.parts, e.types)
		}
	}
}

// see https://lojban.github.io/cll/4/12/
func TestCLLScoring(t *testing.T) {
	examples := []struct {
		string
		int
	}{
		{"zbasai", 5847},
		{"nunynau", 6967},
		{"sairzbata'u", 10385},
		{"zbazbasysarji", 12976}}
	for i, e := range examples {
		s := Score(e.string)
		if s != e.int {
			t.Errorf("'%s' (example #%d): got %d, expected %d",
				e.string, i, s, e.int)
		}
	}
}

// see https://lojban.github.io/cll/4/13/
func TestCLLExamples(t *testing.T) {
	examples := []struct {
		selci [][]string
		lujvo string
		score int
	}{
		{[][]string{{"ger", "ge'u", "gerk"}, {"zda", "zdani"}}, "gerzda", 5878},
		{[][]string{{"lot", "blo", "lo'i"}, {"kle", "lei"}}, "blolei", 5847},
		{[][]string{{"loj", "logj"}, {"ban", "bau", "bang"}, {"gri", "girzu"}}, "lojbaugri", 8796},
		{[][]string{{"loj", "logj"}, {"ban", "bau", "bang"}, {"gir", "girz"}}, "lojbaugir", 8816},
		{[][]string{{"nak", "nakn"}, {"kem"}, {"cin", "cins"}, {"ctu", "ctuca"}}, "nakykemcinctu", 12876}}
	for i, e := range examples {
		s, err := Lujvo(e.selci)
		if err != nil {
			t.Errorf("(example #%d): error %v", i, err)
			continue
		}
		score := Score(s)
		if s != e.lujvo || score != e.score {
			t.Errorf("(example #%d): got %s (%d), expected %s (%d)",
				i, s, score, e.lujvo, e.score)
		}
	}
}

func TestIsTosmabruInitial(t *testing.T) {
	examples := map[string]bool{
		"":          true,
		"jacnal":    true,
		"jacnalsel": false,
		"skebap":    false,
		"skebai":    false,
		"pevrisn":   false,
	}
	for k, v := range examples {
		res := isTosmabruInitial(k)
		if res != v {
			t.Errorf("'%s': got %v",
				k, res)
		}
	}
}
