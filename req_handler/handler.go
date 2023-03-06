package req_handler

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	database "github.com/dkmats/blog-app-sample/db"
	"github.com/jmoiron/sqlx"
)

var PageTemplates = make(map[string]*template.Template)

func LoadTemplate(name string) *template.Template {
	funcMap := template.FuncMap{
		"isEven": func(x int) bool { return x%2 == 0 },
	}
	t := template.Must(template.New(name).Funcs(funcMap).ParseFiles(
		"template/"+name,
		"template/_header.html",
		"template/_footer.html",
		"template/_item_col.html",
	))
	return t
}

func MakeHandler(db *sqlx.DB, fn func(http.ResponseWriter, *http.Request, *sqlx.DB)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, db)
	}
}

// handler of `/`
func IndexHandler(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	switch r.Method {
	case http.MethodGet:
		stmt := "SELECT * FROM article ORDER BY created_at DESC"
		articleList := []database.Article{}
		if err := db.Select(&articleList, stmt); err != nil {
			log.Println(err)
		}
		if len(articleList) > 20 {
			articleList = articleList[:20]
		}
		if err := PageTemplates["index"].Execute(w, articleList); err != nil {
			log.Printf("failed to execute template: %v", err)
			errorTpl := template.Must(template.ParseFiles("template/error.html"))
			errorTpl.Execute(w, nil)
		}
	default:
		log.Println("IndexHandler: unsupported method was sent")
	}
}

// handler of `/new`
func CreateArticleHandler(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	switch r.Method {
	case http.MethodGet:
		if err := PageTemplates["create"].Execute(w, nil); err != nil {
			log.Printf("failed to execute template: %v", err)
			errorTpl := template.Must(template.ParseFiles("template/error.html"))
			errorTpl.Execute(w, nil)
		}
	case http.MethodPost:
		r.ParseForm()
		if r.Form["TitleInput"][0] != "" && r.Form["BodyTextArea"][0] != "" {
			article := database.Article{
				Title:       r.Form["TitleInput"][0],
				Body:        r.Form["BodyTextArea"][0],
				Tag:         r.Form["ArticleTag"][0],
				CreatedTime: time.Now(),
			}
			database.InsertArticle(db, article)
			if err := PageTemplates["createConfirm"].Execute(w, article); err != nil {
				log.Printf("failed to execute template: %v", err)
				errorTpl := template.Must(template.ParseFiles("template/error.html"))
				errorTpl.Execute(w, nil)
			}
		} else {
			http.Redirect(w, r, "/new", http.StatusSeeOther)
		}
	default:
		log.Println("unsupported method was sent:")
	}
}

// handler of `/article/?id=id`
func ReadArticleHandler(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	if r.Method != http.MethodGet {
		return
	}
	query := r.URL.Query()
	stmt := "SELECT * FROM article WHERE id == " + fmt.Sprint(query["id"][0]) + ";"
	article := []database.Article{}
	if err := db.Select(&article, stmt); err != nil {
		log.Println(err)
	}
	if err := PageTemplates["article"].Execute(w, article[0]); err != nil {
		log.Printf("failed to execute template: %v", err)
		errorTpl := template.Must(template.ParseFiles("template/error.html"))
		errorTpl.Execute(w, nil)
	}
}
