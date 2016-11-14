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
	var statusdesc string

	t0 := time.Now()
	if method == "GET" {
		statuscode, statusdesc = getResult(url)
	} else {
		fparams := make(map[string]string)
		statuscode, statusdesc = postResult(url, fparams)
	}
	//calc
	t1 := time.Since(t0)
	pAppData.Millis += int64(t1.Nanoseconds()/1000) / int64(1000)

	if statuscode != 200 || statusdesc == "" {
		pAppData.Summary["Failed"]++
		pStats.setStats("Failed")
	} else {
		pAppData.Summary["Success"]++
		pStats.setStats("Success")
	}
	//send signal -> DONE
	doneFlg <- true
}

//getResult http req a url
func getResult(url string) (int, string) {
	var timeout time.Duration
	timeout = 120
	if pReqTimeout > 0 {
		timeout = time.Duration(pReqTimeout)
	}
	//client
	c := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true, RootCAs: pool},
		Dial: (&net.Dialer{
			Timeout: timeout * time.Second,
		}).Dial,
		//DisableKeepAlives: true,
	},
	}
	res, err := c.Get(url)
	//make sure to free-up
	if res != nil {
		defer res.Body.Close()
	}
	if err != nil {
		log.Println("ERROR: getResult:", err)
		return 0, ""
	}
	//get response
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("ERROR: getResult:", err)
		return 0, ""
	}
	//give
	return res.StatusCode, string(body)
}

//postResult
func postResult(uri string, fparams map[string]string) (int, string) {
	var timeout time.Duration
	timeout = 120
	if pReqTimeout > 0 {
		timeout = time.Duration(pReqTimeout)
	}
	//client
	c := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true, RootCAs: pool},
		Dial: (&net.Dialer{
			Timeout: timeout * time.Second,
		}).Dial,
		//DisableKeepAlives: true,
	},
	}
	form := &url.Values{}
	for xk, xv := range fparams {
		form.Add(xk, xv)
	}
	req, errs := http.NewRequest("POST", uri, strings.NewReader(form.Encode()))
	if errs != nil {
		fmt.Println("ERROR: postResult:", errs)
		return 0, ""
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, errv := c.Do(req)
	//make sure to free-up
	if res != nil {
		defer res.Body.Close()
	}
	if errv != nil {
		log.Println("ERROR: postResult:", errv)
		return 0, ""
	}
	//get response
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("ERROR: postResult:", err)
		return 0, ""
	}
	//give
	return res.StatusCode, string(body)
}
