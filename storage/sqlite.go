package storage

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/ignoxx/podara/poc3/types"
	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite"
)

type SqliteStorage struct {
	db *sql.DB
}

func NewSqliteStorage(file string) *SqliteStorage {
	db, err := sql.Open("sqlite", file)

	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id TEXT PRIMARY KEY UNIQUE,
            email TEXT NOT NULL UNIQUE,
            password TEXT NOT NULL
        );

        CREATE TABLE IF NOT EXISTS podcasts (
            id TEXT PRIMARY KEY UNIQUE,
            title TEXT NOT NULL UNIQUE,
            description TEXT NOT NULL,
            cover_url TEXT NOT NULL,
            user_id TEXT NOT NULL,
            created_at DATETIME NOT NULL,
            updated_at DATETIME NOT NULL,
            FOREIGN KEY(user_id) REFERENCES users(id)
        );

        CREATE TABLE IF NOT EXISTS episodes (
            id TEXT PRIMARY KEY UNIQUE,
            title TEXT NOT NULL UNIQUE,
            description TEXT NOT NULL,
            podcast_id TEXT NOT NULL,
            cover_url TEXT NOT NULL,
            audio_url TEXT NOT NULL,
            duration_ms INTEGER NOT NULL,
            created_at DATETIME NOT NULL,
            updated_at DATETIME NOT NULL,
            FOREIGN KEY(podcast_id) REFERENCES podcasts(id)
        );
    `)

    if err != nil {
        panic(err)
    }

    _, err = db.Exec("PRAGMA journal_mode=WAL")

    if err != nil {
        panic(err)
    }

	return &SqliteStorage{
		db: db,
	}
}

func (s *SqliteStorage) Close() {
	s.db.Close()
}

func (s *SqliteStorage) GetUserByEmail(email string) (*types.User, error) {
	var user types.User

	err := s.db.QueryRow(`
        SELECT id, email, password FROM users WHERE email = ?
    `, email).Scan(&user.Id, &user.Email, &user.Password)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *SqliteStorage) CreateUser(email, name, password string) (*types.User, error) {
	var user types.User

	user.Id = uuid.NewString()
	user.Email = email
	user.Name = name

	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	user.Password = string(hashedPassword)

	_, err = s.db.Exec(`
        INSERT INTO users (id, email, password) VALUES (?, ?, ?)
    `, user.Id, user.Email, user.Password)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *SqliteStorage) CreatePodcast(p *types.Podcast) (*types.Podcast, error) {
	p.Id = uuid.NewString()

	_, err := s.db.Exec(`
        INSERT INTO podcasts (id, title, description, cover_url, user_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
    `, p.Id, p.Title, p.Description, p.CoverImageUrl, p.UserId, p.CreatedAt, p.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return p, nil
}

func (s *SqliteStorage) GetAllPodcasts() ([]*types.Podcast, error) {
	rows, err := s.db.Query(`
        SELECT id, title, description, cover_url, user_id, created_at, updated_at FROM podcasts
    `)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var podcasts []*types.Podcast

	for rows.Next() {
		var p types.Podcast

		err := rows.Scan(&p.Id, &p.Title, &p.Description, &p.CoverImageUrl, &p.UserId, &p.CreatedAt, &p.UpdatedAt)

		if err != nil {
			return nil, err
		}

		podcasts = append(podcasts, &p)
	}

	return podcasts, nil
}

func (s *SqliteStorage) GetPodcastByID(podcast_id string) (*types.Podcast, error) {
	var p types.Podcast

	err := s.db.QueryRow(`
        SELECT id, title, description, cover_url, user_id, created_at, updated_at FROM podcasts WHERE id = ?
    `, podcast_id).Scan(&p.Id, &p.Title, &p.Description, &p.CoverImageUrl, &p.UserId, &p.CreatedAt, &p.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (s *SqliteStorage) UpdatePodcast(p *types.Podcast) (*types.Podcast, error) {
	_, err := s.db.Exec(`
        UPDATE podcasts SET title = ?, description = ?, cover_url = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?
    `, p.Title, p.Description, p.CoverImageUrl, p.UpdatedAt, p.Id)

	if err != nil {
		return nil, err
	}

	return p, nil
}

func (s *SqliteStorage) DeletePodcast(podcastId string) error {
	_, err := s.db.Exec(`
        DELETE FROM podcasts WHERE id = ?
    `, podcastId)

	if err != nil {
		return err
	}

	return nil
}

func (s *SqliteStorage) CreateEpisode(e *types.Episode) (*types.Episode, error) {
	e.Id = uuid.NewString()

	_, err := s.db.Exec(`
        INSERT INTO episodes (id, title, description, podcast_id, cover_url, audio_url, duration_ms, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
    `, e.Id, e.Title, e.Description, e.PodcastId, e.CoverImageUrl, e.AudioUrl, e.DurationMs, e.CreatedAt, e.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return e, nil
}

func (s *SqliteStorage) GetAllEpisodes(podcastId string) ([]*types.Episode, error) {
	rows, err := s.db.Query(`
        SELECT id, title, description, podcast_id, cover_url, audio_url, duration_ms, created_at, updated_at FROM episodes WHERE podcast_id = ?
    `, podcastId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var episodes []*types.Episode

	for rows.Next() {
		var e types.Episode

		err := rows.Scan(&e.Id, &e.Title, &e.Description, &e.PodcastId, &e.CoverImageUrl, &e.AudioUrl, &e.DurationMs, &e.CreatedAt, &e.UpdatedAt)

		if err != nil {
			return nil, err
		}

		episodes = append(episodes, &e)
	}

	return episodes, nil
}

func (s *SqliteStorage) GetEpisodeByID(podcastId, episodeId string) (*types.Episode, error) {
	var e types.Episode

	err := s.db.QueryRow(`
        SELECT id, title, description, podcast_id, cover_url, audio_url, duration_ms, created_at, updated_at FROM episodes WHERE podcast_id = ? AND id = ?
    `, podcastId, episodeId).Scan(&e.Id, &e.Title, &e.Description, &e.PodcastId, &e.CoverImageUrl, &e.AudioUrl, &e.DurationMs, &e.CreatedAt, &e.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &e, nil
}

func (s *SqliteStorage) UpdateEpisode(podcastId string, e *types.Episode) (*types.Episode, error) {
	_, err := s.db.Exec(`
        UPDATE episodes SET title = ?, description = ?, cover_url = ?, audio_url = ?, duration_ms = ?, updated_at = CURRENT_TIMESTAMP WHERE podcast_id = ? AND id = ?
    `, e.Title, e.Description, e.CoverImageUrl, e.AudioUrl, e.DurationMs, e.UpdatedAt, podcastId, e.Id)

	if err != nil {
		return nil, err
	}

	return e, nil
}

func (s *SqliteStorage) DeleteEpisode(podcastId, episodeId string) error {
	_, err := s.db.Exec(`
        DELETE FROM episodes WHERE podcast_id = ? AND id = ?
    `, podcastId, episodeId)

	if err != nil {
		return err
	}

	return nil
}

func (s *SqliteStorage) GetPodcastAndEpisodesByPodcastID(podcastId string) (*types.Podcast, []*types.Episode, error) {
	var p *types.Podcast
	var episodes []*types.Episode

    p, err := s.GetPodcastByID(podcastId)

    if err != nil {
        return nil, nil, err
    }

    episodes, err = s.GetAllEpisodes(podcastId)

    if err != nil {
        return nil, nil, err
    }

    return p, episodes, nil
}
