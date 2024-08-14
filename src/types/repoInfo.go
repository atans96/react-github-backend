package types

import "time"

type RepoInfo struct {
	ID       int    `json:"id" bson:"id"`
	NodeID   string `json:"node_id" bson:"node_id"`
	Name     string `json:"name" bson:"name"`
	FullName string `json:"full_name" bson:"full_name"`
	Owner    struct {
		Login      string `json:"login" bson:"login"`
		ID         int    `json:"id" bson:"id"`
		NodeID     string `json:"node_id" bson:"node_id"`
		AvatarURL  string `json:"avatar_url" bson:"avatar_url"`
		GravatarID string `json:"gravatar_id" bson:"gravatar_id"`
		URL        string `json:"url" bson:"url"`
	} `json:"owner" bson:"owner"`
	HTMLURL         string      `json:"html_url" bson:"html_url"`
	Forks           int         `json:"forks" bson:"forks"`
	Topics          []string    `json:"topics" bson:"topics"`
	Description     interface{} `json:"description" bson:"description"`
	URL             string      `json:"url" bson:"url"`
	CreatedAt       time.Time   `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at" bson:"updated_at"`
	Homepage        interface{} `json:"homepage" bson:"homepage"`
	StargazersCount int         `json:"stargazers_count" bson:"stargazers_count"`
	WatchersCount   int         `json:"watchers_count" bson:"watchers_count"`
	Language        string      `json:"language" bson:"language"`
	Watchers        int         `json:"watchers" bson:"watchers"`
	DefaultBranch   string      `json:"default_branch" bson:"default_branch"`
}
type RepoInfos []RepoInfo
type ReturnRepoInfo struct {
	FullName      string      `json:"fullName" bson:"fullName"`
	Description   interface{} `json:"description" bson:"description"`
	Stars         int         `json:"stars" bson:"stars"`
	Forks         int         `json:"forks" bson:"forks"`
	UpdatedAt     time.Time   `json:"updatedAt" bson:"updatedAt"`
	Language      string      `json:"language" bson:"language"`
	Topics        []string    `json:"topics" bson:"topics"`
	DefaultBranch string      `json:"defaultBranch" bson:"defaultBranch"`
	HtmlUrl       string      `json:"html_url" bson:"html_url"`
	Readme        string      `json:"readme" bson:"readme"`
}
