package main

import (
	"log"
	"time"

	"github.com/aglide100/news-scrap/pkg/db"
	"github.com/aglide100/news-scrap/pkg/logger"
	"github.com/aglide100/news-scrap/pkg/scrap"
	"go.uber.org/zap"
)

// var (
// 	target = flag.String("target", "https://news.naver.com/section/102", "target url")
// 	area = flag.String("area", "li.rl_item._LAZY_LOADING_WRAP", "list area for watching")
// 	article_title = flag.String("article_title", ".media_end_head_title", "article area")
// 	article_date = flag.String("article_date", "._ARTICLE_DATE_TIME", "article date")
// 	article_content = flag.String("article_ct", "#dic_area", "article area")
// 	comment = flag.String("comment", "", "comment area")
// )



func main() {
	db, err := db.NewDB()
	if err != nil {
		log.Fatalf(err.Error())
	}

	duration, _ := time.ParseDuration("20000s")
	ticker := time.NewTicker(duration)
	defer ticker.Stop()
	for range ticker.C {
		_, _, err := scrap.Scrap(db)
		if err != nil {
			logger.Info("err", zap.Any("err", err))
		}
		// for _, article := range articles {
		// 	docs := strings.Split(article.Content, ".")
	
	
		// 	keysents := textrank.TextRankSentences(docs, 2, 0.85, 30, 3)
	
		// 	for _, keysent := range keysents {
		// 		fmt.Printf("Score: %.4f, Sentence: %s\n", keysent.Score, keysent.Sentence)
		// 	}
	
		// }
	}
}

