package githubtrending

import (
	"context"
	"fmt"
	"strings"
	"time"

	ghb "github.com/google/go-github/v32/github"
)

/* -------------------- Exported Functions -------------------- */

func FetchRepos(client *ghb.Client, settings *Settings) ([]Repo, error) {
	queryDate, err := getQueryDate(settings.since)
	if err != nil {
		return nil, err
	}

	response, err := fetch(client, settings.langs, settings.minStars, queryDate)
	if err != nil {
		return nil, err
	}

	limit := settings.limit
	if len(response) < settings.limit {
		limit = len(response)
	}

	items := response[0:limit]
	repos := make([]Repo, limit)

	for idx, repo := range items {
		repos[idx] = Repo{
			Name:        getStr(repo.Name, ""),
			FullName:    getStr(repo.FullName, ""),
			Url:         getStr(repo.HTMLURL, ""),
			Description: getStr(repo.Description, ""),
			Language:    getStr(repo.Language, ""),
			Stars:       getInt(repo.StargazersCount, 0),
			StarsToday:  -1, // Not available via GH Search API
		}
	}
	return repos, nil
}

/* -------------------- Unexported Functions -------------------- */

func getStr(strptr *string, defaultValue string) string {
	if strptr == nil {
		return defaultValue
	}
	return *strptr
}

func getInt(intptr *int, defaultValue int) int {
	if intptr == nil {
		return defaultValue
	}
	return *intptr
}

func getQueryDate(period string) (string, error) {
	currentDate := time.Now()
	switch period {
	case "daily":
		return currentDate.AddDate(0, 0, -1).Format("2006-01-02"), nil
	case "weekly":
		return currentDate.AddDate(0, 0, -7).Format("2006-01-02"), nil
	case "monthly":
		return currentDate.AddDate(0, 0, -31).Format("2006-01-02"), nil
	default:
		return "", fmt.Errorf("unknown since value: %s", period)
	}
}

func fetch(
	client *ghb.Client,
	langs []string,
	minStars int,
	since string) ([]*ghb.Repository, error) {
	var q []string
	q = append(q, fmt.Sprintf("created:>=%s", since))
	q = append(q, fmt.Sprintf("stars:>=%d", minStars))
	for _, lang := range langs {
		q = append(q, fmt.Sprintf("language:%s", lang))
	}

	repositories, _, err :=
		client.Search.Repositories(context.Background(), strings.Join(q, " "), &ghb.SearchOptions{
			Sort:  "stars",
			Order: "desc",
		})
	if err != nil {
		return nil, err
	}
	return repositories.Repositories, nil
}
