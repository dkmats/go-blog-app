package db

import (
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Article struct {
	ID          uint32    `db:"id"`
	Tag         string    `db:"tag"`
	Author      string    `db:"author"`
	Title       string    `db:"title"`
	Body        string    `db:"body"`
	CreatedTime time.Time `db:"created_at"`
	UpdatedTime time.Time `db:"updated_at"`
}

var schema string = `
CREATE TABLE article (
	id         integer not null primary key autoincrement,
	author     text,
	tag        text,
	title      text,
	body       text,
	created_at date,
	updated_at date
);`

func CreateArticleTable(db *sqlx.DB) {
	if _, err := db.Exec(schema); err != nil {
		log.Printf("%q: %s\n", err, schema)
	}
}

func InsertArticle(db *sqlx.DB, article Article) {
	db.NamedExec(`
	INSERT INTO article (author, tag, title, body, created_at, updated_at)
	VALUES (:author, :tag, :title, :body, :created_at, :updated_at);
	`, article)
}

func CountRows(db *sqlx.DB) int {
	var rowNum int
	db.Get(&rowNum, "SELECT COUNT (*) FROM article;")
	return rowNum
}
