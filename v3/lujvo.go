package jvozba

import (
	"bytes"
	"fmt"
)

// Yields the numerical score of a lujvo form, using the algorithm described in
// the Complete Logical Language.
func Score(lujvo []byte) int {
	L := len(lujvo)
	// apostrophe count, hyphen count, rafsi score count, vowel count
	A, H, R, V := 0, 0, 0, 0
	var curr []byte
	for _, b := range lujvo {
		if IsVowel(byte(b)) {
			V++
		} else if b == 'y' {
			H++
		} else if b == '\'' {
			A++
			H++
		}
	}
	for len(lujvo) > 0 {
		curr, lujvo = katna(lujvo)
		tai := rafsiTarmi(curr)
		R += int(tai)
	}
	return 1000*L - 500*A + 100*H - 10*R - V
}

// the tarmi enum directly corresponds to the scores for R
type tarmi int

const (
	hyphen tarmi = iota
	cvccv
	cvcc
	ccvcv
	ccvc
	cvc
	cvhv
	ccv
	cvv
	fuhivla
)

func rafsiTarmi(rafsi []byte) tarmi {
	l := len(rafsi)
	if l == 0 {
		return fuhivla
	} else if l == 2 && rafsi[0] == '\'' && rafsi[1] == 'y' {
		return hyphen
	} else if !IsConsonant(rafsi[0]) && l != 1 {
		return fuhivla
	}
	switch l {
	case 1:
		if len(rafsi) == 1 && (rafsi[0] == 'y' || rafsi[0] == 'r' || rafsi[0] == 'n') {
			return hyphen
		}
	case 3:
		if !IsVowel(rafsi[2]) {
			if IsVowel(rafsi[1]) && IsConsonant(rafsi[2]) {
				return cvc
			}
		} else {
			if IsVowel(rafsi[1]) {
				return cvv
			} else if IsConsonant(rafsi[1]) {
				return ccv
			}
		}
	case 4:
		if IsVowel(rafsi[1]) {
			if IsVowel(rafsi[3]) {
				if rafsi[2] == '\'' {
					return cvhv
				}
			} else if IsConsonant(rafsi[2]) && IsConsonant(rafsi[3]) {
				return cvcc
			}
		} else if IsConsonant(rafsi[1]) && IsVowel(rafsi[2]) && IsConsonant(rafsi[3]) {
			return ccvc
		}
	case 5:
		if IsGismu(rafsi) {
			if IsVowel(rafsi[2]) {
				return ccvcv
			} else {
				return cvccv
			}
		}
	}
	return fuhivla
}

// Whether `lujvo` could lead to tosmabru on appending -y or is one already.
func isTosmabruInitial(lujvo []byte) bool {
	var r []byte
	var lastChar byte
	i := 0
	for len(lujvo) > 0 {
		r, lujvo = katna(lujvo)
		switch t := rafsiTarmi(r); t {
		case cvc:
			if i > 0 && !IsValidInitial(lastChar, r[0]) {
				return false
			}
			lastChar = r[2]
		case cvccv:
			return i > 0 &&
				IsValidInitial(lastChar, r[0]) &&
				IsValidInitial(r[2], r[3])
		case hyphen:
			return i > 1 && bytes.Equal(r, []byte{'y'})
		default:
			return false
		}
		i++
	}
	return true
}

// Make one cut. Atrocious code.
func katna(lujvo []byte) ([]byte, []byte) {
	var point /* of cission */ int
	l := len(lujvo)
	switch {
	case l >= 1 && lujvo[0] == 'y':
		point = 1
	case l >= 2 && lujvo[0] == '\'' && lujvo[1] == 'y':
		point = 2
	case l >= 4 && (lujvo[0] == 'n' || lujvo[0] == 'r' || lujvo[0] == 'y') && IsConsonant(lujvo[1]):
		point = 1
	case l >= 8 && lujvo[3] == '\'' && lujvo[4] == 'y':
		point = 3
	case l >= 8 && lujvo[4] == 'y':
		point = 4
	case l >= 7 && lujvo[3] == 'y':
		point = 3
	case l >= 7 && lujvo[2] == '\'' && lujvo[3] != 'y':
		point = 4
	case l >= 6 && bytes.Index(lujvo[:6], []byte{'y'}) == -1:
		point = 3
	default:
		for point = 0; point < l; point++ {
			if lujvo[point] == 'y' {
				break
			}
		}
	}
	return lujvo[:point], lujvo[point:]
}

type scored struct {
	lujvo    []byte
	score    int
	tosmabru bool
}

// -y-: L += 1 && H += 1 -> score += 1100
const yPenalty = 1100

// Lujvo is the most direct interface to the lujvo maker. The only argument is
// a list of affix forms for each constituent.
func Lujvo(selci [][][]byte) ([]byte, error) {
	if len(selci) < 2 {
		return []byte{}, fmt.Errorf("need at least two tanru words")
	}
	candidates := []scored{{[]byte{}, 0, false}}
	for selciN, cnino := range selci {
		isLast := selciN == len(selci)-1
		newCand := []scored{}
		for _, rafsi := range cnino {
			var best *scored
			var bestTosmabru *scored
			for _, laldo := range candidates {
				hyphen := []byte{}
				if len(laldo.lujvo) > 0 && IsInvalidCluster(laldo.lujvo[len(laldo.lujvo)-1], rafsi[:2]) {
					hyphen = []byte{'y'}
				} else if !isLast && rafsiTarmi(selci[selciN+1][0]) == fuhivla {
					switch rafsiTarmi(rafsi) {
					case cvhv, ccv, cvv:
						rafsi = append(rafsi, '\'')
					}
				} else if selciN == 1 {
					tai := rafsiTarmi(laldo.lujvo)
					if (tai == cvv || tai == cvhv) && !(isLast && rafsiTarmi(rafsi) == ccv) {
						if rafsi[0] == 'r' {
							hyphen = []byte{'n'}
						} else {
							hyphen = []byte{'r'}
						}
					}
				}
				if !isLast && (rafsiTarmi(rafsi) == cvcc || rafsiTarmi(rafsi) == ccvc) && rafsiTarmi(selci[selciN+1][0]) != fuhivla {
					rafsi = append(rafsi, 'y')
				}
				newPart := append(hyphen, rafsi...)
				newLujvo := append(append([]byte{}, laldo.lujvo...), newPart...)
				newScore := laldo.score + Score(newPart)
				tosmabru := isTosmabruInitial(newLujvo)
				if laldo.tosmabru {
					newScore -= yPenalty
				}
				if tosmabru {
					newScore += yPenalty
				}
				newScored := scored{newLujvo, newScore, tosmabru}
				// DRY
				if tosmabru {
					if bestTosmabru == nil || bestTosmabru.score > newScore {
						bestTosmabru = &newScored
					}
				} else {
					if best == nil || best.score > newScore {
						best = &newScored
					}
				}
			}
			if best != nil {
				newCand = append(newCand, *best)
			}
			if bestTosmabru != nil {
				newCand = append(newCand, *bestTosmabru)
			}
		}
		candidates = newCand
	}
	bestOption := candidates[0]
	for _, o := range candidates {
		if bestOption.score > o.score {
			bestOption = o
		}
	}
	result := bestOption.lujvo
	if bestOption.tosmabru && IsVowel(result[len(result)-1]) {
		result = append(result[:3], append([]byte{'y'}, result[3:]...)...)
	}
	return result, nil
}
