package jvozba

import (
	"bytes"
	"fmt"
)

func Katna(lujvo []byte) (result [][]byte) {
	chunk := make([][]byte, 0, len(lujvo)/3)
	fuhivlaTainted, i := false, 0
	var rafsi []byte
	lastRun := false
	for len(lujvo) > 0 || !lastRun {
		if len(lujvo) == 0 {
			lastRun = true
		}
		var tai tarmi
		if lastRun {
			rafsi = []byte{'y'}
			tai = hyphen
		} else {
			rafsi, lujvo = katna(lujvo)
			tai = rafsiTarmi(rafsi)
		}
		switch tai {
		case hyphen:
			if rafsi[len(rafsi)-1] != 'y' {
				if i != 1 {
					chunk = append(chunk, rafsi)
					fuhivlaTainted = true
				}
			} else {
				if fuhivlaTainted {
					res := bytes.Join(chunk, []byte{})
					if res[len(res)-1] == '\'' {
						res = res[:len(res)-1]
					} else if res[len(res)-1] != '-' && !(IsVowel(res[len(res)-1]) && len(lujvo) == 0) {
						res = append(res, '-')
					}
					result = append(result, res)
				} else {
					result = append(result, chunk...)
				}
				chunk = [][]byte{}
				fuhivlaTainted = false
			}
		case fuhivla:
			fuhivlaTainted = true
			fallthrough
		default:
			if rafsi[0] == '\'' {
				rafsi = rafsi[1:]
			}
			if tai == cvcc || tai == ccvc {
				rafsi = append(append(make([]byte, 0, len(rafsi)+1), rafsi...), '-')
			}
			chunk = append(chunk, rafsi)
		}
		i++
	}
	return result
}

func Veljvo(lujvo string) ([]string, error) {
	pieces, err := Normalize(lujvo)
	if err != nil {
		return []string{}, err
	}
	if len(pieces) != 1 {
		return []string{}, fmt.Errorf("unexpected number of space-delimited elements %d", len(pieces))
	}
	rafpoi := Katna(pieces[0])
	tanru := make([]string, len(rafpoi))
	for i, raf := range rafpoi {
		ok := false
		if IsGismu(raf) || (len(raf) == 4 && IsGismu(append(append(make([]byte, 0, 5), raf...), 'a'))) {
			for selrafsi := range Rafsi {
				if len(selrafsi) == 5 && selrafsi[:4] == string(raf[:4]) {
					tanru[i] = selrafsi
					ok = true
					break
				}
			}
		}
		if !ok {
			for selrafsi, rafsiporsi := range Rafsi {
				for _, rafsi := range rafsiporsi {
					if rafsi == string(raf) {
						tanru[i] = selrafsi
						ok = true
						break
					}
				}
				if ok {
					break
				}
			}
		}
		if !ok {
			if len(raf) >= 4 && !IsCmavo(raf) {
				if IsConsonant(raf[len(raf)-1]) {
					tanru[i] = string(raf) + "-"
				} else {
					tanru[i] = string(raf)
				}
			} else {
				tanru[i] = fmt.Sprintf("-%s-", raf)
			}
		}
	}
	return tanru, nil
}
