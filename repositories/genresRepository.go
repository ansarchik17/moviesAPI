package repositories

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"goozinshe/models"
)

type GenresRepository struct {
	db *pgxpool.Pool
}

func NewGenresRepository(conn *pgxpool.Pool) *GenresRepository {
	return &GenresRepository{db: conn}
}

func (r *GenresRepository) FindAllByIds(c context.Context, ids []int) ([]models.Genre, error) {
	rows, err := r.db.Query(c, "SELECT id, title FROM genres WHERE id = ANY($1)", ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	genres := make([]models.Genre, 0)
	for rows.Next() {
		var genre models.Genre
		if err := rows.Scan(&genre.Id, &genre.Title); err != nil {
			return nil, err
		}
		genres = append(genres, genre)
	}
	return genres, nil
}

func (r *GenresRepository) GetGenres(c context.Context) ([]models.Genre, error) {
	rows, err := r.db.Query(c, "SELECT id, title FROM genres")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	genres := make([]models.Genre, 0)
	for rows.Next() {
		var genre models.Genre
		if err := rows.Scan(&genre.Id, &genre.Title); err != nil {
			return nil, err
		}
		genres = append(genres, genre)
	}
	return genres, nil
}

func (r *GenresRepository) CreateGenre(c context.Context, genre models.Genre) (int, error) {
	var id int
	row := r.db.QueryRow(c, "INSERT INTO genres (title) VALUES ($1) RETURNING id", genre.Title)
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *GenresRepository) UpdateGenre(c context.Context, id int, genre models.Genre) error {
	_, err := r.db.Exec(c, "UPDATE genres SET title=$1 WHERE id=$2", genre.Title, id)
	return err
}

func (r *GenresRepository) DeleteGenre(c context.Context, id int) error {
	_, err := r.db.Exec(c, "DELETE FROM genres WHERE id=$1", id)
	return err
}
