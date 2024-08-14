package service

import (
	"backend/src/types"
	"context"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/json-iterator/go"
	"github.com/machinebox/graphql"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongoDB "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
)

type mongo struct {
	DB          *mongoDB.Collection
	DBSuggested *mongoDB.Collection
	Session     *mongoDB.Client
}

var Mongo mongo //global variable to remember internal state of struct
var json = jsoniter.ConfigCompatibleWithStandardLibrary
var GQLClient = graphql.NewClient("https://api.github.com/graphql")
var cfg = elasticsearch.Config{
	Addresses: []string{
		"http://" + os.Getenv("ES_HOST") + ":" + os.Getenv("ES_PORT"),
	},
	Transport: &FastHTTPTransport{},
}
var ESClient, esErr = elasticsearch.NewClient(cfg)

type Fetch struct {
	Mongo      mongo
	collection *mongoDB.Collection
}
type UserInfo struct {
	UserName string `json:"userName" bson:"userName"`
	Token    string `json:"token" bson:"token"`
	FullName string `json:"fullName" bson:"fullName"`
}

func (m *mongo) NewDatastore() *mongo {
	if esErr != nil {
		panic(esErr)
	}
	m.connect()
	if m.DB != nil && m.Session != nil {
		mongoDataStore := new(mongo)
		mongoDataStore.DB = m.DB
		mongoDataStore.DBSuggested = m.DBSuggested
		mongoDataStore.Session = m.Session
		return mongoDataStore
	} else {
		panic("nil DB mongo")
	}
	return nil
}

func (m *mongo) connect() {
	fmt.Println("connected to MongoDB")
	m.DB, m.DBSuggested, m.Session = m.connectToMongo()
	go func(m *mongo) {
		var f Fetch
		f.Mongo = *m
		f.watch(m.DB)
	}(m)
}

func (m *mongo) connectToMongo() (a *mongoDB.Collection, as *mongoDB.Collection, b *mongoDB.Client) {

	var err error
	session, err := mongoDB.NewClient(options.Client().ApplyURI(os.Getenv("DATABASE")))
	if err != nil {
		panic(err)
	}
	err = session.Connect(context.Background())
	if err != nil {
		return nil, nil, nil
	}
	if err != nil {
		panic(err)
	}

	var DB = session.Database("github").Collection("users")
	var DBSuggested = session.Database("github").Collection("suggested")
	return DB, DBSuggested, session
}
func (f *Fetch) watch(collection *mongoDB.Collection) {
	//events := make(chan *types.ChangeEvent)
	go func() {
		cs, err := collection.Watch(context.TODO(), mongoDB.Pipeline{})
		if err != nil {
			panic(err)
		}
		// Whenever there is a new change event, decode the change event and print some information about it
		for cs.Next(context.TODO()) {
			var userInfo UserInfo //notice that we put here, not at above cs, err :=
			event := &types.ChangeEvent{}
			err := cs.Decode(event)
			if err != nil {
				panic(err)
			}
			switch event.OperationType {
			case "update", "replace":
				for key, value := range event.UpdateDescription.UpdatedFields {
					switch value.(type) {
					case primitive.A:
						if key == "starred" && len(value.(primitive.A)) > 0 {
							if len(userInfo.UserName) == 0 {
								opt := options.FindOne().SetProjection(bson.M{"userName": 1, "token": 1})
								err = f.Mongo.DB.FindOne(context.TODO(), bson.M{"_id": event.DocumentKey.ID}, opt).Decode(&userInfo)
								if err != nil {
									panic(err)
								}
							}
							for _, v := range value.(primitive.A) {
								bsonBytes, _ := bson.Marshal(v)
								var s types.Starred
								bson.Unmarshal(bsonBytes, &s)
								userInfo.FullName = s.FullName
								fmt.Println("Processing: " + s.FullName)
								f.fetchStargazersQuery(userInfo)
								//don't read variable from struct but pass it to goroutine for sync goroutines
								//otherwise two or more goroutines are reading the same userInfo.FullName
							}
						}
					}
				}
				//event.Data = &s
				//events <- event
				//if err != nil {
				//	panic(err)
				//}
				break
			case "insert":
				go func(event *types.ChangeEvent) {
					var page = 1
					userInfo.UserName = event.FullDocument.UserName
					userInfo.Token = event.FullDocument.Token
					f.fetchStarred(page, &userInfo)
				}(event)
				fmt.Println("new user created!")
			}
		}
	}()
}
