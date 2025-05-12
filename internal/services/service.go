package services

import (
	"regexp"
)

var namesRelation = map[string]string{
	"habr":     "Habr",
	"russian":  "Russia-Today",
	"lenta":    "Lenta-ru",
	"nytimes":  "New-York-Times",
	"research": "Research-swtch",
}

func getNameSite(url string) string {
	re := regexp.MustCompile(`//(?:www\.)?([^./]+)`)
	match := re.FindStringSubmatch(url)
	return namesRelation[match[1]]
}
