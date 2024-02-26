package application

import (
	"log"
	"sync"

	"github.com/jrgmonsalve/corvette/cmd/cli/internal/domain"
	"github.com/jrgmonsalve/corvette/cmd/cli/pkg/helpers"
)

type Indexer struct {
	collector       domain.Collector
	emailRepository domain.EmailRepository
}

func NewIndexer(collector domain.Collector, emailRepository domain.EmailRepository) *Indexer {
	return &Indexer{
		collector:       collector,
		emailRepository: emailRepository,
	}
}

func (i *Indexer) Start() {
	batchSize := helpers.GetIntValueFromEnv("BATCH_SIZE", 10)
	emailChan := make(chan domain.Email, 4)
	var wg sync.WaitGroup

	wg.Add(1)
	go i.collector.Collect(emailChan, &wg)

	wg.Add(1)
	go func() {
		defer wg.Done()
		var batch []domain.Email

		for email := range emailChan {
			batch = append(batch, email)
			if len(batch) >= batchSize {
				if err := i.emailRepository.CreateBulk(batch); err != nil {
					log.Println("Error creating bulk: ", err)
				}
				batch = []domain.Email{}
			}
		}

		if len(batch) > 0 {
			if err := i.emailRepository.CreateBulk(batch); err != nil {
				log.Println("Error creating bulk: ", err)
			}
		}
	}()

	wg.Wait()
}
