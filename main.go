package main

import (
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
	"goozinshe/config"
	"goozinshe/handlers"
	"goozinshe/repositories"
)

func main() {
	r := gin.Default()
	corsConfig := cors.Config{
		AllowAllOrigins: true,
		AllowHeaders:    []string{"*"},
		AllowMethods:    []string{"*"},
	}
	r.Use(cors.New(corsConfig))

	err := loadConfig()
	if err != nil {
		panic(err)
	}

	conn, err := connetToBd()
	if err != nil {
		panic(err)
	}

	moviesRepository := repositories.NewMoviesRepository()
	genresRepository := repositories.NewGenresRepository(conn)

	moviesHandler := handlers.NewMoviesHandler(moviesRepository, genresRepository)
	genresHandler := handlers.NewGenreHandlers(genresRepository)
	//movies routers
	r.GET("/movies", moviesHandler.FindAll)
	r.GET("/movies/:id", moviesHandler.FindById)
	r.POST("/movies", moviesHandler.Create)
	r.PUT("/movies/:id", moviesHandler.Update)
	r.DELETE("/movies/:id", moviesHandler.Delete)
	//genres routers
	r.GET("/genres", genresHandler.GetGenres)
	r.GET("/genres/:id", genresHandler.FindAllByIds)
	r.POST("/genres", genresHandler.CreateGenre)
	r.PUT("/genres/:id", genresHandler.UpdateGenre)
	r.DELETE("/genres/:id", genresHandler.DeleteGenre)
	r.Run(config.Config.AppHost)
}

func connetToBd() (*pgxpool.Pool, error) {
	conn, err := pgxpool.New(context.Background(), config.Config.DbConnectionString)
	if err != nil {
		return nil, err
	}
	err = conn.Ping(context.Background())
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func loadConfig() error {
	viper.SetConfigName(".env")
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	var mapConfig config.MapConfig
	err = viper.Unmarshal(&mapConfig)
	if err != nil {
		return err
	}
	config.Config = &mapConfig
	return nil
}
