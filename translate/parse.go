package translate

import (
	"fmt"
	"io"
	"strings"
)

var allowedClasses = []string{
	"noun", "verb", "adjective", "adverb", "pronoun", "preposition", "conjunction", "particle", "phrase", "auxiliary verb"}

type token struct {
	value    string
	children []*token
}

func parseToWords(in io.Reader) ([]*TranslatedWord, error) {
	root := &token{}
	if err := parse(in, root); err != nil {
		return nil, err
	}
	var result []*TranslatedWord
	if len(root.children) != 1 {
		return nil, fmt.Errorf("fetched empty result from Google Translate")
	}
	if len(root.children[0].children) < 3 {
		return nil, fmt.Errorf("expect the root to have 3 or more children, found %d", len(root.children[0].children))
	}
	translationBlocks := root.children[0].children[1]
	for _, t := range translationBlocks.children {
		if len(t.children) != 5 {
			return nil, fmt.Errorf("expect the translation block to have 5 children, found %d", len(t.children))
		}
		class := t.children[0].value
		if !allowableClass(class) {
			return nil, fmt.Errorf("expected one of %v, found %s", allowedClasses, class)
		}
		for _, w := range t.children[1].children {
			result = append(result, &TranslatedWord{class, w.value, false})
		}
	}
	return result, nil
}

func allowableClass(c string) bool {
	for _, allowed := range allowedClasses {
		if c == allowed {
			return true
		}
	}
	return false
}

// This method only works if we assume Google Translate always return well formed responses
// TODO: write some test cases
func parse(in io.Reader, parent *token) error {
	var builder strings.Builder
	currentChild := &token{}
	for b, success := next(in); success; b, success = next(in) {
		switch b {
		case '[':
			if builder.Len() > 0 {
				return fmt.Errorf("expected a \", or , not [")
			}
			if err := parse(in, currentChild); err != nil {
				return err
			}
			parent.children = append(parent.children, currentChild)
			currentChild = &token{}
		case ']':
			if builder.Len() > 0 {
				currentChild.value = builder.String()
				parent.children = append(parent.children, currentChild)
			}
			return nil
		case ',':
			if builder.Len() > 0 {
				currentChild.value = builder.String()
				builder = strings.Builder{}
				parent.children = append(parent.children, currentChild)
				currentChild = &token{}
			}
		case '"':
		case '\n':
			break
		default:
			builder.WriteByte(b)
		}
	}
	return nil
}

func next(in io.Reader) (byte, bool) {
	b := make([]byte, 1)
	i, err := in.Read(b)
	return b[0], i == 1 && err == nil
}
