package application

type Indexer struct {
	recolector Recolector
}

func NewIndexer(recolector Recolector) *Indexer {
	return &Indexer{recolector: recolector}
}

func (i *Indexer) Start() {
	i.recolector.Collect()
}
