package jvozba

import (
	"bytes"
	"fmt"
)

func IsCmavo(what []byte) bool {
	if len(what) < 2 || !(IsConsonant(what[0]) || what[0] == 'i' || what[0] == 'u') || !(IsVowel(what[1]) || what[1] == 'y') {
		return false
	}
	previous := what[1]
	for _, c := range what[2:] {
		if previous == '\'' {
			if !bytes.Contains([]byte("aeiouy"), []byte{c}) {
				return false
			}
		} else {
			switch c {
			case 'i':
				if !(previous == 'a' || previous == 'e' || previous == 'o') {
					return false
				}
			case 'u':
				if previous != 'a' {
					return false
				}
			default:
				if c != '\'' {
					return false
				}
			}
		}
		previous = c
	}
	return previous != '\''
}

func IsGismu(what []byte) bool {
	return len(what) == 5 &&
		IsConsonant(what[0]) && IsConsonant(what[3]) && IsVowel(what[4]) &&
		((IsVowel(what[1]) && IsConsonant(what[2])) ||
			(IsConsonant(what[1]) && IsVowel(what[2])))
}

func isLujvoInitial(what []byte) bool {
	raf, acc := []byte{}, what
	for len(acc) > 0 {
		raf, acc = katna(acc)
		if rafsiTarmi(raf) == fuhivla {
			return false
		}
	}
	return true
}

// Convert a string (that is expected to be valid Lojban text) to a []byte
// slice with appropriate replacements applied to it.
func Byter(s string) ([][]byte, error) {
	parts := make([][]byte, 0, len(s)/6)
	base := 0
	for i, r := range s + "\000" {
		switch {
		case (r >= 'a' && r <= 'z' && r != 'q' && r != 'w'):
			continue
		case (r >= 'A' && r <= 'Z' && r != 'Q' && r != 'W'):
			continue
		case r == '\'' || r == ',' || r == '.' || r == '’' || r == '-':
			continue
		case r == ' ' || r == '\t' || r == '\r' || r == '\n' || r == '\000':
			if i > base {
				parts = append(parts, []byte(s[base:i]))
			}
			base = i + 1
		default:
			return [][]byte{}, fmt.Errorf("unexpected character %v", r)
		}
	}
	for i, p := range parts {
		j := 0
		offset := 0
		lenp := len(p)
		for ; j < lenp; j++ {
			switch p[j] {
			// ’ U+2019 = e2 80 99
			case 0xe2:
				p[j] = '\''
				copy(p[j+1:], p[j+3:])
				lenp -= 2
			case '-':
				if j != 0 && j != lenp-1 {
					return [][]byte{}, fmt.Errorf("misplaced hyphen in %s", p[:len(p)+offset])
				}
			case 'h', 'H':
				p[j] = '\''
			}
			if p[j] >= 'A' && p[j] <= 'Z' {
				p[j] += 'a' - 'A'
			}
		}
		parts[i] = p[:lenp]
	}
	return parts, nil
}

func Selci(tanru [][]byte, rafste map[string][]string, config Config) ([][][]byte, error) {
	if config&(Brivla|Cmevla) == 0 {
		return [][][]byte{}, fmt.Errorf("neither Cmevla nor Brivla was specified")
	}

	selci := make([][][]byte, 0, len(tanru))
	for i, p := range tanru {
		final := i == len(tanru)-1
		var r [][]byte

		if p[0] == '-' {
			r = [][]byte{p[1 : len(p)-1]}
		} else if p[len(p)-1] == '-' {
			r = [][]byte{p[:len(p)-1]}
		} else {
			r_ := rafste[string(p)]
			r = make([][]byte, len(r_))
			for i, rafsi := range r_ {
				r[i] = []byte(rafsi)
			}
			if !IsCmavo(p) {
				canShort, canLong := config&Cmevla == Cmevla, config&Brivla == Brivla
				if !final {
					canShort, canLong = true, false
				}
				midPrefix := []byte("y'")
				if i == 0 || IsGismu(p) {
					midPrefix = []byte{}
				} else if IsConsonant(p[0]) || (IsVowel(p[1]) && (p[0] == 'i' || p[0] == 'u')) {
					midPrefix = []byte{'y'}
				}
				if canShort {
					if !IsGismu(p) && ((len(midPrefix) < 2 && isLujvoInitial(p[:len(p)-1])) || config&LongFuhivla == LongFuhivla) {
						if !final {
							r = append(r, bytes.Join([][]byte{midPrefix, p, []byte("'y")}, []byte{}))
						}
					} else if !(IsGismu(p) && bytes.Equal(p[:4], []byte("brod")) && p[4] != 'a') {
						if final || IsGismu(p) {
							r = append(r, bytes.Join([][]byte{midPrefix, p[:len(p)-1]}, []byte{}))
						} else {
							r = append(r, bytes.Join([][]byte{midPrefix, p[:len(p)-1], {'y'}}, []byte{}))
						}
					}
				}
				if canLong {
					r = append(r, bytes.Join([][]byte{midPrefix, p}, []byte{}))
				}
			}
		}

		filtered := make([][]byte, 0, len(r))
		for _, one := range r {
			keep := true
			switch rafsiTarmi(one) {
			case ccv, cvv, cvhv:
				if final {
					keep = config&Brivla == Brivla
				} else if i == 0 {
					keep = IsGismu(tanru[1]) || IsCmavo(tanru[1])
				}
			case cvcc, ccvc, cvc:
				if final {
					keep = config&Cmevla == Cmevla
				}
			}
			if keep {
				filtered = append(filtered, one)
			}
		}
		if len(filtered) == 0 {
			return [][][]byte{}, fmt.Errorf("no applicable rafsi found for %s", p)
		}
		selci = append(selci, filtered)
	}
	return selci, nil
}

// Zbasu is like Jvozba, but it allows you to specify your own list of rafsi.
func Zbasu(tanru string, rafste map[string][]string, config Config) (string, error) {
	byted, err := Byter(tanru)
	if err != nil {
		return "", err
	}
	slemei, err := Selci(byted, rafste, config)
	if err != nil {
		return "", err
	}
	res, err := Lujvo(slemei)
	return string(res), err
}
