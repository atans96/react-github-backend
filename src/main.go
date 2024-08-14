package main

import (
	"backend/src/service"
	"github.com/joho/godotenv"
	"os"
	"path/filepath"
)

func init() {
	dir, err := filepath.Abs("./")
	if err != nil {
		panic(err)
	}
	err = os.Chdir(filepath.Join(dir, "src"))
	if err != nil {
		panic(err)
	}
	err = godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
}
func main() {
	_ = service.Mongo.NewDatastore()
	select {}
}
