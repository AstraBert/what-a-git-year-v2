package cache

import (
	"os"
	"slices"
	"testing"
	"time"

	"github.com/AstraBert/what-a-git-year-v2/gh"
)

func TestNewApiCache(t *testing.T) {
	apiCache, err := NewApiCache("cache_test.db")
	if err != nil {
		t.Errorf("Not expecting errors, got %s", err.Error())
	}
	_, isType := any(apiCache).(Cache)
	if !isType {
		t.Error("Expecting ApiCache to implement Cache, but it does not")
	}
}

func TestApiCacheGetSet(t *testing.T) {
	apiCache, _ := NewApiCache("cache_test.db")
	err := apiCache.Set("hello", gh.NewUserStats(1, 1, 1, 1, []string{"hello", "ai"}, []string{"no-repo", "repo"}, "test", "not a url"), nil)
	if err != nil {
		t.Errorf("Not expecting errors, got %s", err.Error())
	}
	err = apiCache.Set("hello", gh.NewUserStats(1, 1, 1, 1, []string{"hello", "ai"}, []string{"no-repo", "repo"}, "test", "not a url"), nil)
	if err == nil {
		t.Errorf("Expecting errors related to the uniqueness of the key, got none")
	}
	userStats, orgStats, err := apiCache.Get("hello", "user")
	if err != nil {
		t.Errorf("Expecting no error, got %s", err.Error())
	}
	if orgStats != nil {
		t.Errorf("Expecting org stats to be null, got %v", orgStats)
	}
	if userStats.Repositories != 1 || userStats.Commits != 1 || userStats.Forks != 1 || userStats.Stars != 1 || userStats.User != "test" || userStats.AvatarUrl != "not a url" || !slices.Equal(userStats.TopRepositories, []string{"no-repo", "repo"}) || !slices.Equal(userStats.Topics, []string{"hello", "ai"}) {
		t.Error("Expected retrieved stats to be the same as inserted ones, got different ones")
	}
	_, _, err = apiCache.Get("hello1", "user")
	if err == nil {
		t.Error("Expecting an error related to the unexistent key, got none")
	}
	userStats, orgStats, err = apiCache.Get("hello", "org")
	if err != nil {
		t.Errorf("Expecting no error, got %s", err.Error())
	}
	if userStats != nil {
		t.Errorf("Expecting user stats to be null, got %v", orgStats)
	}
	if orgStats.Repositories != 1 || orgStats.Commits != 1 || orgStats.Forks != 1 || orgStats.Stars != 1 || orgStats.Organization != "test" || orgStats.AvatarUrl != "not a url" || !slices.Equal(orgStats.TopRepositories, []string{"no-repo", "repo"}) || !slices.Equal(orgStats.Topics, []string{"hello", "ai"}) {
		t.Error("Expected retrieved stats to be the same as inserted ones, got different ones")
	}
	_ = os.Remove("cache_test.db")
}

func TestApiCacheClean(t *testing.T) {
	apiCache, _ := NewApiCache("cache_test.db")
	_ = apiCache.Set("hello", gh.NewUserStats(1, 1, 1, 1, []string{"hello", "ai"}, []string{"no-repo", "repo"}, "test", "not a url"), nil)
	time.Sleep(2 * time.Second)
	err := apiCache.Clean("test")
	if err != nil {
		t.Errorf("Not expecting any error, got %s", err.Error())
	}
	_, _, err = apiCache.Get("hello", "user")
	if err == nil {
		t.Error("Expecting an error related to the unexistent key, got none")
	}
	_ = os.Remove("cache_test.db")
}
