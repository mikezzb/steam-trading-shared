package database_test

import (
	"testing"
	"time"

	shared "github.com/mikezzb/steam-trading-shared"
	"github.com/mikezzb/steam-trading-shared/database"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var dbUri = "mongodb://localhost:27017"
var dbName = "steam-trading-unit-test"

var dbClient, _ = database.NewDBClient(dbUri, dbName, 10*time.Second)
var err error

func TestTimeConvert(t *testing.T) {
	t.Run("TimeConvert", func(t *testing.T) {
		err = dbClient.ConvertUnixToTime("listings")
		if err != nil {
			t.Fatalf("Failed to convert time: %v", err)
		}
	})
}

func TestInit(t *testing.T) {
	defer dbClient.Disconnect()

	if err = dbClient.Ping(); err != nil {
		t.Fatalf("Failed to ping db: %v", err)
	}

	// init
	if err = dbClient.Init(); err != nil {
		t.Fatalf("Failed to init db: %v", err)
	}
}

func TestMigrate(t *testing.T) {
	dbClient, err := database.NewDBClient(dbUri, dbName, 10*time.Second)
	if err != nil {
		t.Fatalf("Failed to connect to db: %v", err)
	}
	defer dbClient.Disconnect()

	if err = dbClient.Ping(); err != nil {
		t.Fatalf("Failed to ping db: %v", err)
	}

	// migrate
	if err = dbClient.MigrateTransactions("transactions-old", "transactions"); err != nil {
		t.Fatalf("Failed to migrate db: %v", err)
	}
}

func TestDecimal128(t *testing.T) {
	t.Run("Decimal128", func(t *testing.T) {
		a, _ := primitive.ParseDecimal128("0.434223353413354")
		b, _ := primitive.ParseDecimal128("0.434223353413355")

		if shared.DecCompareTo(a, b) != -1 {
			t.Fatalf("Failed to compare Decimal128")
		}

		if shared.DecCompareTo(b, a) != 1 {
			t.Fatalf("Failed to compare Decimal128")
		}
	})
}

func TestConvertToDecimal(t *testing.T) {
	t.Run("ConvertToDecimal", func(t *testing.T) {
		fields := []string{
			"price",
			"paintWear",
		}

		err := dbClient.ConvertToDecimal128("transactions-old", fields)
		if err != nil {
			t.Fatalf("Failed to convert to decimal: %v", err)
		}
	})
}

func TestReformatTransactions(t *testing.T) {
	t.Run("ReformatTransactions", func(t *testing.T) {
		err := dbClient.ReformatTransactionCollection("transactions-old")
		if err != nil {
			t.Fatalf("Failed to reformat transactions: %v", err)
		}
	})
}
