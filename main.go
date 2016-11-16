package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {

	//might help u
	defer func() {
		recvr := recover()
		if recvr != nil {
			fmt.Println("MAIN-RECOV-INIT: ", recvr)
		}
	}()
	//elapsed
	t0 := time.Now()

	handle()

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
		fmt.Println(strings.TrimSpace(k), ":", v)
	}
	t1 := time.Since(t0)
	pAppData.Elapsed = int64(t1.Nanoseconds()/1000) / int64(1000)
	fmt.Println("Elapsed :", pAppData.Elapsed, "millisecs")
	fmt.Println("Requests:", fmt.Sprintf("%.04f", (float64(pAppData.Requests)*float64(1000))/float64(pAppData.Elapsed)), " ( # per sec )")
	os.Exit(0)
}
