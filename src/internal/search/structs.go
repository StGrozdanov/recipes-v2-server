package search

import "github.com/lib/pq"

type CollectionSearch struct {
	Content string `json:"content" db:"content"`
}

type GlobalSearch struct {
	ResultType string         `json:"resultType" db:"collection_name"`
	Content    pq.StringArray `json:"content" db:"results"`
}
