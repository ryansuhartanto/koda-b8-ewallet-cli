package models

import (
	"time"

	"github.com/ryansuhartanto/koda-b8-ewallet-cli/db"
)

type UserSPI struct {
	db.ID

	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time

	VerifiedAt *time.Time

	Ssn       string
	LegalName string
	Dob       time.Time

	TaxId string
}
