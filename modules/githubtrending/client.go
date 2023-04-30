package githubtrending

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	ghb "github.com/google/go-github/v32/github"
)

/* -------------------- Exported Functions -------------------- */

func GetRepos(client *ghb.Client, settings *Settings) ([]Repo, error) {
	queryDate, err := getQueryDate(settings.period)
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
	case "day":
		return currentDate.AddDate(0, 0, -1).Format("2006-01-02"), nil
	case "week":
		return currentDate.AddDate(0, 0, -7).Format("2006-01-02"), nil
	case "month":
		return currentDate.AddDate(0, 0, -31).Format("2006-01-02"), nil
	default:
		// try parse value as int in case there is a custom period
		i, err := strconv.Atoi(period)
		if err != nil {
			return "", fmt.Errorf("unknown period value: %s", period)
		}
		return currentDate.AddDate(0, 0, -i).Format("2006-01-02"), nil
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
