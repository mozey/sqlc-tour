package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/kyleconroy/sqlc-tour/pkg/dbconn"
	"github.com/kyleconroy/sqlc-tour/pkg/logutil"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func run(ctx context.Context, db *sql.DB) error {
	q := &sqlc.Queries{db: db}

	insertedAuthor, err := q.CreateAuthor(ctx, CreateAuthorParams{
		Name: "Brian Kernighan",
		Bio: sql.NullString{
			String: "Co-author of The C Programming Language",
			Valid:  true,
		},
	})
	if err != nil {
		return errors.WithStack(err)
	}

	authors, err := q.ListAuthors(ctx)
	if err != nil {
		return errors.WithStack(err)
	}
	b, err := json.MarshalIndent(authors, "", "  ")
	if err != nil {
		return errors.WithStack(err)
	}
	log.Info().RawJSON("authors", b).Msg("")

	err = q.DeleteAuthor(ctx, insertedAuthor.ID)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func main() {
	logutil.SetupLogger(true)

	db, err := dbconn.GetConnection(&dbconn.ConnectionConfig{
		Host: "localhost",
		User: "postgres",
		Pass: "",
		Port: "5432",
		DB:   "sqlc",
	})
	if err != nil {
		log.Error().Stack().Err(err).Msg("")
	}

	err = db.Ping()
	if err != nil {
		log.Error().Stack().Err(err).Msg("")
	}

	err = run(context.Background(), db.DB)
	if err != nil {
		log.Error().Stack().Err(err).Msg("")
	}

	log.Info().Msg("It works!")
}
