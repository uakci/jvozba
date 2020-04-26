package jvozba

import (
	"fmt"
	"regexp"
	"strings"
)

var cmavoRegex = regexp.MustCompilePOSIX(
	"^[bcdfgjklmnprstvxziu]([aeiouy]|[aeo]i|au)('([aeiouy]|[aeo]i|au))*$")

func isCmavo(what string) bool {
	return cmavoRegex.FindIndex([]byte(what)) != nil
}

func isGismu(what string) bool {
	return len(what) == 5 &&
		isVowel(what[4]) &&
		isConsonant(what[0]) && isConsonant(what[3]) && ((isVowel(what[1]) && isConsonant(what[2])) ||
		(isConsonant(what[1]) && isVowel(what[2])))
}

func isLujvoInitial(what string) bool {
	raf, acc := "", what
	for len(acc) > 0 {
		raf, acc = katna(acc)
		if rafsiTarmi(raf) == fuhivla {
			return false
		}
	}
	return true
}

func selci(tanru string, rafste map[string][]string, config Config) ([][]string, error) {
	if config&(Brivla|Cmevla) == 0 {
		return [][]string{}, fmt.Errorf("neither Cmevla or Brivla was specified")
	}
	dirty_parts := strings.Split(tanru, " ")
	parts := make([]string, 0, len(dirty_parts))
	for _, p := range dirty_parts {
		p = strings.Trim(p, "\n")
		if p != "" {
			parts = append(parts, p)
		}
	}

	selci := make([][]string, 0, len(parts))
	for i, p := range parts {
		final := i == len(parts)-1
		r := rafste[p]

		if !isCmavo(p) {
			canShort, canLong := config&Cmevla == Cmevla, config&Brivla == Brivla
			if !final {
				canShort, canLong = true, false
			}
			midPrefix := "y'"
			if i == 0 || isGismu(p) {
				midPrefix = ""
			} else if isConsonant(p[0]) || (isVowel(p[0]) && isVowel(p[1])) {
				midPrefix = "y"
			}
			if canShort {
				if !isGismu(p) && ((len(midPrefix) < 2 && isLujvoInitial(p[:len(p)-1])) || config&LongFuhivla == LongFuhivla) {
					if !final {
						r = append(r, midPrefix+p+"'y")
					}
				} else if !(isGismu(p) && p[:4] == "brod" && p[4] != 'a') {
					if final || isGismu(p) {
						r = append(r, midPrefix+p[:len(p)-1])
					} else {
						r = append(r, midPrefix+p[:len(p)-1]+"y")
					}
				}
			}
			if canLong {
				r = append(r, midPrefix+p)
			}
		}

		filtered := make([]string, 0, len(r))
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
			return [][]string{}, fmt.Errorf("no applicable rafsi found for %s", p)
		}
		selci = append(selci, filtered)
	}
	return selci, nil
}

// Zbasu is like Jvozba, but it allows you to specify your own list of rafsi.
func Zbasu(tanru string, rafste map[string][]string, config Config) (string, error) {
	slemei, err := selci(tanru, rafste, config)
	if err != nil {
		return "", err
	}
	res, err := Lujvo(slemei)
	return res, err
}
