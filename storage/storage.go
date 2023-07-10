package storage

import "github.com/ignoxx/podara/poc3/types"

type Storage interface {
	GetUserByEmail(string) (*types.User, error)
	CreateUser(string, string, string) (*types.User, error)

	CreatePodcast(p *types.Podcast) (*types.Podcast, error)
	GetAllPodcasts() ([]*types.Podcast, error)
	GetPodcastByID(podcastId string) (*types.Podcast, error)
    GetPodcastByUserID(userId string) ([]*types.Podcast, error)
	UpdatePodcast(p *types.Podcast) (*types.Podcast, error)
    DeletePodcast(podcastId string) error

	CreateEpisode(e *types.Episode) (*types.Episode, error)
	GetAllEpisodes(podcastId string) ([]*types.Episode, error)
	GetEpisodeByID(podcastId, episodeId string) (*types.Episode, error)
	UpdateEpisode(podcastId string, e *types.Episode) (*types.Episode, error)
    DeleteEpisode(podcastId, episodeId string) error

    GetPodcastAndEpisodesByPodcastID(podcastId string) (*types.Podcast, []*types.Episode, error)
}
