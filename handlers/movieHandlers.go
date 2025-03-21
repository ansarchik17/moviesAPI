package handlers

import (
	"github.com/gin-gonic/gin"
	"goozinshe/models"
	"goozinshe/repositories"
	"net/http"
	"strconv"
)

type MoviesHandler struct {
	moviesRepo *repositories.MoviesRepository
	genreRepo  *repositories.GenresRepository
}

func NewMoviesHandler(moviesRepo *repositories.MoviesRepository, genreRepo *repositories.GenresRepository) *MoviesHandler {
	return &MoviesHandler{
		moviesRepo: moviesRepo,
		genreRepo:  genreRepo,
	}
}

func (h *MoviesHandler) FindById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid movie ID"))
		return
	}

	movie, err := h.moviesRepo.FindById(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.NewApiError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, movie)
}

func (h *MoviesHandler) FindAll(c *gin.Context) {
	movies := h.moviesRepo.FindAll(c)
	c.JSON(http.StatusOK, movies)
}

func (h *MoviesHandler) Create(c *gin.Context) {
	var request struct {
		Title       string
		Description string
		ReleaseYear int
		Director    string
		TrailerUrl  string
		GenreIds    []int
	}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid request data"))
		return
	}

	genres, err := h.genreRepo.FindAllByIds(c, request.GenreIds)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError(err.Error()))
		return
	}

	movie := models.Movie{
		Title:       request.Title,
		Description: request.Description,
		ReleaseYear: request.ReleaseYear,
		Director:    request.Director,
		TrailerUrl:  request.TrailerUrl,
		Genres:      genres,
	}

	id := h.moviesRepo.Create(c, movie)
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *MoviesHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid movie ID"))
		return
	}

	var request struct {
		Title       string
		Description string
		ReleaseYear int
		Director    string
		TrailerUrl  string
		GenreIds    []int
	}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid request data"))
		return
	}

	genres, err := h.genreRepo.FindAllByIds(c, request.GenreIds)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError(err.Error()))
		return
	}

	movie := models.Movie{
		Title:       request.Title,
		Description: request.Description,
		ReleaseYear: request.ReleaseYear,
		Director:    request.Director,
		TrailerUrl:  request.TrailerUrl,
		Genres:      genres,
	}

	if err := h.moviesRepo.Update(c, id, movie); err != nil {
		c.JSON(http.StatusNotFound, models.NewApiError(err.Error()))
		return
	}

	c.Status(http.StatusOK)
}

func (h *MoviesHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid movie ID"))
		return
	}

	if err := h.moviesRepo.Delete(c, id); err != nil {
		c.JSON(http.StatusNotFound, models.NewApiError(err.Error()))
		return
	}

	c.Status(http.StatusOK)
}
