package integration

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/happYness-Project/taskManagementGolang/pkg/configs"
	"github.com/happYness-Project/taskManagementGolang/pkg/dbs"
	"github.com/stretchr/testify/require"
)

var testDB *sql.DB

// TestMain sets up the test database connection for all integration tests
func TestMain(m *testing.M) {
	// Load configuration - connects to docker-compose database
	env := configs.InitConfig("")
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable timezone=UTC connect_timeout=5",
		env.DBHost, env.DBPort, env.DBUser, env.DBPwd, env.DBName,
	)

	var err error
	testDB, err = dbs.ConnectToDb(connStr)
	if err != nil {
		log.Fatal("Failed to connect to test database. Make sure Docker Compose is running (make start):", err)
	}

	// Verify connection
	if err := testDB.Ping(); err != nil {
		log.Fatal("Failed to ping test database:", err)
	}

	log.Println("Integration tests: Connected to test database")

	// Run tests
	code := m.Run()

	// Cleanup
	testDB.Close()
	os.Exit(code)
}

// setupTest provides test isolation by cleaning the database before each test
// Returns a cleanup function that can be deferred
func setupTest(t *testing.T) func() {
	t.Helper()
	cleanDatabase(t, testDB)

	return func() {
		// Optional: cleanup after test if needed
		// cleanDatabase(t, testDB)
	}
}

// cleanDatabase truncates all tables in the correct order to respect foreign key constraints
func cleanDatabase(t *testing.T, db *sql.DB) {
	t.Helper()

	ctx := context.Background()

	// Order matters - delete child tables first, then parent tables
	tables := []string{
		"taskcontainer_task", // Join table - task to container relationship
		"usergroup_user",     // Join table - user to group relationship
		"task",               // Tasks
		"taskcontainer",      // Containers
		"usergroup",          // Groups
		"user",               // Users
	}

	for _, table := range tables {
		query := fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", table)
		_, err := db.ExecContext(ctx, query)
		require.NoError(t, err, "Failed to truncate table: %s", table)
	}
}

// withTransaction runs a test function within a database transaction and always rolls back
// This provides complete isolation without affecting other tests
func withTransaction(t *testing.T, db *sql.DB, fn func(*sql.Tx)) {
	t.Helper()

	tx, err := db.Begin()
	require.NoError(t, err, "Failed to begin transaction")

	defer func() {
		// Always rollback - even if test passes, we don't want to commit test data
		if err := tx.Rollback(); err != nil {
			// If rollback fails and transaction was already committed/rolled back, that's okay
			t.Logf("Transaction rollback error (might be already closed): %v", err)
		}
	}()

	fn(tx)
}

// Helper to execute a query and return the count
func countRows(t *testing.T, db *sql.DB, table string) int {
	t.Helper()

	var count int
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", table)
	err := db.QueryRow(query).Scan(&count)
	require.NoError(t, err)

	return count
}
