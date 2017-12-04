package main

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"sync"
	"time"
	"ubica/glumac"
	"ubica/short_link"
)

func main() {
	glumac.helloActor()

	var wg sync.WaitGroup
	httpClient := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	all := [16]int{66, 67, 73, 74, 79, 80, 81, 82, 84, 85, 89, 90, 107, 108, 109, 110}
	stats := make([]int64, 0)
	errors := make([]string, 0)

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
					stats = append(stats, d.Nanoseconds())
					errors = append(errors, err.Error())
					fmt.Printf("%s\t%v\n", err, d)
					return
				}
				d := time.Now().Sub(s)
				stats = append(stats, d.Nanoseconds())
				fmt.Printf("%s\t%v\n", resp.Status, d)
			}
		}()
	}
	wg.Wait()
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
	fmt.Printf("Success:\t%v%%\n", len(errors)/len(stats)*100)
	fmt.Println("-----------------------------------")
}

type ByDuration []int64

func (a ByDuration) Len() int           { return len(a) }
func (a ByDuration) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDuration) Less(i, j int) bool { return a[i] < a[j] }
