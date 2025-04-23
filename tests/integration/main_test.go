package tests

import (
	"log"
	"net/http"
	"os"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/PyMarcus/TCC_SistemasDeInformacao2025/internal/adapters/config"
	adapters "github.com/PyMarcus/TCC_SistemasDeInformacao2025/internal/adapters/http"
	"github.com/PyMarcus/TCC_SistemasDeInformacao2025/internal/adapters/repository"
	"github.com/PyMarcus/TCC_SistemasDeInformacao2025/internal/core/usecase"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Creates a mock database connection using sqlmock
func mockDBConn() (*gorm.DB, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("failed to create sqlmock: %s", err)
	}

	dialector := postgres.New(postgres.Config{Conn: db})
	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to open GORM connection: %s", err)
	}

	cleanup := func() {
		db.Close()
	}

	return gormDB, mock, cleanup
}

func TestMain(m *testing.M) {
	// Load .env config
	_, err := config.LoadConfig("../../.env")
	if err != nil {
		log.Println("[-] Failed to load config:", err)
		os.Exit(1)
	}

	testExternalAPI()
	testDatabaseMock()

	code := m.Run()
	os.Exit(code)
}

// Simulates an external API request
func testExternalAPI() {
	clientService := adapters.NewApiRequestService()
	clientUsecase := usecase.NewAPIRequestUsecase(clientService)

	headers := map[string]string{"Content-Type": "application/json"}

	response, err := clientUsecase.Fetch("https://example.com", headers, "")
	if err != nil || response == nil {
		log.Println("[-] API request error:", err)
		return
	}

	if response.StatusCode == http.StatusBadRequest {
		log.Println("[-] BadRequest error (400)")
	}
}

// Simulates a database fetch and update operation using mock
func testDatabaseMock() {
	db, mock, cleanup := mockDBConn()
	defer cleanup()

	
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "datasets"`)).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "class", "atom", "snippet", "line", "github_link", "status_code",
			"marked_by_agent_one", "marked_by_agent_two",
		}).AddRow(
			1, "security", "atom-1", "snippet-1", "line-1", "http://github.com", "200 OK", true, false,
		))

	sqlDB, err := db.DB()
	if err != nil {
		log.Println("[-] Failed to access underlying DB from GORM:", err)
		os.Exit(1)
	}

	if err = sqlDB.Ping(); err != nil {
		log.Println("[-] Failed to ping database:", err)
		os.Exit(1)
	}

	datasetRepo := repository.NewDatasetRepository(db)
	datasetUsecase := usecase.NewDatasetUsecase(datasetRepo)

	// Test: Fetch all datasets
	datasets, err := datasetUsecase.FindAll()
	if err != nil {
		log.Println("[-] Failed to fetch datasets:", err)
		os.Exit(1)
	}

	if len(datasets) > 0 {
		log.Println("First atom found:", datasets[0].Atom)
	}

	// Test: Update MarkedByAgentOne to true for a specific ID
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "datasets" SET "marked_by_agent_one"=$1 WHERE id = $2`)).
		WithArgs(true, 1).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err = datasetRepo.UpdateMarkedByAgent(1, 1)
	if err != nil {
		log.Println("[-] Failed to update MarkedByAgentOne:", err)
		os.Exit(1)
	}

	log.Println("[+] UpdateMarkedByAgentOne executed successfully")
}
