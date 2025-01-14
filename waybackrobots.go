package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
	"strings"

	"github.com/fatih/color"
)

type WaybackTimestamps [][]string

func main() {
	var domain string
	var fromDate int
	flag.StringVar(&domain, "domain", "", "which domain to find old robots for")
	flag.IntVar(&fromDate, "fd", 2015, "choose date from when to get robots from format: 2015")
	flag.Parse()

	if domain == "" {
		color.Set(color.FgRed)
		fmt.Fprintln(os.Stderr, "supply a domain, -domain")
		color.Unset()
		flag.PrintDefaults()
		os.Exit(1)
	}

	if len(strings.Split(domain, ".")) < 1 {
		color.Set(color.FgRed)
		fmt.Fprintf(os.Stderr, "ERROR not a valid domain: %s\n", domain)
		color.Unset()
		os.Exit(1)
	}

	var wt WaybackTimestamps

	timestamps := fmt.Sprintf("https://web.archive.org/cdx/search/cdx?url=%s/robots.txt&output=json&fl=timestamp,original&filter=statuscode:200&collapse=digest&from=%d", domain, fromDate)
	resp, err := http.Get(timestamps)
	if err != nil {
		color.Set(color.FgRed)
		fmt.Fprintf(os.Stderr, "ERROR getting timestamps response: %v\n", err)
		color.Unset()
		os.Exit(1)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		color.Set(color.FgRed)
		fmt.Fprintf(os.Stderr, "ERROR reading response body: %v\n", err)
		color.Unset()
		os.Exit(1)
	}

	if resp.StatusCode != 200 {
		color.Set(color.FgRed)
		fmt.Fprintf(os.Stderr, "ERROR getting timestamps, status: %d, body: %v\n", resp.StatusCode, string(body))
		color.Unset()
		os.Exit(1)
	}

	if json.Unmarshal(body, &wt); err != nil {
		color.Set(color.FgRed)
		fmt.Fprintf(os.Stderr, "ERROR parsing json: %v\n", err)
		color.Unset()
		os.Exit(1)
	}

	color.Set(color.FgYellow)
	fmt.Fprintf(os.Stderr, "[i] found %d old robots.txt files", len(wt))
	color.Unset()

	var uniquePaths []string
	for _, e := range wt {
		ts := e[0]
		snapshotUrl := fmt.Sprintf("https://web.archive.org/web/%sid_/%s/robots.txt", ts, domain)
		resp, err := http.Get(snapshotUrl)
		if err != nil {
			color.Set(color.FgRed)
			fmt.Fprintf(os.Stderr, "ERROR getting robots resp: %v\n", err)
			color.Unset()
			os.Exit(1)
		}

		defer resp.Body.Close()

		// fmt.Fprintf(os.Stderr, "Sending: %s, Resp status code: %d\n", snapshotUrl, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			color.Set(color.FgRed)
			fmt.Fprintf(os.Stderr, "ERROR reading response body: %v\n", err)
			color.Unset()
			os.Exit(1)
		}

		textContent := string(body)
		if strings.Contains(textContent, "Disallow:") {
			for _, tc := range strings.Split(textContent, "\n") {
				if strings.Contains(tc, "Disallow:") {
					path := strings.Split(tc, "Disallow:")
					if !slices.Contains(uniquePaths, path[1]) && path[1] != "" {
						uniquePaths = append(uniquePaths, path[1])
						fmt.Println(path[1])
					}
				}
			}
		}

	}
}
