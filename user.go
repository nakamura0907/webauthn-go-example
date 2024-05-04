package main

import "github.com/go-webauthn/webauthn/webauthn"

type User struct {
	ID          []byte `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name,omitempty"`

	Credentials []webauthn.Credential `json:"-"`
}

// WebAuthnCredentials implements webauthn.User.
func (u *User) WebAuthnCredentials() []webauthn.Credential {
	return u.Credentials
}

// WebAuthnDisplayName implements webauthn.User.
func (u *User) WebAuthnDisplayName() string {
	return u.DisplayName
}

// WebAuthnID implements webauthn.User.
func (u *User) WebAuthnID() []byte {
	return u.ID
}

// WebAuthnIcon implements webauthn.User.
func (u *User) WebAuthnIcon() string {
	return ""
}

// WebAuthnName implements webauthn.User.
func (u *User) WebAuthnName() string {
	return u.Name
}
