package crawler

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gabrielmoura/WebCrawler/infra/data"
	"github.com/gabrielmoura/WebCrawler/infra/log"
	"golang.org/x/net/html"
	"regexp"
)

// commonStopWords Palavras de parada comuns (personalize conforme necessário, geradas por GPT)
var commonStopWords = map[string][]string{
	"en": {"is", "or", "a", "and", "the", "are", "of", "to"},
	"pt": {
		"a", "à", "ao", "aos", "aquela", "aquelas", "aquele", "aqueles", "aquilo", "as", "até",
		"com", "como", "da", "das", "de", "dela", "delas", "dele", "deles", "desta", "destas",
		"deste", "destes", "do", "dos", "e", "é", "ela", "elas", "ele", "eles", "em", "entre",
		"era", "essa", "essas", "esse", "esses", "esta", "estamos", "estas", "estava", "estavam",
		"estávamos", "este", "estes", "estou", "eu", "foi", "fomos", "for", "foram", "havia",
		"isso", "isto", "já", "lhe", "lhes", "mais", "mas", "me", "mesmo", "meu", "meus", "minha",
		"minhas", "muito", "na", "nas", "nem", "no", "nos", "nossa", "nossas", "nosso", "nossos",
		"num", "numa", "o", "os", "ou", "para", "pela", "pelas", "pelo", "pelos", "por", "qual",
		"quando", "que", "quem", "se", "sem", "seu", "seus", "sua", "suas", "só", "também", "te",
		"tem", "temos", "tenha", "tenham", "teu", "teus", "teve", "tinha", "tinham", "tua", "tuas",
		"um", "uma", "você", "vocês", "vos",
	},
	"ru": {
		"а", "без", "более", "бы", "был", "была", "были", "было", "быть", "в", "вам", "вас", "всё",
		"все", "всего", "всех", "вы", "где", "да", "даже", "для", "до", "его", "ее", "ей", "ему",
		"если", "ест", "есть", "ещё", "ж", "же", "за", "здесь", "и", "из", "или", "им", "их", "к",
		"как", "какая", "какой", "когда", "кое", "кто", "куда", "ли", "либо", "мне", "много", "может",
		"можно", "мой", "моя", "мы", "на", "над", "надо", "наш", "не", "него", "нее", "нет", "ни",
		"них", "но", "ну", "о", "об", "однако", "он", "она", "они", "оно", "от", "очень", "по",
		"под", "после", "при", "про", "с", "сам", "сама", "сами", "само", "свое", "своего", "своей",
		"свои", "себе", "себя", "сейчас", "со", "совсем", "так", "такой", "там", "тебя", "тем",
		"теперь", "то", "тогда", "того", "тоже", "только", "том", "ты", "у", "уже", "хотя", "чего",
		"чей", "чем", "что", "чтобы", "чуть", "эта", "эти", "это", "этого", "этой", "этом", "эту",
		"я",
	},
	"es": {
		"el", "la", "los", "las", "de", "del", "y", "a", "en", "un", "una", "unos", "unas", "con",
		"para", "por", "su", "se", "que", "es", "soy", "eres", "somos", "son", "me", "te", "nos", "le",
		"les", "lo", "mi", "tu", "si", "no", "pero", "porque", "como", "esta", "estoy", "estas", "estamos",
		"estais", "estan", "muy", "poco", "mucho", "todo", "todos", "al", "algo", "alguien", "donde", "cuando",
		"como", "aqui", "ahi", "alli", "ahora", "antes", "despues", "hoy", "ayer", "mañana", "siempre", "nunca",
	},
	"hindi": {
		"का", "के", "की", "में", "है", "और", "यह", "वह", "से", "को", "पर", "इस", "होता", "ही", "हैं", "ये", "वो", "कर", "गया", "लिए",
		"अपना", "अपनी", "अपने", "कुछ", "थी", "थे", "थीं", "हुआ", "जा", "रहा", "रहे", "जाता", "जाती", "जाते", "एक", "दो", "तीन", "चार",
		"पांच", "छह", "सात", "आठ", "नौ", "दस",
	},
	"ch": {
		"的", "了", "在", "是", "我", "有", "和", "就", "不", "人", "这", "那", "中", "来", "上", "大", "为", "个", "国",
		"以", "说", "到", "要", "子", "你", "会", "着", "能", "里", "去", "年", "得", "他", "她", "它", "们", "地", "也",
		"自", "这", "时", "那", "儿", "可", "就", "给", "下", "都", "向", "看", "起", "还", "过", "只", "把", "对", "做",
		"当", "想", "成", "事", "被", "用", "多", "从", "面", "等", "前", "些", "于", "后", "所", "又", "经", "方", "现",
		"没", "吧", "定", "得", "该", "好好", "家", "种", "那", "里", "然", "其", "间", "什", "么", "很", "得", "哪", "些",
		"向", "生", "里", "果", "再", "两", "并", "而", "些", "定",
	},
}

// CountWordsInText Extrai e conta a frequência de palavras do conteúdo HTML, ignorando palavras irrelevantes comuns.
func CountWordsInText(node *html.Node) (*data.Words, error) {

	// Etapa 1: remover tags HTML
	htmlRegex := regexp.MustCompile("<([^>]*)>")
	plainText := htmlRegex.ReplaceAll([]byte(node.Data), []byte("")) // Convert node.Data to []byte

	// Etapa 2: Extraia a primeira frase (opcional)
	// Isso pode ser útil se você quiser analisar apenas uma parte do texto.
	sentenceRegex := regexp.MustCompilePOSIX(".+[\\.\\!\\?] ")
	firstSentence := sentenceRegex.Find(plainText) // Opcionalmente use firstSentence em vez de plainText

	// Etapa 3: Normalizar texto
	normalizedText := bytes.ToLower(firstSentence) // or use plainText if not extracting the first sentence

	// Etapa 4: Remova caracteres especiais e divida em palavras
	wordRegex := regexp.MustCompile("\n|\t|&[a-z]+|[.,]+ |;|\u0009")
	words := bytes.Split(wordRegex.ReplaceAll(normalizedText, []byte(" ")), []byte(" "))

	// Etapa 5: Conte a frequência das palavras (ignorando palavras comuns)
	var wordCounts data.Words
	for _, wordBytes := range words {
		word := string(wordBytes)

		// Skip short words and common stop words
		if len(word) < 2 || containsMap(commonStopWords, word) {
			continue
		}

		wordCounts[word]++
	}

	return &wordCounts, nil
}

// ContainsMap Verifica se uma palavra está em uma lista de stop words comuns,
func containsMap(wordMap map[string][]string, item string) bool {
	for key, slice := range wordMap {
		// Ignora a primeira string do mapa (chave vazia ou primeira chave lexicograficamente)
		if key == "" {
			continue
		}

		for _, a := range slice {
			if a == item {
				return true
			}
		}
	}
	return false
}

func extractData(n *html.Node) (*data.Page, error) {
	var dataPage data.Page

	var extract func(*html.Node)
	extract = func(n *html.Node) {
		if n.Type == html.ElementNode {
			if n.Data == "title" && n.FirstChild != nil {
				dataPage.Title = n.FirstChild.Data
			} else if n.Data == "meta" {
				var isDescription bool
				for _, a := range n.Attr {
					if a.Key == "name" && a.Val == "description" {
						isDescription = true
					}
					if a.Key == "content" {
						if isDescription {
							dataPage.Description = a.Val
						} else {
							dataPage.Meta = append(dataPage.Meta, a.Val)
						}
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extract(c)
		}
	}

	extract(n)
	ok, _ := CountWordsInText(n)
	dataPage.Words = ok
	return &dataPage, nil
}

func extractLinks(parentLink string, n *html.Node) ([]string, error) {
	var links []string

	var extract func(*html.Node)
	extract = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					urlE, err := prepareLink(a.Val)
					if err != nil {
						if errors.Is(invalidSchemaErr, err) {
							preparedLink, err := prepareParentLink(parentLink, a.Val)
							if err != nil {
								continue
							}
							urlE = preparedLink
						}
						log.Logger.Debug(fmt.Sprintf("Error preparing link: %s", err))
						continue
					}
					links = append(links, urlE.String())
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extract(c)
		}
	}

	extract(n)
	return links, nil
}
