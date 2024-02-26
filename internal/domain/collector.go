package domain

import "sync"

type Collector interface {
	Collect(chan<- Email, *sync.WaitGroup) error
}
