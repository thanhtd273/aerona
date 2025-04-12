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

type SearchHistoryRepository struct {
	client *elasticsearch.Client
}

func NewSearchHistoryRepository(client *elasticsearch.Client) *SearchHistoryRepository {
	return &SearchHistoryRepository{
		client: client,
	}
}

func (r *SearchHistoryRepository) CreateHistory(history models.SearchHistory) (*models.SearchHistory, error) {
	searchId, err := utils.GenerateUniqueId()
	if err != nil {
		return nil, fmt.Errorf("failed to generate Snowflake ID, error: %v", err)
	}
	history.SearchId = searchId

	jsonData, err := json.MarshalIndent(history, "", " ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal search history: %v", err)
	}
	res, err := r.client.Index("search histories",
		bytes.NewReader(jsonData),
		r.client.Index.WithDocumentID(history.SearchId),
		r.client.Index.WithContext(context.Background()))
	if err != nil {
		return nil, fmt.Errorf("failed to index search history: %v", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("error response from Elasticsearch: %v", res.String())
	}
	return &history, nil
}

func (r *SearchHistoryRepository) UpdateHistory(history models.SearchHistory) (*models.SearchHistory, error) {
	now := time.Now()
	history.SearchedAt = &now
	script := fmt.Sprintf(`{
		"query": {
			"term": {
				"search_code": "%s"
			}
		},
		"script": {
			"source": "ctx._source.searched_at = params.searched_at",
			"lang": "painless",
			"param": {
				"searched_at": %s
			}
		}
	}`, history.SearchCode, now)

	res, err := r.client.UpdateByQuery(
		[]string{"search_histories"},
		r.client.UpdateByQuery.WithContext(context.Background()),
		r.client.UpdateByQuery.WithBody(strings.NewReader(script)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update search history")
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
		return nil, fmt.Errorf("no search history found ")
	}

	return &history, nil
}

func (r *SearchHistoryRepository) GetRecentSearches() ([]models.SearchHistory, error) {
	const LIMIT = 3
	query := fmt.Sprintf(`{
		"sort": [
			{
			"searched_at": {
				"order": "desc"
			}
			}
		],
		"size": %d
		}`, LIMIT)

	res, err := r.client.Search(
		r.client.Search.WithContext(context.Background()),
		r.client.Search.WithIndex("search_histories"),
		r.client.Search.WithBody(strings.NewReader(query)),
		r.client.Search.WithSize(3),
		r.client.Search.WithPretty(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get 3 most popular search histories, error: %v", err)
	}
	defer res.Body.Close()
	if res.IsError() {
		return nil, fmt.Errorf("elasticsearch response had error: %v", res.String())
	}

	var result models.ESResult[models.SearchHistory]
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode Elasticsearch response, error: %v", err)
	}

	searchHistorys := make([]models.SearchHistory, 0, len(result.Hits.Hits))
	for _, hit := range result.Hits.Hits {
		searchHistorys = append(searchHistorys, hit.Source)
	}
	return searchHistorys, nil
}

func (r *SearchHistoryRepository) DeleteBySearchId(searchId string) error {
	query := fmt.Sprintf(`{
		"query": {
			"match": {
				"search_id": "%s"
			}
		}
	}`, searchId)
	res, err := r.client.DeleteByQuery(
		[]string{"search_histories"},
		strings.NewReader(query),
		r.client.DeleteByQuery.WithContext(context.Background()))

	if err != nil {
		return fmt.Errorf("failed to delete search history by search_id: %v", searchId)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("elasticsearch error: %v", res.String())
	}

	return nil
}

func (r *SearchHistoryRepository) IsHistoryExisting(searchCode string) (bool, error) {
	query := fmt.Sprintf(`{
		"query": {
			"term": {
				"search_code": {
					"value": "%s"
				}
			}
		},
		"size": 0
	}`, searchCode)

	res, err := r.client.Search(
		r.client.Search.WithContext(context.Background()),
		r.client.Search.WithIndex("airports"),
		r.client.Search.WithBody(strings.NewReader(query)),
	)
	if err != nil {
		return false, fmt.Errorf("failed to check the existing of airport, error: %v", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return false, fmt.Errorf("response have error: %v", res.String())
	}

	var result models.ESResult[models.Airport]
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return false, fmt.Errorf("failed to decode response: %v", err)
	}

	return result.Hits.Total.Value > 0, nil
}
