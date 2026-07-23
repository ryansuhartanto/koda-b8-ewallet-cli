package model

import "time"

type UserSPI struct {
	Id `db:"id"`

	CreatedAt *time.Time `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`

	VerifiedAt *time.Time `db:"verified_at"`

	Ssn       string    `db:"ssn"`
	LegalName string    `db:"legal_name"`
	Dob       time.Time `db:"dob"`

	TaxID string `db:"tax_id"`
}
