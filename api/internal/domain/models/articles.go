package models

type ArticleResponse struct {
	Content []Article `json:"content"`
}

type UpsertArticle struct {
	ExternalID int `json:"externalID"`
	Article
}

type Article struct {
	ExternalID  int    `json:"externalID"`
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Date        string `json:"date"`
	Body        string `json:"body"`
	Summary     string `json:"summary"`
	LeadMedia   Media  `json:"leadMedia"`
	Tags        []Tag  `json:"tags"`
}

type Media struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Type  string `json:"type"`
}

type Tag struct {
	ID    int    `json:"id"`
	Label string `json:"label"`
}

type PaginatedArticles struct {
	PageInfo PageInfo  `json:"pageInfo"`
	Content  []Article `json:"content"`
}

type PageInfo struct {
	Page       int `json:"page"`
	NumPages   int `json:"numPages"`
	PageSize   int `json:"pageSize"`
	NumEntries int `json:"numEntries"`
}
