package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	cortex "github.com/ilyaglow/go-cortex"
)

const (
	apiURL    = "https://hacked-emails.com/api?q="
	namespace = "HackedEmails"
	predicate = "ResultsCount"
)

type reply struct {
	Status  string      `json:"status"`
	Query   string      `json:"query"`
	Results int         `json:"results"`
	Data    []replyData `json:"data"`
}

type replyData struct {
	SourceID       string    `json:"source_id"`
	SourceURL      string    `json:"source_url"`
	SourceLines    int       `json:"source_lines"`
	SourceSize     int       `json:"source_size"`
	SourceNetwork  string    `json:"source_network"`
	SourceProvider string    `json:"source_provider"`
	Verified       bool      `json:"verified"`
	Title          string    `json:"title"`
	Author         string    `json:"author"`
	DateCreated    time.Time `json:"date_created"`
	DateLeaked     time.Time `json:"date_leaked"`
	EmailsCount    int       `json:"emails_count"`
	Details        string    `json:"details"`
}

func main() {
	i, err := cortex.NewInput()
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.Get(apiURL + i.Data)
	if err != nil {
		cortex.SayError(i, err.Error())
	}
	defer resp.Body.Close()

	var r reply
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&r); err != nil {
		cortex.SayError(i, err.Error())
	}

	var txs []cortex.Taxonomy
	if r.Results == 0 {
		txs = append(txs, cortex.Taxonomy{
			Namespace: namespace,
			Predicate: predicate,
			Level:     cortex.TxSafe,
			Value:     strconv.FormatInt(int64(r.Results), 10),
		})
	} else {
		txs = append(txs, cortex.Taxonomy{
			Namespace: namespace,
			Predicate: predicate,
			Level:     cortex.TxSuspicious,
			Value:     strconv.FormatInt(int64(r.Results), 10),
		})
	}

	cortex.SayReport(r, txs)
}
