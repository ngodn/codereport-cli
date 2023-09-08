package git_test

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/ngodn/codereport-cli/extensions"
	"github.com/ngodn/codereport-cli/extensions/options"
	"github.com/ngodn/codereport-cli/pkg/locator"
	_ "github.com/ngodn/codereport-cli/pkg/sqlite"
	"go.riyazali.net/sqlite"
)

func init() {
	// register sqlite extension when this package is loaded
	sqlite.Register(extensions.RegisterFn(
		options.WithExtraFunctions(), options.WithRepoLocator(locator.CachedLocator(locator.MultiLocator(nil))),
	))
}

// tests' entrypoint that registers the extension
// automatically with all loaded database connections
func TestMain(m *testing.M) { os.Exit(m.Run()) }

// Memory represents a uri to an in-memory database
const Memory = "file:testing.db?mode=memory"

// Connect opens a connection with the sqlite3 database using
// the given data source address and pings it to check liveliness.
func Connect(t *testing.T, dataSourceName string) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		t.Fatalf("failed to open connection: %v", err.Error())
	}

	if err = db.Ping(); err != nil {
		t.Fatalf("failed to open connection: %v", err.Error())
	}

	return db
}
