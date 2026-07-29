package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	dsn := os.Getenv("DB_DSN")

	//for sanity check
	if dsn == "" {
		log.Fatal("DB_DSN enviroment variable is required")
	}

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("unable to connect: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("unable to ping: %v", err)
	}

	if err := runMigrations(ctx, pool); err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	if err := seedRoles(ctx, pool); err != nil {
		log.Fatalf("seed failed: %v", err)
	}

	fmt.Println("Migrations completed successfully")

}

func runMigrations(ctx context.Context, pool *pgxpool.Pool) error {

	files, err := filepath.Glob("migrations/*.up.sql")
	if err != nil {
		return fmt.Errorf("reading migration files: %v", err)
	}

	sort.Strings(files)

	for _, file := range files {
		fmt.Printf("Running migration: %s\n", filepath.Base(file))

		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("reading %s: %w", file, err)
		}

		statements := strings.Split(string(content), ";")
		for _, stmt := range statements {
			stmt = strings.TrimSpace(stmt)
			if stmt == "" {
				continue
			}

			_, err := pool.Exec(ctx, stmt)
			if err != nil {
				return fmt.Errorf("executing %s: %w", filepath.Base(file), err)
			}
		}

	}

	return nil

}
