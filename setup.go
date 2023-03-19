package main

import (
	"log"
	"os"

	"github.com/ffardo/user-crud/controllers"
	"github.com/ffardo/user-crud/infrastructures"
	"github.com/ffardo/user-crud/repositories"
	"github.com/ffardo/user-crud/routes"
	"github.com/ffardo/user-crud/services"
	"github.com/gin-gonic/gin"
)

func setup() *gin.Engine {

	mongoUri := os.Getenv("MONGO_URI")
	mongoUsername := os.Getenv("MONGO_USERNAME")
	mongoPassword := os.Getenv("MONGO_PASSWORD")
	apiKey := os.Getenv("API_KEY")

	client, err := infrastructures.CreateMongoClient(mongoUri, mongoUsername, mongoPassword)

	if err != nil {
		log.Fatal(err)
	}

	ur := repositories.UserRepository{
		Client: client,
	}

	ur.Init()

	uc := controllers.UserController{
		UserService: services.UserService{
			UserRepository: ur,
		},
	}

	return routes.InitRouter(&uc, apiKey)

}
