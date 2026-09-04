package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"ntocean-feast/backend/internal/config"
)

type codeRecord struct {
	Value string `json:"value"`
}

func main() {
	seedPath := flag.String("seed", "../random_codes_60000.json", "path to the random-code JSON file")
	direction := flag.String("direction", "up", "migration direction: up or down")
	flag.Parse()

	cfg := config.Load()
	migrationsPath, err := filepath.Abs("migrations")
	if err != nil {
		log.Fatal(err)
	}
	m, err := migrate.New("file://"+filepath.ToSlash(migrationsPath), "mysql://"+cfg.DSN)
	if err != nil {
		log.Fatal(err)
	}
	defer m.Close()

	switch *direction {
	case "up":
		err = m.Up()
	case "down":
		err = m.Steps(-1)
	default:
		log.Fatalf("unknown direction %q", *direction)
	}
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal(err)
	}
	if *direction == "down" {
		log.Println("migration rolled back")
		return
	}

	if err := seedCodes(cfg.DSN, *seedPath); err != nil {
		log.Fatal(err)
	}
}

func seedCodes(dsn, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open seed file: %w", err)
	}
	defer file.Close()

	var records []codeRecord
	if err := json.NewDecoder(file).Decode(&records); err != nil {
		return fmt.Errorf("decode seed file: %w", err)
	}
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	for i, record := range records {
		if record.Value == "" {
			tx.Rollback()
			return fmt.Errorf("empty code at JSON item %d", i+1)
		}
	}
	const batchSize = 500
	for start := 0; start < len(records); start += batchSize {
		end := min(start+batchSize, len(records))
		placeholders := make([]string, 0, end-start)
		args := make([]any, 0, end-start)
		for _, record := range records[start:end] {
			placeholders = append(placeholders, "(?)")
			args = append(args, record.Value)
		}
		query := "INSERT IGNORE INTO lottery_codes (code) VALUES " + strings.Join(placeholders, ",")
		if _, err := tx.Exec(query, args...); err != nil {
			tx.Rollback()
			return fmt.Errorf("insert code batch starting at item %d: %w", start+1, err)
		}
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	log.Printf("migration complete; %d codes checked/imported", len(records))
	return nil
}
