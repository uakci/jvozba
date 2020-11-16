// jvozba is an efficient lujvo (Lojban compound word) constructor.  It can
// devise the optimal forms for a given tanru (full-form words which underlie
// the lujvo) in linear time, where ‘lowest score’ is in the sense defined in
// section 11.4 of the Complete Logical Language. The algorithm used resembles
// an indirect A* search.
//
// jvozba also experimentally supports so-called fu'ivla rafsi, or affix forms
// of foreign words, as of v2.
package jvozba

// Config controls the word forms which the algorithm is allowed to output.
// If both Cmevla and Brivla are chosen (bitwise-or'd), only one option (the
// one that has a lower score) will be output anyway. It is illegal to specify
// a null value. If you're unsure what this all means, use the Brivla value.
type Config int

const (
	// Brivla = name word form; word form which ends in a consonant.
	Brivla Config = 1 << iota
	// Cmevla = predicate word form; word form which ends in a vowel.
	Cmevla

	// Always use full affix forms for fu'ivla. For example, with this option
	// turned on, -almy- will become -alma'y-, exposing the final vowel. This
	// flag does not affect fu'ivla which would be ambiguous without the vowel,
	// such as ‘barduku’ or ‘spatula’ (-barduky-, -spatuly- look like two rafsi
	// each).
	LongFuhivla
)

// This is the easiest lujvo-making function you can use; the possible rafsi
// forms will be constructed for you. For more control, use Zbasu or Lujvo.
func Jvozba(tanru string, config Config) (string, error) {
	return Zbasu(tanru, Rafsi, config)
}

// Zbasu is like Jvozba, but it allows you to specify your own list of rafsi.
func Zbasu(tanru string, rafste map[string][]string, config Config) (string, error) {
	byted, err := Normalize(tanru)
	if err != nil {
		return "", err
	}
	slemei, err := Selci(byted, rafste, config)
	if err != nil {
		return "", err
	}
	res, err := Lujvo(slemei)
	return string(res), err
}
