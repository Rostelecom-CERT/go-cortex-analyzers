package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	cortex "github.com/ilyaglow/go-cortex"
)

const (
	apiv2Breaches = "https://haveibeenpwned.com/api/v2/breachedaccount/"
	apiv2Pastes   = "https://haveibeenpwned.com/api/v2/pasteaccount/"
)

type report struct {
	Results []result `json:"results"`
}

type result []interface{}

func main() {
	i, err := cortex.NewInput()
	if err != nil {
		log.Fatal(err)
	}

	br, btxs, err := getBreaches(i.Data)
	if err != nil {
		cortex.SayError(i, err.Error())
	}

	pr, ptxs, err := getPastes(i.Data)
	if err != nil {
		cortex.SayError(i, err.Error())
	}

	r := report{}
	if br != nil {
		r.Results = append(r.Results, *br)
	}

	if pr != nil {
		r.Results = append(r.Results, *pr)
	}
	cortex.SayReport(r, append(btxs, ptxs...))
}

func getBreaches(acc string) (*result, []cortex.Taxonomy, error) {
	resp, err := http.Get(apiv2Breaches + acc)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	var txs []cortex.Taxonomy
	predicate := "Breaches"
	namespace := "HaveIBeenPwned"
	switch resp.StatusCode {
	case 404:
		txs = append(txs, cortex.Taxonomy{
			Namespace: namespace,
			Predicate: predicate,
			Level:     cortex.TxSafe,
			Value:     "0",
		})
		return nil, txs, nil
	case 200:
		var r result
		dec := json.NewDecoder(resp.Body)
		if err := dec.Decode(&r); err != nil {
			return nil, nil, err
		}

		txs = append(txs, cortex.Taxonomy{
			Namespace: namespace,
			Predicate: predicate,
			Level:     cortex.TxSuspicious,
			Value:     strconv.FormatInt(int64(len(r)), 10),
		})

		var vf int
		for i := range r {
			m := r[i].(map[string]interface{})
			if m["IsVerified"].(bool) {
				vf++
			}
		}

		predicate = "Verified"
		txs = append(txs, cortex.Taxonomy{
			Namespace: namespace,
			Predicate: predicate,
			Level:     cortex.TxSuspicious,
			Value:     strconv.FormatInt(int64(vf), 10),
		})

		return &r, txs, nil
	default:
		return nil, nil, fmt.Errorf("Unexpected status code from haveibeenpwned.com %s", resp.Status)
	}
}

func getPastes(acc string) (*result, []cortex.Taxonomy, error) {
	resp, err := http.Get(apiv2Pastes + acc)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	var txs []cortex.Taxonomy
	predicate := "Pastes"
	namespace := "HaveIBeenPwned"
	switch resp.StatusCode {
	case 404:
		txs = append(txs, cortex.Taxonomy{
			Namespace: namespace,
			Predicate: predicate,
			Level:     cortex.TxSafe,
			Value:     "0",
		})
		return nil, txs, nil
	case 200:
		var r result
		dec := json.NewDecoder(resp.Body)
		if err := dec.Decode(&r); err != nil {
			return nil, nil, err
		}

		txs = append(txs, cortex.Taxonomy{
			Namespace: namespace,
			Predicate: predicate,
			Level:     cortex.TxSuspicious,
			Value:     strconv.FormatInt(int64(len(r)), 10),
		})

		return &r, txs, nil
	default:
		return nil, nil, fmt.Errorf("Unexpected status code from haveibeenpwned.com %s", resp.Status)
	}
}
