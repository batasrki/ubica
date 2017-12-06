package ubica

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"net/http"
	"net/url"
)

type Record struct {
	Url       string
	Timestamp int
	Client    string
	Title     string
}

func shortenAllLinks() {
	var record Record
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	all, err := client.ZRange("posted:urls", 0, 1).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(all)

	for _, rec := range all {
		json.Unmarshal([]byte(rec), &record)
		fmt.Println(record.Url)
		addShortLink(record.Url, client)

	}
}

func addShortLink(canonicalUrl string, client *redis.Client) {
	vals := url.Values{}
	vals.Set("url", canonicalUrl)
	vals.Add("csrfmiddlewaretoken", "rD2g5SaXUQ5ZiYnjovgoB6im0PTzYUzSCBUuFhduQjDz4Bb6f2OTOzFWVl3bGjf0")
	resp, err := http.PostForm("http://138.197.170.57/GenerateNewShortLink", vals)

	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}
