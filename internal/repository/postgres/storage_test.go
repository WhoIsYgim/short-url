package postgres

import (
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"short-link/internal/entities"
	"short-link/pkg/errs"
	"testing"
)

func TestLinkStorage_GetLink(t *testing.T) {
	type dbBehaviour struct {
		error error
		data  *sqlmock.Rows
	}
	tests := []struct {
		name  string
		token string

		expectedError error
		expectedLink  *entities.Link

		dbBehaviour
		setupMock func(mock sqlmock.Sqlmock, token string, link *entities.Link, behaviour dbBehaviour)
	}{
		{
			name:          "Success",
			token:         "qwerty123_",
			expectedError: nil,
			expectedLink: &entities.Link{
				OriginalLink: "http://wikipedia.org",
				Token:        "qwerty123_",
				ExpiresAt:    "2024-01-02 15:04:05",
			},
			dbBehaviour: dbBehaviour{
				data: sqlmock.NewRows([]string{
					"original_link", "token", "expires_at",
				}).
					AddRow(
						"http://wikipedia.org", "qwerty123_", "2024-01-02 15:04:05",
					),
			},
			setupMock: func(mock sqlmock.Sqlmock, token string, link *entities.Link, behaviour dbBehaviour) {
				mock.ExpectQuery(GetLinkByToken).WithArgs(token).WillReturnRows(behaviour.data)

			},
		},
		{
			name:          "Not Found",
			token:         "qwerty123_",
			expectedError: errs.NotFoundError(),
			dbBehaviour: dbBehaviour{
				error: sql.ErrNoRows,
			},
			setupMock: func(mock sqlmock.Sqlmock, token string, link *entities.Link, behaviour dbBehaviour) {
				mock.ExpectQuery(GetLinkByToken).WithArgs(token).WillReturnError(behaviour.error)
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			require.NoError(t, err)
			defer func(db *sql.DB) {
				_ = db.Close()
			}(db)

			test.setupMock(mock, test.token, test.expectedLink, test.dbBehaviour)

			sqlxDB := sqlx.NewDb(db, "sqlmock")
			repo := &LinkStorage{
				db: sqlxDB,
			}

			link, err := repo.GetLink(test.token)
			if test.expectedError != nil {
				require.ErrorIs(t, errors.Unwrap(err), errors.Unwrap(test.expectedError))
			} else {
				require.NoError(t, err)
				require.Equal(t, test.expectedLink, link)
			}
		})
	}
}
