package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

func b(s string) *bufio.Reader { return bufio.NewReader(strings.NewReader(s)) }

//handle main entry
func handle() {

	uFlag := make(chan bool)
	uwg := new(sync.WaitGroup)

	//go parsql.Save(uFlag, uwg)
	for k := 0; k < pAppData.Requests; k++ {
		uwg.Add(1)
		go process(uFlag, uwg, pAppData.Method, pAppData.URL)
		if k%pAppData.Concurrent == 0 && k > 0 {
			uwg.Wait()
		}
	}
	uwg.Wait()
	close(uFlag)
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
func process(doneFlg chan bool, wg *sync.WaitGroup, method, url string) {

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
	if method == "GET" {
		statuscode = getResult(url)
	} else {
		statuscode = postResult(url, pFormData)
	}
	//calc
	t1 := time.Since(t0)
	pAppData.Millis += int64(t1.Nanoseconds()/1000) / int64(1000)

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
	c := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true, RootCAs: pool},
		Dial: (&net.Dialer{
			Timeout: getTimeoutCfg() * time.Second,
		}).Dial,
		//DisableKeepAlives: true,
	},
	}

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

//postResult
func postResult(uri string, form *url.Values) int {
	//client
	c := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true, RootCAs: pool},
		Dial: (&net.Dialer{
			Timeout: getTimeoutCfg() * time.Second,
		}).Dial,
		//DisableKeepAlives: true,
	},
	}

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
