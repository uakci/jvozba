package jvozba

import (
	"bytes"
	"fmt"
	"strings"
)

func isCmavo(what []byte) bool {
	if len(what) < 2 || !bytes.Contains([]byte("bcdfgjklmnprstvxziu"), []byte{what[0]}) || !bytes.Contains([]byte("aeiouy"), []byte{what[1]}) {
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

func isGismu(what []byte) bool {
	return len(what) == 5 &&
		isVowel(what[4]) &&
		isConsonant(what[0]) && isConsonant(what[3]) && ((isVowel(what[1]) && isConsonant(what[2])) ||
		(isConsonant(what[1]) && isVowel(what[2])))
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

func Selci(tanru string, rafste map[string][]string, config Config) ([][][]byte, error) {
	if config&(Brivla|Cmevla) == 0 {
		return [][][]byte{}, fmt.Errorf("neither Cmevla nor Brivla was specified")
	}
	dirty_parts := strings.Fields(tanru)
	parts := make([][]byte, 0, len(dirty_parts))
	for _, p := range dirty_parts {
		if len(p) > 0 {
			neu := make([]byte, 0, len(p))
			for i, c := range p {
				if c == '’' || c == 'h' {
					neu = append(neu, '\'')
				} else if (c >= 'a' && c <= 'z' && c != 'q' && c != 'w') || c == '\'' {
					neu = append(neu, byte(c))
				} else if !(i == 0 && c == '.') {
					return [][][]byte{}, fmt.Errorf("invalid character: %v", c)
				}
			}
			parts = append(parts, neu)
		}
	}

	selci := make([][][]byte, 0, len(parts))
	for i, p := range parts {
		final := i == len(parts)-1
		r_ := rafste[string(p)]
		r := make([][]byte, len(r_))
		for i, rafsi := range r_ {
			r[i] = []byte(rafsi)
		}

		if !isCmavo(p) {
			canShort, canLong := config&Cmevla == Cmevla, config&Brivla == Brivla
			if !final {
				canShort, canLong = true, false
			}
			midPrefix := []byte("y'")
			if i == 0 || isGismu(p) {
				midPrefix = []byte{}
			} else if isConsonant(p[0]) || (isVowel(p[1]) && (p[0] == 'i' || p[0] == 'u')) {
				midPrefix = []byte{'y'}
			}
			if canShort {
				if !isGismu(p) && ((len(midPrefix) < 2 && isLujvoInitial(p[:len(p)-1])) || config&LongFuhivla == LongFuhivla) {
					if !final {
						r = append(r, append(append(append([]byte{}, midPrefix...), p...), "'y"...))
					}
				} else if !(isGismu(p) && bytes.Equal(p[:4], []byte("brod")) && p[4] != 'a') {
					if final || isGismu(p) {
						r = append(r, append(append([]byte{}, midPrefix...), p[:len(p)-1]...))
					} else {
						r = append(r, append(append(append([]byte{}, midPrefix...), p[:len(p)-1]...), 'y'))
					}
				}
			}
			if canLong {
				r = append(r, append(midPrefix, p...))
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
					keep = isGismu(parts[1]) || isCmavo(parts[1])
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
	slemei, err := Selci(tanru, rafste, config)
	if err != nil {
		return "", err
	}
	res, err := Lujvo(slemei)
	return string(res), err
}
