package api

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	vercel "github.com/tbxark/g4vercel"
	"net/http"
	"os"
	"regexp"
	"strconv"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	server := vercel.New()
	db := InitRepo()

	server.GET("/", func(context *vercel.Context) {
		context.JSON(200, vercel.H{"endpoints": "article and article/:id"})
	})

	server.GET("/article", func(context *vercel.Context) {
		articles, err := GetAllArticles(context, db)
		if err != nil {
			context.JSON(400, vercel.H{"message": "Something is not yes."})
		}
		context.JSON(200, articles)
	})

	server.GET("/article/:id", func(context *vercel.Context) {
		article, err := GetArticleFromPage(context, db)
		if err != nil {
			context.JSON(400, vercel.H{"message": "Something is not yes."})
		}

		context.JSON(200, article)
	})

	server.Handle(w, r)
}

type DatabaseCredentialsStruct struct {
	HOST     string
	PORT     int
	USER     string
	PASSWORD string
	NAME     string
	REPO     string
}

func InitRepo() *sql.DB {
	PORT, _ := strconv.Atoi(os.Getenv("DB_PORT"))
	dbc := &DatabaseCredentialsStruct{
		HOST:     os.Getenv("DB_HOST"),
		PORT:     PORT,
		USER:     os.Getenv("DB_USER"),
		PASSWORD: os.Getenv("DB_PASSWORD"),
		NAME:     os.Getenv("DB_NAME"),
		REPO:     "SUPABASE",
	}

	psql := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbc.HOST, dbc.PORT, dbc.USER, dbc.PASSWORD, dbc.NAME)
	db, err := sql.Open("postgres", psql)

	if err != nil {
		fmt.Println(err.Error())
		panic(fmt.Sprintf("%s DB Failed", dbc.REPO))
	}

	return db
}

type ArticleStruct struct {
	Id          string   `json:"id"`
	Page        int      `json:"page"`
	Content     string   `json:"content"`
	ReadingTime int      `json:"readingTime"`
	Header      string   `json:"header"`
	AuthorId    []string `json:"authorId"`
}

func GetAllArticles(context *vercel.Context, db *sql.DB) ([]ArticleStruct, error) {
	rows, err := db.Query(`SELECT * FROM articles`)
	var articles []ArticleStruct
	if err != nil {
		return articles, errors.New("Cannot get from DB")
	}

	defer rows.Close()

	for rows.Next() {
		var article ArticleStruct
		rows.Scan(&article.Id, &article.Page, &article.Content, &article.ReadingTime, &article.Header, &article.AuthorId)

		articles = append(articles, article)
	}

	return articles, nil
}

func GetArticleFromPage(context *vercel.Context, db *sql.DB) (ArticleStruct, error) {
	paramId := context.Param("id")
	id, _ := strconv.Atoi(paramId)
	isInt := regexp.MustCompile("[0-9]+").MatchString(paramId)
	var article ArticleStruct

	if isInt == false {
		return article, errors.New("ID isn't correct")
	}

	rows, _ := db.Query(`SELECT * FROM articles WHERE page=$1`, id)
	defer rows.Close()

	for rows.Next() {
		rows.Scan(&article.Id, &article.Page, &article.Content, &article.ReadingTime, &article.Header, &article.AuthorId)
	}

	return article, nil
}
