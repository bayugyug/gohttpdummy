package main

import (
	"crypto/rand"
	"fmt"
	mt "math/rand"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const randChars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

//memDmp dump the current Mem in MBytes
func memDmp() {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	m, _ := strconv.Atoi(strconv.FormatUint(mem.Alloc, 10))
	infoLog.Println("Mem", fmt.Sprintf("%.04f", float64(m)/(1024*1024)), " MB")
}

//statsDmp dump all the stats summary
func statsDmp() {
	slist := pStats.getStatsList()
	infoLog.Println("Stats Summary")
	for k, v := range slist {
		infoLog.Println(k, " -> ", v)
	}
}

//createTempStr uniq uuid generator
func createTempStr(pfx string) string {

	var uniqid string

	if len(pfx) > 0 {
		uniqid = fmt.Sprintf("%s%04X%04X%16X", pfx, mt.Intn(9999), mt.Intn(9999), time.Now().UTC().UnixNano())
	} else {
		uniqid = fmt.Sprintf("%s%04X%04X%16X", "tmf", mt.Intn(9999), mt.Intn(9999), time.Now().UTC().UnixNano())
	}

	return strings.ToUpper(uniqid + randStr(8))
}

//randStr more than 1 way to do random chars
func randStr(strSize int) string {
	var bytes = make([]byte, strSize)
	rand.Read(bytes)
	for k, v := range bytes {
		bytes[k] = randChars[v%byte(len(randChars))]
	}
	return string(bytes)
}

//timeElapsed display the time elapsed since t0
func timeElapsed(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Println(fmt.Sprintf("%s took %d ms", name, elapsed.Nanoseconds()/1000))
}
