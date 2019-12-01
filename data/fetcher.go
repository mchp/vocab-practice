package data

import (
	"database/sql"
	"fmt"
	"log"
	"sort"
	"time"

	// MySQL driver
	_ "github.com/go-sql-driver/mysql"
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

// TranslationAndTest
type TranslationAndTest struct {
	Translation string    `json:"translation"`
	LastTested  time.Time `json:"lastTested"`
}

// DB handles interfacing with vocab storage
type DB struct {
	db *sql.DB
}

// Init returns a usable instance of DB
func Init() (*DB, error) {
	db, err := sql.Open("mysql", "root:rainstop@tcp(:3306)/vocabpractice?parseTime=true")
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

// FetchNext get the least recently tested vocab/translation pair
func (d *DB) FetchNext() (*Word, error) {
	row := d.db.QueryRow("SELECT vocab FROM vocabs ORDER BY last_test ASC LIMIT 1")
	var vocab string
	if err := row.Scan(&vocab); err != nil {
		log.Fatal(err)
	}

	rows, err := d.db.Query("SELECT translation, last_test FROM vocabs WHERE vocab=?", vocab)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	w := &Word{Vocab: vocab}
	for rows.Next() {
		var trans string
		var lastTested sql.NullTime
		if err := rows.Scan(&trans, &lastTested); err != nil {
			log.Fatal(err)
		}
		w.Translations = append(w.Translations, &TranslationAndTest{trans, lastTested.Time})
	}
	return w, nil
}

func (d *DB) checkExist(vocab, translation string) (bool, error) {
	rows, err := d.db.Query("SELECT * FROM vocabs WHERE vocab=? AND translation=?", vocab, translation)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	for rows.Next() {
		return true, nil
	}
	return false, nil
}

// Pass should be called when the user correctly identified a vocab-translation pair
func (d *DB) Pass(vocab, translation string) error {
	exist, err := d.checkExist(vocab, translation)
	if err != nil || !exist {
		return fmt.Errorf("could not find %s -> %s: %v", vocab, translation, err)
	}
	result, err := d.db.Exec("UPDATE vocabs SET last_test=? WHERE vocab=? AND translation=?", time.Now(), vocab, translation)
	if err != nil {
		return fmt.Errorf("could not update test time for %s -> %s: %s", vocab, translation, err.Error())
	}
	if affected, _ := result.RowsAffected(); affected != 1 {
		return fmt.Errorf("updating %s -> %s test time affected %d rows", vocab, translation, affected)
	}
	return nil
}

// Input submits a new vocab translation pair into the database
func (d *DB) Input(vocab, translation string) error {
	exist, err := d.checkExist(vocab, translation)
	if err != nil {
		return err
	}
	if exist {
		return fmt.Errorf("translation pair %s -> %s already exists", vocab, translation)
	}
	_, err = d.db.Exec("INSERT INTO vocabs (vocab, translation) VALUES(?,?)", vocab, translation)
	if err != nil {
		return err
	}
	return nil
}

// List returns all the vocab and translations in storage
func (d *DB) List() ([]*Word, error) {
	rows, err := d.db.Query("SELECT vocab, translation, last_test FROM vocabs")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	m := map[string]*Word{}
	for rows.Next() {
		var vocab, translation string
		var lastTest sql.NullTime
		if err := rows.Scan(&vocab, &translation, &lastTest); err != nil {
			log.Fatal(err)
		}
		if m[vocab] == nil {
			m[vocab] = &Word{Vocab: vocab}
		}
		m[vocab].Translations = append(m[vocab].Translations, &TranslationAndTest{translation, lastTest.Time})
	}

	// sort the vocabs alphabetically
	var vocabs []string
	for k := range m {
		vocabs = append(vocabs, k)
	}
	sort.Strings(vocabs)
	var words []*Word
	for _, k := range vocabs {
		words = append(words, m[k])
	}
	return words, nil
}
