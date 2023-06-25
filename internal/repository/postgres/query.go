package postgres

var (
	GetLinkByToken = `
		select s.original_link, s.token, s.expires_at from short_link s
			where s.token = $1
		`

	GetLinkByOriginalUrl = `
		select s.original_link, s.token, s.expires_at from short_link s
			where s.original_link = $1
		`

	InsertLink = `
		insert into short_link (original_link, token, expires_at) 
			values ($1, $2, $3)
		`

	DeleteOldRecords = `
         delete
             from short_link
			 where expires_at < $1
         returning token
	`
)
