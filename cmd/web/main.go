package main

import (
	"database/sql"
	"flag"
	"github.com/aifuxi/snippetbox/internal/models"
	"github.com/go-playground/form/v4"
	"html/template"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type config struct {
	addr      string
	staticDir string
	dsn       string
}

type application struct {
	infoLog       *log.Logger
	errorLog      *log.Logger
	cfg           config
	snippets      *models.SnippetModel
	templateCache map[string]*template.Template
	formDecoder   *form.Decoder
}

func main() {
	var cfg config

	// 定义一个标志 addr
	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.staticDir, "staticDir", "./ui/static", "Static files path")
	flag.StringVar(&cfg.dsn, "dsn", "root:123456@tcp(127.0.0.1:33060)/snippetbox?parseTime=true", "MySQL data source name")

	flag.Parse()

	f, err := os.OpenFile("./logs/info.log", os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	//infoLog := log.New(f, "[INFO]\t", log.Ldate|log.Ltime)
	infoLog := log.New(os.Stdout, "[INFO]\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "[ERROR]\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(cfg.dsn)
	if err != nil {
		errorLog.Fatalln(err)
	}

	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatalln(err)
	}

	defer db.Close()

	formDecoder := form.NewDecoder()

	app := &application{
		infoLog:  infoLog,
		errorLog: errorLog,
		cfg:      cfg,
		snippets: &models.SnippetModel{
			DB: db,
		},
		templateCache: templateCache,
		formDecoder:   formDecoder,
	}

	infoLog.Printf("Starting server on %v", cfg.addr)

	srv := http.Server{
		Addr:     cfg.addr,
		Handler:  app.routes(),
		ErrorLog: errorLog,
	}

	err = srv.ListenAndServe()
	if err != nil {
		errorLog.Fatalln(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
