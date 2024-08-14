package types

type ApiRateLimitRule struct {
	Limit     int   `json:"limit" bson:"limit"`
	Remaining int   `json:"remaining" bson:"remaining"`
	Reset     int64 `json:"reset" bson:"reset"`
}

type ApiRateLimitResource struct {
	Core                ApiRateLimitRule `json:"core" bson:"core"`
	Search              ApiRateLimitRule `json:"search" bson:"search"`
	GraphQL             ApiRateLimitRule `json:"graphql" bson:"graph_ql"`
	IntegrationManifest ApiRateLimitRule `json:"integration_manifest" bson:"integration_manifest"`
}
type ApiRateLimit struct {
	Resources ApiRateLimitResource `json:"resources" bson:"resources"`
	Rate      ApiRateLimitRule     `json:"rate" bson:"rate"`
}
