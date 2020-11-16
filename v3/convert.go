package jvozba

import (
	"fmt"
)

// Convert a string to an ASCII []byte slice for faster processing, applying
// replacements along the way. Fails if input is not comprehensible as Lojban.
func Normalize(s string) ([][]byte, error) {
	parts := make([][]byte, 0, len(s)/6)
	base := 0
	for i, r := range s + "\000" {
		switch {
		case (r >= 'a' && r <= 'z' && r != 'q' && r != 'w'):
			continue
		case (r >= 'A' && r <= 'Z' && r != 'Q' && r != 'W'):
			continue
		case r == '\'' || r == ',' || r == '.' || r == '’' || r == '-':
			continue
		case r == ' ' || r == '\t' || r == '\r' || r == '\n' || r == '\000':
			if i > base {
				parts = append(parts, []byte(s[base:i]))
			}
			base = i + 1
		default:
			return [][]byte{}, fmt.Errorf("unexpected character %v", r)
		}
	}
	for i, p := range parts {
		j := 0
		offset := 0
		lenp := len(p)
		for ; j < lenp; j++ {
			switch p[j] {
			// ’ U+2019 = e2 80 99
			case 0xe2:
				p[j] = '\''
				copy(p[j+1:], p[j+3:])
				lenp -= 2
			case '-':
				if j != 0 && j != lenp-1 {
					return [][]byte{}, fmt.Errorf("misplaced hyphen in %s", p[:len(p)+offset])
				}
			case 'h', 'H':
				p[j] = '\''
			}
			if p[j] >= 'A' && p[j] <= 'Z' {
				p[j] += 'a' - 'A'
			}
		}
		parts[i] = p[:lenp]
	}
	return parts, nil
}
