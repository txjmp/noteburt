package main

import (
	"flag"
	"fmt"
	"lib"
	"log"
	"math/rand"
	d "nb/data"
	"nb/schedule"
	"nb/web"
	"os"
	"os/signal"
	"time"
)

func main() {
	prod := flag.Bool("prod", false, "include -prod if running in production mode")
	traceLevel := flag.Int("trace", 0, "trace level 0-2, default is 0")
	flag.Parse()

	errLog, err := os.OpenFile("error.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.ModeAppend)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(errLog)

	if *prod { // running in production mode
		prodStart()
	} else {
		testStart()
	}
	lib.TraceLevel(*traceLevel)

	setInterrupt() // allows app to be stopped gracefully with ctrl-c

	rand.Seed(time.Now().UnixNano()) // rand used by common.RandCode

	web.WebStart(":http")

	log.Println("pgm stopped unexpectedly")
}

func prodStart() {
	fmt.Println("starting in production mode")
	now := time.Now()
	_, mth, day := now.Date()
	traceFile := fmt.Sprintf("logs/%v%d_%02d%02d%02d.log", mth.String()[0:3], day, now.Hour(), now.Minute(), now.Second())
	lib.TraceStart(traceFile)
	d.DataStart("db/prod.db") // init global memory data vars, load data, start data dispatch goroutine
	go schedule.Scheduler()
}

func testStart() {
	fmt.Println("starting in test mode - test.db database recreated")
	lib.TraceStart("logs/test.log")
	// lib.TraceStart("stdout")
	os.Remove("db/test.db")
	d.DataStart("db/test.db")
}

func setInterrupt() {
	interruptChan := make(chan os.Signal, 1) // stop pgm if ctrl-c entered
	signal.Notify(interruptChan, os.Interrupt)
	waitForInterrupt := func() {
		<-interruptChan
		log.Println("interrupt received, pgm shutting down gracefully")
		d.Data("shutdown", &d.SimpleRequest{nil})
		log.Println("waiting")
		time.Sleep(3 * time.Second)
		os.Exit(0)
	}
	go waitForInterrupt()
}
