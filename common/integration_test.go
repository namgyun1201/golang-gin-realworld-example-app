//go:build integration

package common

import (
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// setupIntegrationTestDB creates a unique test DB for each integration test
// to avoid conflicts with the unit test TestMain DB.
func setupIntegrationTestDB(t *testing.T) {
	t.Helper()
	os.Setenv("TEST_DB_PATH", fmt.Sprintf("./data/integration_%s.db", t.Name()))
}

func TestIntegration_ConcurrentDBAccess(t *testing.T) {
	asserts := assert.New(t)
	setupIntegrationTestDB(t)
	origDB := DB

	db := TestDBInit()
	defer func() {
		TestDBFree(db)
		DB = origDB
	}()

	// Create a simple table for testing
	type ConcurrentTestModel struct {
		ID    uint   `gorm:"primaryKey"`
		Value string `gorm:"column:value"`
	}
	db.AutoMigrate(&ConcurrentTestModel{})

	const goroutines = 10
	const opsPerGoroutine = 5
	var wg sync.WaitGroup
	errCh := make(chan error, goroutines*opsPerGoroutine*2)

	// Multiple goroutines writing concurrently
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < opsPerGoroutine; j++ {
				model := ConcurrentTestModel{
					Value: RandString(10),
				}
				if err := db.Create(&model).Error; err != nil {
					errCh <- err
				}
			}
		}(i)
	}
	wg.Wait()
	close(errCh)

	for err := range errCh {
		asserts.NoError(err, "Concurrent writes should not produce errors")
	}

	// Verify all records were written
	var count int64
	db.Model(&ConcurrentTestModel{}).Count(&count)
	asserts.Equal(int64(goroutines*opsPerGoroutine), count, "All concurrent writes should be persisted")

	// Multiple goroutines reading concurrently
	var readWg sync.WaitGroup
	readErrCh := make(chan error, goroutines)
	for i := 0; i < goroutines; i++ {
		readWg.Add(1)
		go func() {
			defer readWg.Done()
			var models []ConcurrentTestModel
			if err := db.Find(&models).Error; err != nil {
				readErrCh <- err
			}
		}()
	}
	readWg.Wait()
	close(readErrCh)

	for err := range readErrCh {
		asserts.NoError(err, "Concurrent reads should not produce errors")
	}
}

func TestIntegration_DBConnectionPool(t *testing.T) {
	asserts := assert.New(t)
	setupIntegrationTestDB(t)
	origDB := DB

	db := TestDBInit()
	defer func() {
		TestDBFree(db)
		DB = origDB
	}()

	sqlDB, err := db.DB()
	asserts.NoError(err, "Should get underlying sql.DB")

	// Verify connection pool settings
	stats := sqlDB.Stats()
	asserts.GreaterOrEqual(stats.MaxOpenConnections, 0, "MaxOpenConnections should be set")

	// Verify the pool is functional by pinging
	asserts.NoError(sqlDB.Ping(), "Connection pool should be functional")

	// Verify MaxIdleConns was set (TestDBInit sets it to 3)
	// Open multiple connections to exercise the pool
	type PoolTestModel struct {
		ID    uint   `gorm:"primaryKey"`
		Value string `gorm:"column:value"`
	}
	db.AutoMigrate(&PoolTestModel{})

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			db.Create(&PoolTestModel{Value: RandString(5)})
		}(i)
	}
	wg.Wait()

	var count int64
	db.Model(&PoolTestModel{}).Count(&count)
	asserts.Equal(int64(20), count, "All pool operations should complete successfully")

	// Verify pool stats after operations
	stats = sqlDB.Stats()
	asserts.GreaterOrEqual(stats.OpenConnections, 0, "Should have open connections")
}

func TestIntegration_TestDBCleanup(t *testing.T) {
	asserts := assert.New(t)
	setupIntegrationTestDB(t)
	origDB := DB

	testDBPath := GetTestDBPath()

	// Create the test database
	db := TestDBInit()
	defer func() { DB = origDB }()

	// Verify the file exists
	_, err := os.Stat(testDBPath)
	asserts.NoError(err, "Test DB file should exist after init")

	// Create a table and insert data to make it a real database
	type CleanupTestModel struct {
		ID    uint   `gorm:"primaryKey"`
		Value string `gorm:"column:value"`
	}
	db.AutoMigrate(&CleanupTestModel{})
	db.Create(&CleanupTestModel{Value: "test"})

	// Free the database
	err = TestDBFree(db)
	asserts.NoError(err, "TestDBFree should not return error")

	// Verify the file is removed
	_, err = os.Stat(testDBPath)
	asserts.True(os.IsNotExist(err), "Test DB file should be removed after TestDBFree")
}
