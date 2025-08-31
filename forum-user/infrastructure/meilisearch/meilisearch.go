package meilisearch

import (
	"encoding/json"
	"github.com/meilisearch/meilisearch-go"
	"log/slog"
)

func New(host string, port string, apiKey string) *ClientMeiliSearch {
	client := meilisearch.New(
		"http://"+host+":"+port,
		meilisearch.WithAPIKey(apiKey),
	)
	return &ClientMeiliSearch{
		client: client,
	}
}

type ClientMeiliSearch struct {
	client meilisearch.ServiceManager
}

func (m *ClientMeiliSearch) CreateIndexes(indexes []string) error {
	for _, index := range indexes {
		err := m.createIndex(index)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *ClientMeiliSearch) createIndex(indexName string) error {
	_, err := m.client.CreateIndex(&meilisearch.IndexConfig{Uid: indexName, PrimaryKey: "id"})
	if err != nil {
		slog.Error(err.Error())
	}

	index, err := m.client.GetIndex(indexName)
	if err != nil {
		return err
	}

	_, err = index.UpdateSettings(&meilisearch.Settings{
		SearchableAttributes: []string{"login:prefix"},
	})

	if err != nil {
		return err
	}

	return nil
}

func (m *ClientMeiliSearch) AddDocuments(indexName string, documents []any) error {
	index, err := m.client.GetIndex(indexName)
	if err != nil {
		return err
	}

	_, err = index.AddDocuments(documents, "id")
	if err != nil {
		return err
	}
	return nil
}

func (m *ClientMeiliSearch) SimpleSearch(indexName, query string) ([]byte, error) {
	index, err := m.client.GetIndex(indexName)
	if err != nil {
		return nil, err
	}

	searchResult, err := index.Search(
		query,
		&meilisearch.SearchRequest{
			Query: query,
		},
	)
	if err != nil {
		return nil, err
	}

	searchResultBytes, err := json.Marshal(searchResult.Hits)
	if err != nil {
		return nil, err
	}
	return searchResultBytes, nil
}
