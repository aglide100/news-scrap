package scrap

import (
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/aglide100/news-scrap/pkg/logger"
	"github.com/aglide100/news-scrap/pkg/model"
	"github.com/aglide100/news-scrap/pkg/textrank"
	"go.uber.org/zap"
)


var (
	target = flag.String("target", "https://news.naver.com/section/102", "target url")
	area = flag.String("area", "li.rl_item._LAZY_LOADING_WRAP", "list area for watching")
	article_title = flag.String("article_title", ".media_end_head_title", "article area")
	article_date = flag.String("article_date", "._ARTICLE_DATE_TIME", "article date")
	article_content = flag.String("article_ct", "#dic_area", "article area")
	comment = flag.String("comment", "", "comment area")
)

func Scrap() ([]*model.Article, []*model.Comment, error) {
	flag.Parse()

	body, err := CreateHttpReq(*target)
	if err != nil {
		logger.Error(err.Error())
	}
	
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		logger.Error(err.Error())
	}

	links := []string{}
	doc.Find(*area).Each(func(i int, s *goquery.Selection) {
		s.Find("a").Each(func(j int, s2 *goquery.Selection) {
			href, exists := s2.Attr("href")
			if exists && strings.HasPrefix(href, "http") {
				links = append(links, href)
				fmt.Println(href)
			}
		})
	})

	articles := []*model.Article{}
	comments := []*model.Comment{}

	for _, link := range links {
		body, err := CreateHttpReq(link)
		if err != nil {
			logger.Error(err.Error())
		}

		doc, err := goquery.NewDocumentFromReader(body)
		if err != nil {
			logger.Error(err.Error())
		}
		
		var title, date, content string

		doc.Find(*article_title).Each(func(i int, s *goquery.Selection) {
			if (s.Text() != "" ) {
				title += s.Text() + " "
			}
		})
	
		doc.Find(*article_date).Each(func(i int, s *goquery.Selection) {
			if s.Text() != "" {
				date += s.Text() + " "
			}
		})
	
		doc.Find(*article_content).Each(func(i int, s *goquery.Selection) {
			if s.Text() != "" {
				content += s.Text() + " "
			}
		})

		title = Preprocessing(title)

		date = Preprocessing(date)

		content = Preprocessing(content)

		logger.Info("title", zap.Any("title", title))
		logger.Info("date", zap.Any("date", date))
		logger.Info("content", zap.Any("content", content))

		if len(title) > 5 && len(content) > 5 {
			newArticle := &model.Article{
				Title: title,
				Written_date: date,
				Content: content,
				Url: link,
			}

			articles = append(articles, newArticle)

			docs := strings.Split(newArticle.Content, ".")
			keysents := textrank.TextRankSentences(docs, 2, 0.85, 30, 4)

			for _, keysent := range keysents {
				fmt.Printf("Score: %.4f, Sentence: %s\n", keysent.Score, keysent.Sentence)
			}
		}

		time.Sleep(5 * time.Second)
	}

	defer body.Close()

	return articles, comments, nil
}