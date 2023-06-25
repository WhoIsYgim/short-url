package postgres

import (
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"short-link/internal/entities"
	"short-link/internal/utils"
	"short-link/pkg/errs"
	"time"
)

type LinkStorage struct {
	db *sqlx.DB
}

func (ls *LinkStorage) GetLink(token string) (*entities.Link, error) {
	link := &entities.Link{}
	err := ls.db.QueryRowx(GetLinkByToken, token).StructScan(link)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.NotFoundError()
		}
		return nil, errs.InternalError(err)
	}
	return link, nil
}

func (ls *LinkStorage) GetLinkByOriginal(origLink string) (*entities.Link, error) {
	link := &entities.Link{}
	err := ls.db.QueryRowx(GetLinkByOriginalUrl, origLink).StructScan(link)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NotFoundError()
		}
		return nil, errs.InternalError(err)
	}
	return link, nil
}

func (ls *LinkStorage) StoreLink(link *entities.Link) error {
	stmt, err := ls.db.Prepare(InsertLink)
	if err != nil {
		return errs.InternalError(err)
	}
	_, err = stmt.Exec(link.OriginalLink, link.Token, link.ExpiresAt)
	if err != nil {
		return errs.InternalError(err)
	}
	return nil
}

func (ls *LinkStorage) StartRecalculation(interval time.Duration, deleted chan []string) {
	ticker := time.NewTicker(interval)
	go func() {
		for {
			<-ticker.C
			rows, err := ls.db.Queryx(DeleteOldRecords, utils.CurrentTimeString())
			if err != nil {
				continue
			}
			var del []string
			for rows.Next() {
				var deletedToken string
				err := rows.Scan(&deletedToken)
				if err != nil {
					continue
				}
				del = append(del, deletedToken)
			}
			deleted <- del
		}
	}()
}

func (ls *LinkStorage) ShutDown() error {
	return ls.db.Close()
}

func NewLinkStorage(db *sqlx.DB) *LinkStorage {
	return &LinkStorage{
		db: db,
	}
}
