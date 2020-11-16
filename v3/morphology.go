package jvozba

import (
	"bytes"
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

// Matches aeiou (not y!).
func IsVowel(o byte) bool {
	return o == 'a' || o == 'e' || o == 'i' || o == 'o' || o == 'u'
}

// Matches bcdfgjklmnprstvxz.
func IsConsonant(o byte) bool {
	return !(o < 'a' || o > 'z' || o == 'a' || o == 'e' || o == 'i' || o == 'o' || o == 'u' || o == 'y')
}

// Returns true if twoBytes is a valid initial consonant cluster.
func IsValidInitial(twoBytes ...byte) bool {
	if len(twoBytes) != 2 {
		return false
	}
	for _, validInitial := range ValidInitials[twoBytes[0]] {
		if validInitial == twoBytes[1] {
			return true
		}
	}
	return false
}

var ValidInitials = map[byte][]byte{
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

func isCjsz(one byte) bool {
	switch one {
	case 'c', 'j', 's', 'z':
		return true
	default:
		return false
	}
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

// `IsInvalidCluster` checks if an y-hyphen needs to be inserted
func IsInvalidCluster(previous byte, current []byte) bool {
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
