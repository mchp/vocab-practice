package data

import (
	"fmt"
	"time"
)

// Word represents all the translations assigned to a vocab and the last time each translation was tested and passed
type Word struct {
	Vocab        string                `json:"vocab"`
	Translations []*TranslationAndTest `json:"translations"`
}

func (w *Word) String() string {
	str := w.Vocab + ": "
	for _, o := range w.Translations {
		str = fmt.Sprintf("%s %s(%s)", str, o.Translation, o.LastTested)
	}
	return str
}

// TranslationAndTest indicates when this translation was last tested
type TranslationAndTest struct {
	Translation string    `json:"translation"`
	LastTested  time.Time `json:"lastTested"`
}

// Database is the interface of database
type Database interface {
	FetchNext() (*Word, error)
	List() ([]*Word, error)
	QueryWord(string) (*Word, error)
	Input(string, string) error
	Pass(string, string) error
}
