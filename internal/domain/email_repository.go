package domain

type EmailRepository interface {
	Store(email Email) error
}
