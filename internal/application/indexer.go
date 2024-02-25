package application

import "fmt"

type Indexer struct {
	recolector Recolector
	uploader   Uploader
}

func NewIndexer(recolector Recolector, uploader Uploader) *Indexer {
	return &Indexer{
		recolector: recolector,
		uploader:   uploader,
	}
}

func (i *Indexer) Start() {
	err := i.recolector.Collect()
	if err != nil {
		fmt.Println(err)
		return
	}
}
