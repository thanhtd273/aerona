package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"aerona.thanhtd.com/flight-search-service/internal/api/models"
	"aerona.thanhtd.com/flight-search-service/internal/utils"
	"github.com/elastic/go-elasticsearch/v7"
)

type AirlineRepository struct {
	client *elasticsearch.Client
}

func NewAirlineRepository(client *elasticsearch.Client) *AirlineRepository {
	return &AirlineRepository{client: client}
}

func (r *AirlineRepository) CreateAirline(airline models.Airline) (*models.Airline, error) {
	// Assign ID for airline
	airlineId, err := utils.GenerateUniqueId()
	if err != nil {
		return nil, fmt.Errorf("failed to generate Snowflake ID, error: %v", err)
	}
	airline.AirlineId = airlineId

	jsonData, err := json.MarshalIndent(airline, "", " ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal airline: %v", err)
	}
	res, err := r.client.Index("airlines",
		bytes.NewReader(jsonData),
		r.client.Index.WithDocumentID(airline.AirlineId),
		r.client.Index.WithContext(context.Background()))
	if err != nil {
		return nil, fmt.Errorf("failed to index airline: %v", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("error response from Elasticsearch: %v", res.String())
	}
	return &airline, nil
}

func (r *AirlineRepository) GetAllAirline() ([]models.Airline, error) {
	query := `{
        "query": {
            "match_all": {}
        },
        "size": 1000
    }`

	res, err := r.client.Search(
		r.client.Search.WithContext(context.Background()),
		r.client.Search.WithIndex("airlines"),
		r.client.Search.WithBody(strings.NewReader(query)),
		r.client.Search.WithSize(1000),
		r.client.Search.WithPretty(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to search airlines: %v", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("elasticsearch error: %s", res.String())
	}

	// Parse response
	var result models.ESResult[models.Airline]

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	airlines := make([]models.Airline, 0, len(result.Hits.Hits))

	for _, hit := range result.Hits.Hits {
		airlines = append(airlines, hit.Source)
	}
	return airlines, nil
}

func (r *AirlineRepository) FindByAirlineId(airlineId string) (*models.Airline, error) {
	query := fmt.Sprintf(`{
		"query": {
			"match": {
				"airline_id": "%s"
			}
		},
		"size": 1
	}`, airlineId)

	res, err := r.client.Search(
		r.client.Search.WithContext(context.Background()),
		r.client.Search.WithIndex("airlines"),
		r.client.Search.WithBody(strings.NewReader(query)),
		r.client.Search.WithSize(1),
		r.client.Search.WithPretty(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch airline by id=%v, error: %v", airlineId, err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("response have error: %v", res.String())
	}

	var result models.ESResult[models.Airline]
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	if len(result.Hits.Hits) == 0 {
		return nil, fmt.Errorf("airline with airline_id=%s not found", airlineId)
	}

	return &result.Hits.Hits[0].Source, nil

}

func (r *AirlineRepository) DeleteByAirlineId(airlineId string) error {

	query := fmt.Sprintf(`{
		"query": {
			"match": {
				"airline_id": "%s"
			}
		}
	}`, airlineId)
	res, err := r.client.DeleteByQuery(
		[]string{"airlines"},
		strings.NewReader(query),
		r.client.DeleteByQuery.WithContext(context.Background()))

	if err != nil {
		return fmt.Errorf("failed to delete airline by airline_id: %v", airlineId)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("elasticsearch error: %v", res.String())
	}

	return nil
}
