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

type FlightRepository struct {
	client *elasticsearch.Client
}

func NewFlightRepository(client *elasticsearch.Client) *FlightRepository {
	return &FlightRepository{client: client}
}

func (r *FlightRepository) CreateFlight(flight models.Flight) (*models.Flight, error) {
	// Asign ID for flight
	flight.FlightId = uuid.New().String()

	jsonData, err := json.MarshalIndent(flight, "", " ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal flight: %v", err)
	}
	res, err := r.client.Index("flights",
		bytes.NewReader(jsonData),
		r.client.Index.WithDocumentID(flight.FlightId),
		r.client.Index.WithContext(context.Background()))
	if err != nil {
		return nil, fmt.Errorf("failed to index flight: %v", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("error response from Elasticsearch: %v", res.String())
	}
	return &flight, nil
}

func (r *FlightRepository) GetAllFlight() ([]models.Flight, error) {
	query := `{
        "query": {
            "match_all": {}
        },
        "size": 1000
    }`

	res, err := r.client.Search(
		r.client.Search.WithContext(context.Background()),
		r.client.Search.WithIndex("flights"),
		r.client.Search.WithBody(strings.NewReader(query)),
		r.client.Search.WithSize(1000),
		r.client.Search.WithPretty(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to search flights: %v", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("elasticsearch error: %s", res.String())
	}

	// Parse response
	var result models.ESResult[models.Flight]

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	flights := make([]models.Flight, 0, len(result.Hits.Hits))

	for _, hit := range result.Hits.Hits {
		flights = append(flights, hit.Source)
	}
	return flights, nil
}

func (r *FlightRepository) FindByFlightId(flightId string) (*models.Flight, error) {
	query := fmt.Sprintf(`{
		"query": {
			"match": {
				"flight_id": "%s"
			}
		},
		"size": 1
	}`, flightId)

	res, err := r.client.Search(
		r.client.Search.WithContext(context.Background()),
		r.client.Search.WithIndex("flights"),
		r.client.Search.WithBody(strings.NewReader(query)),
		r.client.Search.WithSize(1),
		r.client.Search.WithPretty(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch flight by id=%v, error: %v", flightId, err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("response have error: %v", res.String())
	}

	var result models.ESResult[models.Flight]
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	if len(result.Hits.Hits) == 0 {
		return nil, fmt.Errorf("flight with flight_id=%s not found", flightId)
	}

	return &result.Hits.Hits[0].Source, nil

}

func (r *FlightRepository) UpdateFlight(flightId string, rawData models.RawFlightData) (*models.Flight, error) {
	updatedFlight := models.ParseFlight(rawData)
	now := time.Now()
	updatedFlight.UpdatedAt = &now

	updatedData, err := json.Marshal(updatedFlight)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal updated flight data: %v", err)
	}

	script := fmt.Sprintf(`{
		"query": {
			"term": {
				"flight_id": "%s"
			}
		},
		"script": {
			"source": "ctx._source = params.updatedFlight",
			"lang": "painless",
			"param": {
				"updatedFlight": %s
			}
		}
	}`, flightId, string(updatedData))

	res, err := r.client.UpdateByQuery(
		[]string{"flights"},
		r.client.UpdateByQuery.WithContext(context.Background()),
		r.client.UpdateByQuery.WithBody(strings.NewReader(script)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update flight by flight_id: %v", flightId)
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
		return nil, fmt.Errorf("no flight found with flight_id: %s", flightId)
	}

	return &updatedFlight, nil
}

func (r *FlightRepository) DeleteByFlightId(flightId string) error {

	query := fmt.Sprintf(`{
		"query": {
			"match": {
				"flight_id": "%s"
			}
		}
	}`, flightId)
	res, err := r.client.DeleteByQuery(
		[]string{"flights"},
		strings.NewReader(query),
		r.client.DeleteByQuery.WithContext(context.Background()))

	if err != nil {
		return fmt.Errorf("failed to delete flight by flight_id: %v", flightId)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("elasticsearch error: %v", res.String())
	}

	return nil
}
