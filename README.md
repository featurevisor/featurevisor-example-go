# Featurevisor Go example

A small application showing how to use the [Featurevisor Go SDK](https://github.com/featurevisor/featurevisor-go).

It fetches a production datafile, creates a Featurevisor instance with customer context, then evaluates:

- the `commerce_platform` feature as a flag;
- the `checkout_experience` variation;
- integer and array variables from `checkout_experience`;
- the `serviceEndpoints` global object variable into a Go struct;
- the `supportContact` global string variable.

The datafile comes from the [Featurevisor Cloudflare example](https://github.com/featurevisor/featurevisor-example-cloudflare).

## Requirements

- Go 1.21 or newer
- internet access for downloading the SDK and example datafile

## Run the example

Install the dependencies:

```bash
go mod download
```

Run the application:

```bash
go run .
```

Expected output:

```text
Commerce platform enabled: true
Checkout variation: express
Maximum checkout items: 25
Payment methods: [card wallet]
Service endpoint: https://api.eu.example.com (timeout: 1200 ms, retries: 4)
Support contact: support-nl@example.com
```

The evaluations use this context:

```go
featurevisor.Context{
    "userId":      "customer-123",
    "country":     "nl",
    "locale":      "nl-NL",
    "accountPlan": "pro",
}
```

Change these values to see how Featurevisor selects different rules, variations, and global variable overrides.

Learn more in the [Featurevisor documentation](https://featurevisor.com/docs/sdks/go/).
