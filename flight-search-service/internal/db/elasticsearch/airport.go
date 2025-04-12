package elasticsearch

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"aerona.thanhtd.com/flight-search-service/internal/api/models"
	"github.com/elastic/go-elasticsearch/v7"
)

type AirportRepository struct {
	client *elasticsearch.Client
}

func NewAirportRepository(client *elasticsearch.Client) *AirportRepository {
	return &AirportRepository{
		client: client,
	}
}

func (r *AirportRepository) GetPopularAirports() ([]models.Airport, error) {
	const LIMIT int = 3

	query := fmt.Sprintf(`{
	"sort": [
		{
		"popularity": {
			"order": "desc"
		}
		}
	],
	"size": %d
	}`, LIMIT)

	res, err := r.client.Search(
		r.client.Search.WithContext(context.Background()),
		r.client.Search.WithIndex("airports"),
		r.client.Search.WithBody(strings.NewReader(query)),
		r.client.Search.WithSize(3),
		r.client.Search.WithPretty(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get 3 most popular airports, error: %v", err)
	}
	defer res.Body.Close()
	if res.IsError() {
		return nil, fmt.Errorf("elasticsearch response had error: %v", res.String())
	}

	var result models.ESResult[models.Airport]
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode Elasticsearch response, error: %v", err)
	}

	airports := make([]models.Airport, 0, len(result.Hits.Hits))
	for _, hit := range result.Hits.Hits {
		airports = append(airports, hit.Source)
	}
	return airports, nil

}

func (r *AirportRepository) IncreasePopularity(airportId string) error {
	query := fmt.Sprintf(`{
		"script": {
			"source": "ctx._source.popularity = ctx._source.popularity + 1; ctx._source.updated_at = params.now",
			"lang": "painless",
			"params": {
				"now": "%v"
			}
		},
		"query": {
			"match": {
				"airport_id": "%s"
			}
		}
	}`, time.Now().UTC().Format(time.RFC3339), airportId)

	res, err := r.client.UpdateByQuery(
		[]string{"airports"},
		r.client.UpdateByQuery.WithContext(context.Background()),
		r.client.UpdateByQuery.WithBody(strings.NewReader(query)),
	)
	if err != nil {
		return fmt.Errorf("failed to update popularity of airport, airportId: %s, error: %v", airportId, err)
	}
	defer res.Body.Close()
	if res.IsError() {
		return fmt.Errorf("elasticsearch response had error: %v", err)
	}
	return nil
}

func (r *AirportRepository) FindByAirportId(airportId string) (*models.Airport, error) {
	query := fmt.Sprintf(`{
		"query": {
			"match": {
				"airport_id": "%s"
			}
		},
		"size": 1
	}`, airportId)

	res, err := r.client.Search(
		r.client.Search.WithContext(context.Background()),
		r.client.Search.WithIndex("airports"),
		r.client.Search.WithBody(strings.NewReader(query)),
		r.client.Search.WithSize(1),
		r.client.Search.WithPretty(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch airport by id=%v, error: %v", airportId, err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("response have error: %v", res.String())
	}

	var result models.ESResult[models.Airport]
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	if len(result.Hits.Hits) == 0 {
		return nil, fmt.Errorf("airport with airport_id=%s not found", airportId)
	}

	return &result.Hits.Hits[0].Source, nil

}

func (r *AirportRepository) FindByAirportCode(airportCode string) (*models.Airport, error) {
	query := fmt.Sprintf(`{
		"query": {
			"match": {
				"code": "%s"
			}
		},
		"size": 1
	}`, airportCode)

	res, err := r.client.Search(
		r.client.Search.WithContext(context.Background()),
		r.client.Search.WithIndex("airports"),
		r.client.Search.WithBody(strings.NewReader(query)),
		r.client.Search.WithSize(1),
		r.client.Search.WithPretty(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch airport by code=%v, error: %v", airportCode, err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("response have error: %v", res.String())
	}

	var result models.ESResult[models.Airport]
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	if len(result.Hits.Hits) == 0 {
		return nil, fmt.Errorf("airport with code=%s not found", airportCode)
	}

	return &result.Hits.Hits[0].Source, nil

}
