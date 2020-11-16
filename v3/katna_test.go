package jvozba

import (
	"bytes"
	"regexp"
	"strings"
	"testing"
)

func TestKatna(t *testing.T) {
	const want = "co'a/zuk/barduku/arc-/iarcii-/u'enca/ciz/tong-/da'u/dat/to'a/barduku"
	res := Katna([]byte("co'arzukybarduku'y'arcyiarciiy'u'enca'ycizytongyda'udatyto'a'ybarduku"))
	have := string(bytes.Join(res, []byte{'/'}))
	if have != want {
		t.Errorf("want %s", want)
		t.Errorf("have %s", have)
	}
}

func TestVeljvo(t *testing.T) {
	for i, n := range jvozbaTests {
		res, err := Veljvo(n.lujvo)
		if err != nil {
			t.Errorf("on example #%d (%s): got error: %v", i, n.lujvo, err)
		} else if regexp.MustCompile("^"+strings.ReplaceAll(strings.Join(res, " "), "-", ".")+"$").FindStringIndex(n.tanru) == nil {
			t.Errorf("on example #%d: expected %v, got %v (input: %s)", i, n.lujvo, res, n.tanru)
		}
	}
}

func BenchmarkNormalizeWithBloti1000(b *testing.B) {
	bloti := strings.Repeat("blo", 1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Normalize(bloti)
	}
}

func BenchmarkKatnaWithBloti1000(b *testing.B) {
	bloti_, _ := Normalize(strings.Repeat("blo", 1000))
	bloti := bloti_[0]
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Katna(bloti)
	}
}

func BenchmarkVeljvoWithBloti1000(b *testing.B) {
	bloti := strings.Repeat("blo", 1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Veljvo(bloti)
	}
}
