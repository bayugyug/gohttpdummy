package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"regexp"
	"strings"
	"time"
)

//os.Stdout, os.Stdout, os.Stderr
func initLogger(i, w, e io.Writer) {
	//just in case
	infoLog = log.New(i,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lmicroseconds)
	warnLog = log.New(w,
		"WARN: ",
		log.Ldate|log.Ltime|log.Lshortfile)
	errorLog = log.New(e,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

//formatLogger try to init all filehandles for logs
func formatLogger(fdir, fname, pfx string) string {
	t := time.Now()
	r := regexp.MustCompile("[^a-zA-Z0-9]")
	p := t.Format("2006-01-02") + "-" + r.ReplaceAllString(strings.ToLower(pfx), "")
	s := path.Join(pLogDir, fdir)
	if _, err := os.Stat(s); os.IsNotExist(err) {
		//mkdir -p
		os.MkdirAll(s, os.ModePerm)
	}
	return path.Join(s, p+"-"+fname+".log")
}

//makeLogger initialize the logger either via file or console
func makeLogger(w io.Writer, ldir, fname, pfx string) *log.Logger {
	logFile := w
	if !pShowConsole {
		var err error
		logFile, err = os.OpenFile(formatLogger(ldir, fname, pfx), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
		if err != nil {
			log.Println(err)
		}
	}
	//give it
	return log.New(logFile,
		pfx,
		log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)

}

//dumpW log into warning
func dumpW(s ...interface{}) {
	warnLog.Println(s...)
}

//dumpE log into error
func dumpE(s ...interface{}) {
	errorLog.Println(s...)
}

//dumpI log into info
func dumpI(s ...interface{}) {
	infoLog.Println(s...)
}

//Write override the log.print
func (w logOverride) Write(bytes []byte) (int, error) {
	//return fmt.Print(w.Prefix + time.Now().UTC().Format("2006-01-02 15:04:05.999") + " " + string(bytes))
	return fmt.Print(string(bytes))
}

//overrideLogger reset the log.print to customized
func overrideLogger(pfx string) {
	log.SetFlags(0)
	log.SetOutput(&logOverride{Prefix: pfx})
}
