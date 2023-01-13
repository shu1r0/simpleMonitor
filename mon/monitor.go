package mon

import (
	"encoding/json"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gocarina/gocsv"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

// one interval stats
type OneStat struct {
	// Unix time
	Time int64

	// cpu
	GlobalCPUPercent float64
	CPUPercent       []float64 `csv:"CPUPercent"`

	// memory
	TotalMemory   uint64
	TotalMemoryMB float64 `json:"TotalMemory(MB)" csv:"TotalMemory(MB)"`
	UsedMemory    uint64
	UsedMemoryMB  float64 `json:"UsedMemory(MB)" csv:"UsedMemory(MB)"`
	UsedPercent   float64
}

// monitor
func NewOneStat() (oneStat *OneStat, err error) {
	s := OneStat{}
	interval := time.Microsecond * 500 // 0.5s

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
	stat_csv := slice[1]

	var new_stat []string
	for i, s := range strings.Split(stat_csv, ",") {
		if i == 2 {
			for _, s := range stat.CPUPercent {
				new_stat = append(new_stat, strconv.FormatFloat(s, 'f', -1, 64))
			}
		} else {
			new_stat = append(new_stat, s)
		}
	}

	return strings.Join(new_stat, ","), err
}

func (stat OneStat) CSVHead() (s string, err error) {
	stats := []*OneStat{&stat}
	str, err := gocsv.MarshalString(stats)
	slice := strings.Split(str, "\n")
	head_csv := slice[0]

	var new_head []string
	for _, s := range strings.Split(head_csv, ",") {
		if s == "CPUPercent" {
			for j, _ := range stat.CPUPercent {
				h := s + "(" + strconv.Itoa(j+1) + ")"
				new_head = append(new_head, h)
			}
		} else {
			new_head = append(new_head, s)
		}
	}

	return strings.Join(new_head, ","), err
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
