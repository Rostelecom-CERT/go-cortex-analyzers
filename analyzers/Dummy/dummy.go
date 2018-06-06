package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

type report struct {
	Body interface{} `json:"body"`
}

type empty struct {
	Field interface{} `json:"field,omitempty"`
}

type cortexReport struct {
	FullReport report   `json:"full"`
	Success    bool     `json:"success"`
	Summary    empty    `json:"summary"`
	Artifacts  []string `json:"artifacts"`
}

func main() {
	input, err := parseInput(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	printReport(input)
}

func parseInput(f io.Reader) ([]byte, error) {
	in, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return in, nil
}

func printReport(body []byte) {
	r := &cortexReport{
		FullReport: report{string(body)},
		Success:    true,
		Artifacts:  []string{},
	}
	b, err := json.Marshal(r)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(string(b))
	os.Exit(0)
}
