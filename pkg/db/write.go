package db

import (
	"errors"

	"github.com/aglide100/news-scrap/pkg/logger"
	"github.com/aglide100/news-scrap/pkg/model"
	"go.uber.org/zap"
)

func (db *Database) SaveArticle(article *model.Article) error {
	const q1 = `
	SELECT title FROM article WHERE url = ?
	`
	
	var title string

	err := db.conn.QueryRow(q1, article.Url).Scan(&title)
	if err.Error() != "sql: no rows in result set" {
		return err
	}

	if len(title) > 1 {
		return errors.New("article is already exist")
	}

	const q2 = `
	INSERT INTO article(title, content, url, written_date)
   		VALUES (?, ?, ?, ?)
	`
	
	_, err = db.conn.Exec(q2, article.Title, article.Content, article.Url, article.Written_date)
	if err != nil {
		logger.Error("Can't insert article", zap.Error(err))
		return err
	}

	return nil
}

func (db *Database) SaveComment(comment *model.Comment) error {
	const q1 = `
	SELECT title FROM article WHERE url = ?
	`
	
	title := ""

	err := db.conn.QueryRow(q1, comment.ArticleURL).Scan(&title)
	if err != nil {
		return err
	}

	if len(title) <= 1 {
		return errors.New("can't find article")
	}

	const q2 = `
	INSERT INTO comment(content, like_cnt, dislike_cnt, articleURL, written_date)
   		VALUES (?, ?, ?, ?, ?)
	`
	
	_, err = db.conn.Exec(q2, comment.Content, comment.Like, comment.Dislike, comment.ArticleURL, comment.Written_date)
	if err != nil {
		logger.Error("Can't insert comment", zap.Error(err))
		return err
	}

	return nil
}


