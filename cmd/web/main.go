package main

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/nt2311-vn/snippetbox/internal/models"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	snippets      *models.SnippetModel
	templateCache map[string]*template.Template
}

func main() {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	if err := godotenv.Load(); err != nil {
		errorLog.Fatalln("cannot read .env  file", err)
	}
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", os.Getenv("MYSQL_CONNSTR"), "MySQL data source name")

	flag.Parse()

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatalln("cannot connect to db", err)
	}

	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatalln("cannot init new template cache: ", err)
	}

	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		snippets:      &models.SnippetModel{DB: db},
		templateCache: templateCache,
	}

	server := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s\n", *addr)
	if err := server.ListenAndServe(); err != nil {
		errorLog.Fatalln("error on starting server", err)
	}
}

type neuteredFileSystem struct {
	fs http.FileSystem
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if err != nil {
		return nil, err
	}

	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := nfs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}
			return nil, err
		}
	}

	return f, nil
}

func initTLSConfig() error {
	certPool := x509.NewCertPool()
	certFile := filepath.Join("ca.pem")

	pem, err := os.ReadFile(certFile)
	if err != nil {
		return err
	}

	if !certPool.AppendCertsFromPEM(pem) {
		return err
	}

	mysql.RegisterTLSConfig("aiven", &tls.Config{
		RootCAs: certPool,
	})

	return nil
}

func openDB(dsn string) (*sql.DB, error) {
	if err := initTLSConfig(); err != nil {
		return nil, err
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
