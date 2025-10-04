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

// Comprehensive Indonesian stop words
var indonesianStopWords []string = []string{
	// Common function words
	"di", "dari", "juga", "lalu", "dengan", "ke", "ini", "itu", "dia", "dan", "aku", "saya", "kamu", "anda", "kita", "mereka", "yang", "adalah", "walaupun", "jika", "jadi", "akan", "tetapi", "begitupun", "bilamana", "bagaimanapun", "apa", "untuk", "kepada", "menurut", "siapa", "dimana", "kapan", "bagaimana", "kenapa", "mengapa", "pada", "dalam", "ada", "adapun", "apapun", "ya", "tidak", "bukan",

	// Personal pronouns
	"aku", "saya", "engkau", "kamu", "anda", "beliau", "ia", "dia", "kita", "kami", "mereka",

	// Demonstratives and relatives
	"ini", "itu", "tersebut", "sini", "situ", "sana", "yang", "dimana", "kemana", "darimana",

	// Conjunctions
	"dan", "atau", "tetapi", "melainkan", "namun", "serta", "lalu", "kemudian", "lalu", "lalu", "lantas", "lalu", "lagi", "juga", "juga", "pun", "pula", "serta", "selanjutnya",

	// Prepositions
	"di", "ke", "dari", "pada", "dalam", "atas", "bawah", "antara", "dengan", "tanpa", "untuk", "bagi", "tentang", "hingga", "sampai", "sejak", "selama", "selama", "waktu",

	// Question words
	"apa", "siapa", "mengapa", "kenapa", "bagaimana", "berapa", "kapan", "dimana", "manakah", "bila",

	// Negation and affirmation
	"tidak", "bukan", "jangan", "tak", "jangan", "belum", "sudah", "telah", "pernah", "mungkin", "mestinya", "hendaknya", "semestinya",

	// Auxiliary verbs
	"dapat", "bisa", "akan", "mau", "ingin", "harus", "perlu", "mesti", "wajib", "boleh", "seharusnya", "sepatutnya",

	// Common verbs (generic)
	"adalah", "ialah", "yakni", "yaitu", "memiliki", "mempunyai", "terdapat", "ada", "ada", "terdapat", "wujud", "ada",

	// Time expressions
	"sekarang", "tadi", "tadinya", "nanti", "besok", "kemarin", "lusa", "tadi", "dulu", "dahulu", "sedang", "akan", "telah", "sudah", "masih", "belum", "pernah", "juga", "lagi", "juga", "sementara", "sementara", "sewaktu", "ketika", "tatkala",

	// Place expressions
	"sini", "situ", "sana", "dalam", "luar", "atas", "bawah", "depan", "belakang", "samping", "tengah", "pinggir", "ujung",

	// Quantity and degree
	"semua", "seluruh", "para", "segala", "berbagai", "macam", "jenis", "rupa", "banyak", "sedikit", "beberapa", "hampir", "kurang", "lebih", "paling", "sangat", "amat", "sekali", "cukup", "hanya", "saja", "melainkan", "selain", "kecuali",

	// Modality and certainty
	"mungkin", "barangkali", "tentu", "pasti", "jelas", "nyata", "jelas", "tentu", "pasti", "yakin", "percaya", "rasa", "kira", "agak", "agaknya", "kira-kira", "kayaknya", "sepertinya",

	// Logical connectors
	"karena", "sebab", "akibat", "demi", "guna", "agar", "supaya", "jika", "jikalau", "apabila", "bilamana", "manakala", "walaupun", "meskipun", "kendati", "sekalipun", "biarpun", "walaupun", "meski",

	// Additional common function words
	"bahwa", "bahwa", "bahwa", "hanya", "saja", "melainkan", "selain", "kecuali", "selain", "terkecuali", "justru", "justru", "justru", "malah", "malahan", "bahkan", "lebih-lebih", "terlebih", "apalagi",

	// Numerical and temporal markers
	"pertama", "kedua", "ketiga", "keempat", "kelima", "keenam", "ketujuh", "kedelapan", "kesembilan", "kesepuluh", "terakhir", "awal", "mula-mula", "sejak", "semenjak", "sedari",

	// Discourse markers
	"nah", "lho", "kok", "toh", "dong", "kan", "tuh", "nih", "deh", "yah", "yak", "yoi", "iya", "ya", "betul", "benar", "tentu", "pasti",

	// Common Indonesian particles and interjections
	"dong", "kan", "kok", "lho", "toh", "nih", "tuh", "deh", "yah", "sih", "dong", "yah", "yak", "ya", "oh", "ah", "aduh", "wah", "asyik", "asyik", "mantap", "cihuy",

	// Additional generic action words
	"melakukan", "membuat", "menjadi", "terjadi", "memberikan", "menggunakan", "mengatakan", "mendapatkan", "menjalankan", "melihat", "mendapat", "memberi", "mengambil", "menuju", "terhadap", "mengalami", "pergi",

	// Common filler and transition words
	"umumnya", "biasanya", "pada umumnya", "pada dasarnya", "intinya", "singkatnya", "ringkasnya", "garis besarnya", "sekali lagi", "lagi-lagi", "berulang kali", "berkali-kali", "beberapa kali", "sekian",

	// Indonesian specific connectors
	"sedangkan", "padahal", "biarpun", "sekalipun", "kendatipun", "meskipun", "walaupun", "daripada", "alih-alih", "justru", "malahan", "bahkan", "terlebih", "apalagi",

	// Common discourse adverbs
	"sering", "kerap", "acap", "kerapkali", "seringkali", "kerapkali", "biasanya", "umumnya", "lazimnya", "umumnya", "pada umumnya", "pada dasarnya", "pada hakikatnya",

	// Additional common Indonesian words
	"tersebut", "tersebut", "tersebut", "berikut", "berikut", "berikut", "berikutnya", "selanjutnya", "berikutnya", "seterusnya", "kemudian", "lalu", "lantas", "langsun", "langsung",

	// Common Indonesian expressions
	"hal", "sesuatu", "segala", "macam", "jenis", "bentuk", "rupa", "kondisi", "keadaan", "situasi", "kedudukan", "posisi", "tempat", "lokasi", "area", "wilayah", "daerah",

	// More function words
	"sebab", "karena", "akibat", "dampak", "pengaruh", "efek", "imbas", "akibatnya", "oleh karena itu", "oleh sebab itu", "maka", "maka dari itu", "karenanya", "oleh sebab itu",

	// Comparative and superlative markers
	"lebih", "paling", "ter-", "ber-", "per-", "tersebut", "terkait", "berkaitan", "menyangkut", "mengenai", "tentang", "soal", "urusan", "hal", "masalah",

	// Common Indonesian modal particles
	"lah", "kah", "tah", "dong", "kan", "toh", "kok", "sih", "nih", "tuh", "deh", "yah", "yuk", "ayo", "mari", "silakan", "pergilah", "datanglah",

	// EXPANDED STOP WORDS - Additional comprehensive list

	// Additional common function words
	"berikutnya", "sebelumnya", "selanjutnya", "berikut", "bermacam-macam", "berlainan", "serupa", "sama", "beda",

	// Extended prepositions and locational words
	"sepanjang", "melalui", "bersama", "mengenai", "seputar", "berkaitan", "terkait",

	// Additional temporal expressions
	"kemudian", "selanjutnya", "sewaktu", "sepanjang", "setiap", "masing-masing", "berturut-turut", "berurutan", "sesekali",

	// Indonesian-specific particles and discourse markers
	"toh", "wong", "lha", "gini", "gitu", "gimana", "begini", "begitu", "sekali",

	// Extended question and relative words
	"bagaimanakah", "siapakah", "apakah", "dimanakah", "kapanpun", "dimanapun", "siapapun", "apapun",

	// Additional negation and affirmation words
	"belum", "jangan", "tak usah", "tidak usah", "sudahlah", "tentulah", "pastilah",

	// Extended conjunctions and connectors
	"meskipun", "kendatipun", "biarpun", "sekalipun", "selama", "semenjak", "sejak", "hingga", "sampai",

	// Common colloquial and informal words
	"gitu", "gini", "gimana", "bener", "banget", "dong", "kan", "lah", "deh",

	// Additional demonstrative and reference words
	"yakni", "yaitu", "ialah",

	// Extended numerical and quantitative words
	"beberapa", "segala", "berbagai", "macam-macam", "jenis-jenis", "aneka", "pelbagai",

	// Common academic and formal words
	"terkait", "berkaitan", "mengenai", "menyangkut", "sehubungan", "berkenaan",

	// Additional verbal function words
	"melakukan", "menjalankan", "mengimplementasikan", "mengaplikasikan", "memproses", "mengolah",

	// Indonesian-specific cultural and contextual words
	"tolong", "mohon", "silakan", "mari", "ayo", "yuk",

	// Common Indonesian abbreviations
	"dll", "dsb", "ybs",

	// Technical and academic function words
	"data", "informasi", "sistem", "proses", "metode", "analisis", "hasil", "penelitian",

	// Time-specific additional words
	"pagi", "siang", "sore", "malam", "dini hari", "tengah malam",

	// Extended place-related words
	"sekitar", "sekeliling",

	// Common Indonesian idiomatic expressions
	"pada dasarnya", "pada umumnya", "pada hakikatnya", "intinya", "singkatnya",

	// Additional connector words
	"sebab", "akibat", "dampak", "pengaruh", "efek", "imbas", "oleh karena itu", "oleh sebab itu", "maka", "maka dari itu", "karenanya",

	// Additional discourse adverbs
	"umumnya", "biasanya", "pada umumnya", "lazimnya", "pada dasarnya", "pada hakikatnya", "sering", "kerap", "acap", "kerapkali", "seringkali",

	// Additional prepositions
	"terhadap", "mengenai", "seputar", "berkaitan", "terkait", "menyangkut", "sehubungan", "berkenaan", "mengenai", "tentang", "soal", "urusan", "hal", "masalah",

	// Additional temporal words
	"setiap", "masing-masing", "berturut-turut", "berurutan", "sesekali", "sewaktu", "sepanjang", "selama", "semenjak", "sedari",

	// Additional quantity words
	"aneka", "pelbagai", "berbagai", "macam-macam", "jenis-jenis", "segala", "beberapa", "pelbagai", "aneka ragam", "bagai macam",

	// Additional modal words
	"mestinya", "hendaknya", "seharusnya", "sepatutnya", "sepantasnya", "seyogyanya", "alangkah", "betapa",

	// Additional colloquial words
	"bener", "banget", "kok", "sih", "dong", "kan", "toh", "nih", "tuh", "deh", "yah", "gitu", "gini", "gimana", "begini", "begitu", "gituan",

	// Additional formal words
	"terkait", "berkaitan", "mengenai", "menyangkut", "sehubungan", "berkenaan", "terkait", "berkenaan", "sehubungan", "menyangkut", "mengenai", "mengenai", "seputar", "soal", "urusan", "hal", "masalah", "persoalan",

	// Additional common verbs
	"memberikan", "menggunakan", "mengatakan", "mendapatkan", "menjalankan", "melihat", "mendapat", "memberi", "mengambil", "menuju", "mengalami", "melakukan", "membuat", "menjadi", "terjadi", "mengimplementasikan", "mengaplikasikan", "memproses", "mengolah",

	// Additional particles
	"lah", "kah", "tah", "dong", "kan", "toh", "kok", "sih", "nih", "tuh", "deh", "yah", "yuk", "ayo", "mari", "silakan", "pergilah", "datanglah",

	// Additional time-related words
	"tadi", "tadinya", "nanti", "besok", "kemarin", "lusa", "dulu", "dahulu", "sedang", "akan", "telah", "sudah", "masih", "belum", "pernah", "lagi", "sementara", "sewaktu", "ketika", "tatkala", "sekarang",

	// Additional place-related words
	"dalam", "luar", "atas", "bawah", "depan", "belakang", "samping", "tengah", "pinggir", "ujung", "sekitar", "sekeliling", "sepanjang", "melalui", "menurut", "bersama", "terhadap",

	// Additional quantity words
	"semua", "seluruh", "para", "segala", "berbagai", "macam", "jenis", "rupa", "banyak", "sedikit", "beberapa", "hampir", "kurang", "lebih", "paling", "sangat", "amat", "sekali", "cukup", "hanya", "saja", "melainkan", "selain", "kecuali",

	// Additional certainty words
	"mungkin", "barangkali", "tentu", "pasti", "jelas", "nyata", "yakin", "percaya", "rasa", "kira", "agak", "agaknya", "kira-kira", "kayaknya", "sepertinya", "mestinya", "hendaknya", "semestinya", "seharusnya", "sepatutnya",

	// Additional connectors
	"karena", "sebab", "akibat", "demi", "guna", "agar", "supaya", "jika", "jikalau", "apabila", "bilamana", "manakala", "walaupun", "meskipun", "kendati", "sekalipun", "biarpun", "meski", "sedangkan", "padahal", "daripada", "alih-alih", "justru", "malahan", "bahkan", "terlebih", "apalagi",

	// Additional demonstrative words
	"ini", "itu", "tersebut", "sini", "situ", "sana", "berikut", "berikutnya", "selanjutnya", "seterusnya", "kemudian", "lalu", "lantas", "langsung", "yakni", "yaitu", "ialah",

	// Additional question words
	"apa", "siapa", "mengapa", "kenapa", "bagaimana", "berapa", "kapan", "dimana", "manakah", "bila", "bagaimanakah", "siapakah", "apakah", "dimanakah", "kapanpun", "dimanapun", "siapapun", "apapun",

	// Additional negation words
	"tidak", "bukan", "jangan", "tak", "belum", "sudahlah", "tentulah", "pastilah", "tak usah", "tidak usah", "jangan", "janganlah",

	// Additional auxiliary words
	"dapat", "bisa", "akan", "mau", "ingin", "harus", "perlu", "mesti", "wajib", "boleh", "seharusnya", "sepatutnya", "mestinya", "hendaknya", "semestinya", "sepantasnya", "seyogyanya",

	// Additional common expressions
	"hal", "sesuatu", "segala", "macam", "jenis", "bentuk", "rupa", "kondisi", "keadaan", "situasi", "kedudukan", "posisi", "tempat", "lokasi", "area", "wilayah", "daerah", "pada dasarnya", "pada umumnya", "pada hakikatnya", "intinya", "singkatnya", "ringkasnya", "garis besarnya", "sekali lagi", "lagi-lagi", "berulang kali", "berkali-kali", "beberapa kali", "sekian",
}

// Comprehensive English stop words
var englishStopWords []string = []string{
	// Articles and determiners
	"a", "an", "the", "this", "that", "these", "those", "another", "other", "others", "some", "any", "no", "every", "each", "either", "neither", "both", "all", "several", "various", "certain", "such", "what", "whatever", "which", "whichever", "many", "much", "more", "most", "few", "fewer", "fewest", "little", "less", "least", "some", "somebody", "someone", "something", "somewhere", "anybody", "anyone", "anything", "anywhere", "everybody", "everyone", "everything", "everywhere", "nobody", "nothing", "nowhere",

	// Personal pronouns
	"i", "me", "my", "mine", "myself", "you", "your", "yours", "yourself", "yourselves", "he", "him", "his", "himself", "she", "her", "hers", "herself", "it", "its", "itself", "we", "us", "our", "ours", "ourselves", "they", "them", "their", "theirs", "themselves",

	// Possessive pronouns
	"my", "your", "his", "her", "its", "our", "their", "mine", "yours", "hers", "ours", "theirs",

	// Demonstratives
	"this", "that", "these", "those", "here", "there", "now", "then", "ago", "before", "after", "since", "until", "when", "whenever", "while", "during",

	// Prepositions
	"about", "above", "across", "after", "against", "along", "amid", "among", "amongst", "around", "as", "at", "atop", "before", "behind", "below", "beneath", "beside", "besides", "between", "betwixt", "beyond", "by", "despite", "down", "during", "except", "for", "from", "in", "inside", "into", "like", "near", "next", "of", "off", "on", "onto", "out", "outside", "over", "per", "since", "through", "throughout", "till", "to", "toward", "towards", "under", "underneath", "until", "unto", "up", "upon", "via", "with", "within", "without", "worth",

	// Conjunctions
	"and", "but", "or", "nor", "for", "so", "yet", "after", "although", "as", "because", "before", "even", "if", "lest", "once", "provided", "rather", "since", "so", "than", "that", "though", "till", "unless", "when", "whenever", "where", "whereas", "wherever", "whether", "while", "why",

	// Auxiliary verbs
	"be", "am", "is", "are", "was", "were", "been", "being", "have", "has", "had", "having", "do", "does", "did", "doing", "done", "shall", "should", "will", "would", "can", "could", "may", "might", "must", "ought", "need", "dare", "used", "suppose", "seem", "appear", "happen", "chance",

	// Common verbs (generic)
	"be", "become", "been", "being", "have", "has", "had", "having", "do", "does", "did", "doing", "done", "go", "goes", "went", "going", "gone", "come", "comes", "came", "coming", "get", "gets", "got", "gotten", "getting", "make", "makes", "made", "making", "take", "takes", "took", "taken", "taking", "give", "gives", "gave", "given", "giving", "see", "sees", "saw", "seen", "seeing", "look", "looks", "looked", "looking", "say", "says", "said", "saying", "tell", "tells", "told", "telling", "ask", "asks", "asked", "asking", "work", "works", "worked", "working", "try", "tries", "tried", "trying", "need", "needs", "needed", "needing", "feel", "feels", "felt", "feeling", "seem", "seems", "seemed", "seeming", "leave", "leaves", "left", "leaving", "call", "calls", "called", "calling", "experience", "experiences", "experienced", "experiencing",

	// Common adverbs
	"also", "always", "never", "often", "sometimes", "usually", "rarely", "seldom", "occasionally", "frequently", "constantly", "continuously", "continually", "again", "once", "twice", "seldom", "rarely", "hardly", "scarcely", "barely", "almost", "nearly", "approximately", "roughly", "about", "around", "here", "there", "everywhere", "nowhere", "somewhere", "anywhere", "abroad", "away", "back", "forward", "backward", "up", "down", "upward", "downward", "inward", "outward", "north", "south", "east", "west", "northward", "southward", "eastward", "westward", "abroad", "outside", "inside", "indoors", "outdoors", "upstairs", "downstairs", "across", "along", "aside", "away", "back", "by", "down", "forth", "home", "near", "off", "on", "out", "over", "past", "round", "through", "under", "up",

	// Time adverbs
	"now", "then", "today", "tomorrow", "yesterday", "tonight", "soon", "later", "early", "late", "before", "after", "since", "until", "when", "whenever", "while", "during", "ago", "already", "yet", "still", "even", "just", "now", "presently", "recently", "currently", "lately", "immediately", "instantly", "quickly", "rapidly", "slowly", "gradually", "suddenly", "abruptly", "unexpectedly",

	// Degree adverbs
	"very", "quite", "rather", "fairly", "pretty", "somewhat", "slightly", "a little", "a bit", "extremely", "incredibly", "amazingly", "surprisingly", "remarkably", "notably", "particularly", "especially", "mainly", "mostly", "primarily", "chiefly", "principally", "largely", "generally", "usually", "typically", "normally", "commonly", "frequently", "often", "regularly", "repeatedly", "constantly", "continually", "always", "never", "rarely", "seldom", "hardly", "scarcely", "barely",

	// Question words
	"who", "whom", "whose", "what", "which", "where", "when", "why", "how", "whatever", "whichever", "whoever", "whomever", "whenever", "wherever", "however", "whyever", "howsoever",

	// Negation
	"not", "no", "none", "nothing", "nowhere", "neither", "never", "no one", "nobody", "hardly", "scarcely", "barely", "rarely", "seldom", "without", "lack", "lacking", "absence", "absent", "deny", "denied", "refuse", "refused", "reject", "rejected", "negative", "negatively",

	// Affirmation
	"yes", "yeah", "yep", "yup", "indeed", "certainly", "definitely", "absolutely", "surely", "truly", "really", "actually", "in fact", "of course", "naturally", "obviously", "clearly", "plainly", "evidently", "apparently", "seemingly", "supposedly", "presumably", "allegedly", "reportedly",

	// Logical connectors
	"therefore", "thus", "hence", "consequently", "accordingly", "so", "then", "henceforth", "thereby", "thereupon", "wherefore", "why", "because", "since", "as", "for", "due to", "owing to", "thanks to", "because of", "on account of", "by virtue of", "in view of", "in light of", "considering", "given", "granted", "assuming", "supposing", "provided", "providing", "if", "unless", "whether", "either", "or", "neither", "nor", "both", "and", "but", "yet", "however", "nevertheless", "nonetheless", "notwithstanding", "although", "though", "even though", "despite", "in spite of", "regardless of", "irrespective of", "disregarding",

	// Additional common function words
	"also", "too", "either", "neither", "both", "all", "every", "each", "any", "some", "none", "no", "not", "only", "just", "merely", "simply", "purely", "solely", "exclusively", "particularly", "especially", "specifically", "exactly", "precisely", "approximately", "roughly", "about", "around", "nearly", "almost", "virtually", "practically", "essentially", "basically", "fundamentally", "primarily", "mainly", "chiefly", "principally", "largely", "mostly", "generally", "usually", "typically", "normally", "commonly", "frequently", "often", "regularly", "repeatedly", "constantly", "continually", "always", "never", "rarely", "seldom", "hardly", "scarcely", "barely",

	// Modal and auxiliary expressions
	"can", "could", "may", "might", "must", "shall", "should", "will", "would", "ought", "need", "dare", "used to", "have to", "had to", "must", "should", "ought to", "supposed to", "going to", "about to", "bound to", "likely to", "expected to", "required to", "forced to", "compelled to", "made to",

	// Common discourse markers
	"well", "now", "so", "then", "anyway", "anyhow", "besides", "moreover", "furthermore", "additionally", "also", "too", "as well", "what's more", "in addition", "further", "again", "once more", "anyway", "anyhow", "regardless", "nevertheless", "nonetheless", "however", "but", "yet", "still", "though", "although", "even so", "all the same", "in any case", "in any event", "at any rate", "in either case", "whatever happens", "come what may",

	// Filler and hesitation words
	"uh", "um", "er", "ah", "oh", "well", "like", "you know", "I mean", "kind of", "sort of", "I guess", "I suppose", "I think", "I believe", "I feel", "I reckon", "I suspect", "I imagine", "I expect", "I hope", "I wish", "I want", "I need", "I'd like", "I'd prefer", "I'd rather", "I'd better", "I should", "I must", "I have to", "I've got to",

	// Contractions and informal forms
	"im", "youre", "hes", "shes", "its", "were", "theyre", "ive", "youve", "weve", "theyve", "id", "youd", "hed", "shed", "wed", "theyd", "ill", "youll", "hell", "shell", "well", "theyll", "isnt", "arent", "wasnt", "werent", "havent", "hasnt", "hadnt", "dont", "doesnt", "didnt", "wont", "wouldnt", "cant", "couldnt", "shouldnt", "shouldnt", "mustnt", "mightnt", "neednt", "darent", "didnt", "couldnt", "shouldnt", "wouldnt", "werent", "wasnt", "arent", "isnt", "havent", "hasnt", "hadnt", "doesnt", "dont", "didnt", "wont", "wouldnt", "cant", "couldnt", "shouldnt", "mustnt", "mightnt", "neednt",

	// Additional common expressions
	"etc", "et cetera", "and so on", "and so forth", "etcetera", "vice versa", "per se", "de facto", "de jure", "status quo", "modus operandi", "modus vivendi", "sine qua non", "terra incognita", "tabula rasa", "raison d'etre", "tour de force", "piece de resistance", "cri de coeur", "bete noire", "enfant terrible", "nom de plume", "nom de guerre", "tour de force", "force majeure", "fait accompli", "coup de grace", "coup d'etat", "pièce de résistance", "raison d'être",

	// Numbers and numerical expressions
	"one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "ten", "eleven", "twelve", "thirteen", "fourteen", "fifteen", "sixteen", "seventeen", "eighteen", "nineteen", "twenty", "thirty", "forty", "fifty", "sixty", "seventy", "eighty", "ninety", "hundred", "thousand", "million", "billion", "trillion", "first", "second", "third", "fourth", "fifth", "sixth", "seventh", "eighth", "ninth", "tenth", "once", "twice", "thrice", "again", "single", "double", "triple", "quadruple", "multiple", "several", "various", "numerous", "countless", "innumerable", "many", "few", "several", "some", "any", "no", "none",

	// Days, months, seasons
	"sunday", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "january", "february", "march", "april", "may", "june", "july", "august", "september", "october", "november", "december", "spring", "summer", "autumn", "fall", "winter", "morning", "afternoon", "evening", "night", "midnight", "noon", "dawn", "dusk", "sunrise", "sunset", "daytime", "nighttime",

	// Directions and locations
	"north", "south", "east", "west", "northeast", "northwest", "southeast", "southwest", "up", "down", "left", "right", "forward", "backward", "upward", "downward", "inward", "outward", "across", "along", "around", "through", "between", "among", "amid", "behind", "below", "beneath", "beside", "beyond", "inside", "outside", "upon", "within", "without", "aboard", "about", "above", "across", "after", "against", "along", "amid", "among", "anti", "around", "as", "at", "before", "behind", "below", "beneath", "beside", "besides", "between", "beyond", "but", "by", "concerning", "considering", "despite", "down", "during", "except", "excepting", "excluding", "following", "for", "from", "in", "inside", "into", "like", "minus", "near", "of", "off", "on", "onto", "opposite", "outside", "over", "past", "per", "plus", "regarding", "round", "save", "since", "than", "through", "to", "toward", "towards", "under", "underneath", "unlike", "until", "up", "upon", "versus", "via", "with", "within", "without",

	// Colors and basic adjectives
	"red", "blue", "green", "yellow", "orange", "purple", "pink", "brown", "black", "white", "gray", "grey", "big", "small", "large", "little", "tiny", "huge", "enormous", "giant", "massive", "microscopic", "miniature", "minuscule", "long", "short", "tall", "high", "low", "deep", "shallow", "wide", "narrow", "thick", "thin", "fat", "skinny", "slim", "heavy", "light", "hot", "cold", "warm", "cool", "freezing", "boiling", "scorching", "icy", "new", "old", "young", "fresh", "stale", "modern", "ancient", "recent", "current", "past", "future", "present", "good", "bad", "nice", "nasty", "pleasant", "unpleasant", "beautiful", "ugly", "pretty", "handsome", "attractive", "repulsive", "clean", "dirty", "tidy", "messy", "neat", "disorganized", "organized", "orderly", "chaotic", "easy", "difficult", "hard", "simple", "complex", "complicated", "straightforward", "basic", "advanced", "elementary", "sophisticated",

	// Basic shapes and materials
	"circle", "square", "triangle", "rectangle", "oval", "round", "flat", "curved", "straight", "bent", "wood", "metal", "plastic", "glass", "paper", "stone", "rock", "sand", "water", "oil", "fire", "air", "earth", "soil", "dirt", "clay", "concrete", "steel", "iron", "gold", "silver", "copper", "bronze", "brass", "lead", "tin", "aluminum", "plastic", "rubber", "leather", "cloth", "fabric", "silk", "cotton", "wool", "linen", "nylon", "polyester",

	// Food and drink basics
	"food", "eat", "drink", "water", "milk", "bread", "rice", "meat", "fish", "chicken", "beef", "pork", "vegetable", "fruit", "apple", "banana", "orange", "grape", "strawberry", "potato", "tomato", "onion", "carrot", "lettuce", "cabbage", "corn", "wheat", "flour", "sugar", "salt", "pepper", "oil", "butter", "cheese", "egg", "coffee", "tea", "juice", "soda", "beer", "wine", "alcohol", "cocktail", "dessert", "cake", "cookie", "pie", "ice cream", "chocolate", "candy", "sweet", "sour", "bitter", "spicy", "hot", "mild", "fresh", "cooked", "raw", "baked", "fried", "boiled", "grilled", "roasted", "steamed", "frozen", "canned", "dried", "preserved",

	// Body parts and health
	"head", "face", "eye", "ear", "nose", "mouth", "lip", "tongue", "tooth", "teeth", "hair", "neck", "throat", "shoulder", "arm", "elbow", "wrist", "hand", "finger", "thumb", "chest", "back", "stomach", "waist", "hip", "leg", "knee", "ankle", "foot", "toe", "skin", "bone", "muscle", "blood", "heart", "lung", "brain", "liver", "kidney", "stomach", "health", "healthy", "sick", "ill", "disease", "illness", "pain", "ache", "hurt", "injury", "wound", "cut", "bruise", "medicine", "drug", "pill", "tablet", "capsule", "treatment", "therapy", "cure", "heal", "recover", "rest", "sleep", "wake", "tired", "exhausted", "energetic", "strong", "weak", "fit", "exercise", "workout", "diet", "nutrition",

	// Family and relationships
	"family", "mother", "father", "parent", "brother", "sister", "sibling", "son", "daughter", "child", "kid", "baby", "infant", "toddler", "teenager", "adult", "elderly", "grandmother", "grandfather", "grandparent", "uncle", "aunt", "cousin", "nephew", "niece", "husband", "wife", "spouse", "partner", "boyfriend", "girlfriend", "friend", "enemy", "neighbor", "colleague", "coworker", "boss", "employee", "employer", "manager", "leader", "follower", "teacher", "student", "mentor", "protégé", "role model", "hero", "villain",

	// Education and learning
	"school", "college", "university", "academy", "institute", "class", "course", "lesson", "subject", "topic", "study", "learn", "teach", "educate", "train", "practice", "exercise", "homework", "assignment", "project", "research", "experiment", "test", "exam", "quiz", "grade", "score", "mark", "degree", "diploma", "certificate", "qualification", "skill", "ability", "talent", "gift", "knowledge", "information", "data", "fact", "truth", "wisdom", "intelligence", "smart", "clever", "brilliant", "stupid", "dumb", "ignorant", "uneducated", "illiterate", "literate", "educated", "learned", "scholarly", "academic", "theoretical", "practical", "applied", "basic", "advanced", "elementary", "intermediate", "beginner", "expert", "master", "professional", "amateur", "novice", "specialist", "generalist",

	// Work and career
	"work", "job", "career", "profession", "occupation", "business", "company", "corporation", "organization", "institution", "agency", "department", "office", "factory", "shop", "store", "market", "industry", "trade", "commerce", "economy", "finance", "money", "cash", "currency", "dollar", "pound", "euro", "yen", "price", "cost", "expense", "budget", "income", "salary", "wage", "pay", "earn", "profit", "loss", "gain", "investment", "saving", "spending", "buying", "selling", "shopping", "purchase", "sale", "discount", "bargain", "deal", "offer", "contract", "agreement", "negotiation", "meeting", "conference", "presentation", "report", "document", "file", "record", "database", "computer", "technology", "internet", "website", "email", "phone", "call", "message", "communication", "conversation", "discussion", "argument", "debate", "dispute", "conflict", "resolution", "solution", "answer", "question", "problem", "issue", "challenge", "opportunity", "success", "failure", "achievement", "accomplishment", "goal", "objective", "target", "purpose", "mission", "vision", "strategy", "plan", "project", "task", "duty", "responsibility", "obligation", "commitment", "dedication", "effort", "attempt", "try", "performance", "result", "outcome", "consequence", "effect", "impact", "influence", "change", "development", "progress", "improvement", "growth", "expansion", "increase", "decrease", "reduction", "decline", "rise", "fall", "growth", "shrink", "expand", "contract",

	// EXPANDED STOP WORDS - Additional comprehensive list

	// Technical and scientific terms
	"abstract", "acknowledge", "analysis", "analyze", "approach", "appropriate", "approximately", "aspect", "assess", "assessment", "assume", "assumption", "available", "basis", "brief", "case", "chapter", "chart", "clear", "clearly", "compare", "comparison", "concept", "concern", "concerning", "condition", "conduct", "consequently", "consider", "considerable", "consist", "consistent", "context", "contrast", "correspond", "data", "define", "definition", "definite", "demonstrate", "depart", "derive", "design", "determine", "device", "different", "dimension", "distinct", "distribute", "distribution", "diverse", "document", "domain", "due", "effect", "effective", "efficiency", "efficient", "element", "else", "energy", "enforce", "engage", "engine", "enhance", "enormous", "ensure", "enter", "entire", "environment", "equal", "equally", "error", "especially", "establish", "estimate", "evaluate", "event", "evident", "exactly", "example", "exceed", "except", "exist", "expand", "expect", "expensive", "experience", "experiment", "expert", "explain", "explore", "express", "extend", "extent", "factor", "feature", "field", "figure", "final", "follow", "following", "form", "formal", "former", "formula", "foundation", "function", "functional", "fundamental", "further", "future", "generate", "given", "goal", "grade", "group", "growth", "guideline", "hence", "however", "identify", "illustrate", "impact", "implement", "imply", "impose", "improve", "improvement", "include", "incorporate", "increase", "indicate", "individual", "inference", "influence", "information", "initial", "initiate", "injury", "input", "inquiry", "insert", "inside", "insight", "inspect", "instance", "institute", "institution", "instruction", "instrument", "insurance", "integrate", "intend", "intense", "interact", "interest", "interface", "intermediate", "internal", "interpret", "interval", "intervene", "introduce", "invest", "investigate", "investment", "involve", "issue", "item", "job", "journal", "justification", "lack", "large", "largely", "latter", "lead", "leadership", "learn", "lecture", "legal", "legislate", "less", "level", "library", "license", "life", "lifestyle", "limit", "link", "locate", "location", "logic", "logical", "long", "maintain", "major", "manage", "management", "mandate", "manual", "manufacture", "manufacturing", "margin", "mark", "market", "master", "match", "material", "matrix", "matter", "maximum", "means", "measure", "mechanism", "media", "mediate", "medium", "meet", "mention", "message", "method", "methodology", "metric", "micro", "might", "migration", "military", "million", "mind", "minimum", "minor", "model", "modify", "module", "monitor", "motivation", "multiple", "municipal", "museum", "music", "nation", "national", "natural", "nature", "negative", "network", "neutral", "nevertheless", "news", "nonetheless", "normal", "normally", "north", "notable", "note", "noteworthy", "notion", "novel", "nuclear", "objective", "observation", "observe", "obtain", "obvious", "occasion", "occur", "offer", "office", "official", "okay", "opening", "operate", "operation", "operational", "opinion", "opportunity", "opposite", "optimal", "option", "order", "ordinary", "organization", "organize", "orient", "origin", "original", "outcome", "output", "overall", "overlap", "overseas", "owner", "package", "page", "paragraph", "parameter", "parent", "part", "partial", "participate", "particular", "particularly", "partner", "passage", "passion", "passive", "past", "pattern", "pay", "payment", "pension", "people", "perceive", "percent", "perfect", "perform", "performance", "perhaps", "period", "permit", "persist", "person", "personal", "perspective", "phase", "phenomenon", "philosophy", "phone", "photo", "phrase", "physical", "pick", "picture", "piece", "place", "plan", "plane", "planning", "plastic", "plot", "plus", "point", "police", "policy", "political", "politics", "pool", "popular", "population", "portion", "portrait", "pose", "position", "positive", "possess", "possibility", "possible", "post", "potentially", "practical", "practice", "practitioner", "precede", "precedent", "precise", "precisely", "predict", "prefer", "preference", "preliminary", "premise", "premium", "prepare", "presence", "present", "president", "pressure", "presume", "previous", "price", "primary", "prime", "principal", "principle", "print", "prior", "priority", "private", "probably", "problem", "procedure", "proceed", "process", "produce", "product", "profile", "program", "progress", "project", "promote", "prompt", "proper", "property", "proposal", "prospect", "protect", "protocol", "psychology", "public", "publish", "purchase", "purpose", "pursue", "push", "quality", "query", "question", "quit", "quote", "radical", "range", "rate", "rather", "ratio", "rational", "reach", "read", "ready", "real", "reality", "realize", "reason", "recall", "receive", "recent", "recognition", "recommend", "record", "recover", "recruit", "refer", "reference", "reflect", "regime", "region", "register", "regular", "regulate", "regulation", "reinforce", "relate", "relation", "relationship", "relative", "release", "relevant", "reliability", "reliable", "relief", "religion", "reluctant", "remain", "remember", "remind", "remote", "remove", "repeat", "replace", "reply", "report", "represent", "republican", "request", "require", "research", "resemble", "residence", "resident", "resolve", "resource", "respond", "response", "responsibility", "rest", "restore", "retain", "retire", "retreat", "reveal", "revenue", "reverse", "review", "revise", "revolution", "reward", "rhythm", "rise", "risk", "robust", "role", "roman", "room", "root", "rough", "route", "routine", "row", "royal", "rural", "satisfy", "save", "scale", "scan", "schedule", "scheme", "scope", "score", "screen", "script", "search", "season", "seat", "second", "secret", "section", "sector", "secure", "seek", "seem", "segment", "seize", "select", "selection", "sell", "send", "senior", "sense", "sensitive", "sentence", "separate", "sequence", "series", "serious", "serve", "service", "session", "set", "setting", "settle", "seven", "shadow", "shape", "share", "sharp", "sheet", "shift", "shine", "ship", "shop", "short", "shot", "shoulder", "show", "side", "sight", "sign", "signal", "significance", "significance", "significant", "silence", "similar", "simple", "simulate", "simulation", "since", "single", "site", "situate", "skill", "skin", "slave", "sleep", "slide", "slight", "slightly", "slope", "slow", "small", "smart", "smile", "soil", "solar", "solid", "solution", "solve", "somebody", "somehow", "someone", "something", "sometime", "somewhere", "song", "sophisticated", "sorry", "sort", "sound", "source", "south", "space", "speak", "special", "specialist", "species", "specific", "specify", "spectrum", "spend", "spice", "spin", "split", "spokesperson", "sport", "staff", "stage", "stake", "stand", "standard", "start", "state", "statement", "station", "statistical", "statistics", "status", "steel", "stick", "still", "stock", "stomach", "stone", "stop", "store", "storm", "story", "straight", "strange", "strategic", "strategy", "stream", "street", "strength", "stress", "stretch", "strike", "string", "strip", "stroke", "strong", "structure", "struggle", "student", "study", "stuff", "style", "subject", "submit", "subsequent", "substance", "substantial", "succeed", "success", "successful", "succession", "such", "sudden", "sufficiently", "suggest", "suitable", "summary", "summer", "supply", "support", "suppose", "supreme", "sure", "surface", "surgeon", "surgery", "surprise", "surround", "survey", "survive", "suspect", "sustain", "swing", "switch", "symbol", "sympathetic", "sympathy", "system", "table", "tackle", "tactic", "tail", "take", "talent", "talk", "tank", "tape", "target", "task", "taste", "tax", "technical", "technique", "technology", "teen", "telephone", "telescope", "television", "tell", "temperature", "temporary", "tend", "tension", "tent", "term", "terminate", "terrible", "territory", "test", "testimony", "text", "texture", "theme", "theory", "therapy", "thereby", "therefore", "these", "thesis", "thick", "thing", "think", "third", "thorough", "though", "thought", "thousand", "threat", "three", "threshold", "thrill", "through", "throughout", "throw", "thumb", "thus", "ticket", "tiger", "tight", "time", "tiny", "title", "toast", "today", "together", "tomorrow", "tone", "tongue", "tonight", "tool", "tooth", "topic", "tornado", "total", "totally", "touch", "tough", "tour", "tourist", "toward", "tower", "town", "trace", "track", "trade", "traffic", "tragedy", "tragic", "trail", "train", "transfer", "transform", "transition", "translate", "transport", "travel", "treat", "treatment", "treaty", "tree", "trend", "trial", "triangle", "tribe", "trigger", "trim", "trip", "troop", "tropical", "trouble", "truck", "truly", "trust", "truth", "try", "tube", "tune", "turn", "tutor", "twelve", "twice", "twenty", "twist", "type", "typical", "ultimate", "unable", "uncertain", "uncle", "under", "undergo", "undergraduate", "underground", "underlie", "underlying", "undermine", "underneath", "understand", "undertake", "unemployment", "unexpected", "unfair", "unfold", "unfortunately", "uniform", "unify", "union", "unique", "unit", "unite", "university", "unknown", "unless", "unlike", "unlikely", "unlock", "unnecessary", "unpack", "unpleasant", "unreasonable", "unstable", "until", "unusual", "unveil", "update", "upgrade", "uphold", "upon", "upper", "urban", "urge", "usage", "use", "used", "useful", "user", "usual", "usually", "utility", "utilize", "utmost", "utter", "vacant", "vacation", "vague", "valid", "valley", "valuable", "value", "van", "variable", "variation", "varied", "variety", "various", "vary", "vast", "vehicle", "venture", "venue", "verbal", "verify", "version", "vertical", "very", "vessel", "veteran", "viable", "vibrant", "victim", "victory", "video", "view", "village", "violate", "violence", "violent", "virtually", "virtue", "virus", "visible", "vision", "visit", "visitor", "visual", "vital", "vivid", "vocabulary", "vocal", "voice", "volume", "voluntary", "volunteer", "vote", "voyage", "vulnerable", "wage", "wait", "walk", "wall", "want", "war", "warm", "warn", "wash", "waste", "watch", "water", "wave", "weak", "wealth", "weapon", "wear", "weather", "week", "weekend", "weekly", "weigh", "weight", "weird", "welcome", "welfare", "well", "west", "western", "whale", "whatsoever", "wheel", "whereas", "whereby", "wherein", "whether", "which", "whichever", "whichever", "whip", "whisper", "whistle", "white", "whoever", "whole", "wholesale", "whom", "whomsoever", "whose", "why", "wicked", "wide", "widely", "widespread", "widow", "width", "wife", "wild", "wilderness", "will", "willing", "win", "wind", "window", "wine", "wing", "winner", "winter", "wipe", "wire", "wisdom", "wise", "wish", "with", "withdraw", "within", "without", "witness", "wolf", "woman", "wonder", "wood", "wooden", "wool", "word", "work", "worker", "workforce", "workshop", "world", "worried", "worry", "worth", "would", "wound", "wrap", "wreckage", "wrestle", "wrist", "write", "writer", "writing", "written", "wrong", "yard", "year", "yellow", "yes", "yesterday", "yet", "yield", "young", "younger", "yourself", "youth", "zero", "zone",

	// Modern internet and digital terms
	"click", "scroll", "swipe", "tap", "download", "upload", "refresh", "link", "url", "http", "https", "www", "web", "site", "page", "post", "blog", "vlog", "tweet", "retweet", "like", "share", "comment", "subscribe", "follow", "unfollow", "block", "mute", "notification", "dm", "pm", "tag", "hashtag", "trending", "viral", "meme", "gif", "emoji", "sticker", "filter", "story", "reel", "shorts", "live", "stream", "podcast", "playlist", "algorithm", "feed", "timeline", "profile", "bio", "avatar", "handle", "username", "password", "login", "logout", "signup", "register", "verify", "verification", "auth", "oauth", "captcha", "cookie", "cache", "history", "bookmark", "tab", "window", "browser", "chrome", "safari", "firefox", "edge", "search", "query", "result", "keyword", "seo", "ads", "popup", "banner", "spam", "phishing", "malware", "virus", "hack", "hacker", "cyber", "data", "analytics", "metrics", "stats", "dashboard", "report", "export", "import", "sync", "backup", "cloud", "server", "host", "domain", "ip", "vpn", "proxy", "firewall", "encryption", "security", "privacy", "terms", "conditions", "policy", "guidelines", "rules", "community", "forum", "thread", "reply", "edit", "delete", "archive", "pin", "sticky", "lock", "unlock", "ban", "suspend", "appeal", "mod", "admin", "op", "user", "member", "guest", "anonymous", "anon",

	// Additional contractions and informal forms
	"ain't", "gonna", "gotta", "kinda", "sorta", "lotsa", "gotcha", "yep", "yup", "nah", "uh-huh", "mm-hmm", "uh-uh", "coulda", "shoulda", "woulda", "musta", "mighta", "lovin'", "wan", "finna", "tryna", "kinda", "sorta", "lotsa", "gotta", "gonna", "wanna", "hafta", "sup", "yo", "bruh", "fam", "lit", "savage", "bae", "ghosting", "flex", "stan", "clapback", "shade", "drag", "salty", "triggered", "woke", "basic", "extra", "gucci", "bet", "cap", "no cap", "drip", "rizz", "iykyk", "ngl", "fr", "frfr", "sus", "chad", "simp", "poggers", "dank", "cringe", "based", "redpilled", "wojak", "pepe", "doge", "stonks", "tendies", "fud", "hodl", "gm", "gn", "wagmi", "ngmi", "lol", "omg", "btw", "idk", "tbh", "smh", "rofl", "lmao", "stfu", "gtfo", "fyi", "asap", "aka", "tgif", "hbd", "hmu", "wyd", "wfh", "fomo", "yolo", "swag", "yass", "slay", "periodt", "sis", "bro", "dude", "bruh", "oof", "yeet", "vibe", "vibes", "iykyk", "ngl", "fr", "frfr", "no cap", "bet", "sus", "chad", "simp", "poggers", "dank", "meme", "cringe", "based", "redpilled", "wojak", "pepe", "doge", "stonks", "tendies", "fud", "hodl", "diamond hands", "paper hands", "ape", "gm", "gn", "wagmi", "ngmi", "tldr", "imho", "ama", "nsfw", "tldr", "imho", "ama", "nsfw", "lol", "omg", "btw", "idk", "tbh", "smh", "rofl", "lmao", "stfu", "gtfo", "fyi", "asap", "aka", "tgif", "hbd", "hmu", "wyd", "wfh", "fomo", "yolo", "swag", "yass", "slay", "periodt", "sis", "bro", "dude", "bruh", "oof", "yeet", "vibe", "vibes",
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
			// Skip empty strings
			if word != "" {
				sentence = append(sentence, word)
			}
			if len(sentence) > 0 {
				sentences = append(sentences, sentence)
			}
			sentence = []string{}
		} else {
			// sanitize them FIX 2
			word = sanitizeWord(word)
			// Skip empty strings
			if word != "" {
				sentence = append(sentence, word)
			}
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
		// Skip empty strings and stop words
		if v != "" && !stopWordsMap[v] {
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

	// Exclude pure numbers from tagging
	if word != "" && isNumeric(word) {
		return ""
	}

	return word
}

// isNumeric checks if a string contains only digits (and optionally hyphens for number ranges)
func isNumeric(s string) bool {
	if s == "" {
		return false
	}

	hasDigit := false
	for _, r := range s {
		if unicode.IsDigit(r) {
			hasDigit = true
		} else if r != '-' {
			// If any character other than digit or hyphen, it's not a pure number
			return false
		}
	}

	return hasDigit
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
		// Sanitize the word to exclude numbers
		word = sanitizeWord(word)
		// Skip empty strings
		if word != "" && dict[word] == 0 {
			dict[word] = i
			i++
		}
	}
	return dict
}
