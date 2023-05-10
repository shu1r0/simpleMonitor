package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"time"

	mon "github.com/shu1r0/simpleMonitor/pkg/monitor"
)

func measureCSV(filename string, printStat bool) {
	fmt.Println("Monitor start ...")

	stat, _ := mon.NewOneStat()
	s, _ := stat.CSVHead()

	// write header
	if err := mon.WriteLine(filename, []string{s}); err != nil {
		fmt.Errorf("Write file Error: %s\n", err)
		os.Exit(1)
	}

	// write stat
	for true {
		stat, _ := mon.NewOneStat()
		if printStat {
			jsonstat, err := stat.JSON()
			if err != nil {
				fmt.Errorf("%s", err)
			}
			fmt.Println(jsonstat)
		}
		s, _ := stat.CSV()
		if err := mon.WriteLine(filename, []string{s}); err != nil {
			fmt.Errorf("Write file Error: %s\n", err)
		}
	}
}

func main() {
	var (
		printStat = flag.Bool("out", false, "stat output")
		filename  = flag.String("file", "", "result csv file path")
	)
	flag.Parse()

	if *filename == "" {
		*filename = "./stats-" + strconv.FormatInt(time.Now().UnixNano(), 10) + ".csv"
	}

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	go measureCSV(*filename, *printStat)

	<-quit
	fmt.Println("Monitor Stopped.")
}
