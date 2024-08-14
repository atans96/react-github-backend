package service

import (
	"backend/src/types"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	mongoDB "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (f *Fetch) processContributionByRepository(ct types.StargazersNode, node []types.ContributionByRepository, key string, userInfo UserInfo) {
	for _, rep := range node {
		var contributionByRepo ContributionByRepo
		contributionByRepo.Languages = ProcessLanguages(rep.Repository.Languages.Edges)
		contributionByRepo.FullName = rep.Repository.Owner.Login + "/" + rep.Repository.Name

		opts := options.FindOneAndUpdate().SetUpsert(true)
		query := bson.M{"userName": ct.Login}
		data, err := bson.Marshal(contributionByRepo)
		if err != nil {
			panic(err)
		}

		err = bson.Unmarshal(data, &contributionByRepo)
		update := bson.D{
			{Key: "$addToSet",
				Value: bson.D{
					{Key: key, Value: contributionByRepo},
					{Key: "linkage", Value: userInfo.UserName},
				},
			},
			{Key: "$set",
				Value: bson.D{
					{Key: "bio", Value: ct.Bio},
					{Key: "id", Value: ct.Id},
					{Key: "location", Value: ct.Location},
					{Key: "isGithubStar", Value: ct.IsGithubStar},
				},
			},
		}
		err = f.Mongo.DBSuggested.FindOneAndUpdate(nil, query, update, opts).Err()
		if err != nil {
			// ErrNoDocuments means that the filter did not match any documents in the collection
			if err == mongoDB.ErrNoDocuments {
				fmt.Println(err)
			}
			fmt.Println(err)
		}
	}
}
