package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"ryanlawton.art/photospace-api/auth"
	"ryanlawton.art/photospace-api/photo"

	authhttp "ryanlawton.art/photospace-api/auth/delivery/http"
	authmongo "ryanlawton.art/photospace-api/auth/repository/mongo"
	authusecase "ryanlawton.art/photospace-api/auth/usecase"
	photohttp "ryanlawton.art/photospace-api/photo/delivery/http"
	photomongo "ryanlawton.art/photospace-api/photo/repository/mongo"
	photousecase "ryanlawton.art/photospace-api/photo/usecase"
)

type App struct {
	httpServer *http.Server

	photoUC photo.UseCase
	authUC  auth.UseCase
}

func NewApp() *App {
	db := initDB()

	userRepo := authmongo.NewUserRepository(db, viper.GetString("mongo.user_collection"))
	photoRepo := photomongo.NewPhotoRepository(db, viper.GetString("mongo.photo_collection"))

	return &App{
		photoUC: photousecase.NewPhotoUseCase(photoRepo),
		authUC: authusecase.NewAuthUseCase(
			userRepo,
			viper.GetString("auth.hash_salt"),
			[]byte(viper.GetString("auth.signing_key")),
			viper.GetDuration("auth.token_ttl"),
		),
	}
}

func (a *App) Run(port string) error {
	// Init gin handler
	router := gin.Default()
	router.Use(
		gin.Recovery(),
		gin.Logger(),
	)

	// Set up http handlers
	// SignUp/SignIn endpoints
	authhttp.RegisterHTTPEndpoints(router, a.authUC)

	// API endpoints
	authMiddleware := authhttp.NewAuthMiddleware(a.authUC)
	api := router.Group("/api", authMiddleware)

	photohttp.RegisterHTTPEndpoints(api, a.photoUC)

	// HTTP Server
	a.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := a.httpServer.ListenAndServe(); err != nil {
			log.Fatalf("Failed to listen and serve: %+v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	return a.httpServer.Shutdown(ctx)
}

func initDB() *mongo.Database {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Printf("Connecting to mongoDB at uri: %s", viper.GetString("mongo.uri"))

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(viper.GetString("mongo.uri")))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	return client.Database(viper.GetString("mongo.name"))
}
