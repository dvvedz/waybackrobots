package main

import (
	"bufio"
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

func getRobotsContent(domain, ts string) []string {
	// fmt.Fprintf(os.Stderr, "[DEBUG] downloading %v\n", ts)

	snapshotUrl := fmt.Sprintf("https://web.archive.org/web/%sid_/%s/robots.txt", ts, domain)
	resp, err := http.Get(snapshotUrl)
	if err != nil {
		defer resp.Body.Close()
		color.Set(color.FgRed)
		fmt.Fprintf(os.Stderr, "ERROR getting robots resp: %v\n", err)
		color.Unset()
	} else {

		defer resp.Body.Close()
		// fmt.Fprintf(os.Stderr, "Sending: %s, Resp status code: %d\n", snapshotUrl, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			color.Set(color.FgRed)
			fmt.Fprintf(os.Stderr, "ERROR reading response body: %v\n", err)
			color.Unset()
		}

		if resp.StatusCode != 200 {
			fmt.Fprintf(os.Stderr, "Error response code was not 200, got %d, boyd: %s\n", resp.StatusCode, string(body))
			return nil
		}

		textContent := string(body)
		var paths []string
		if strings.Contains(textContent, "Disallow:") {
			for _, tc := range strings.Split(textContent, "\n") {
				if strings.Contains(tc, "Disallow:") {
					path := strings.Split(tc, "Disallow:")
					paths = append(paths, path[1])
				}
			}
		}
		return paths
	}
	return nil
}

func getTimestamps(domain, strategy string, fromDate int) {

	var wt WaybackTimestamps
	var timestamps string
	if strategy == "digest" {
		timestamps = fmt.Sprintf("https://web.archive.org/cdx/search/cdx?url=%s/robots.txt&output=json&fl=timestamp,original&filter=statuscode:200&collapse=digest&from=%d", domain, fromDate)
	} else if strategy == "day" {
		timestamps = fmt.Sprintf("https://web.archive.org/cdx/search/cdx?url=%s/robots.txt&output=json&fl=timestamp,original&filter=statuscode:200&collapse=timestamp:7&from=%d", domain, fromDate)
	} else if strategy == "month" {
		timestamps = fmt.Sprintf("https://web.archive.org/cdx/search/cdx?url=%s/robots.txt&output=json&fl=timestamp,original&filter=statuscode:200&collapse=timestamp:6&from=%d", domain, fromDate)
	} else {
		return
	}
	resp, err := http.Get(timestamps)
	if err != nil {
		color.Set(color.FgRed)
		fmt.Fprintf(os.Stderr, "ERROR getting timestamps response: %v\n", err)
		color.Unset()
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		color.Set(color.FgRed)
		fmt.Fprintf(os.Stderr, "ERROR reading response body: %v\n", err)
		color.Unset()
	}

	if resp.StatusCode != 200 {
		color.Set(color.FgRed)
		fmt.Fprintf(os.Stderr, "ERROR getting timestamps, status: %d, body: %v\n", resp.StatusCode, string(body))
		color.Unset()
	}

	if json.Unmarshal(body, &wt); err != nil {
		color.Set(color.FgRed)
		fmt.Fprintf(os.Stderr, "ERROR parsing json: %v\n", err)
		color.Unset()
	}

	color.Set(color.FgYellow)
	fmt.Fprintf(os.Stderr, "[i] found %d old robots.txt files\n", len(wt))
	color.Unset()

	var uniquePaths []string
	sem := make(chan struct{}, 3)
	for _, e := range wt {
		sem <- struct{}{}
		go func() {
			ts := e[0]
			urls := getRobotsContent(domain, ts)
			for _, url := range urls {
				if !slices.Contains(uniquePaths, url) && url != "" {
					fmt.Println(url)
					uniquePaths = append(uniquePaths, url)
				}
			}
			<-sem
		}()
	}
}

func main() {

	var domain string
	var fromDate int
	var strategy string
	strategyValues := []string{"digest", "day", "month"}

	flag.StringVar(&domain, "domain", "", "which domain to find old robots for")
	flag.StringVar(&strategy, "strat", "digest", "interval to get robots for, possible values: digest, day, month")
	flag.IntVar(&fromDate, "fd", 2015, "choose date from when to get robots from format: 2015")

	flag.Parse()

	if !slices.Contains(strategyValues, strategy) {
		color.Set(color.FgRed)
		fmt.Fprintf(os.Stderr, "%s is not a valid value for strat flag\n", strategy)
		color.Unset()
		flag.PrintDefaults()
		os.Exit(1)
	}

	if domain != "" {
		if len(strings.Split(domain, ".")) < 1 {
			fmt.Fprintf(os.Stderr, "ERROR not a valid domain: %s\n", domain)
			os.Exit(1)
		}
		getTimestamps(domain, strategy, fromDate)
	} else {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			domain := scanner.Text()
			if len(strings.Split(domain, ".")) < 1 {
				fmt.Fprintf(os.Stderr, "Skipping %s no a valid domain\n", domain)
			} else {
				getTimestamps(domain, strategy, fromDate)
			}
		}
	}

}
