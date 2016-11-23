package main

import (
	"fmt"
	"os"
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
	//timing
	t0 := time.Now()
	//do
	handle()
	//stats
	showSummary(t0)
	//good ;-)
	os.Exit(0)
}
