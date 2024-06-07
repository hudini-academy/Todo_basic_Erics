package main

import (
	"database/sql"
	"flag"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golangcollege/sessions"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
	"todo/pkg/models/mysql"
)

// Config struct stores configuration settings for the application.
type Config struct {
	Addr      string
	StaticDir string
	Dsn       string
	secretKey string
}

// Application struct stores application-wide dependencies.
type Application struct {
	infolog       *log.Logger 
	errorlog      *log.Logger
	config        *Config 
	todo          *mysql.TodoModel 
	templateCache map[string]*template.Template 
	session       *sessions.Session	
	users         *mysql.UserModel
}

func main() {
	// Creating new instance and setting default settings.
	config := new(Config)
	flag.StringVar(&config.Addr, "addr", ":4000", "Default server address")
	flag.StringVar(&config.StaticDir, "static-dir", "./ui/static", "Path to static directory")
	flag.StringVar(&config.Dsn, "dsn", "root:root@/todo?parseTime=true", "Mysql connection")
	flag.StringVar(&config.secretKey, "secretKey", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret Key for session storage")
	flag.Parse()

	// Open the file to log the errors.
	f, err := os.OpenFile("./tmp/error.log", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}

	// Open the file to log the information messages.
	infoFile, err := os.OpenFile("./tmp/info.log", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}

	// Creating custom loggers.
	infoLog := log.New(infoFile, "INFO\t", log.Ldate|log.Ltime|log.Lshortfile)
	errorLog := log.New(f, "ERROR\t", log.Ldate|log.Ltime|log.Llongfile)

	// Initialize session manager.
	session := sessions.New([]byte(config.secretKey))
	session.Lifetime = 2 * time.Minute

	// initialize templateCache.
	ts, errTemplateCaching := newTemplateCache("./ui/html/")
	if errTemplateCaching != nil {
		log.Fatalln(errTemplateCaching)
		errorLog.Println(errTemplateCaching)
	}

	// Initialize database connection.
	db, errDB := openDB(config.Dsn)
	if errDB != nil {
		errorLog.Fatal(errDB)
		log.Println(errDB)
	}

	// Initialize dependencies.
	app := &Application{
		errorlog:      errorLog,
		infolog:       infoLog,
		config:        config,
		todo:          &mysql.TodoModel{DB: db},
		templateCache: ts,
		session:       session,
		users:         &mysql.UserModel{DB: db},
	}

	defer db.Close()
	log.Println("Starting server on", app.config.Addr)
	infoLog.Println("Starting server on", app.config.Addr)
	errorLog.Fatal(http.ListenAndServe(app.config.Addr, app.routes()))
}

// Establishes a connection to the database and test the connection.
func openDB(dsn string) (*sql.DB, error) {
	db, errDBOpen := sql.Open("mysql", dsn)
	if errDBOpen != nil {
		log.Println(errDBOpen)
		return nil, errDBOpen
	}
	if err := db.Ping(); err != nil {
		log.Println(errDBOpen)
		return nil, err
	}
	log.Println(db)
	return db, nil
}
