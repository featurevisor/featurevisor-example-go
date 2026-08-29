package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/featurevisor/featurevisor-go/v3"
)

const datafileURL = "https://featurevisor-example-cloudflare.pages.dev/production/featurevisor-sdk-v3.json"

type serviceEndpoints struct {
	BaseURL   string `json:"baseUrl"`
	TimeoutMS int    `json:"timeoutMs"`
	Retries   int    `json:"retries"`
}

func main() {
	client := &http.Client{Timeout: 10 * time.Second}
	response, err := client.Get(datafileURL)
	if err != nil {
		panic(fmt.Errorf("fetch datafile: %w", err))
	}
	defer response.Body.Close()

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		panic(fmt.Errorf("fetch datafile: unexpected HTTP status %s", response.Status))
	}

	datafileJSON, err := io.ReadAll(response.Body)
	if err != nil {
		panic(fmt.Errorf("read datafile: %w", err))
	}

	var datafile featurevisor.DatafileContent
	if err := datafile.FromJSON(string(datafileJSON)); err != nil {
		panic(fmt.Errorf("parse datafile: %w", err))
	}

	logLevel := featurevisor.LogLevelError
	f := featurevisor.CreateFeaturevisor(featurevisor.FeaturevisorOptions{
		Datafile: datafile,
		LogLevel: &logLevel,
		Context: featurevisor.Context{
			"userId":      "customer-123",
			"country":     "nl",
			"locale":      "nl-NL",
			"accountPlan": "pro",
		},
	})
	defer f.Close()

	commerceEnabled := f.IsEnabled("commerce_platform")
	checkoutVariation := f.GetVariation("checkout_experience")
	maxItems := f.GetVariableInteger("checkout_experience", "max_items")
	paymentMethods := f.GetVariableArray("checkout_experience", "payment_methods")
	supportContact := f.GetGlobalVariableString("supportContact")

	var endpoints serviceEndpoints
	if err := f.GetGlobalVariableObjectInto("serviceEndpoints", &endpoints); err != nil {
		panic(fmt.Errorf("evaluate service endpoints: %w", err))
	}

	fmt.Println("Commerce platform enabled:", commerceEnabled)
	fmt.Println("Checkout variation:", valueOrUnavailable(checkoutVariation))
	fmt.Println("Maximum checkout items:", valueOrUnavailable(maxItems))
	fmt.Println("Payment methods:", paymentMethods)
	fmt.Printf("Service endpoint: %s (timeout: %d ms, retries: %d)\n", endpoints.BaseURL, endpoints.TimeoutMS, endpoints.Retries)
	fmt.Println("Support contact:", valueOrUnavailable(supportContact))
}

func valueOrUnavailable[T any](value *T) any {
	if value == nil {
		return "unavailable"
	}
	return *value
}
