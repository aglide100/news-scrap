package db

import (
	"github.com/aglide100/news-scrap/pkg/logger"
	"github.com/aglide100/news-scrap/pkg/model"
	"go.uber.org/zap"
)

func (db *Database) SaveArticle(article *model.Article) error {
	const q = `
	INSERT INTO article(title, content, url, written_date)
   		VALUES (?, ?, ?, ?)
	`
	
	_, err := db.conn.Exec(q, article.Title, article.Content, article.Url, article.Written_date)
	if err != nil {
		logger.Error("Can't insert article", zap.Error(err))
		return err
	}

	return nil
}

