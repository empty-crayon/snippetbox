package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"

	"github.com/empty-crayon/snippetbox/internal/models"
	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	logger *slog.Logger
	// Allows us to make the snippet model object available to our handlers
	snippets      *models.SnippetModel
	templateCache map[string]*template.Template
}

func main() {

	addr := flag.String("addr", ":4000", "HTTP Network address")
	dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")

	flag.Parse()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,            // adds caller location in logs, line number
		Level:     slog.LevelDebug, // sets minimum level to debug
	}))

	db, err := openDB(*dsn)

	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer db.Close()

	templateCache, err := newTemplateCache()

	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	app := &application{
		logger:        logger,
		snippets:      &models.SnippetModel{DB: db},
		templateCache: templateCache,
	}

	logger.Info("Starting Server", "addr", *addr)

	err = http.ListenAndServe(*addr, app.routes())

	logger.Error(err.Error())
	os.Exit(1)
}

func openDB(dsn string) (*sql.DB, error) {
	// first parameter is the driver name,
	// second parameter is the dsn: data source name
	// returns a sql.DB object: not a db connection but a pool for many connections
	// go manages connections via this pool, auto opening and closing
	// concurrency safe
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil

}
