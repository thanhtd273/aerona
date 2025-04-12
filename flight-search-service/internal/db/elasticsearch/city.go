package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"aerona.thanhtd.com/flight-search-service/internal/api/models"
	"aerona.thanhtd.com/flight-search-service/internal/utils"
	"github.com/elastic/go-elasticsearch/v7"
)

type CityRepository struct {
	client *elasticsearch.Client
}

func NewCityRepository(client *elasticsearch.Client) *CityRepository {
	return &CityRepository{client: client}
}

func (r *CityRepository) FindByCityCode(code string) (*models.City, error) {
	query := fmt.Sprintf(`{
		"query": {
			"match": {
				"code": "%s"
			}
		},
		"size": 1
	}`, code)

	res, err := r.client.Search(
		r.client.Search.WithContext(context.Background()),
		r.client.Search.WithIndex("cities"),
		r.client.Search.WithBody(strings.NewReader(query)),
		r.client.Search.WithSize(1),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch city by code=%v, error: %v", code, err)
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
		return nil, fmt.Errorf("city with code=%s not found", code)
	}

	return &result.Hits.Hits[0].Source, nil
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

func (r *CityRepository) CreateCity(city models.City) (*models.City, error) {
	// Assign ID for city
	citId, err := utils.GenerateUniqueId()
	if err != nil {
		return nil, fmt.Errorf("failed to generate Snowflake ID, error: %v", err)
	}
	city.CityId = citId

	jsonData, err := json.MarshalIndent(city, "", " ")
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

func (r *CityRepository) GetPopularCities() ([]models.City, error) {
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
		r.client.Search.WithIndex("cities"),
		r.client.Search.WithBody(strings.NewReader(query)),
		r.client.Search.WithSize(3),
		r.client.Search.WithPretty(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get 3 most popular cities, error: %v", err)
	}
	defer res.Body.Close()
	if res.IsError() {
		return nil, fmt.Errorf("elasticsearch response had error: %v", res.String())
	}

	var result models.ESResult[models.City]
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode Elasticsearch response, error: %v", err)
	}

	cities := make([]models.City, 0, len(result.Hits.Hits))
	for _, hit := range result.Hits.Hits {
		cities = append(cities, hit.Source)
	}
	return cities, nil

}

func (r *CityRepository) IncreasePopularity(cityId string) error {
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
				"city_id": "%s"
			}
		}
	}`, time.Now().UTC().Format(time.RFC3339), cityId)

	res, err := r.client.UpdateByQuery(
		[]string{"cities"},
		r.client.UpdateByQuery.WithContext(context.Background()),
		r.client.UpdateByQuery.WithBody(strings.NewReader(query)),
	)
	if err != nil {
		return fmt.Errorf("failed to update popularity of city, cityId: %s, error: %v", cityId, err)
	}
	defer res.Body.Close()
	if res.IsError() {
		return fmt.Errorf("elasticsearch response had error: %v", err)
	}
	return nil
}
