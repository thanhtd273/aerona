package elasticsearch

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v7"

	"aerona.thanhtd.com/flight-search-service/internal/api/dto"
	"aerona.thanhtd.com/flight-search-service/internal/api/models"
)

type FlightRepository struct {
	client *elasticsearch.Client
}

func NewFlightRepository(client *elasticsearch.Client) *FlightRepository {
	return &FlightRepository{client: client}
}

func (r *FlightRepository) SearchFlights(searchInfo dto.SearchInfo, offset int, limit int) ([]models.Flight, error) {
	query := fmt.Sprintf(`{
	"from": %d,
	"size": %d,
	"query": {
	  "bool": {
		"filter": [
		  {"term": { "departure.airport_code": "%s" }},
		  {"term": { "arrival.airport_code": "%s" }},
		  {"term": { "flight_date": "%s" }},
		  {"term": { "status": "active" } }
		],
		"must": [
			{ "range": {"seats_available": { "gte": %d }} }
			]
	  	}
		}
	}`, (offset-1)*limit, limit, searchInfo.From, searchInfo.To, searchInfo.DepartureDate, searchInfo.NumOfPassengers)

	res, err := r.client.Search(
		r.client.Search.WithContext(context.Background()),
		r.client.Search.WithIndex("flights"),
		r.client.Search.WithBody(strings.NewReader(query)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to search flights: %v", err)
	}
	defer res.Body.Close()
	if res.IsError() {
		return nil, fmt.Errorf("elasticsearch error: %s", res.String())
	}

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

func (r *FlightRepository) FilterFlights(searchInfo dto.SearchInfo) ([]models.Flight, error) {
	var filterClauses []string
	var mustClauses []string

	// Parse flight_date to use in range query
	flightDate, err := time.Parse("2006-01-02", searchInfo.DepartureDate)
	if err != nil {
		return nil, fmt.Errorf("invalid flight_date format, expected YYYY-MM-DD: %v", err)
	}
	dateStr := flightDate.Format("2006-01-02")

	filterClauses = append(filterClauses,
		fmt.Sprintf(`{"term": {"departure.airport_code": "%s"}}`, searchInfo.From),
		fmt.Sprintf(`{"term": {"arrival.airport_code": "%s"}}`, searchInfo.To),
		fmt.Sprintf(`{"term": {"flight_date": "%s"}}`, searchInfo.DepartureDate),
		`{"term": {"status": "active"}}`,
	)
	mustClauses = append(mustClauses,
		fmt.Sprintf(`{"range": {"seat_available": {"gte": %d}}}`, searchInfo.NumOfPassengers),
	)

	filters := searchInfo.Filters
	if len(filters.AirlineCode) > 0 {
		airlineTerms := make([]string, len(filters.AirlineCode))
		for i, code := range filters.AirlineCode {
			airlineTerms[i] = fmt.Sprintf(`"%s"`, code)
		}
		airlineClause := fmt.Sprintf(`{"terms": {"airport_code": [%s]}}`, strings.Join(airlineTerms, ","))
		filterClauses = append(filterClauses, airlineClause)
	}

	if len(filters.DepartureRange) > 0 {
		for _, r := range filters.DepartureRange {
			fromTime := fmt.Sprintf("%sT%s:00Z", dateStr, r.From)
			toTime := fmt.Sprintf("%sT%s:00Z", dateStr, r.To)
			rangeClause := fmt.Sprintf(`{"range": {"departure.scheduled": {"gte": "%s", "lte": "%s"}}}`, fromTime, toTime)
			filterClauses = append(filterClauses, rangeClause)
		}
	}

	if len(filters.ArrivalRange) > 0 {
		for _, r := range filters.ArrivalRange {
			fromTime := fmt.Sprintf("%sT%s:00Z", dateStr, r.From)
			toTime := fmt.Sprintf("%sT%s:00Z", dateStr, r.To)
			rangeClause := fmt.Sprintf(`{"range": {"arrival.scheduled": {"gte": "%s", "lte": "%s"}}}`, fromTime, toTime)
			filterClauses = append(filterClauses, rangeClause)
		}
	}

	query := fmt.Sprintf(`{
		"query": {
			"bool": {
				"filter": [%s],
				"must": [%s]
			}
		}
	}`, strings.Join(filterClauses, ","), strings.Join(mustClauses, ","))

	fmt.Printf("filter query: %s\n\n\n", query)

	res, err := r.client.Search(
		r.client.Search.WithContext(context.Background()),
		r.client.Search.WithIndex("flights"),
		r.client.Search.WithBody(strings.NewReader(query)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to search flights: %v", err)
	}
	defer res.Body.Close()
	if res.IsError() {
		return nil, fmt.Errorf("elasticsearch error: %s", res.String())
	}

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
