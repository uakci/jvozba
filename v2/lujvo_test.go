package jvozba

import (
	"bytes"
	"strings"
	"testing"
)

func TestHyphenation(t *testing.T) {
	examples := []struct {
		whole string
		parts string
		types string
	}{
		{"zuky'u'enca'yzukte", "zuk-y-'u'e-n-ca'-y-zukte", "5090901"},
		{"barduku'ynei", "bar-duku'-y-nei", "5908"},
		{"zbasai", "zba-sai", "78"},
		{"nunynau", "nun-y-nau", "508"},
		{"sairzbata'u", "sai-r-zba-ta'u", "8076"},
		{"zbazbasysarji", "zba-zbas-y-sarji", "7401"},
		{"lojbaugri", "loj-bau-gri", "587"}}
	for i, e := range examples {
		total := [][]byte{}
		types := []byte{}
		s := []byte(e.whole)
		var sle []byte
		for len(s) > 0 {
			sle, s = katna([]byte(s))
			total = append(total, sle)
			types = append(types, '0'+byte(rafsiTarmi(sle)))
		}
		attempt := bytes.Join(total, []byte{'-'})
		if !bytes.Equal([]byte(e.parts), attempt) || !bytes.Equal([]byte(e.types), types) {
			t.Errorf("'%s' (example #%d): got %s %s, expected %s %s",
				e.whole, i, string(attempt), string(types), e.parts, e.types)
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
		s := Score([]byte(e.string))
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
		newSelci := make([][][]byte, len(e.selci))
		for i, u := range e.selci {
			newUnit := make([][]byte, len(u))
			for j, u := range e.selci[i] {
				newUnit[j] = []byte(u)
			}
			newSelci[i] = newUnit
		}
		s, err := Lujvo(newSelci)
		if err != nil {
			t.Errorf("(example #%d): error %v", i, err)
			continue
		}
		score := Score(s)
		if string(s) != e.lujvo || score != e.score {
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
		res := isTosmabruInitial([]byte(k))
		if res != v {
			t.Errorf("'%s': got %v",
				k, res)
		}
	}
}

func TestKatna(t *testing.T) {
	res := Katna([]byte("co'arzukybarduku'y'arcyiarciiy'u'enca'ycizda'u"))
	got := string(bytes.Join(res, []byte{'-'}))
	if got != "co'a-zuk-barduku-arca-iarciia-u'enca-ciz-da'u" {
		t.Errorf("%s", got)
	}
}

func BenchmarkLujvoWithBloti1000(b *testing.B) {
	bloti := strings.Repeat("bloti", 1000)
	selci, _ := Selci(bloti, Rafsi, Brivla)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Lujvo(selci)
	}
}

func BenchmarkLujvoWithBlaci1000(b *testing.B) {
	blaci := strings.Repeat("blaci", 1000)
	selci, _ := Selci(blaci, Rafsi, Brivla)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Lujvo(selci)
	}
}
