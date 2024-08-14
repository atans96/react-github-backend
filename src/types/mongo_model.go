package types

import "time"

type languagePreference struct {
	Language string `json:"language" bson:"language"`
	Checked  bool   `json:"checked" bson:"checked"`
}
type renderImages struct {
	ID    int      `json:"id" bson:"id"`
	Value []string `json:"value" bson:"value"`
}
type starRanking struct {
	Id     int `json:"id" bson:"id"`
	Trends struct {
		Daily     int `json:"daily" bson:"daily"`
		Weekly    int `json:"weekly" bson:"weekly"`
		Monthly   int `json:"monthly" bson:"monthly"`
		Quarterly int `json:"quarterly" bson:"quarterly"`
		Yearly    int `json:"yearly" bson:"yearly"`
	} `json:"trends" bson:"trends"`
	TimeSeries struct {
		Daily   []int `json:"daily" bson:"daily"`
		Monthly struct {
			Year     int `json:"year" bson:"year"`
			Months   int `json:"months" bson:"months"`
			FirstDay int `json:"firstDay" bson:"firstDay"`
			LastDay  int `json:"lastDay" bson:"lastDay"`
			Delta    int `json:"delta" bson:"delta"`
		}
	} `json:"timeSeries" bson:"timeSeries"`
}
type repoInfo struct {
	FullName      string   `json:"full_name" bson:"full_name"`
	Description   string   `json:"description" bson:"description"`
	Stars         int      `json:"stars" bson:"stars"`
	Forks         int      `json:"forks" bson:"forks"`
	UpdatedAt     string   `json:"updatedAt" bson:"updatedAt"`
	Language      string   `json:"language" bson:"language"`
	Topics        []string `json:"topics" bson:"topics"`
	DefaultBranch string   `json:"default_branch" bson:"default_branch"`
	HtmlUrl       string   `json:"html_url" bson:"html_url"`
	Readme        string   `json:"readme" bson:"readme"`
}
type repoContributions struct {
	FullName     int `json:"full_name" bson:"full_name"`
	Contributors []struct {
		Login         string `json:"login" bson:"login"`
		AvatarUrl     string `json:"avatar_url" bson:"avatar_url"`
		Contributions int    `json:"contributions" bson:"contributions"`
	} `json:"contributors" bson:"contributors"`
}
type seenCards struct {
	StargazersCount int    `json:"stargazers_count" bson:"stargazers_count"`
	FullName        string `json:"full_name" bson:"full_name"`
	Owner           struct {
		Login     string `json:"login" bson:"login"`
		AvatarUrl string `json:"avatar_url" bson:"avatar_url"`
		HtmlUrl   string `json:"html_url" bson:"html_url"`
	} `json:"owner" bson:"owner"`
	Description   string   `json:"description" bson:"description"`
	Language      string   `json:"language" bson:"language"`
	Topics        []string `json:"topics" bson:"topics"`
	HtmlUrl       []string `json:"html_url" bson:"html_url"`
	Name          string   `json:"name" bson:"name"`
	Id            int      `json:"id" bson:"id"`
	DefaultBranch string   `json:"default_branch" bson:"default_branch"`
	IsQueried     bool     `json:"is_queried" bson:"is_queried"`
}
type searches struct {
	Search    string `json:"search" bson:"search"`
	Count     int    `json:"count" bson:"count"`
	UpdatedAt string `json:"updatedAt" bson:"updatedAt"`
}
type repoInfoSuggested struct {
	From            string `json:"from" bson:"from"`
	IsSeen          bool   `json:"is_seen" bson:"is_seen"`
	StargazersCount int    `json:"stargazers_count" bson:"stargazers_count"`
	FullName        string `json:"full_name" bson:"full_name"`
	DefaultBranch   string `json:"default_branch" bson:"default_branch"`
	Owner           struct {
		Login     string `json:"login" bson:"login"`
		AvatarUrl string `json:"avatar_url" bson:"avatar_url"`
		HtmlUrl   string `json:"html_url" bson:"html_url"`
	} `json:"owner" bson:"owner"`
	Description string   `json:"description" bson:"description"`
	Language    string   `json:"language" bson:"language"`
	Topics      []string `json:"topics" bson:"topics"`
	HtmlUrl     string   `json:"html_url" bson:"html_url"`
	Id          string   `json:"id" bson:"id"`
	Name        string   `json:"name" bson:"name"`
}
type clicked struct {
	IsQueried bool   `json:"is_queried" bson:"is_queried"`
	FullName  string `json:"full_name" bson:"full_name"`
	Owner     struct {
		Login string `json:"login" bson:"login"`
	} `json:"owner" bson:"owner"`
}
type UserModel struct {
	UserName           string               `json:"userName" bson:"userName"`
	Avatar             string               `json:"avatar" bson:"avatar"`
	JoinDate           time.Time            `json:"joinDate" bson:"joinDate"`
	Token              string               `json:"token" bson:"token"`
	LanguagePreference []languagePreference `json:"languagePreference" bson:"languagePreference"`
	Starred            []Starred            `json:"starred" bson:"starred"`
	RenderImages       []renderImages       `json:"renderImages" bson:"renderImages"`
	Languages          []string             `json:"languages" bson:"languages"`
	RepoContributions  []repoContributions  `json:"repoContributions" bson:"repoContributions"`
	RepoInfo           []repoInfo           `json:"repoInfo" bson:"repoInfo"`
	StarRanking        []starRanking        `json:"starRanking" bson:"starRanking"`
	SeenCards          []seenCards          `json:"seenCards" bson:"seenCards"`
	Searches           []searches           `json:"searches" bson:"searches"`
	RepoInfoSuggested  []repoInfoSuggested  `json:"repoInfoSuggested" bson:"repoInfoSuggested"`
	Clicked            []clicked            `json:"clicked" bson:"clicked"`
	Rss                []string             `json:"rss" bson:"rss"`
	LastSeen           []string             `json:"lastSeen" bson:"lastSeen"`
}
