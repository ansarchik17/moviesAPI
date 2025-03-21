package handlers

import (
	"github.com/gin-gonic/gin"
	"goozinshe/models"
	"goozinshe/repositories"
	"net/http"
	"strconv"
	"strings"
)

type GenreHandlers struct {
	genreRepo *repositories.GenresRepository
}

func NewGenreHandlers(genreRepo *repositories.GenresRepository) *GenreHandlers {
	return &GenreHandlers{
		genreRepo: genreRepo,
	}
}

// GET: Retrieve all genres
func (g *GenreHandlers) GetGenres(c *gin.Context) {
	genres, err := g.genreRepo.GetGenres(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewApiError("Failed to retrieve genres"))
		return
	}
	c.JSON(http.StatusOK, genres)
}

// POST: Create a new genre
func (g *GenreHandlers) CreateGenre(c *gin.Context) {
	var genre models.Genre
	if err := c.BindJSON(&genre); err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid request format"))
		return
	}
	id, err := g.genreRepo.CreateGenre(c.Request.Context(), genre)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewApiError(err.Error()))
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// GET: Retrieve multiple genres by IDs
func (g *GenreHandlers) FindAllByIds(c *gin.Context) {
	idStr := c.Query("ids") // Expecting a query parameter like "?ids=1,2,3"
	if idStr == "" {
		c.JSON(http.StatusBadRequest, models.NewApiError("No IDs provided"))
		return
	}

	idStrings := strings.Split(idStr, ",")
	var ids []int
	for _, idString := range idStrings {
		id, err := strconv.Atoi(strings.TrimSpace(idString))
		if err != nil {
			c.JSON(http.StatusBadRequest, models.NewApiError("Invalid ID format"))
			return
		}
		ids = append(ids, id)
	}

	genres, err := g.genreRepo.FindAllByIds(c.Request.Context(), ids)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewApiError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, genres)
}

// PUT: Update an existing genre
func (g *GenreHandlers) UpdateGenre(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid ID format"))
		return
	}

	_, err = g.genreRepo.FindAllByIds(c.Request.Context(), []int{id})
	if err != nil {
		c.JSON(http.StatusNotFound, models.NewApiError("Genre not found"))
		return
	}

	var updatedGenre models.Genre
	if err := c.BindJSON(&updatedGenre); err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid request format"))
		return
	}

	if err := g.genreRepo.UpdateGenre(c.Request.Context(), id, updatedGenre); err != nil {
		c.JSON(http.StatusInternalServerError, models.NewApiError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Genre updated successfully",
		"updated_genre": updatedGenre,
	})
}

// DELETE: Remove a genre
func (g *GenreHandlers) DeleteGenre(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid ID format"))
		return
	}

	_, err = g.genreRepo.FindAllByIds(c.Request.Context(), []int{id})
	if err != nil {
		c.JSON(http.StatusNotFound, models.NewApiError("Genre not found"))
		return
	}

	if err := g.genreRepo.DeleteGenre(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, models.NewApiError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Genre deleted successfully"})
}
