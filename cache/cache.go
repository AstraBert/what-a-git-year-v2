package cache

import (
	"context"
	"database/sql"
	"strings"

	"github.com/AstraBert/what-a-git-year-v2/cacheutils"
	"github.com/AstraBert/what-a-git-year-v2/gh"
)

type Cache interface {
	GetDb() *sql.DB
	Get(string, string) (*gh.UserStats, *gh.OrgStats, error)
	Set(string, *gh.UserStats, *gh.OrgStats) error
	Clean(string) error
}

type ApiCache struct {
	DatabaseFile string
	db           *sql.DB
}

func (c *ApiCache) GetDb() *sql.DB {
	return c.db
}

func (c *ApiCache) Get(key string, searchType string) (*gh.UserStats, *gh.OrgStats, error) {
	c.Clean("app")
	ctx := context.Background()
	queries := cacheutils.New(c.GetDb())
	val, err := queries.GetStatsByKey(ctx, key)
	if err != nil {
		return nil, nil, err
	} else {
		if searchType == "org" {
			return nil, gh.NewOrgStats(int(val.Repositories), int(val.Commits), int(val.Stars), int(val.Forks), strings.Split(val.Topics.String, ","), strings.Split(val.TopRepositories.String, ","), val.User, val.AvatarUrl), nil
		} else {
			return gh.NewUserStats(int(val.Repositories), int(val.Commits), int(val.Stars), int(val.Forks), strings.Split(val.Topics.String, ","), strings.Split(val.TopRepositories.String, ","), val.User, val.AvatarUrl), nil, nil
		}
	}
}

func (c *ApiCache) Set(key string, userStats *gh.UserStats, orgStats *gh.OrgStats) error {
	ctx := context.Background()
	queries := cacheutils.New(c.GetDb())
	if userStats != nil {
		_, err := queries.SetStats(ctx, cacheutils.SetStatsParams{K: key, User: userStats.User, AvatarUrl: userStats.AvatarUrl, Repositories: int64(userStats.Repositories), Commits: int64(userStats.Commits), Forks: int64(userStats.Forks), Stars: int64(userStats.Stars), TopRepositories: sql.NullString{Valid: true, String: strings.Join(userStats.TopRepositories, ",")}, Topics: sql.NullString{Valid: true, String: strings.Join(userStats.Topics, ",")}})
		return err
	} else {
		_, err := queries.SetStats(ctx, cacheutils.SetStatsParams{K: key, User: orgStats.Organization, AvatarUrl: orgStats.AvatarUrl, Repositories: int64(orgStats.Repositories), Commits: int64(orgStats.Commits), Forks: int64(orgStats.Forks), Stars: int64(orgStats.Stars), TopRepositories: sql.NullString{Valid: true, String: strings.Join(orgStats.TopRepositories, ",")}, Topics: sql.NullString{Valid: true, String: strings.Join(orgStats.Topics, ",")}})
		return err
	}
}

func (c *ApiCache) Clean(mode string) error {
	ctx := context.Background()
	queries := cacheutils.New(c.GetDb())
	if mode == "test" {
		return queries.CleanCacheTest(ctx)
	} else {
		return queries.CleanCache(ctx)
	}

}

func NewApiCache(dbFile string) (*ApiCache, error) {
	db, err := CreateNewDb(dbFile)
	if err != nil {
		return nil, err
	}
	return &ApiCache{
		DatabaseFile: dbFile,
		db:           db,
	}, nil
}
