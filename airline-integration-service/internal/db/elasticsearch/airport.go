package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"aerona.thanhtd.com/airline-integration-service/internal/api/models"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/google/uuid"
)

type AirportRepository struct {
	client *elasticsearch.Client
}

func NewAirportRepository(client *elasticsearch.Client) *AirportRepository {
	return &AirportRepository{client: client}
}

func (r *AirportRepository) CreateAirport(airport models.Airport) (*models.Airport, error) {
	airport.AirportId = uuid.New().String()
	jsonData, err := json.Marshal(airport)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal airport: %v", err)
	}

	res, err := r.client.Index("airports",
		bytes.NewReader(jsonData),
		r.client.Index.WithDocumentID(airport.AirportId),
		r.client.Index.WithContext(context.Background()))
	if err != nil {
		return nil, fmt.Errorf("failed to index airport: %v", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("error response from Elasticsearch: %v", res.String())
	}
	return &airport, nil
}

func (r *AirportRepository) GetAllAirport() ([]models.Airport, error) {
	query := `{
		"query": {
			"match_all": {}	
		},
		"size": 1000
	}`

	res, err := r.client.Search(
		r.client.Search.WithContext(context.Background()),
		r.client.Search.WithIndex("airports"),
		r.client.Search.WithBody(strings.NewReader(query)),
		r.client.Search.WithSize(1000),
		r.client.Search.WithPretty(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to search airports: %v", err)
	}
	defer res.Body.Close()
	if res.IsError() {
		return nil, fmt.Errorf("elasticsearch error: %s", res.String())
	}

	var result models.ESResult[models.Airport]
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	airports := make([]models.Airport, 0, len(result.Hits.Hits))
	for _, hit := range result.Hits.Hits {
		airports = append(airports, hit.Source)
	}
	return airports, nil
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

func (r *AirportRepository) UpdateAirport(airportId string, rawData models.RawAirport) (*models.Airport, error) {
	updatedAirport := models.ParseAirport(rawData)
	now := time.Now()
	updatedAirport.UpdatedAt = &now

	updatedData, err := json.Marshal(updatedAirport)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal updated airport: %v", err)
	}

	script := fmt.Sprintf(`{
        "query": {
            "term": {
                "airport_id": "%s"
            }
        },
        "script": {
            "source": "for (entry in params.updates.entrySet()) { ctx._source[entry.getKey()] = entry.getValue() }",
            "lang": "painless",
            "params": {
                "updates": %s
            }
        }
    }`, airportId, string(updatedData))

	res, err := r.client.UpdateByQuery(
		[]string{"airports"},
		r.client.UpdateByQuery.WithContext(context.Background()),
		r.client.UpdateByQuery.WithBody(strings.NewReader(script)))
	if err != nil {
		return nil, fmt.Errorf("failed to update airport by airport_id: %v", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("elasticsearch error: %s", res.String())
	}

	var result struct {
		Updated int64 `json:"updated"`
	}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode update response: %v", err)
	}
	if result.Updated == 0 {
		return nil, fmt.Errorf("no airport found with airport_id: %s", airportId)
	}
	return &updatedAirport, nil
}

func (r *AirportRepository) DeleteByAirportId(airportId string) error {
	query := fmt.Sprintf(`{
		"query": {
			"match": {
				"airport_id": "%s"
			}
		}
	}`, airportId)

	res, err := r.client.DeleteByQuery(
		[]string{"airports"},
		strings.NewReader(query),
		r.client.DeleteByQuery.WithContext(context.Background()))

	if err != nil {
		return fmt.Errorf("failed to delete airport by airport_id: %v", airportId)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("elasticsearch error: %v", res.String())
	}

	return nil
}
