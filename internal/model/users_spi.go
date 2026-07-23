package model

import (
	"time"

	"github.com/ryansuhartanto/koda-b8-ewallet-cli/internal/db"
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

	TaxID string
}

type RepoUserSPI db.RepoSoftDelete[UserSPI, db.ID]
