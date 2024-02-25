package model_zincsearch

import "github.com/jrgmonsalve/corvette/cmd/cli/internal/domain"

type EmailBulkRequest struct {
	Index   string         `json:"index"`
	Records []domain.Email `json:"records"`
}
