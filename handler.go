package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	gor "github.com/gorilla/http"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"runtime"
	"strings"
	"sync"
	"time"
)

func b(s string) *bufio.Reader { return bufio.NewReader(strings.NewReader(s)) }

//handle main entry
func handle() {

	uFlag := make(chan bool)
	vFlag := make(chan bool)
	uwg := new(sync.WaitGroup)
	vwg := new(sync.WaitGroup)

	runtime.GOMAXPROCS(runtime.NumCPU())
	http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = 10000

	pTunnelURL = make(chan string, pAppData.Concurrent)
	if false {
		for k := 0; k < pAppData.Requests; k++ {
			uwg.Add(1)
			go process(uFlag, uwg)
			if k%pAppData.Concurrent == 0 && k > 0 {
				uwg.Wait()
			}
		}
	}
	if true {
		uwg.Add(1)
		go getTunnel(uFlag, uwg)

		vwg.Add(1)
		go setTunnel(vFlag, vwg)

		vwg.Wait()
		close(pTunnelURL)
	}

	uwg.Wait()
	close(uFlag)
	close(vFlag)

}

//setTunnel add the data at the end-of-tunnel
func setTunnel(doneFlg chan bool, wg *sync.WaitGroup) {
	go func() {
		for {
			select {
			//wait till doneFlag has value ;-)
			case <-doneFlg:
				//done already ;-)
				wg.Done()
				return
			}
		}
	}()

	//go parsql.Save(uFlag, uwg)
	for k := 0; k < pAppData.Requests; k++ {
		//sig-check
		if !pStillRunning {
			log.Println("Signal detected ...")
			doneFlg <- true
			return
		}
		pTunnelURL <- pAppData.URL
	}

	//send signal -> DONE
	doneFlg <- true
}

//getTunnel process it here
func getTunnel(doneFlg chan bool, wg *sync.WaitGroup) {

	go func() {
		for {
			select {
			//wait till doneFlag has value ;-)
			case <-doneFlg:
				//done already ;-)
				wg.Done()
				return
			}
		}
	}()

	uFlag := make(chan bool)
	uwg := new(sync.WaitGroup)
	for {
		_, ok := <-pTunnelURL
		if !ok {
			break
		}
		//sig-check
		if !pStillRunning {
			log.Println("Signal detected ...")
			break
		}
		uwg.Add(1)
		go process(uFlag, uwg)
		runtime.Gosched()
	}
	uwg.Wait()
	close(uFlag)
	//send signal -> DONE
	doneFlg <- true
}

//showSummary list all the results statistics
func showSummary(t0 time.Time) {
	//elapsed
	msg := `

Version ` + pVersion + `

Benchmarking is now in progress .... 

Please be patient!

Statistics :
`
	fmt.Println(msg)
	fmt.Println()
	h := strings.SplitN(pAppData.URLInfo.Host, ":", 2)
	fmt.Println("Server Hostname:", h[0])
	if len(h) > 1 {
		fmt.Println("Server Port    :", h[1])
	}
	fmt.Println("Document Path  :", pAppData.URLInfo.Path)
	fmt.Println("Request Method :", pAppData.Method)
	fmt.Println()
	slist := pStats.getStatsList()
	for k, v := range slist {
		fmt.Println(strings.TrimSpace(k), ": ", v)
	}
	t1 := time.Since(t0)
	pAppData.Elapsed = float64(t1.Nanoseconds()) / 1000000 //millis = nanos / 1,000,000
	fmt.Println("Elapsed : ", fmt.Sprintf("%.06f", pAppData.Elapsed), "( millisecs )")
	fmt.Println("Requests: ", fmt.Sprintf("%.06f", (float64(pAppData.Requests)*float64(1000))/float64(pAppData.Elapsed)), "( # per sec )")
	fmt.Println("App Time: ", t1.String())
	fmt.Println("Sys Time: ", time.Since(t0).String())
}

//process
func process(doneFlg chan bool, wg *sync.WaitGroup) {

	go func() {
		for {
			select {
			//wait till doneFlag has value ;-)
			case <-doneFlg:
				//done already ;-)
				wg.Done()
				return
			}
		}
	}()
	var statuscode int

	t0 := time.Now()
	if pAppData.Method == "GET" {
		statuscode = getResultFast(pAppData.URL)
	} else {
		statuscode = postResultFast(pAppData.URL, pFormData)
	}
	//calc
	t1 := time.Since(t0)
	pAppData.Millis += int64(t1.Nanoseconds() / 1000000)

	//http.StatusText(statuscode)
	if statuscode != http.StatusOK {
		pStats.setStats("Failed")
	} else {
		pStats.setStats("Success")
	}
	//send signal -> DONE
	doneFlg <- true
}

//getResult http req a url
func getResult(url string) int {
	//client
	c := setupHttpClient()

	//init
	var res *http.Response
	var err error

	//Get
	res, err = c.Get(url)
	//make sure to free-up
	if res != nil {
		defer res.Body.Close()
	}
	if err != nil {
		log.Println("ERROR: getResult:", err)
		return 0
	}
	//get response
	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("ERROR: getResult:", err)
		return 0
	}
	//give
	return res.StatusCode
}

//getResultFast http req a url
func getResultFast(url string) int {

	//init
	statRet, _, r, err := gor.DefaultClient.Get(url, nil)
	if err != nil {
		log.Println("ERROR: getResultFast:", err)
		return 0
	}
	if r != nil {
		defer r.Close()
		/**
		var body []byte
			body, err = ioutil.ReadAll(r)
			if err != nil {
				log.Println("ERROR: getResultFast:", err)
				return 0
			}**/
	}
	//give
	if !statRet.IsSuccess() {
		log.Println("ERROR: getResultFast: Invalid status code", statRet.String())
		return 0
	}
	return statRet.Code
}

//postResultFast http req a url
func postResultFast(url string, form *url.Values) int {

	//init
	phdrs := make(map[string][]string)
	pbody := strings.NewReader(form.Encode())
	statRet, _, r, err := gor.DefaultClient.Post(url, phdrs, pbody)
	if err != nil {
		log.Println("ERROR: postResultFast:", err)
		return 0
	}
	if r != nil {
		defer r.Close()
		/**
		var body []byte
			body, err = ioutil.ReadAll(r)
			if err != nil {
				log.Println("ERROR: getResultFast:", err)
				return 0
			}**/
	}
	//give
	if !statRet.IsSuccess() {
		log.Println("ERROR: postResultFast: Invalid status code", statRet.String())
		return 0
	}
	return statRet.Code
}

//postResult
func postResult(uri string, form *url.Values) int {
	//client
	c := setupHttpClient()

	//init
	var req *http.Request
	var res *http.Response
	var err error

	//Post
	req, err = http.NewRequest("POST", uri, strings.NewReader(form.Encode()))
	if err != nil {
		fmt.Println("ERROR: postResult:", err)
		return 0
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, err = c.Do(req)
	//make sure to free-up
	if res != nil {
		defer res.Body.Close()
	}
	if err != nil {
		log.Println("ERROR: postResult:", err)
		return 0
	}
	//get response
	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("ERROR: postResult:", err)
		return 0
	}
	//give
	return res.StatusCode
}

//getTimeoutCfg get timeout settings
func getTimeoutCfg() time.Duration {
	var timeout time.Duration
	timeout = 120
	if pReqTimeout > 0 {
		timeout = time.Duration(pReqTimeout)
	}
	return timeout
}

func setupHttpClient() *http.Client {
	//client
	return &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true, RootCAs: pool},
		Dial: (&net.Dialer{
			Timeout: getTimeoutCfg() * time.Second,
		}).Dial,
		//DisableKeepAlives: true,
		MaxIdleConnsPerHost: 10000,
	},
	}
}
