package errors

import (
	"database/sql"
	"errors"
	"log/slog"

	"github.com/lib/pq"
)

var (
    ErrAlreadyExists = errors.New("already exists")
    ErrNotFound      = errors.New("not found")
)

func MapPgError(err error) error {
    if err == nil {
        return nil
    }

    if e, ok := pgErr(err); ok {
        switch e.Code {
        case "23505": // unique_violation
            return ErrAlreadyExists
		default:
			slog.Warn("Unhandled PostgreSQL error", 
				slog.String("code", string(e.Code)), 
				slog.String("message", e.Message),
			)
			return err
        }
		
    }

    if errors.Is(err, sql.ErrNoRows) {
        return ErrNotFound
    }

    return err
}

func pgErr(err error) (*pq.Error, bool) {
    var e *pq.Error
    if errors.As(err, &e) {
        return e, true
    }
    return nil, false
}
