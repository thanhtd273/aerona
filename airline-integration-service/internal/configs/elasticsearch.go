package configs

import (
	"fmt"
	"os"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
)

type ElasticClient struct {
	*elasticsearch.Client
}

func NewElasticClient() (*elasticsearch.Client, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{
			os.Getenv("ELASTICSEARCH_URL"), // e.g., "http://localhost:9200"
		},
		// Username:             os.Getenv("ELASTICSEARCH_USERNAME"), // Optional
		// Password:             os.Getenv("ELASTICSEARCH_PASSWORD"), // Optional
		RetryOnStatus:        []int{502, 503, 504, 429}, // Retry on these HTTP statuses
		RetryBackoff:         func(i int) time.Duration { return time.Duration(i) * 100 * time.Millisecond },
		MaxRetries:           5,    // Max retry attempts
		EnableRetryOnTimeout: true, // Retry on timeout
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create Elasticsearch client: %v", err)
	}

	// Test connection
	// ctx := context.Background()
	res, err := client.Info()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Elasticsearch: %v", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("elasticsearch returned an error: %s", res.String())
	}

	return client, nil
}
