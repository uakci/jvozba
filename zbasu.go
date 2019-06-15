package jvozba

import (
	"fmt"
	"strings"
)

func isGismu(what string) bool {
	return len(what) == 5 &&
		!isVowel(what[0]) &&
		!isVowel(what[3]) &&
		isVowel(what[4]) &&
		isVowel(what[1]) != isVowel(what[2])
}

func selci(tanru string, rafste map[string][]string, isCmene bool) ([][]string, error) {
	parts := strings.Split(tanru, " ")
	count := 0
	selci := make([][]string, len(parts))
	for i, p := range parts {
		if p == "" {
			continue
		}
		r := rafste[p]
		if r == nil && !isGismu(p) {
			return [][]string{}, fmt.Errorf("no rafsi found for %s", p)
		}
		if isGismu(p) {
			r = append(r, p, p[:4])
		}
		filtered := make([]string, len(r))
		c := 0
		for _, one := range r {
			var keep bool
			switch rafsiTarmi(one) {
			case CVCCV, CCVCV:
				keep = i == len(parts)-1 && !isCmene
			case CCV, CVV, CVhV:
				keep = !(i == len(parts)-1 && isCmene)
			case CVCC, CCVC, CVC:
				keep = !(i == len(parts)-1 && !isCmene)
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

func Zbasu(tanru string, rafste map[string][]string, isCmene bool) (string, error) {
	slemei, err := selci(tanru, rafste, isCmene)
	if err != nil {
		return "", err
	}
	res, err := Lujvo(slemei)
	if err != nil {
		return "", err
	}
	return res, nil
}
