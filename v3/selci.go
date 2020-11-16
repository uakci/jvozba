package jvozba

import (
	"bytes"
	"fmt"
)

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
