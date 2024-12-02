package movieservice

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DisneyMovieDataAccess struct {
	host     string
	port     string
	user     string
	password string
	dbname   string
	pool     *pgxpool.Pool
}

type MovieDataAccess interface {
	GetMovies() (*[]MovieDBObject, error)
	InsertMovies(movies *[]MovieDBObject) error
}

type MovieDBObject struct {
	id                    int64  `db:"id"`
	title                 string `db:"title"`
	overview              string `db:"overview"`
	genres                []string
	horizontalPosterw1080 string `db:"horizontal_poster_w1080"`
	verticalPosterw720    string `db:"vertical_poster_w720"`
}

func CreateMovieDataAccess(host string, port string, user string, password string, dbname string) (MovieDataAccess, error) {
	result := &DisneyMovieDataAccess{
		host:     host,
		port:     port,
		user:     user,
		password: password,
		dbname:   dbname,
	}
	pool, err := pgxpool.New(context.Background(), result.getDatabaseUrl())
	if err != nil {
		return nil, err
	}
	result.pool = pool
	return result, nil
}

func (da *DisneyMovieDataAccess) getDatabaseUrl() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?search_path=disneyschema", da.user, da.password, da.host, da.port, da.dbname)
}

func (da *DisneyMovieDataAccess) GetMovies() (*[]MovieDBObject, error) {
	rows, err := da.pool.Query(context.Background(), "SELECT * FROM movies")
	if err != nil {
		return nil, err
	}
	movies, err := pgx.CollectRows(rows, pgx.RowToStructByName[MovieDBObject])
	if err != nil {
		return nil, err
	}
	return &movies, nil
}

func (da *DisneyMovieDataAccess) InsertMovies(movies *[]MovieDBObject) error {
	query := `
	INSERT INTO disneyschema.movies (title, overview, horizontal_poster_w1080, vertical_poster_w720)
	VALUES (@title, @overview, @hp1080, @vp720) ON CONFLICT ON CONSTRAINT unique_title
	DO UPDATE
	SET overview = EXCLUDED.overview,
	horizontal_poster_w1080 = EXCLUDED.horizontal_poster_w1080,
	vertical_poster_w720 = EXCLUDED.vertical_poster_w720`
	batch := &pgx.Batch{}
	for _, movie := range *movies {
		args := pgx.NamedArgs{
			"title":    movie.title,
			"overview": movie.overview,
			"hp1080":   movie.horizontalPosterw1080,
			"vp720":    movie.verticalPosterw720,
		}
		batch.Queue(query, args)
	}
	results := da.pool.SendBatch(context.Background(), batch)
	defer results.Close()
	for range *movies {
		_, err := results.Exec()
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
				continue
			}

			return fmt.Errorf("unable to insert row: %w", err)
		}
	}
	return results.Close()
}
