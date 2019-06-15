package "github.com/ciuak/jvozba"

import (
	"strings"
)

func Score(lujvo string) int {
	L := len(lujvo)
	// apostrophe count, hyphen count, rafsi score count, vowel count
	A := 0
	H := 0
	R := 0
	V := 0
	var curr string
	for len(lujvo) > 0 {
		curr, lujvo = katna(lujvo)
		tai := rafsiTarmi(curr)
		R += int(tai)
		switch tai {
		case Hyphen:
			H++
		case CVhV:
			V += 2
			A++
		case CVCCV, CCVCV, CVV:
			V += 2
		default:
			V += 1
		}
	}
	return 1000*L - 500*A + 100*H - 10*R - V
}

// the tarmi enum directly corresponds to the scores for R
type tarmi int

const (
	Hyphen tarmi = iota
	CVCCV
	CVCC
	CCVCV
	CCVC
	CVC
	CVhV
	CCV
	CVV
)

func rafsiTarmi(rafsi string) tarmi {
	l := len(rafsi)
	switch l {
	case 1:
		return Hyphen
	case 5:
		if isVowel(rafsi[2]) {
			return CCVCV
		} else {
			return CVCCV
		}
	case 4:
		if rafsi[2] == '\'' {
			return CVhV
		} else if !isVowel(rafsi[3]) {
			if isVowel(rafsi[2]) {
				return CCVC
			} else {
				return CVCC
			}
		} else {
			return Hyphen
		}
	// case 3:
	default:
		if !isVowel(rafsi[2]) {
			return CVC
		} else {
			if isVowel(rafsi[1]) {
				return CVV
			} else {
				return CCV
			}
		}
	}
}

type zunsna int

const (
	Unvoiced zunsna = iota
	Voiced
	Liquid
)

func zunsnaType(one byte) zunsna {
	switch one {
	case 'b', 'd', 'g', 'v', 'j', 'z':
		return Voiced
	case 'p', 't', 'k', 'f', 'c', 's', 'x':
		return Unvoiced
	default:
		return Liquid
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
func needsY(previous byte, current string) bool {
	if isVowel(previous) {
		return false
	}
	head := current[0]
	prevType := zunsnaType(previous)
	headType := zunsnaType(head)
	if (prevType == Voiced && headType == Unvoiced) || (headType == Voiced && prevType == Unvoiced) || previous == head || (isCjsz(previous) && isCjsz(head)) {
		return true
	}
	comb := string([]byte{previous, head})
	switch comb {
	case "cx", "kx", "xc", "xk", "mz":
		return true
	}
	switch string(previous) + current {
	case "ndj", "ndz", "ntc", "nts":
		return true
	}
	return false
}

func isValidInitial(twoBytes ...byte) bool {
	return len(twoBytes) == 2 &&
		strings.Contains("bl~br~cf~ck~cl~cm~cn~cp~cr~ct~dj~dr~dz~fl~fr~gl~gr~jb~jd~jg~jm~jv~kl~kr~ml~mr~pl~pr~sf~sk~sl~sm~sn~sp~sr~st~tc~tr~ts~vl~vr~xl~xr~zb~zd~zg~zm~zv", string(twoBytes))
}

// Whether `lujvo` could lead to tosmabru on appending -y or is one already.
func isTosmabruInitial(lujvo string) bool {
	var r string
	var lastChar byte
	i := 0
	for lujvo != "" {
		r, lujvo = katna(lujvo)
		switch t := rafsiTarmi(r); t {
		case CVC:
			if i > 0 && !isValidInitial(lastChar, r[0]) {
				return false
			}
			lastChar = r[2]
		case CVCCV:
			return i > 0 &&
				isValidInitial(lastChar, r[0]) &&
				isValidInitial(r[2], r[3])
		case Hyphen:
			return i > 1 && r == "y"
		default:
			return false
		}
		i++
	}
	return true
}

// Mind the lack of y.
func isVowel(one byte) bool {
	switch one {
	case 'a', 'e', 'i', 'o', 'u':
		return true
	default:
		return false
	}
}

// Make one cut. Atrocious code.
func katna(lujvo string) (string, string) {
	var point /* of cission */ int
	l := len(lujvo)
	switch {
	case l >= 4 && (lujvo[0] == 'n' || lujvo[0] == 'r' || lujvo[0] == 'y') && !isVowel(lujvo[1]):
		point = 1
	case l >= 8 && lujvo[4] == 'y':
		point = 4
	case l >= 7 && lujvo[3] == 'y':
		point = 3
	case l >= 7 && lujvo[2] == '\'' && isVowel(lujvo[3]):
		point = 4
	case l >= 6:
		point = 3
	default:
		point = l
	}
	return lujvo[:point], lujvo[point:]
}

type scored struct {
	lujvo    string
	score    int
	tosmabru bool
}

// -y-: L += 1 && H += 1 -> score += 1100
const yPenalty = 1100

func Lujvo(selci [][]string) (string, error) {
	candidates := []scored{{"", 0, false}}
	for selciN, cnino := range selci {
		isLast := selciN == len(selci)-1
		newCand := []scored{}
		for _, rafsi := range cnino {
			var best *scored
			var bestTosmabru *scored
			for _, laldo := range candidates {
				l := len(laldo.lujvo)
				hyphen := ""
				if l > 0 && needsY(laldo.lujvo[l-1], rafsi[:2]) {
					hyphen = "y"
				} else if selciN == 1 {
					tai := rafsiTarmi(laldo.lujvo)
					if (tai == CVV || tai == CVhV) && !(isLast && rafsiTarmi(rafsi) == CCV) {
						if rafsi[0] == 'r' {
							hyphen = "n"
						} else {
							hyphen = "r"
						}
					}
				}
				if !isLast && (rafsiTarmi(rafsi) == CVCC || rafsiTarmi(rafsi) == CCVC) {
					rafsi += "y"
				}
				newPart := hyphen + rafsi
				newLujvo := laldo.lujvo + newPart
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
	if bestOption.tosmabru {
		result = result[:3] + "y" + result[3:]
	}
	return result, nil
}
