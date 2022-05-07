package main

import (
	"fmt"
	"strconv"
	"time"

	"simpleMonitor/mon"
)

func measureCSV() {
	fmt.Println("monitor start")
	filename := "./" + strconv.FormatInt(time.Now().UnixNano(), 10) + ".csv"

	stat, _ := mon.NewOneStat()
	s, _ := stat.CSVHead()
	mon.WriteLine(filename, []string{s})
	for true {
		stat, _ := mon.NewOneStat()
		s, _ := stat.CSV()
		mon.WriteLine(filename, []string{s})
	}
}

func main() {
	measureCSV()
}
