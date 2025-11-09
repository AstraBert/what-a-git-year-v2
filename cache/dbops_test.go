package cache

import (
	"os"
	"testing"
)

func TestCreateDb(t *testing.T) {
	_, err := CreateNewDb("cache_test.db")
	if err != nil {
		t.Errorf("Not expecting an error when creating a new database instance, got %s", err.Error())
	}
	os.Remove("cache_test.db")
}
