package scrap

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/aglide100/news-scrap/pkg/db"
	"github.com/aglide100/news-scrap/pkg/logger"
	"github.com/aglide100/news-scrap/pkg/model"
	"github.com/aglide100/news-scrap/pkg/textrank"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
)


var (
	target = flag.String("target", "https://news.naver.com/section/102", "target url")

	article_area = flag.String("article_area", "li.rl_item._LAZY_LOADING_WRAP", "list area for watching")
	article_title = flag.String("article_title", ".media_end_head_title", "article area")
	article_date = flag.String("article_date", "._ARTICLE_DATE_TIME", "article date")
	article_content = flag.String("article_ct", "#dic_area", "article area")
	
	comment_link = flag.String("comment_link", "https://apis.naver.com/commentBox/cbox/web_naver_list_jsonp.json?ticket=news&templateId=default_politics_m3&pool=cbox5&lang=ko&country=KR&categoryId=&pageSize=20&indexSize=10&groupId=&listType=OBJECT&pageType=more&page=1&initialize=true&followSize=5&userType=&useAltSort=true&replyPageSize=20&sort=favorite&includeAllStatus=true&objectId=news", "comment link")
	comment_area = flag.String("comment_area", "result.commentList", "comment area")
	comment_content = flag.String("comment_content", "contents", "content")
	comment_like = flag.String("comment_like", "sympathyCount", "comment_like")
	comment_dislike = flag.String("comment_dislike", "antipathyCount", "comment_dislike")
	comment_date = flag.String("comment_date", "regTime", "comment_date")
)

func Scrap(db *db.Database) ([]*model.Article, []*model.Comment, error) {
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
	doc.Find(*article_area).Each(func(i int, s *goquery.Selection) {
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

		article_doc, err := goquery.NewDocumentFromReader(body)
		if err != nil {
			logger.Error(err.Error())
		}
		
		var title, date, content string

		article_doc.Find(*article_title).Each(func(i int, s *goquery.Selection) {
			if (s.Text() != "" ) {
				title += s.Text() + " "
			}
		})
	
		article_doc.Find(*article_date).Each(func(i int, s *goquery.Selection) {
			if s.Text() != "" {
				date += s.Text() + " "
			}
		})
	
		article_doc.Find(*article_content).Each(func(i int, s *goquery.Selection) {
			if s.Text() != "" {
				content += s.Text() + " "
			}
		})

		title = Preprocessing(title)

		date = Preprocessing(date)

		content = Preprocessing(content)

		// logger.Info("article_doc", zap.Any("doc", Preprocessing(article_doc.Text())))
		if len(title) > 5 && len(content) > 5 {
			newArticle := &model.Article{
				Title: title,
				Written_date: date,
				Content: content,
				Url: link,
			}

			logger.Info("link", zap.Any("link", link))
			logger.Info("title", zap.Any("title", title))
			logger.Info("date", zap.Any("date", date))
			logger.Info("content", zap.Any("content", content))

			
			err = db.SaveArticle(newArticle)
			if err != nil {
				log.Printf("can't save article %s", err)
				continue
			}
			
			p1, p2 := extractID(link)
			commentLink := *comment_link + p1 + "%2C" +p2 
			
			body, err := CreateHttpReqWithReferer(commentLink, link)
			if err != nil {
				logger.Error(err.Error())
			}

			raw := extractJson(body)
			vals := gjson.Get(raw, *comment_area)
			
			for _, val := range vals.Array() {
				content := val.Get(*comment_content)
				like := val.Get(*comment_like)
				like_cnt, err := strconv.Atoi(like.String())
				if err != nil {
					log.Printf("can't convert to int %s", err)
				}

				dislike := val.Get(*comment_dislike)
				dislike_cnt, err := strconv.Atoi(dislike.String())
				if err != nil {
					log.Printf("can't convert to int %s", err)
				}

				date := val.Get(*comment_date)

				newComment := &model.Comment{
					Content: content.String(),
					Like: like_cnt,
					Dislike: dislike_cnt,
					Written_date: date.String(),
					ArticleURL: link,
				}
				comments = append(comments, newComment)


				err = db.SaveComment(newComment)
				if err != nil {
					log.Printf("can't save comment %s", err)
				}
				
			}
			
			logger.Info("comments", zap.Any("len", len(comments)))

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

func extractID(url string) (string, string) {
	parts := strings.Split(url, "/")
	p1 := parts[len(parts)-2]
	p2 := parts[len(parts)-1]

	if idx := strings.Index(p2, "?"); idx != -1 {
		p2 = p2[:idx]
	}

	return p1, p2
}

func extractJson(text string) string {
	front := 0
	rear := 0

	for i := 0; i < len(text); i++ {
		if (text[i] == '{') {
			front = i
			break
		}
	}

	for i := len(text)-1; i>=0; i-- {
		if (text[i] == '}') {
			rear = i
			break
		}
	}

	return text[front:rear+1]
}
