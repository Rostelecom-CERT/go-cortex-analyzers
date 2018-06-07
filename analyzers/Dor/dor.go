package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/ilyaglow/dor"
	"gopkg.ilya.app/ilyaglow/go-cortex.v1"
)

type dorResponse struct {
	RequestData string             `json:"data"`
	Hits        []dor.ExtendedRank `json:"ranks"`
	Timestamp   time.Time          `json:"timestamp"`
}

func main() {
	// Grab stdin to JobInput structure
	input, err := cortex.NewInput()
	if err != nil {
		log.Fatal(err)
	}

	// Get url parameter from analyzer config
	url, err := input.Config.GetString("url")
	if err != nil {
		// Report an error if something went wrong
		cortex.SayError(input, err.Error())
	}

	// You get somehow report struct from JobInput.Data
	rep, err := do(input.Data, url)
	if err != nil {
		cortex.SayError(input, err.Error())
	}

	// Make taxonomies
	var txs []cortex.Taxonomy
	namespace := "DomainRank"
	if len(rep.Hits) != 0 {
		for i := range rep.Hits {
			txs = append(txs, cortex.Taxonomy{
				Namespace: namespace,
				Predicate: rep.Hits[i].Source,
				Level:     cortex.TxInfo,
				Value:     strconv.FormatInt(int64(rep.Hits[i].GetRank()), 10),
			})
		}
	} else {
		txs = append(txs, cortex.Taxonomy{
			Namespace: namespace,
			Predicate: "Report",
			Level:     cortex.TxSuspicious,
			Value:     "unranked",
		})
	}

	// Report accept marshallable struct and taxonomies
	cortex.SayReport(rep, txs)
}

func do(domain string, url string) (*dorResponse, error) {
	resp, err := http.Get(url + "/rank/" + domain)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var f dorResponse

	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&f); err != nil {
		return nil, err
	}

	return &f, nil
}
