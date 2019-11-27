package translate

import (
	"fmt"
	"net/http"
	"net/url"
)

const urlFormat = "https://translate.googleapis.com/translate_a/single?client=gtx&sl=%s&tl=%s&dt=bd&q=%s"
const english = "en"
const spanish = "es"

// TranslatedWord is a translation and the class of word
type TranslatedWord struct {
	class       string
	translation string
}

// Lookup fetches translation results from Google Translate
func Lookup(input string) ([]TranslatedWord, error) {
	url := fmt.Sprintf(urlFormat, spanish, english, url.PathEscape(input))
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	return parseToWords(resp.Body)
}
