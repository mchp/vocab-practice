CREATE DATABASE vocabpractice;
USE vocabpractice;

CREATE TABLE vocabs (
    vocab       VARCHAR(150) NOT NULL,
    translation VARCHAR(150) NOT NULL,
    last_test   DATETIME,
    PRIMARY KEY (vocab, translation)
)
