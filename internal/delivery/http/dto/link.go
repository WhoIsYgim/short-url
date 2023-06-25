package dto

//go:generate easyjson --all link.go


type CreateLinkRequest struct {
	Link string `json:"link"`
}

type CreateLinkResponse struct {
	ShortLink string `json:"short_link"`
	ExpiresAt string `json:"expires_at"`
}
