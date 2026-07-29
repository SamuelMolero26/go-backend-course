package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

func seedRoles(ctx context.Context, pool *pgxpool.Pool) error {

	query := `INSERT INTO roles (name) VALUES ($1), ($2) on CONFLICT DO NOTHING`
	_, err := pool.Exec(ctx, query, "user", "admin")

	if err != nil {
		return err
	}

	fmt.Println("Seeded roles Executed")
	return nil
}
