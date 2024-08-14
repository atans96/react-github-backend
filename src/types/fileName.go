package types

type Tree struct {
	Path string `json:"path" bson:"path"`
	Type string `json:"type" bson:"type"`
	Sha  string `json:"sha" bson:"sha"`
	URL  string `json:"url" bson:"url"`
}

type FileNameResponse struct {
	Tree []Tree `json:"tree" bson:"tree"`
}
