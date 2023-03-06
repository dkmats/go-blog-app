package main

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"

	database "github.com/dkmats/blog-app-sample/db"
	"github.com/dkmats/blog-app-sample/req_handler"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var Setting setting

func init() {
	Setting = getSetting()

	req_handler.PageTemplates["index"] = req_handler.LoadTemplate("index.html")
	req_handler.PageTemplates["create"] = req_handler.LoadTemplate("create.html")
	req_handler.PageTemplates["article"] = req_handler.LoadTemplate("article.html")
	req_handler.PageTemplates["createConfirm"] = req_handler.LoadTemplate("create_confirm.html")
}

func main() {
	db, err := sqlx.Open("sqlite3", Setting.dbPath)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer db.Close()

	if f, err := os.Stat(Setting.dbPath); errors.Is(err, fs.ErrNotExist) {
		database.CreateArticleTable(db)
	} else if f.IsDir() {
		log.Printf("same name directory already exist on database path: %s\n", Setting.dbPath)
		os.Exit(1)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", req_handler.MakeHandler(db, req_handler.IndexHandler))
	mux.HandleFunc("/new/", req_handler.MakeHandler(db, req_handler.CreateArticleHandler))
	mux.HandleFunc("/article", req_handler.MakeHandler(db, req_handler.ReadArticleHandler))

	fmt.Printf("listening on port %s...\n", Setting.portNum)
	http.ListenAndServe(":"+Setting.portNum, mux)
}
