package limoBookingApp

import (
	"html/template"
	"io"
	"log"
	"os"
	"sync"
)

//variable declarations for templates and WaitGroup.
var (
	tpl *template.Template
	wg  sync.WaitGroup
)

//variable declaration for log files.
var (
	errLog   *os.File
	userLog  *os.File
	adminLog *os.File
)

//variable declaration for logger types.
var (
	ErrorLogger *log.Logger
	UserLogger  *log.Logger
	AdminLogger *log.Logger
)

//funcMap declaration.
var funcMap = template.FuncMap{
	"add": Add,
}

func init() {
	var err error
	tpl = template.Must(template.New("").Funcs(funcMap).ParseGlob("templates/*"))
	errLog, err = os.OpenFile("logs/errors.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open error log file:", err)
	}
	userLog, err = os.OpenFile("logs/userLoginAndLogout.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open user log file:", err)
	}
	adminLog, err = os.OpenFile("logs/adminLoginAndLogout.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open admin log file:", err)
	}

	flags := log.LstdFlags | log.Lshortfile
	ErrorLogger = log.New(io.MultiWriter(errLog, os.Stderr), "ERROR: ", flags)
	UserLogger = log.New(io.MultiWriter(userLog, os.Stderr), "USER LOG: ", flags)
	AdminLogger = log.New(io.MultiWriter(adminLog, os.Stderr), "ADMIN LOG: ", flags)
}
