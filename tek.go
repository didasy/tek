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
