package main

import (
	"crypto/x509"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"os"
	"regexp"
	"time"
)

const (
	usageMethod        = "method       Method to use during the http request"
	usageReqTotal      = "requests     Number of requests to perform"
	usageReqConcurrent = "concurrency  Number of multiple requests to make at a time"
	usageReqTimeout    = "timeout      Seconds to max. wait for each response"
)

type AppData struct {
	URL        string         `json:"url,omitempty"`
	Method     string         `json:"method,omitempty"`
	URLInfo    *url.URL       `json:"urlinfo,omitempty"`
	Concurrent int            `json:"concurrent,omitempty"`
	Requests   int            `json:"requests,omitempty"`
	Timeout    int            `json:"timeout,omitempty"`
	Summary    map[string]int `json:"summary,omitempty"`
	Millis     int64          `json:"millis,omitempty"`
	Elapsed    int64          `json:"elapsed,omitempty"`
}

var (
	pLogDir = "."
	//loggers
	infoLog  *log.Logger
	warnLog  *log.Logger
	errorLog *log.Logger

	//signal flag
	pStillRunning = true
	pShowConsole  = true

	pBuildTime = "0"
	pVersion   = "0.1.0" + "-" + pBuildTime

	//envt
	pEnvVars = map[string]*string{
		"PARASQL_LDIR": &pLogDir,
	}

	//ssl certs
	pool *x509.CertPool

	//stats
	pStats *StatsHelper

	//params
	pHTTPMethod    = "GET"
	pReqTotal      = 1
	pReqConcurrent = 1
	pReqTimeout    = 60
	pReqURI        = ""

	pAppData *AppData
)

type logOverride struct {
	Prefix string `json:"prefix,omitempty"`
}

func init() {
	//uniqueness
	rand.Seed(time.Now().UnixNano())

	pAppData = &AppData{Summary: make(map[string]int),
		Method: pHTTPMethod}
	//evt
	initEnvParams()

	//loggers
	initLogger(os.Stdout, os.Stdout, os.Stderr)

	//global vars

	pStats = StatsHelperNew()
	//init certs
	pool = x509.NewCertPool()
	pool.AppendCertsFromPEM(pemCerts)

}

//initRecov is for dumpIng segv in
func initRecov() {
	//might help u
	defer func() {
		recvr := recover()
		if recvr != nil {
			fmt.Println("MAIN-RECOV-INIT: ", recvr)
		}
	}()
}

//initEnvParams enable all OS envt vars to reload internally
func initEnvParams() {
	//just in-case, over-write from ENV
	for k, v := range pEnvVars {
		if os.Getenv(k) != "" {
			*v = os.Getenv(k)
		}
	}

	//fmt
	flag.StringVar(&pHTTPMethod, "m", pHTTPMethod, usageMethod)
	flag.IntVar(&pReqTotal, "r", pReqTotal, usageReqTotal)
	flag.IntVar(&pReqConcurrent, "c", pReqConcurrent, usageReqConcurrent)
	flag.IntVar(&pReqTimeout, "t", pReqTimeout, usageReqTimeout)

	flag.Parse()

	if pReqTotal <= 0 || pReqConcurrent <= 0 || pReqConcurrent > pReqTotal || len(os.Args) <= 1 {
		showMessage()
	}
	//last param is URL
	pReqURI = os.Args[len(os.Args)-1]
	if ok, _ := regexp.MatchString("^http(s)?://", pReqURI); !ok {
		showMessage()
	}

	//URL:&url.URL{Scheme:"https", Opaque:"", User:(*url.Userinfo)(nil), Host:"google.com:443", Path:"/search", RawPath:"", ForceQuery:false, RawQuery:"q=golang", Fragment:""}
	var err error
	pAppData.URLInfo, err = url.Parse(pReqURI)
	if err != nil {
		showMessage()
	}
	if pAppData.URLInfo.Scheme == "" || pAppData.URLInfo.Host == "" {
		showMessage()
	}
	pAppData.URL = pReqURI
	pAppData.Concurrent = pReqConcurrent
	pAppData.Method = pHTTPMethod
	pAppData.Requests = pReqTotal
}

func showMessage() {

	msg := `
Version ` + pVersion + `

Usage: gohttpdummy [options] [http[s]://]hostname[:port]/path

	   Options are:
	
`
	fmt.Println()
	fmt.Println()
	fmt.Println(msg)
	flag.PrintDefaults()
	fmt.Println()
	os.Exit(0)
}
