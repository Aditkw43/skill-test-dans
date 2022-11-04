package internal

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
)

var (
	page  int = 1
	limit int = 1
)

func GetData() []*Job {
	url := os.Getenv("API_DATA")
	resp, err := http.Get(url)

	if err != nil {
		log.Fatal("Failed get data from url: ", err)
	}

	defer resp.Body.Close()

	var dataJob []*Job
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	respByte := buf.Bytes()
	if err := json.Unmarshal(respByte, &dataJob); err != nil {
		log.Fatal("unmarshal json error:", err)
	}

	return dataJob
}

func Paginate(data []*Job, page int, limit int) []*Job {
	start := (page - 1) * limit
	stop := start + limit

	if start > len(data) {
		return nil
	}

	if stop > len(data) {
		stop = len(data)
	}

	return data[start:stop]
}
