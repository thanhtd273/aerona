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

type CityRepository struct {
	client *elasticsearch.Client
}

func NewCityRepository(client *elasticsearch.Client) *CityRepository {
	return &CityRepository{client: client}
}

func (r *CityRepository) CreateCity(city models.City) (*models.City, error) {
	city.CityId = uuid.New().String()
	jsonData, err := json.Marshal(city)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal city: %v", err)
	}

	res, err := r.client.Index("cities",
		bytes.NewReader(jsonData),
		r.client.Index.WithDocumentID(city.CityId),
		r.client.Index.WithContext(context.Background()))
	if err != nil {
		return nil, fmt.Errorf("failed to index city: %v", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("error response from Elasticsearch: %v", res.String())
	}
	return &city, nil
}

func (r *CityRepository) GetAllCity() ([]models.City, error) {
	query := `{
		"query": {
			"match_all": {}	
		},
		"size": 1000
	}`

	res, err := r.client.Search(
		r.client.Search.WithContext(context.Background()),
		r.client.Search.WithIndex("cities"),
		r.client.Search.WithBody(strings.NewReader(query)),
		r.client.Search.WithSize(1000),
		r.client.Search.WithPretty(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to search cities: %v", err)
	}
	defer res.Body.Close()
	if res.IsError() {
		return nil, fmt.Errorf("elasticsearch error: %s", res.String())
	}

	var result models.ESResult[models.City]
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	cities := make([]models.City, 0, len(result.Hits.Hits))
	for _, hit := range result.Hits.Hits {
		cities = append(cities, hit.Source)
	}
	return cities, nil
}

func (r *CityRepository) FindByCityId(cityId string) (*models.City, error) {
	query := fmt.Sprintf(`{
		"query": {
			"match": {
				"city_id": "%s"
			}
		},
		"size": 1
	}`, cityId)

	res, err := r.client.Search(
		r.client.Search.WithContext(context.Background()),
		r.client.Search.WithIndex("cities"),
		r.client.Search.WithBody(strings.NewReader(query)),
		r.client.Search.WithSize(1),
		r.client.Search.WithPretty(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch city by id=%v, error: %v", cityId, err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("response have error: %v", res.String())
	}

	var result models.ESResult[models.City]
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	if len(result.Hits.Hits) == 0 {
		return nil, fmt.Errorf("city with city_id=%s not found", cityId)
	}

	return &result.Hits.Hits[0].Source, nil

}

func (r *CityRepository) UpdateCity(cityId string, rawData models.RawCity) (*models.City, error) {
	updatedCity := models.ParseCity(rawData)
	now := time.Now()
	updatedCity.UpdatedAt = &now

	updatedData, err := json.Marshal(updatedCity)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal updated city: %v", err)
	}

	script := fmt.Sprintf(`{
        "query": {
            "term": {
                "city_id": "%s"
            }
        },
        "script": {
            "source": "for (entry in params.updates.entrySet()) { ctx._source[entry.getKey()] = entry.getValue() }",
            "lang": "painless",
            "params": {
                "updates": %s
            }
        }
    }`, cityId, string(updatedData))

	res, err := r.client.UpdateByQuery(
		[]string{"cities"},
		r.client.UpdateByQuery.WithContext(context.Background()),
		r.client.UpdateByQuery.WithBody(strings.NewReader(script)))
	if err != nil {
		return nil, fmt.Errorf("failed to update city by city_id: %v", err)
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
		return nil, fmt.Errorf("no city found with city_id: %s", cityId)
	}
	return &updatedCity, nil
}

func (r *CityRepository) DeleteByCityId(cityId string) error {
	query := fmt.Sprintf(`{
		"query": {
			"match": {
				"city_id": "%s"
			}
		}
	}`, cityId)

	res, err := r.client.DeleteByQuery(
		[]string{"cities"},
		strings.NewReader(query),
		r.client.DeleteByQuery.WithContext(context.Background()))

	if err != nil {
		return fmt.Errorf("failed to delete city by city_id: %v", cityId)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("elasticsearch error: %v", res.String())
	}

	return nil
}
