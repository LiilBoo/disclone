// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0

package db

import ()

type UserProviderMapping struct {
	ID             int32
	ProviderID     string
	ProviderUserID string
	EmailAddress   string
	DisplayName    string
}
