package entities

type Link struct {
	OriginalLink string `db:"original_link"`
	ShortLink    string
	Token        string `db:"token"`
	ExpiresAt    string `db:"expires_at"`
}

func (l *Link) Expired(now string) bool {
	return now > l.ExpiresAt
}
