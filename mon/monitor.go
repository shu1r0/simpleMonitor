package mon

import (
	"encoding/json"
	"math"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gocarina/gocsv"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

type OneStat struct {
	// Unix time
	Time int64

	// cpu
	GlobalCPUPercent float64
	CPUPercent       []float64 `csv:"CPUPercent(not support)"`

	// memory
	TotalMemory   uint64
	TotalMemoryMB float64 `json:"TotalMemory(MB)" csv:"TotalMemory(MB)"`
	UsedMemory    uint64
	UsedMemoryMB  float64 `json:"UsedMemory(MB)" csv:"UsedMemory(MB)"`
	UsedPercent   float64
}

func NewOneStat() (oneStat *OneStat, err error) {
	s := OneStat{}
	interval := time.Second // 1s

	// timestamp
	s.Time = time.Now().Unix()

	// cpu
	var percentWaitGroup sync.WaitGroup
	percentWaitGroup.Add(2)

	go func(stat *OneStat, wg *sync.WaitGroup) {
		defer wg.Done()
		p, _ := cpu.Percent(interval, true)
		stat.CPUPercent = p
	}(&s, &percentWaitGroup)

	go func(stat *OneStat, wg *sync.WaitGroup) {
		defer wg.Done()
		p, _ := cpu.Percent(interval, false)
		stat.GlobalCPUPercent = p[0]
	}(&s, &percentWaitGroup)

	percentWaitGroup.Wait()

	// memory
	m, _ := mem.VirtualMemory()
	s.TotalMemory = m.Total
	s.TotalMemoryMB = float64(s.TotalMemory) * math.Pow10(-6)
	s.UsedMemory = m.Used
	s.UsedMemoryMB = float64(s.UsedMemory) * math.Pow10(-6)
	s.UsedPercent = m.UsedPercent

	return &s, nil
}

func (stat OneStat) JSON() (s string, err error) {
	b, err := json.Marshal(stat)
	return string(b), err
}

func (stat OneStat) CSV() (s string, err error) {
	stats := []*OneStat{&stat}
	str, err := gocsv.MarshalString(stats)
	slice := strings.Split(str, "\n")
	return slice[1], err
}

func (stat OneStat) CSVHead() (s string, err error) {
	stats := []*OneStat{&stat}
	str, err := gocsv.MarshalString(stats)
	slice := strings.Split(str, "\n")
	return slice[0], err
}

func WriteLine(filename string, lines []string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)

	if err != nil {
		return err
	}

	defer file.Close()

	for _, line := range lines {
		_, err := file.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}
