package model_zincsearch

import (
	"encoding/json"
	"log"

	"github.com/jrgmonsalve/corvette/cmd/cli/internal/domain"
)

type EmailBulkRequest struct {
	IndexName string         `json:"index"`
	Records   []domain.Email `json:"records"`
}

func (emailBulk *EmailBulkRequest) MappingToJson() []byte {
	jsonData, err := json.Marshal(emailBulk)
	if err != nil {
		log.Println("Error marshaling JSON:", err)
	}
	return jsonData
}
