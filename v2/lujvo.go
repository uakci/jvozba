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
	if l == 0 || (!IsConsonant(rafsi[0]) && l != 1) {
		return fuhivla
	}
	switch l {
	case 1:
		if len(rafsi) == 1 && (rafsi[0] == 'y' || rafsi[0] == 'r' || rafsi[0] == 'n') {
			return hyphen
		}
	case 2:
		if bytes.Equal(rafsi, []byte("'y")) {
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
		if isGismu(rafsi) {
			if IsVowel(rafsi[2]) {
				return ccvcv
			} else {
				return cvccv
			}
		}
	}
	return fuhivla
}

type zunsna int

const (
	unvoiced zunsna = iota
	voiced
	liquid
)

func zunsnaType(one byte) zunsna {
	switch one {
	case 'b', 'd', 'g', 'v', 'j', 'z':
		return voiced
	case 'p', 't', 'k', 'f', 'c', 's', 'x':
		return unvoiced
	default:
		return liquid
	}
}

func isCjsz(one byte) bool {
	switch one {
	case 'c', 'j', 's', 'z':
		return true
	default:
		return false
	}
}

// `needsY` checks if an y-hyphen needs to be inserted
func needsY(previous byte, current []byte) bool {
	if IsVowel(previous) {
		return false
	}
	head := current[0]
	prevType := zunsnaType(previous)
	headType := zunsnaType(head)
	if (prevType == voiced && headType == unvoiced) || (headType == voiced && prevType == unvoiced) || previous == head || (isCjsz(previous) && isCjsz(head)) {
		return true
	}
	switch previous {
	case 'c', 'k':
		if head == 'x' {
			return true
		}
	case 'x':
		if head == 'c' || head == 'k' {
			return true
		}
	case 'm':
		if head == 'z' {
			return true
		}
	case 'n':
		if len(current) < 2 {
			break
		}
		switch current[0] {
		case 't':
			if current[1] == 's' || current[1] == 'c' {
				return true
			}
		case 'd':
			if current[1] == 'z' || current[1] == 'j' {
				return true
			}
		}
	}
	return false
}

var invalidClusters = [][2]byte{
	'c': {'x'},
	'k': {'x'},
	'x': {'c', 'k'},
	'm': {'z'},
}

var validInitials = map[byte][]byte{
	'b': {'l', 'r'},
	'c': {'f', 'k', 'l', 'm', 'n', 'p', 'r', 't'},
	'd': {'j', 'r', 'z'},
	'f': {'l', 'r'},
	'g': {'l', 'r'},
	'j': {'b', 'd', 'g', 'm', 'v'},
	'k': {'l', 'r'},
	'm': {'l', 'r'},
	'p': {'l', 'r'},
	's': {'f', 'k', 'l', 'm', 'n', 'p', 'r', 't'},
	't': {'c', 'r', 's'},
	'v': {'l', 'r'},
	'x': {'l', 'r'},
	'z': {'b', 'd', 'g', 'm', 'v'},
}

func isValidInitial(twoBytes ...byte) bool {
	if len(twoBytes) != 2 {
		return false
	}
	for _, validInitial := range validInitials[twoBytes[0]] {
		if validInitial == twoBytes[1] {
			return true
		}
	}
	return false
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
			if i > 0 && !isValidInitial(lastChar, r[0]) {
				return false
			}
			lastChar = r[2]
		case cvccv:
			return i > 0 &&
				isValidInitial(lastChar, r[0]) &&
				isValidInitial(r[2], r[3])
		case hyphen:
			return i > 1 && bytes.Equal(r, []byte{'y'})
		default:
			return false
		}
		i++
	}
	return true
}

// Mind the lack of y.
func IsVowel(one byte) bool {
	switch one {
	case 'a', 'e', 'i', 'o', 'u':
		return true
	default:
		return false
	}
}

func IsConsonant(one byte) bool {
	return !IsVowel(one) && one != 'y' && one != '\''
}

// Make one cut. Atrocious code.
func katna(lujvo []byte) ([]byte, []byte) {
	var point /* of cission */ int
	l := len(lujvo)
	switch {
	case l >= 1 && lujvo[0] == 'y':
		point = 1
	case l >= 2 && bytes.Equal(lujvo[:2], []byte("'y")):
		point = 2
	case l >= 4 && (lujvo[0] == 'n' || lujvo[0] == 'r' || lujvo[0] == 'y') && IsConsonant(lujvo[1]):
		point = 1
	case l >= 8 && lujvo[4] == 'y':
		point = 4
	case l >= 7 && lujvo[3] == 'y':
		point = 3
	case l >= 7 && lujvo[2] == '\'' && IsVowel(lujvo[3]):
		point = 4
	case l >= 6 && bytes.Index(lujvo[:6], []byte{'y'}) == -1:
		point = 3
	default:
		point = bytes.Index(lujvo, []byte{'y'})
		if point == -1 {
			point = l
		}
	}
	return lujvo[:point], lujvo[point:]
}

func Katna(lujvo []byte) (result [][]byte) {
	chunk := make([][]byte, 0, len(lujvo)/3)
	fuhivlaTainted, i := false, 0
	var rafsi []byte
	for len(lujvo) > 0 {
		rafsi, lujvo = katna(lujvo)
		tai := rafsiTarmi(rafsi)
		switch tai {
		case hyphen:
			if rafsi[0] != 'y' {
				if i != 1 {
					chunk = append(chunk, rafsi)
					fuhivlaTainted = true
				}
			} else {
				if fuhivlaTainted {
					res := bytes.Join(chunk, []byte{})
					if IsConsonant(res[len(res)-1]) {
						res = append(res, 'a')
					} else if res[len(res)-1] == '\'' {
						res = res[:len(res)-1]
					} else {
						res = append(res, 'a')
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
			chunk = append(chunk, rafsi)
		}
		i++
	}
	if fuhivlaTainted {
		result = append(result, bytes.Join(chunk, []byte{}))
	} else {
    result = append(result, chunk...)
	}
	return result
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
				if len(laldo.lujvo) > 0 && needsY(laldo.lujvo[len(laldo.lujvo)-1], rafsi[:2]) {
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
