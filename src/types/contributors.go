package types

type ContributorResponse struct {
	Login         string `json:"login" bson:"login"`
	AvatarUrl     string `json:"avatar_url" bson:"avatar_url"`
	Contributions int    `json:"contributions" bson:"contributions"`
}

type ReturnContributors struct {
	FullName     string                 `json:"fullName" bson:"fullName"`
	Contributors *[]ContributorResponse `json:"contributors" bson:"contributors"`
}
