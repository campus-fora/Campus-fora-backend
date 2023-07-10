package posts

import (
	"encoding/json"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	elastic "gopkg.in/olivere/elastic.v5"
)

type Document struct {
	gorm.Model
	Content map[string]interface{} `gorm:"column:content" sql:"type:jsonb"`
}

func main() {
	// Establish PostgreSQL connection using GORM
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=yourdb password=yourpassword sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer db.Close()

	// Initialize Elasticsearch client
	esClient, err := elastic.NewClient(elastic.SetURL("http://localhost:9200"))
	if err != nil {
		log.Fatalf("Failed to connect to Elasticsearch: %v", err)
	}

	// Create Elasticsearch index and mapping
	indexName := "documents"
	mapping := `
	{
		"mappings": {
			"properties": {
				"content": {
					"type": "text"
				}
			}
		}
	}`

	_, err = esClient.CreateIndex(indexName).BodyString(mapping).Do()
	if err != nil {
		log.Fatalf("Failed to create Elasticsearch index: %v", err)
	}

	// Enable automatic indexing of Document records in Elasticsearch
	db.Callback().Create().After("gorm:create").Register("index_document", func(scope *gorm.Scope) {
		if doc, ok := scope.Value.(*Document); ok {
			err := indexDocument(esClient, indexName, doc)
			if err != nil {
				log.Printf("Failed to index document: %v", err)
			}
		}
	})

	db.Callback().Update().After("gorm:update").Register("index_document", func(scope *gorm.Scope) {
		if doc, ok := scope.Value.(*Document); ok {
			err := indexDocument(esClient, indexName, doc)
			if err != nil {
				log.Printf("Failed to update document index: %v", err)
			}
		}
	})

	db.Callback().Delete().After("gorm:delete").Register("delete_document_index", func(scope *gorm.Scope) {
		if doc, ok := scope.Value.(*Document); ok {
			err := deleteDocumentIndex(esClient, indexName, doc)
			if err != nil {
				log.Printf("Failed to delete document index: %v", err)
			}
		}
	})

	// Perform a keyword-based search
	keyword := "example"
	results, err := searchDocuments(esClient, indexName, keyword)
	if err != nil {
		log.Fatalf("Failed to search documents: %v", err)
	}

	log.Printf("Search results:")
	for _, result := range results {
		log.Printf("Document ID: %d", result.ID)
	}
}

// Index a document in Elasticsearch
func indexDocument(client *elastic.Client, indexName string, doc *Document) error {
	body, err := json.Marshal(doc.Content)
	if err != nil {
		return err
	}

	_, err = client.Index().
		Index(indexName).
		Type("document").
		Id(doc.ID).
		BodyJson(string(body)).
		Refresh("wait_for").
		Do()
	return err
}

// Delete a document index from Elasticsearch
func deleteDocumentIndex(client *elastic.Client, indexName string, doc *Document) error {
	_, err := client.Delete().
		Index(indexName).
		Type("document").
		Id(doc.ID).
		Refresh("wait_for").
		Do()
	return err
}

// Search documents in Elasticsearch based on a keyword
func searchDocuments(client *elastic.Client, indexName, keyword string) ([]*Document, error) {
	query := elastic.NewMatchQuery("content", keyword)

	searchResult, err := client.Search().
		Index(indexName).
		Query(query).
		Do()
	if err != nil {
		return nil, err
	}

	var results []*Document
	for _, hit := range searchResult.Hits.Hits {
		var doc Document
		err := json.Unmarshal(*hit.Source, &doc)
		if err != nil {
			log.Printf("Failed to unmarshal search result: %v", err)
			continue
		}
		results = append(results, &doc)
	}

	return results, nil
}
