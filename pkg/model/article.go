package model


type Article struct {
	Title string
	Written_date string
	Content string
	Url string
}

type Comment struct {
	Content string
	Like int
	Dislike int
	Written_date string
	ArticleURL string
}
