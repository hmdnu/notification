package collector

import (
	"fmt"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/hmdnu/bot/internal/client"
)

const (
	subjectCard = ".gallery_grid_item"
	anchorTag   = "a"
)

type Subjects struct {
	Name string
	Link string
}

func Collector() {
	// get subjects
	fmt.Println()
	fmt.Println("Fetching subjects")
	html, err := client.FetchSubjectContent()

	if err != nil {
		log.Fatalln("Failed fetching subjects", err.Error())
	}

	fmt.Println("Get subjects information")
	subjects, err := getSubjects(html)

	if err != nil {
		log.Fatalln("Failed getting subjects information", err.Error())
	}

	fmt.Println("Get LMS contents")
	content, err := getContent(subjects)

	if err != nil {
		log.Fatalln("Failed getting lms content", err.Error())
	}

	fmt.Println(content)
}

func getSubjects(html string) ([]Subjects, error) {

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))

	if err != nil {
		return nil, err
	}

	var subjects []Subjects

	doc.Find(subjectCard).Each(func(i int, s *goquery.Selection) {
		title, exist := s.Attr("title")

		if !exist {
			fmt.Println("title not found")
		}

		href, exist := s.Find(anchorTag).Attr("href")

		if !exist {
			fmt.Println("href not found")
		}

		subjects = append(subjects, Subjects{Name: title, Link: href})
	})

	return subjects, nil
}

func getContent(subjects []Subjects) ([]string, error) {
	var titles []string

	for _, subject := range subjects {
		html, err := client.FetchLmsContent(subject.Link)

		if err != nil {
			return nil, err
		}
		// do dom extraction
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))

		if err != nil {
			return nil, err
		}

		title := doc.Find("title").Text()

		titles = append(titles, title)
	}

	return titles, nil
}
