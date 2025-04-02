package main

import (
	"database/sql"
	"flag"
	"github.com/aifuxi/snippetbox/internal/models"
	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type config struct {
	addr      string
	staticDir string
	dsn       string
}

type application struct {
	infoLog        *log.Logger
	errorLog       *log.Logger
	cfg            config
	snippets       *models.SnippetModel
	users          *models.UserModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
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

	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true // 设置为true时，只会在使用https时才发送cookie到服务器

	app := &application{
		infoLog:  infoLog,
		errorLog: errorLog,
		cfg:      cfg,
		snippets: &models.SnippetModel{
			DB: db,
		},
		users: &models.UserModel{
			DB: db,
		},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}

	infoLog.Printf("Starting server on %v", cfg.addr)

	srv := http.Server{
		Addr:         cfg.addr,
		Handler:      app.routes(),
		ErrorLog:     errorLog,
		IdleTimeout:  time.Minute,      // 单个链接最大空闲超时时间
		ReadTimeout:  5 * time.Second,  // 请求超时时间
		WriteTimeout: 10 * time.Second, // 返回响应内容的超时时间
	}

	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
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
