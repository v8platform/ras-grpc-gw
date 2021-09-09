package domain

type AccessToken string

type RefreshToken string

type Tokens struct {
	Access  AccessToken
	Refresh RefreshToken
}
