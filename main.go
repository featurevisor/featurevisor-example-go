package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/featurevisor/featurevisor-go"
)

func main() {
	/**
	 * Fetch datafile
	 */
	datafileURL := "https://featurevisor-example-cloudflare.pages.dev/production/featurevisor-tag-all.json"

	resp, err := http.Get(datafileURL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	datafileBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var datafileContent featurevisor.DatafileContent
	if err := datafileContent.FromJSON(string(datafileBytes)); err != nil {
		panic(err)
	}

	/**
	 * Create Featurevisor instance
	 */
	f := featurevisor.CreateFeaturevisor(featurevisor.FeaturevisorOptions{
		Datafile: datafileContent,
	})
	f.SetContext(featurevisor.Context{
		"userId":   "123",
		"deviceId": "device-23456",
		"country":  "nl",
	})

	/**
	 * Evaluate values
	 */
	featureIsEnabled := f.IsEnabled("my_feature")
	fmt.Println("Feature is enabled:", featureIsEnabled)
}
