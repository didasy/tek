/*
tek is an automatic tagging library for Go.
*/
package tek

import (
	"math"
	"runtime"
	"sort"
	"strings"
	"unicode"
)

const (
	VERSION = "0.1.1"
)

func init() {
	// Initialize default stop words map for English
	stopWordsMap = make(map[string]bool, len(englishStopWords))
	for _, word := range englishStopWords {
		stopWordsMap[word] = true
	}
}

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
	"mengalami",
	"pergi",
	"dapat",
	"bisa",
	"melakukan",
	"membuat",
	"menjadi",
	"terjadi",
	"memberikan",
	"memiliki",
	"menggunakan",
	"mengatakan",
	"mendapatkan",
	"menjalankan",
	"melihat",
	"mendapat",
	"memberi",
	"menerapkan",
	"melakukan",
	"mengambil",
	"menuju",
	"terhadap",
	"sebagai",
	"oleh",
	"karena",
	"saat",
	"selain",
	"tersebut",
	"yaitu",
	"yakni",
	"ialah",
	"bahwa",
	"dimana",
	"kemudian",
	"setelah",
	"sebelum",
	"sejak",
	"hingga",
	"sampai",
	"daripada",
	"sering",
	"kerap",
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
	"i",
	"or",
	"of",
	"other",
	"others",
	"so",
	"than",
	"though",
	"to",
	"too",
	"they",
	"through",
	"until",
	"go",
	"went",
	"going",
	"gone",
	"came",
	"coming",
	"get",
	"getting",
	"see",
	"saw",
	"seeing",
	"seen",
	"know",
	"knew",
	"knowing",
	"known",
	"take",
	"took",
	"taking",
	"taken",
	"give",
	"gave",
	"giving",
	"given",
	"make",
	"made",
	"making",
	"can",
	"could",
	"may",
	"might",
	"must",
	"shall",
	"should",
	"being",
	"been",
	"do",
	"doing",
	"done",
	"say",
	"said",
	"saying",
	"tell",
	"told",
	"telling",
	"ask",
	"asked",
	"asking",
	"work",
	"worked",
	"working",
	"seem",
	"seemed",
	"seeming",
	"leave",
	"left",
	"leaving",
	"call",
	"called",
	"calling",
	"try",
	"tried",
	"trying",
	"need",
	"needed",
	"needing",
	"feel",
	"felt",
	"feeling",
	"become",
	"became",
	"becoming",
	"experience",
	"experienced",
	"experiencing",
	"experiences",
}

var lang string = "en"

var stopWords []string = englishStopWords
var stopWordsMap map[string]bool

func SetStopWords(s []string) {
	stopWords = s
	stopWordsMap = make(map[string]bool, len(s))
	for _, word := range s {
		stopWordsMap[word] = true
	}
}

// need to tweak these values later
// var modifier map[string]float64 = map[string]float64{ "nama": 2.5, "nomina" : 1.75, "verba" : 1, "adjektiva" : 0.5, "adverbia" : 0.75, "numeralia" : 0.5 }
var modifier map[string]float64 = map[string]float64{"nama": 3.5, "nomina": 3.0, "verba": 2.0, "adjektiva": 1.0, "adverbia": 0.25, "numeralia": 0.5}

type Vocab struct {
	Id   int    `json:"id"`
	Word string `json:"word"`
	Type string `json:"type"`
}

var pos []*Vocab
var posMap map[string]*Vocab

// Set language used, defaulted to english if not called. If argument is not "id" or "en", empty stop words will be used
// For now only support Indonesian and English stop words
func SetLang(l string) error {
	switch l {
	case "id":
		stopWords = indonesianStopWords
		pos = indonesianPos
		// Build POS map for O(1) lookup
		posMap = make(map[string]*Vocab, len(indonesianPos))
		for _, vocab := range indonesianPos {
			posMap[vocab.Word] = vocab
		}
		// Build stop words map
		stopWordsMap = make(map[string]bool, len(indonesianStopWords))
		for _, word := range indonesianStopWords {
			stopWordsMap[word] = true
		}
		break
	case "en":
		stopWords = englishStopWords
		// Build stop words map for English
		stopWordsMap = make(map[string]bool, len(englishStopWords))
		for _, word := range englishStopWords {
			stopWordsMap[word] = true
		}
		posMap = nil
		break
	default:
		// if undefined language, use empty stopwords
		stopWords = []string{}
		stopWordsMap = make(map[string]bool)
		posMap = nil
		break
	}
	lang = l
	return nil
}

func findIdf(idx int, termsInfo []*Info, sentences [][]string, termsCount float64, term string) {
	count := 0.0
	for _, sen := range sentences {
		found := false
		for _, word := range sen {
			if term == word {
				found = true
				// Don't break here - the original logic didn't break
			}
		}
		if found {
			count++
		}
	}
	idf := math.Log(termsCount / count)
	termsInfo[idx] = &Info{term, idf, 0.0, 0.0}
}

func findTfidf(idx int, termsInfo []*Info, termsCount float64, sentences [][]string) {
	count := 0.0
	term := termsInfo[idx].Term
	for _, sen := range sentences {
		for _, word := range sen {
			word = sanitizeWord(word) // Restore sanitizeWord call for correctness
			if term == word {
				count++
			}
		}
	}
	termsInfo[idx].Tf = count / termsCount
	termsInfo[idx].Tfidf = termsInfo[idx].Tf * termsInfo[idx].Idf
}

func modifyTfidfId(idx int, termsInfo []*Info, posMap map[string]*Vocab) {
	term := termsInfo[idx].Term
	found := false
	for _, vocab := range pos { // Use original pos array for exact same behavior
		if term != vocab.Word {
			termsInfo[idx].Tfidf += termsInfo[idx].Tfidf * modifier["nama"]
			found = true
			break
		}
		if term == vocab.Word {
			// Apply the exact same logic as original
			if vocab.Type != "lain-lain" {
				termsInfo[idx].Tfidf += termsInfo[idx].Tfidf * modifier[vocab.Type]
			}
			if vocab.Type != "pronomina" {
				termsInfo[idx].Tfidf += termsInfo[idx].Tfidf * modifier[vocab.Type]
			}
			if vocab.Type != "interjeksi" {
				termsInfo[idx].Tfidf += termsInfo[idx].Tfidf * modifier[vocab.Type]
			}
			if vocab.Type != "preposisi" {
				termsInfo[idx].Tfidf += termsInfo[idx].Tfidf * modifier[vocab.Type]
			}
			found = true
			break
		}
	}
	// If word not found in POS dictionary at all
	if !found {
		termsInfo[idx].Tfidf += termsInfo[idx].Tfidf * modifier["nama"]
	}
}

type Info struct {
	Term  string
	Idf   float64
	Tf    float64
	Tfidf float64
}

// The main method of this package, return a slice of *Info struct, sorted by their weight descending.
func GetTags(text string, num int) []*Info {
	return GetTagsWithWorkers(text, num, runtime.NumCPU())
}

// GetTagsWithWorkers allows specifying the number of workers for concurrent processing.
// If numWorkers is 0 or negative, it defaults to the number of available CPU cores.
func GetTagsWithWorkers(text string, num int, numWorkers int) []*Info {
	// sequential ops, cannot go parallel
	dict := createDictionary(text)
	seq := createSeqDict(dict)
	// we could go concurrent here
	rmStopWordsChan := make(chan []string)
	createSentencesChan := make(chan [][]string)
	defer close(rmStopWordsChan)
	defer close(createSentencesChan)
	go removeStopWords(seq, stopWordsMap, rmStopWordsChan)
	go createSentences(text, createSentencesChan)
	sens := <-createSentencesChan
	seq = <-rmStopWordsChan
	// end
	termsCount := float64(len(flatten(sens)))

	// Use worker pools for better concurrency
	if numWorkers <= 0 {
		numWorkers = runtime.NumCPU()
	}
	if len(seq) < numWorkers {
		numWorkers = len(seq)
	}

	// Parallel IDF calculation with worker pool
	termsInfo := make([]*Info, len(seq))
	idfJobs := make(chan int, len(seq))
	idfDone := make(chan bool, numWorkers)

	// Start IDF workers
	for w := 0; w < numWorkers; w++ {
		go func() {
			for idx := range idfJobs {
				findIdf(idx, termsInfo, sens, termsCount, seq[idx])
			}
			idfDone <- true
		}()
	}

	// Send jobs
	for i := range seq {
		idfJobs <- i
	}
	close(idfJobs)

	// Wait for workers to complete
	for i := 0; i < numWorkers; i++ {
		<-idfDone
	}

	// Parallel TF-IDF calculation with worker pool
	tfidfJobs := make(chan int, len(termsInfo))
	tfidfDone := make(chan bool, numWorkers)

	// Start TF-IDF workers
	for w := 0; w < numWorkers; w++ {
		go func() {
			for idx := range tfidfJobs {
				findTfidf(idx, termsInfo, termsCount, sens)
			}
			tfidfDone <- true
		}()
	}

	// Send jobs
	for i := range termsInfo {
		tfidfJobs <- i
	}
	close(tfidfJobs)

	// Wait for workers to complete
	for i := 0; i < numWorkers; i++ {
		<-tfidfDone
	}

	if lang == "id" {
		// Parallel Indonesian POS modification with worker pool
		posJobs := make(chan int, len(termsInfo))
		posDone := make(chan bool, numWorkers)

		// Start POS workers
		for w := 0; w < numWorkers; w++ {
			go func() {
				for idx := range posJobs {
					modifyTfidfId(idx, termsInfo, posMap)
				}
				posDone <- true
			}()
		}

		// Send jobs
		for i := range termsInfo {
			posJobs <- i
		}
		close(posJobs)

		// Wait for workers to complete
		for i := 0; i < numWorkers; i++ {
			<-posDone
		}
	}

	// Sort only once using sort.SliceStable (remove the insertion sort)
	sort.SliceStable(termsInfo, func(i, j int) bool {
		return termsInfo[i].Tfidf > termsInfo[j].Tfidf
	})

	// out of range error guard
	if num >= len(termsInfo) {
		num = len(termsInfo)
	}

	// return only N number of tags
	result := make([]*Info, num)
	copy(result, termsInfo[:num])
	return result
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
			word = strings.Map(func(r rune) rune {
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
	z := make([]string, len(sentences))
	for i, v := range sentences {
		j := strings.Join(v, " ")
		z[i] = j
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
	unique := make([][]string, len(uniq))
	for i, v := range uniq {
		unique[i] = strings.Fields(v)
	}
	return unique
}

func removeStopWords(seq []string, stopWordsMap map[string]bool, rmStopWordsChan chan<- []string) {
	// Pre-allocate result slice with estimated capacity
	res := make([]string, 0, len(seq))
	for _, v := range seq {
		if !stopWordsMap[v] {
			res = append(res, v)
		}
	}
	rmStopWordsChan <- res
}

func sanitizeWord(word string) string {
	word = strings.ToLower(word)
	var prev rune
	word = strings.Map(func(r rune) rune {
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
	text = strings.Map(func(r rune) rune {
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
