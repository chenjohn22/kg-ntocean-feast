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

type restaurantRecord struct {
	Category  string `json:"category"`
	Name      string `json:"name"`
	Location  string `json:"location"`
	SortOrder uint   `json:"sort_order"`
}

func main() {
	seedPath := flag.String("seed", "../random_codes_60000.json", "path to the random-code JSON file")
	restaurantsPath := flag.String("restaurants", "data/activity_restaurants.json", "path to the activity-restaurant JSON file")
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
	if err := seedRestaurants(cfg.DSN, *restaurantsPath); err != nil {
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

func seedRestaurants(dsn, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open restaurant seed file: %w", err)
	}
	defer file.Close()

	var records []restaurantRecord
	if err := json.NewDecoder(file).Decode(&records); err != nil {
		return fmt.Errorf("decode restaurant seed file: %w", err)
	}
	if len(records) == 0 {
		return errors.New("restaurant seed file is empty")
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
	if _, err := tx.Exec("UPDATE activity_restaurants SET active = FALSE"); err != nil {
		tx.Rollback()
		return fmt.Errorf("deactivate restaurants: %w", err)
	}
	for index, record := range records {
		record.Category = strings.TrimSpace(record.Category)
		record.Name = strings.TrimSpace(record.Name)
		record.Location = strings.TrimSpace(record.Location)
		if record.Category == "" || record.Name == "" {
			tx.Rollback()
			return fmt.Errorf("invalid restaurant at JSON item %d", index+1)
		}
		if _, err := tx.Exec(`
			INSERT INTO activity_restaurants (category, name, location, sort_order, active)
			VALUES (?, ?, ?, ?, TRUE)
			ON DUPLICATE KEY UPDATE sort_order = VALUES(sort_order), active = TRUE`,
			record.Category, record.Name, record.Location, record.SortOrder); err != nil {
			tx.Rollback()
			return fmt.Errorf("insert restaurant at JSON item %d: %w", index+1, err)
		}
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	log.Printf("migration complete; %d activity restaurants imported", len(records))
	return nil
}
