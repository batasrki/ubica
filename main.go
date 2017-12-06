package main

import (
	"fmt"
	// "github.com/batasrki/ubica/glumac"
	// "github.com/batasrki/ubica/short_link"
	"net/http"
	"sort"
	"strconv"
	"sync"
	"time"
)

func main() {
	// glumac.HelloActor()

	var wg sync.WaitGroup
	httpClient := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	all := [16]int{66, 67, 73, 74, 79, 80, 81, 82, 84, 85, 89, 90, 107, 108, 109, 110}
	statsChan := make(chan int64, 16000)
	errorsChan := make(chan string, 16000)

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, rec := range all {
				target := "http://138.197.170.57/"
				target += strconv.Itoa(rec)
				s := time.Now()
				resp, err := httpClient.Get(target)
				if err != nil {
					// panic(err)
					d := time.Now().Sub(s)
					statsChan <- d.Nanoseconds()
					errorsChan <- err.Error()
					fmt.Printf("%s\t%v\n", err, d)
					return
				}
				d := time.Now().Sub(s)
				statsChan <- d.Nanoseconds()
				fmt.Printf("%s\t%v\n", resp.Status, d)
			}
		}()
	}
	wg.Wait()
	stats := make([]int64, 0)
	errors := make([]string, 0)

	stats = collectFromStatsChan(statsChan, stats)
	errors = collectFromErrorsChan(errorsChan, errors)
	sort.Sort(ByDuration(stats))
	var sum int64

	for i := 0; i < len(stats)-1; i++ {
		sum += stats[i]
	}
	avg := sum / int64(len(stats))

	fmt.Println("-----------------------------------")
	fmt.Printf("Min:\t%v\n", time.Duration(stats[0]))
	fmt.Printf("Max:\t%v\n", time.Duration(stats[len(stats)-1]))
	fmt.Printf("Avg:\t%v\n", time.Duration(avg))
	fmt.Printf("Med:\t%v\n", time.Duration(stats[len(stats)/2]))
	fmt.Printf("95th:\t%v\n", time.Duration(stats[len(stats)*95/100]))
	fmt.Printf("99th:\t%v\n", time.Duration(stats[len(stats)*99/100]))
	fmt.Printf("Total:\t%v\n", len(stats))
	fmt.Printf("Errors:\t%v\n", len(errors))
	fmt.Printf("Success: %v%%\n", len(stats)*100/16000)
	fmt.Println("-----------------------------------")
}

func collectFromStatsChan(statsChan chan int64, arr []int64) []int64 {
	for {
		select {
		case d := <-statsChan:
			arr = append(arr, d)
		default:
			return arr
		}
	}
}

func collectFromErrorsChan(errorsChan chan string, arr []string) []string {
	for {
		select {
		case e := <-errorsChan:
			arr = append(arr, e)
		default:
			return arr
		}
	}
}

type ByDuration []int64

func (a ByDuration) Len() int           { return len(a) }
func (a ByDuration) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDuration) Less(i, j int) bool { return a[i] < a[j] }
