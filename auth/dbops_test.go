package auth

import "testing"

func TestCreateDb(t *testing.T) {
	_, err := CreateNewDb()
	if err != nil {
		t.Errorf("Not expecting an error when creating a new database instance, got %s", err.Error())
	}
}
