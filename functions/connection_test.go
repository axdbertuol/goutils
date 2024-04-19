package funcs_test

import (
	"fmt"

	"testing"
	"time"

	utils "github.com/axdbertuol/goutils/functions"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

var viperInstance *viper.Viper

func init() {
	// Initialize viperInstance
	viperInstance = viper.New()
}

type MockGormOpen struct {
	mock.Mock
}

const defaultRetryDelay = 1 * time.Second
const defaultMaxRetries = 3

func (m *MockGormOpen) Open(dsn string) (*gorm.DB, error) {
	args := m.Called(dsn)
	return args.Get(0).(*gorm.DB), args.Error(1)
}

type MockDBConnector func(dbConnStr string) (*gorm.DB, error)

func (m MockDBConnector) OpenDatabaseConnection(dbConnStr string) (*gorm.DB, error) {
	return m(dbConnStr)
}

func TestConnectToDb_SuccessfulConnection(t *testing.T) {
	dbConnStr := "host=localhost port=5432 user=user password=password dbname=test_db sslmode=disable "

	mockDb := &gorm.DB{} // Mock database connection object

	// Create a mock DBConnector function
	mockDBConnector := func(dbConnStr string) (*gorm.DB, error) {
		return mockDb, nil
	}

	// Call the utils.ConnectToDb function with the mock DBConnector
	result, err := utils.ConnectToDb(
		dbConnStr,
		mockDBConnector,
		defaultRetryDelay,
		defaultMaxRetries,
	)
	assert.NoError(t, err)
	// Verify that the utils.ConnectToDb function returns the expected result
	assert.NotNil(t, result)
	assert.Equal(t, mockDb, result)
}

func TestConnectToDb_FailedConnection(t *testing.T) {
	dbConnStr := "mock_db_connection_string"

	// Create a mock DBConnector function that always returns an error
	mockDBConnector := func(dbConnStr string) (*gorm.DB, error) {
		return nil, fmt.Errorf("connection error")
	}

	// Call the utils.ConnectToDb function with the mock DBConnector
	result, err := utils.ConnectToDb(
		dbConnStr,
		mockDBConnector,
		defaultRetryDelay,
		defaultMaxRetries,
	)

	// Verify that the utils.ConnectToDb function returns an error
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetConnection(t *testing.T) {
	// Mock Viper configuration
	viperMock := viper.New()
	viperMock.Set("DATABASE_HOST", "localhost")
	viperMock.Set("DATABASE_PORT", 5432)
	viperMock.Set("DATABASE_USERNAME", "user")
	viperMock.Set("DATABASE_PASSWORD", "password")
	viperMock.Set("DATABASE_DBNAME", "test_db")

	// Override viper default instance with mock
	viper := viperInstance
	viperInstance = viperMock
	defer func() {
		viperInstance = viper
	}()

	// Call utils.GetConnection function
	connStr, err := utils.GetConnection(viperInstance)

	// Check if an error occurred
	assert.NoError(t, err)

	// Expected database connection string
	expectedConnStr := "host=localhost port=5432 user=user password=password dbname=test_db sslmode=disable "

	// Assert that the returned connection string matches the expected connection string
	assert.Equal(t, expectedConnStr, *connStr)
}

func TestGetConnection_Error(t *testing.T) {
	// Mock Viper configuration without setting required values
	viperMock := viper.New()

	// Override viper default instance with mock
	oldViper := viperInstance
	viperInstance = viperMock
	defer func() {
		viperInstance = oldViper
	}()

	// Call utils.GetConnection function
	_, err := utils.GetConnection(viperInstance)

	// Check if an error occurred (missing configuration)
	assert.Error(t, err)
}
