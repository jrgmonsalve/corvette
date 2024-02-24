package application

type Recolector interface {
	Collect() error
}
