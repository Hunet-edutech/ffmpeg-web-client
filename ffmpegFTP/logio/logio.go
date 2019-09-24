package logio

import (
	"log"
	"os"
	"strings"
)

var (
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

func Init() {
	// now := time.Now()
	// custom := now.Format("2000-01-01")

	exePath, _ := os.Executable()
	exePaths := []rune(exePath)
	dirPath := string(exePaths[0:strings.LastIndex(exePath, "\\")])

	fpLog, _ := os.OpenFile(dirPath+"\\logfile.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	Trace = log.New(fpLog,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Info = log.New(fpLog,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(fpLog,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(fpLog,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)

}
