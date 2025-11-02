package gh

import (
	"context"
	"sort"
	"time"

	"github.com/google/go-github/v76/github"
)

type UserStats struct {
	User            string
	Repositories    int
	TopRepositories []string
	Commits         int
	Stars           int
	Forks           int
	Topics          []string
}

type OrgStats struct {
	Organization    string
	Repositories    int
	TopRepositories []string
	Commits         int
	Stars           int
	Forks           int
	Topics          []string
}

func NewUserStats(repos, commits, stars, forks int, topics, topRepos []string, user string) *UserStats {
	return &UserStats{
		User:            user,
		Repositories:    repos,
		Commits:         commits,
		Stars:           stars,
		Forks:           forks,
		TopRepositories: topRepos,
		Topics:          topics,
	}
}

func NewOrgStats(repos, commits, stars, forks int, topics, topRepos []string, organization string) *OrgStats {
	return &OrgStats{
		Organization:    organization,
		Repositories:    repos,
		Commits:         commits,
		Stars:           stars,
		Forks:           forks,
		TopRepositories: topRepos,
		Topics:          topics,
	}
}

type GitHubClient interface {
	GetClient() *github.Client
	GetUserStats(string) (*UserStats, error)
	GetOrgStats(string) (*OrgStats, error)
}

type GitYearClient struct {
	AuthToken string
}

func NewGitYearClient(token string) *GitYearClient {
	return &GitYearClient{
		AuthToken: token,
	}
}

func (c *GitYearClient) GetClient() *github.Client {
	return github.NewClient(nil).WithAuthToken(c.AuthToken)
}

func chooseTopTen(dict map[string]int) []string {
	keys := make([]string, 0, len(dict))
	for k := range dict {
		keys = append(keys, k)
	}
	sort.SliceStable(keys, func(i, j int) bool {
		return dict[keys[i]] > dict[keys[j]]
	})
	if len(keys) > 10 {
		return keys[:10]
	} else {
		return keys
	}
}

func isWithinYear(timestamp time.Time) bool {
	now := time.Now()
	oneYearAgo := now.AddDate(-1, 0, 0)
	tomorrow := now.AddDate(0, 0, 1)

	return timestamp.After(oneYearAgo) && timestamp.Before(tomorrow)
}

func (c *GitYearClient) GetUserStats(user string) (*UserStats, error) {
	ghClient := c.GetClient()
	opt := &github.RepositoryListByUserOptions{Type: "public", ListOptions: github.ListOptions{Page: 0, PerPage: 100}, Sort: "stargazers_count", Direction: "desc"}
	repos, _, err := ghClient.Repositories.ListByUser(context.Background(), user, opt)
	if err != nil {
		return nil, err
	} else {
		repoStars := map[string]int{}
		totalRepos := 0
		totalStars := 0
		topics := map[string]int{}
		totalForks := 0
		totalCommits := 0
		for _, repo := range repos {
			if repo.CreatedAt != nil && isWithinYear(*repo.CreatedAt.GetTime()) && repo.Owner != nil && *repo.Owner.Login == user {
				if repo.FullName != nil {
					totalRepos += 1
					repoStars[*repo.FullName] = *repo.StargazersCount
					totalStars += *repo.StargazersCount
					totalForks += *repo.ForksCount
					now := time.Now()
					oneYearAgo := now.AddDate(-1, 0, 0)
					tomorrow := now.AddDate(0, 0, 1)
					optsCommits := &github.CommitsListOptions{
						ListOptions: github.ListOptions{PerPage: 200},
						Since:       oneYearAgo,
						Until:       tomorrow,
					}
					repoCommits, _, err := ghClient.Repositories.ListCommits(context.Background(), user, *repo.FullName, optsCommits)
					if err == nil {
						totalCommits += len(repoCommits)
					}
					repoTopics := repo.Topics
					for _, topic := range repoTopics {
						val, ok := topics[topic]
						if !ok {
							topics[topic] = 1
						} else {
							topics[topic] = val + 1
						}
					}
				}
			}
		}
		return NewUserStats(totalRepos, totalCommits, totalStars, totalForks, chooseTopTen(topics), chooseTopTen(repoStars), user), nil
	}
}

func (c *GitYearClient) GetOrgStats(organization string) (*OrgStats, error) {
	ghClient := c.GetClient()
	opts := &github.RepositoryListByOrgOptions{Type: "public", ListOptions: github.ListOptions{PerPage: 100}, Sort: "stargazers_count", Direction: "desc"}
	repos, _, err := ghClient.Repositories.ListByOrg(context.Background(), organization, opts)
	if err != nil {
		return nil, err
	} else {
		repoStars := map[string]int{}
		totalRepos := 0
		totalStars := 0
		topics := map[string]int{}
		totalForks := 0
		totalCommits := 0
		for _, repo := range repos {
			if repo.CreatedAt != nil && isWithinYear(*repo.CreatedAt.GetTime()) && repo.Owner != nil && *repo.Owner.Login == organization {
				if repo.FullName != nil {
					totalRepos += 1
					repoStars[*repo.FullName] = *repo.StargazersCount
					totalStars += *repo.StargazersCount
					totalForks += *repo.ForksCount
					now := time.Now()
					oneYearAgo := now.AddDate(-1, 0, 0)
					tomorrow := now.AddDate(0, 0, 1)
					optsCommits := &github.CommitsListOptions{
						ListOptions: github.ListOptions{PerPage: 200},
						Since:       oneYearAgo,
						Until:       tomorrow,
					}
					repoCommits, _, err := ghClient.Repositories.ListCommits(context.Background(), organization, *repo.FullName, optsCommits)
					if err == nil {
						totalCommits += len(repoCommits)
					}
					repoTopics := repo.Topics
					for _, topic := range repoTopics {
						val, ok := topics[topic]
						if !ok {
							topics[topic] = 1
						} else {
							topics[topic] = val + 1
						}
					}
				}
			}
		}
		return NewOrgStats(totalRepos, totalCommits, totalStars, totalForks, chooseTopTen(topics), chooseTopTen(repoStars), organization), nil
	}
}
