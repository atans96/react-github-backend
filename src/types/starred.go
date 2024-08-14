package types

type Starred struct {
	FullName  string `json:"full_name" bson:"full_name"`
	IsQueried bool   `json:"is_queried" bson:"is_queried"`
}
