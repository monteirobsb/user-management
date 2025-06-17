package services

import (
	"log"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/monteirobsb/user-management/backend/database"
	"github.com/monteirobsb/user-management/backend/models"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var testDB *gorm.DB // Holds a global test DB instance for this package

// setupTestSQLiteDB initializes an in-memory SQLite database for testing.
// It migrates the User schema.
func setupTestSQLiteDB(t *testing.T) {
	// Only initialize testDB once for the package to avoid re-opening connections unnecessarily.
	// TestMain can call this, or individual test suites/tests can call it.
	// For simplicity here, each test suite (like TestCreateUser_PasswordHashing_WithSQLite) will ensure it's setup.
	if testDB == nil {
		var err error
		// Using "file::memory:?cache=shared" allows the in-memory database to be potentially shared
		// between connections if needed, though for these tests, one connection is likely sufficient.
		// A simple "file::memory:" would also work for a private in-memory DB per gorm.Open call.
		testDB, err = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
		if err != nil {
			// Use t.Fatalf to stop the test immediately if DB setup fails.
			t.Fatalf("FATAL: Failed to connect to test SQLite database: %v", err)
		}

		// AutoMigrate the schema for the User model.
		err = testDB.AutoMigrate(&models.User{})
		if err != nil {
			t.Fatalf("FATAL: Failed to migrate test database schema: %v", err)
		}
		log.Println("INFO: Test SQLite database initialized and migrated.")
	}
}

// TestMain provides a hook to run setup and teardown code for the package.
func TestMain(m *testing.M) {
	// Optional: Global setup can happen here, e.g. calling setupTestSQLiteDB(nil) if t is not needed
	// or if setup doesn't involve fatal assertions specific to a test.
	// However, it's often better to let tests manage their own DB setup for clarity.

	exitVal := m.Run() // Run all tests in the package.

	// Global teardown: Close the test database connection if it was opened.
	if testDB != nil {
		sqlDB, err := testDB.DB()
		if err == nil && sqlDB != nil {
			errClose := sqlDB.Close()
			if errClose == nil {
				log.Println("INFO: Test SQLite database connection closed.")
			} else {
				log.Printf("WARN: Error closing test SQLite database: %v", errClose)
			}
		}
		testDB = nil // Reset for any potential future runs in a long-lived process.
	}
	os.Exit(exitVal)
}

// TestCreateUser_PasswordHashing_WithSQLite tests the CreateUser service function,
// focusing on password hashing and database interaction using an in-memory SQLite DB.
func TestCreateUser_PasswordHashing_WithSQLite(t *testing.T) {
	setupTestSQLiteDB(t) // Ensure testDB is initialized for this test.
	assert := assert.New(t)

	// Temporarily replace the global database.DB instance with our testDB.
	// This is a common workaround when dealing with global variables in tests.
	// Ensure to restore it afterwards to avoid affecting other tests or packages.
	originalGlobalDB := database.DB
	database.DB = testDB
	defer func() { database.DB = originalGlobalDB }()

	plainPassword := "securePassword123"
	uniqueEmail := "hash.test." + uuid.NewString() + "@example.com" // Ensure unique email for each test run
	user := &models.User{
		Name:  "Hashing Test User",
		Email: uniqueEmail,
	}

	// Call the CreateUser service function.
	err := CreateUser(user, plainPassword)
	assert.NoError(err, "CreateUser should succeed with the test SQLite database")

	// --- Assertions about password hashing ---
	assert.NotEmpty(user.PasswordHash, "PasswordHash should not be empty after CreateUser")

	// Verify the hash against the original password.
	errCompare := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(plainPassword))
	assert.NoError(errCompare, "Hashed password should match the original password")

	assert.NotEqual(plainPassword, user.PasswordHash, "PasswordHash should be different from the plaintext password")

	// models.User.Password is a transient field (gorm:"-") and should not store the plaintext password.
	// The CreateUser function takes plainPassword as a parameter and does not set user.Password from it.
	assert.Empty(user.Password, "models.User.Password field (transient) should remain empty")

	// --- Assertions about database interaction (optional but good for this integration-style service test) ---
	var dbUser models.User
	// Query the database for the created user by ID (GORM populates user.ID upon successful creation).
	dbResult := database.DB.First(&dbUser, "id = ?", user.ID)
	assert.NoError(dbResult.Error, "User should be found in the database after creation using its ID")
	assert.Equal(user.Email, dbUser.Email, "Emails should match for the user retrieved from DB")
	assert.Equal(user.PasswordHash, dbUser.PasswordHash, "PasswordHashes should match for the user retrieved from DB")
	assert.NotEmpty(dbUser.ID, "User ID retrieved from DB should not be empty")
	assert.Equal(user.ID, dbUser.ID, "User ID from CreateUser call and DB retrieval should match")
}
