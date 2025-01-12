package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
	"strings"
)

type WaybackTimestamps [][]string

func main() {
	host := "www.toyota-europe.com"
	var wt WaybackTimestamps

	timestamps := fmt.Sprintf("https://web.archive.org/cdx/search/cdx?url=%s/robots.txt&output=json&fl=timestamp,original&filter=statuscode:200&collapse=digest&from=2012", host)
	resp, err := http.Get(timestamps)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR getting timestamps response: %v\n", err)
		os.Exit(1)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR reading response body: %v\n", err)
		os.Exit(1)
	}

	if resp.StatusCode != 200 {
		fmt.Fprintf(os.Stderr, "ERROR getting timestamps, status: %d, body: %v\n", resp.StatusCode, string(body))
		os.Exit(1)
	}

	if json.Unmarshal(body, &wt); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR parsing json: %v\n", err)
		os.Exit(1)
	}
	var uniquePaths []string
	for _, e := range wt {
		ts := e[0]
		snapshotUrl := fmt.Sprintf("https://web.archive.org/web/%sid_/%s/robots.txt", ts, host)
		resp, err := http.Get(snapshotUrl)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR getting robots resp: %v\n", err)
			os.Exit(1)
		}

		defer resp.Body.Close()

		// fmt.Fprintf(os.Stderr, "Sending: %s, Resp status code: %d\n", snapshotUrl, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR reading response body: %v\n", err)
			os.Exit(1)
		}

		textContent := string(body)
		if strings.Contains(textContent, "Disallow:") {
			for _, tc := range strings.Split(textContent, "\n") {
				if strings.Contains(tc, "Disallow:") {
					path := strings.Split(tc, "Disallow:")
					if slices.Contains(uniquePaths, path[1]) && path[1] != " " {
						uniquePaths = append(uniquePaths, path[1])
						fmt.Println(path[1])
					}
				}
			}
		}

	}
}
