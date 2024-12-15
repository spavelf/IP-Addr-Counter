package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	start := time.Now() // Start time

	db, err := setupDatabase()
	if err != nil {
		log.Fatalf("Failed to set up database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}()

	uniqueCount, err := processIPs(db, "ip_addresses.txt")
	if err != nil {
		log.Fatalf("Failed to process IPs: %v", err)
	}

	// Measure memory usage
	memStats := getMemoryUsage()
	elapsed := time.Since(start) // Time elapsed
	fmt.Println("Number of unique IPs:", uniqueCount)
	fmt.Printf("Memory used: %0.2f MB\n", float64(memStats.Alloc)/1024.0/1024.0)
	fmt.Printf("Time taken: %0.3f ms\n", float64(elapsed.Nanoseconds())/1000000.0)
}

func setupDatabase() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "ip_addresses.db")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Drop table if it exists to start fresh
	_, err = db.Exec(`DROP TABLE IF EXISTS ips;`)
	if err != nil {
		return nil, fmt.Errorf("failed to drop table: %w", err)
	}

	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS ips (
            ip TEXT NOT NULL UNIQUE
        );
    `)
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	return db, nil
}

func processIPs(db *sql.DB, filePath string) (int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, fmt.Errorf("failed to open file: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("Error closing file: %v", err)
		}
	}()

	scanner := bufio.NewScanner(file)
	uniqueCount := 0

	for scanner.Scan() {
		ip := scanner.Text()
		_, err := db.Exec("INSERT INTO ips(ip) VALUES (?)", ip)

		if err != nil {
			// Ignore duplicate errors
			if !isDuplicateError(err) {
				return 0, fmt.Errorf("failed to insert IP %s: %w", ip, err)
			}
		} else {
			uniqueCount++
		}

	}

	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("error reading file: %w", err)
	}

	return uniqueCount, nil
}

func isDuplicateError(err error) bool {
	// In SQLite, duplicate errors have this error string
	if err != nil {
		return strings.Contains(err.Error(), "UNIQUE constraint failed")
	}
	return false
}

func getMemoryUsage() runtime.MemStats {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return memStats
}
