package domain

type EmailRepository interface {
	CreateBulk([]Email) error
}
