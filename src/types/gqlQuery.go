package types

type User struct {
	Id           string `json:"id" bson:"id"`
	Login        string `json:"login" bson:"login"`
	Bio          string `json:"bio" bson:"bio"`
	IsGithubStar bool   `json:"isGithubStar" bson:"isGithubStar"`
	Location     string `json:"location" bson:"location"`
}
type PageInfo struct {
	EndCursor   string `json:"endCursor" bson:"endCursor"`
	HasNextPage bool   `json:"hasNextPage" bson:"hasNextPage"`
	StartCursor string `json:"startCursor" bson:"startCursor"`
}
type LanguagesEdges struct {
	Node struct {
		Name string `json:"name" bson:"name"`
	} `json:"node" bson:"node"`
	Size int `json:"size" bson:"size"`
}
type Languages struct {
	Edges []LanguagesEdges `json:"edges" bson:"edges"`
}
type Reactions struct {
	Content string `json:"content" bson:"content"`
	User    User   `json:"user" bson:"user"`
}
type RepositoryDiscussionsNodes struct {
	UpvoteCount      int  `json:"upvoteCount" bson:"upvoteCount"`
	ViewerHasUpvoted bool `json:"viewerHasUpvoted" bson:"viewerHasUpvoted"`
	Reactions        struct {
		PageInfo   PageInfo    `json:"pageInfo" bson:"pageInfo"`
		TotalCount int         `json:"totalCount" bson:"totalCount"`
		Nodes      []Reactions `json:"nodes" bson:"nodes"`
	} `json:"reactions" bson:"reactions"`
	Discussion struct {
		Category struct {
			Name string `json:"name" bson:"name"`
		} `json:"category" bson:"category"`
	} `json:"discussion" bson:"discussion"`
	Repository struct {
		Name      string    `json:"name" bson:"name"`
		Languages Languages `json:"languages" bson:"languages"`
	} `json:"repository" bson:"repository"`
}
type RepositoryDiscussions struct {
	PageInfo PageInfo                     `json:"pageInfo" bson:"pageInfo"`
	Nodes    []RepositoryDiscussionsNodes `json:"nodes" bson:"nodes"`
}
type StarredRepositoriesNodes struct {
	Languages  Languages `json:"languages" bson:"languages"`
	Stargazers struct {
		PageInfo PageInfo         `json:"pageInfo" bson:"pageInfo"`
		Nodes    []StargazersNode `json:"nodes" bson:"nodes"`
	} `json:"stargazers" bson:"stargazers"`
}
type StarredRepositories struct {
	PageInfo PageInfo                   `json:"pageInfo" bson:"pageInfo"`
	Nodes    []StarredRepositoriesNodes `json:"nodes" bson:"nodes"`
}
type Repository struct {
	DefaultBranchRef struct {
		Name string `json:"name" bson:"name"`
	} `bson:"defaultBranchRef" json:"defaultBranchRef"`
	Name  string `json:"name" bson:"name"`
	Owner struct {
		Id    string `json:"id" bson:"id"`
		Login string `json:"login" bson:"login"`
	} `json:"owner" bson:"owner"`
	Languages Languages `json:"languages" bson:"languages"`
}
type Popular struct {
	Repository   Repository `json:"repository" bson:"repository"`
	Participants struct {
		PageInfo PageInfo         `json:"pageInfo" bson:"pageInfo"`
		Nodes    []StargazersNode `json:"nodes" bson:"nodes"`
	} `json:"participants" bson:"participants"`
}
type ContributionByRepository struct {
	Repository Repository `json:"repository" bson:"repository"`
}
type ContributionsCollection struct {
	TotalIssueContributions             int `json:"totalIssueContributions" bson:"totalIssueContributions"`
	TotalPullRequestContributions       int `json:"totalPullRequestContributions" bson:"totalPullRequestContributions"`
	TotalPullRequestReviewContributions int `json:"totalPullRequestReviewContributions" bson:"totalPullRequestReviewContributions"`
	PopularIssueContribution            struct {
		Issue Popular `json:"issue" bson:"issue"`
	} `json:"popularIssueContribution" bson:"popularIssueContribution"`
	PopularPullRequestContribution struct {
		PullRequest Popular `json:"pullRequest" bson:"pullRequest"`
	} `json:"popularPullRequestContribution" bson:"popularPullRequestContribution"`
	PullRequestReviewContributionsByRepository []ContributionByRepository `json:"pullRequestReviewContributionsByRepository" bson:"pullRequestReviewContributionsByRepository"`
	PullRequestContributionsByRepository       []ContributionByRepository `json:"pullRequestContributionsByRepository" bson:"pullRequestContributionsByRepository"`
}
type StargazersNode struct {
	Id                           string                  `json:"id" bson:"id"`
	Login                        string                  `json:"login" bson:"login"`
	Bio                          string                  `json:"bio" bson:"bio"`
	IsGithubStar                 bool                    `json:"isGithubStar" bson:"isGithubStar"`
	Location                     string                  `json:"location" bson:"location"`
	RepositoryDiscussions        RepositoryDiscussions   `json:"repositoryDiscussions" bson:"repositoryDiscussions"`
	RepositoryDiscussionComments RepositoryDiscussions   `json:"repositoryDiscussionComments" bson:"repositoryDiscussionComments"`
	StarredRepositories          StarredRepositories     `json:"starredRepositories" bson:"starredRepositories"`
	ContributionsCollection      ContributionsCollection `json:"contributionsCollection" bson:"contributionsCollection"`
}
type GraphqlQueryResponse struct {
	Repository struct {
		Stargazers struct {
			Nodes    []StargazersNode `json:"nodes" bson:"nodes"`
			PageInfo PageInfo         `json:"pageInfo" bson:"pageInfo"`
		} `json:"stargazers" bson:"stargazers"`
	} `json:"repository" bson:"repository"`
}
type GraphqlQueryUserResponse struct {
	User struct {
		Id                           string                  `json:"id" bson:"id"`
		Login                        string                  `json:"login" bson:"login"`
		Bio                          string                  `json:"bio" bson:"bio"`
		IsGithubStar                 bool                    `json:"isGithubStar" bson:"isGithubStar"`
		Location                     string                  `json:"location" bson:"location"`
		RepositoryDiscussions        RepositoryDiscussions   `json:"repositoryDiscussions" bson:"repositoryDiscussions"`
		RepositoryDiscussionComments RepositoryDiscussions   `json:"repositoryDiscussionComments" bson:"repositoryDiscussionComments"`
		StarredRepositories          StarredRepositories     `json:"starredRepositories" bson:"starredRepositories"`
		ContributionsCollection      ContributionsCollection `json:"contributionsCollection" bson:"contributionsCollection"`
	} `json:"user" bson:"user"`
}
