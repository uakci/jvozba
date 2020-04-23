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
	parts := strings.Split(tanru, " ")
	count := 0
	selci := make([][]string, len(parts))
	for i, p := range parts {
		p = strings.Trim(p, "\n")
		if p == "" {
			continue
		}
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
				if !isGismu(p) && len(midPrefix) < 2 && isLujvoInitial(p[:len(p)-1]) {
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

		filtered := make([]string, len(r))
		c := 0
		for _, one := range r {
			keep := true
			if final {
				switch rafsiTarmi(one) {
				case ccv, cvv, cvhv:
					keep = config&Brivla == Brivla
				case cvcc, ccvc, cvc:
					keep = config&Cmevla == Cmevla
				}
			}
			if keep {
				filtered[c] = one
				c++
			}
		}
		if c == 0 {
			return [][]string{}, fmt.Errorf("no applicable rafsi found for %s", p)
		}
		selci[count] = filtered[:c]
		count++
	}
	return selci[:count], nil
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
