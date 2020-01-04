package data

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"sort"
	"time"

	// MySQL driver
	_ "github.com/go-sql-driver/mysql"
)

type structuredDB struct {
	db *sql.DB
}

// InitStructured returns a usable instance of database with a relationship database underneath
func InitStructured(host, username, password string) (Database, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:3306)/vocabpractice?parseTime=true", username, password, host))
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &structuredDB{db}, nil
}

// FetchNext get one of the least recently tested vocab/translation pair
func (d *structuredDB) FetchNext() (*Word, error) {
	// Look up 10 to add some randomness to it.
	rows, err := d.db.Query("SELECT vocab FROM vocabs ORDER BY last_test ASC LIMIT 10")
	if err != nil {
		log.Fatal(err)
	}
	var vocabs []string
	for rows.Next() {
		var vocab string
		if err := rows.Scan(&vocab); err != nil {
			log.Fatal(err)
		}
		vocabs = append(vocabs, vocab)
	}
	vocab := vocabs[rand.Intn(len(vocabs))]
	return d.QueryWord(vocab)
}

// QueryWord fetches all the translations of a vocab and the last time the translations are tested
func (d *structuredDB) QueryWord(vocab string) (*Word, error) {
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

func (d *structuredDB) checkExist(vocab, translation string) (bool, error) {
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
func (d *structuredDB) Pass(vocab, translation string) error {
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
func (d *structuredDB) Input(vocab, translation string) error {
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
func (d *structuredDB) List() ([]*Word, error) {
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
