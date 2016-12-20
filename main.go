package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	//timing
	t0 := time.Now()
	//might help u
	defer func() {
		recvr := recover()
		if recvr != nil {
			fmt.Println("MAIN-RECOV-INIT: ", recvr)
		}
	}()
	//do
	handle()
	//stats
	showSummary(t0)
	//good ;-)
	os.Exit(0)
}
