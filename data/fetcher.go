package data

import (
  "database/sql"
  "fmt"
  "log"
  "sort"
  "time"

  _ "github.com/go-sql-driver/mysql"
)

type Word struct {
  Input string
  Outputs []*outputAndTest
}

func (w *Word) String() string {
  str := w.Input + ": ";
  for _, o := range w.Outputs {
    str = fmt.Sprintf("%s %s(%s)", str, o.output, o.tested)
  }
  return str
}

type outputAndTest struct {
  output string
  tested time.Time
}

type DB struct {
  db *sql.DB
}

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

func (d *DB) FetchNext() (*Word, error) {
  row := d.db.QueryRow("SELECT vocab FROM vocabs ORDER BY last_test DESC LIMIT 1")
  var vocab string
  if err := row.Scan(&vocab); err != nil {
    log.Fatal(err)
  }

  rows, err := d.db.Query("SELECT translation, last_test FROM vocabs WHERE vocab=?", vocab)
  if err != nil {
    log.Fatal(err)
  }
  defer rows.Close()
  w := &Word{Input: vocab}
  for rows.Next() {
    var trans string
    var lastTested sql.NullTime
    if err := rows.Scan(&trans, &lastTested); err != nil {
      log.Fatal(err)
    }
    w.Outputs = append(w.Outputs, &outputAndTest{trans, lastTested.Time})
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
    return fmt.Errorf("updating %s -> %s test time affected %i rows. Please debug.", vocab, translation, affected)
  }
  return nil
}
func (d *DB) Input(vocab, translation string) error {
  exist, err := d.checkExist(vocab, translation)
  if err != nil {
    return err
  }
  if exist {
    return nil
  }
  _, err = d.db.Exec("INSERT INTO vocabs (vocab, translation) VALUES(?,?)", vocab, translation)
  if err != nil {
    return err
  }
  return nil
}

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
      m[vocab] = &Word{Input: vocab}
    }
    m[vocab].Outputs = append(m[vocab].Outputs, &outputAndTest{translation, lastTest.Time})
  }

  var vocabs []string
  for k := range m {
    vocabs = append(vocabs, k)
  }
  sort.Strings(vocabs)
  var orderedWords []*Word
  for _, k := range vocabs {
    orderedWords = append(orderedWords, m[k])
  }
  return orderedWords, nil
}
