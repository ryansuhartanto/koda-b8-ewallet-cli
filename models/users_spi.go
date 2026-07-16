package models

import "time"

type UserSpi struct {
	Id int64

	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time

	IdUser     int64
	VerifiedAt *time.Time

	Ssn       string
	LegalName string
	Dob       time.Time

	TaxId string
}
