package models

// Define models for the external API response
type ArticleResponse struct {
	PageInfo struct {
		Page       int `json:"page"`
		NumPages   int `json:"numPages"`
		PageSize   int `json:"pageSize"`
		NumEntries int `json:"numEntries"`
	} `json:"pageInfo"`
	Content []Article `json:"content"`
}

type Article struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Date        string `json:"date"`
	Body        string `json:"body"`
	Summary     string `json:"summary"`
	LeadMedia   Media  `json:"leadMedia"`
	Tags        []Tag  `json:"tags"`
	// Add other fields as needed
}

type Media struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Type  string `json:"type"`
	// Add image URLs/variants as needed
}

type Tag struct {
	ID    int    `json:"id"`
	Label string `json:"label"`
}
