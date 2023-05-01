package githubtrending

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"jaytaylor.com/html2text"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// The Scraper works in a following way
// - 1 query per lang, so [ts, js, python] will result in N=3 queries
// - 1 query per spoken lang so [en, zh] will result in M=2 queries per programming lang
// - Total queries would be N * M in the worst case when programming langs and spoken langs both explicitly defined
// - Aggregate and sort results by trending factor: daily stars / total stars
// - Slice results up to the specified limit

const (
	Url         = "https://github.com"
	TrendingUrl = Url + "/trending"
)

/* -------------------- Exported Functions -------------------- */

func ScrapeRepos(c *colly.Collector, settings *Settings) ([]Repo, error) {
	// Callback for when a visited HTML element is found
	since := settings.since
	limit := settings.limit
	var repos []Repo

	if len(settings.spokenLangs) > 0 {
		if len(settings.langs) > 0 {
			for _, spokenLang := range settings.spokenLangs {
				for _, lang := range settings.langs {
					res, err := scrape(c, createUrl(lang, spokenLang, since))
					if err != nil {
						return nil, err
					}
					repos = append(repos, res[:]...)
				}
			}
		} else {
			for _, spokenLang := range settings.spokenLangs {
				res, err := scrape(c, createUrl("", spokenLang, since))
				if err != nil {
					return nil, err
				}
				repos = append(repos, res[:]...)
			}
		}
	} else {
		var err error
		repos, err = scrape(c, createUrl("", "", since))
		if err != nil {
			return nil, err
		}
		return limitResults(repos, limit), nil
	}

	return process(repos, limit), nil
}

/* -------------------- Unexported Functions -------------------- */

func createUrl(lang string, spokenLang string, since string) string {
	baseUrl := TrendingUrl
	queryParams := url.Values{}
	if lang != "" {
		queryParams.Set("language", lang)
	}
	if spokenLang != "" {
		queryParams.Set("spoken_language_code", spokenLang)
	}
	if since != "" {
		queryParams.Set("since", since)
	}
	return baseUrl + "?" + queryParams.Encode()
}

func process(repos []Repo, limit int) []Repo {
	repos = sortByStars(repos)
	repos = limitResults(repos, limit)
	return repos
}

func sortByStars(repos []Repo) []Repo {
	// Sort by trending factor: StarsToday/Stars
	sort.Slice(repos, func(i, j int) bool {
		return float64(repos[i].StarsToday)/float64(repos[i].Stars) >
			float64(repos[j].StarsToday)/float64(repos[j].Stars)
	})
	return repos
}

func limitResults(repos []Repo, limit int) []Repo {
	if limit >= len(repos) {
		return repos
	}
	return repos[0:limit]
}

func scrape(c *colly.Collector, url string) ([]Repo, error) {
	var repos []Repo

	c.OnHTML("article.Box-row", func(e *colly.HTMLElement) {
		fullName := sanitize(e.ChildText("h2.h3 a"))
		name := sanitize(strings.SplitN(fullName, "/", 2)[1])
		link := sanitize(e.ChildAttr("h2.h3 a", "href"))
		desc := sanitize(e.ChildText("p"))

		var lang string
		var stars int
		var starsToday int

		nodes := e.DOM.Find("div.f6").Children().Nodes
		for i, node := range nodes {
			child := goquery.NewDocumentFromNode(node)
			text := sanitize(child.Text())
			switch i {
			case 0:
				lang = text
			case 1:
				stars, _ = extractStars(text)
			case 4:
				starsToday, _ = extractStars(text)
			}
		}
		repo := Repo{
			Name:        name,
			FullName:    fullName,
			Url:         Url + link,
			Stars:       stars,
			StarsToday:  starsToday,
			Language:    lang,
			Description: desc,
		}
		repos = append(repos, repo)
	})

	err := c.Visit(url)
	if err != nil {
		return nil, err
	}

	return repos, nil
}

func extractStars(text string) (int, error) {
	re := regexp.MustCompile(`\d+`)
	matches := re.FindAllString(text, -1)
	numStr := strings.Join(matches, "")
	num, err := strconv.Atoi(numStr)
	if err != nil {
		return -1, err
	}
	return num, nil
}

func sanitize(text string) string {
	stripped, _ := html2text.FromString(text)
	return strings.TrimSpace(stripped)
}
