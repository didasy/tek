/*
tek is an automatic tagging library for Go.
*/
package tek

import (
	"os"
	"math"
	"strings"
	"unicode"
	"io/ioutil"
	"encoding/json"
)

const (
	VERSION = "0.0.2"
)

// need to expand more and rearranged
var indonesianStopWords []string = []string{
	"di",
	"dari",
	"juga",
	"lalu",
	"dengan",
	"ke",
	"ini",
	"itu",
	"dia",
	"dan",
	"aku",
	"saya",
	"kamu",
	"anda",
	"kita",
	"mereka",
	"yang",
	"adalah",
	"walaupun",
	"jika",
	"jadi",
	"akan",
	"tetapi",
	"begitupun",
	"bilamana",
	"bagaimanapun",
	"apa",
	"untuk",
	"kepada",
	"menurut",
	"siapa",
	"dimana",
	"kapan",
	"bagaimana",
	"kenapa",
	"mengapa",
	"pada",
	"dalam",
	"ada",
	"adapun",
	"apapun",
	"ya",
	"tidak",
	"bukan",
}
var englishStopWords []string = []string{
	"a",
	"an",
	"are",
	"arent",
	"about",
	"alone",
	"also",
	"am",
	"and",
	"as",
	"at",
	"after",
	"all",
	"another",
	"any",
	"be",
	"because",
	"before",
	"beside",
	"besides",
	"between",
	"but",
	"by",
	"come",
	"does",
	"doesnt",
	"did",
	"didnt",
	"do",
	"dont",
	"we",
	"for",
	"his",
	"him",
	"himself",
	"himselves",
	"her",
	"herself",
	"herselves",
	"how",
	"our",
	"ours",
	"yours",
	"your",
	"with",
	"my",
	"you",
	"the",
	"in",
	"that",
	"thats",
	"out",
	"on",
	"off",
	"if",
	"will",
	"these",
	"there",
	"theres",
	"those",
	"he",
	"she",
	"it",
	"its",
	"us",
	"is",
	"would",
	"wouldnt",
	"was",
	"wasnt",
	"have",
	"havent",
	"were",
	"werent",
	"has",
	"hasnt",
	"wont",
	"not",
	"had",
	"hadnt",
	"isnt",
	"etc",
	"for",
	"i",
	"or",
	"of",
	"on",
	"other",
	"others",
	"so",
	"than",
	"that",
	"though",
	"to",
	"too",
	"they",
	"through",
	"until",
}

var lang string = "en"

var stopWords []string = englishStopWords

// Define your own stop words by providing a slice of string of stop words
func SetStopWords(s []string) {
	stopWords = s
}

// need to tweak these values later
// var modifier map[string]float64 = map[string]float64{ "nama": 2.5, "nomina" : 1.75, "verba" : 1, "adjektiva" : 0.5, "adverbia" : 0.75, "numeralia" : 0.5 }
var modifier map[string]float64 = map[string]float64{ "nama": 3.5, "nomina" : 3.0, "verba" : 2.0, "adjektiva" : 1.0, "adverbia" : 0.25, "numeralia" : 0.5 }

type Vocab struct {
	Id int `json:id"`
	Word string `json:"word"`
	Type string `json:"type"`
}

var pos []*Vocab

// Set language used, defaulted to english if not called. If argument is not "id" or "en", empty stop words will be used
// For now only support Indonesian and English stop words
func SetLang(l string) error {
	switch l {
	case "id":
		stopWords = indonesianStopWords
		fb, err := ioutil.ReadFile(os.Getenv("GOPATH") + "/src/github.com/JesusIslam/tek/pos_id.json")
		if err != nil {
			return err
		}
		err = json.Unmarshal(fb, &pos)
		if err != nil {
			return err
		}
	break
	case "en":
		stopWords = englishStopWords
	break
	default:
		// if undefined language, use empty stopwords
		stopWords = []string{}
	break
	}
	lang = l
	return nil
}

// The main method of this package, return a slice of *Info struct, sorted by their weight descending.
func GetTags(text string, num int) []*Info {
	// sequential ops, cannot go parallel
	dict := createDictionary(text)
	seq := createSeqDict(dict)
	// we could go concurrent here
	rmStopWordsChan := make(chan []string)
	createSentencesChan := make(chan [][]string)
	defer close(rmStopWordsChan)
	defer close(createSentencesChan)
	go removeStopWords(seq, stopWords, rmStopWordsChan)
	go createSentences(text, createSentencesChan)
	sens := <- createSentencesChan
	seq = <- rmStopWordsChan
	// end
	termsCount := float64(len(flatten(sens)))
	var termsInfo []*Info 
	for _, term := range seq {
		// find idf of each term by counting word occurence first
		count := 0.0
		for _, sen := range sens {
			found := false
			for _, word := range sen {
				if term == word {
					found = true
				}
			}
			if found {
				count++
			}
		}
		idf := math.Log(termsCount / count)
		termsInfo = append(termsInfo, &Info{term, idf, 0.0, 0.0})
	}
	// find each word their tf-idf
	for i, term := range termsInfo {
		var count float64
		for _, sen := range sens {
			for _, word := range sen {
				word = sanitizeWord(word)
				if term.Term == word {
					count++
				}	 
			}
		}
		termsInfo[i].Tf = count / termsCount
		termsInfo[i].Tfidf = termsInfo[i].Tf * term.Idf
	}
	if lang == "id" {
		for i, term := range termsInfo {
			for _, vocab := range pos {
				if term.Term != vocab.Word {
					termsInfo[i].Tfidf += termsInfo[i].Tfidf * modifier["nama"]
					break
				}
				if vocab.Word == term.Term {
					if vocab.Type != "lain-lain" || vocab.Type != "pronomina" || vocab.Type != "interjeksi" || vocab.Type != "preposisi" {
						termsInfo[i].Tfidf += termsInfo[i].Tfidf * modifier[vocab.Type]
						break
					}
					break
				}
			} 
		}
	}
	// sort descending by tfidf
	for i, v := range termsInfo {
		j := i - 1
		for j >= 0 && termsInfo[j].Tfidf < v.Tfidf {
			termsInfo[j+1] = termsInfo[j]
			j -= 1
		}
		termsInfo[j+1] = v
	}
	// out of range error guard
	if num >= len(termsInfo) {
		num = len(termsInfo)
	}
	// return only N number of tags
	result := make([]*Info, num)
	copy(result, termsInfo[:num])
	// empty termsInfo
	termsInfo = []*Info{}
	return result
}

type Info struct {
	Term string
	Idf float64
	Tf float64
	Tfidf float64
}

func flatten(sens [][]string) []string {
	var flat []string
	for _, v := range sens {
		flat = append(flat, v...)
	}
	return flat
}

func createSentences(text string, createSentencesChan chan<- [][]string) {
	text = strings.TrimSpace(text)
	words := strings.Fields(text)
	var sentence []string
	var sentences [][]string
	for _, word := range words {
		// lowercase them FIX 1
		word = strings.ToLower(word)
		// if there isn't . ? or !, append to sentence. If found, also append (and remove the non alphanumerics) but reset the sentence
		if strings.ContainsRune(word, '.') || strings.ContainsRune(word, '!') || strings.ContainsRune(word, '?') {
			word = strings.Map(func (r rune) rune {
				if r == '.' || r == '!' || r == '?' {
					return -1
				}
				return r
			}, word)
			// sanitize them FIX 2
			word = sanitizeWord(word)
			sentence = append(sentence, word)
			sentences = append(sentences, sentence)
			sentence = []string{}
		} else {
			// sanitize them FIX 2
			word = sanitizeWord(word)
			sentence = append(sentence, word)
		}
	}
	if len(sentence) > 0 {
		sentences = append(sentences, sentence)
	}
	sentences = uniqSentences(sentences)
	createSentencesChan <- sentences
}

func uniqSentences(sentences [][]string) [][]string {
	var z []string
	for _, v := range sentences {
		j := strings.Join(v, " ")
		z = append(z, j)
	}
	m := make(map[string]bool)
	var uniq []string
	for _, v := range z {
		if m[v] {
			continue
		}
		uniq = append(uniq, v)
		m[v] = true
	}
	var unique [][]string
	for _, v := range uniq {
		unique = append(unique, strings.Fields(v))
	}
	return unique
}

func removeStopWords(seq []string, StopWords []string, rmStopWordsChan chan<- []string) {
	var res []string
	for _, v := range seq {
		stopWord := false
		for _, x := range StopWords {
			if v == x {
				stopWord = true
				break
			}
		}
		if !stopWord {
			res = append(res, v)
		}
	}
	rmStopWordsChan <- res
}

func sanitizeWord(word string) string {
	word = strings.ToLower(word)
	var prev rune
	word = strings.Map(func (r rune) rune {
		// don't remove '-' if it exists after alphanumerics
		if r == '-' && ((prev >= '0' && prev <= '9') || (prev >= 'a' && prev <= 'z') || prev == 'ä' || prev == 'ö' || prev == 'ü' || prev == 'ß' || prev == 'é') {
			return r
		}
		if !unicode.IsDigit(r) && !unicode.IsLetter(r) && !unicode.IsSpace(r) {
			return -1
		}
		prev = r
		return r
	}, word)
	return word
}

func createSeqDict(dict map[string]int) []string {
	var seq []string
	for term, _ := range dict {
		seq = append(seq, term)
	}
	return seq
}

func createDictionary(text string) map[string]int {
	// trim all spaces
	text = strings.TrimSpace(text)
	// lowercase the text
	text = strings.ToLower(text)
	// remove all non alphanumerics but spaces
	var prev rune
	text = strings.Map(func (r rune) rune {
		// don't remove '-' if it exists after alphanumerics
		if r == '-' && ((prev >= '0' && prev <= '9') || (prev >= 'a' && prev <= 'z') || prev == 'ä' || prev == 'ö' || prev == 'ü' || prev == 'ß' || prev == 'é') {
			return r
		}
		if !unicode.IsDigit(r) && !unicode.IsLetter(r) && !unicode.IsSpace(r) {
			return -1
		}
		prev = r
		return r
	}, text)
	// TRYING TO FIX BUG : remove all double spaces left
	text = strings.Replace(text, "  ", " ", -1)
	// turn it into bag of words
	words := strings.Fields(text)
	// turn it into dictionary
	dict := make(map[string]int)
	i := 1
	for _, word := range words {
		if dict[word] == 0 {
			dict[word] = i
			i++
		}
	}
	return dict
}