package database

type Track struct {
	ID       uint   `json:"id"`
	Title    string `json:"title"`
	Artist   string `json:"artist"`
	Album    string `json:"album"`
	Genre    string `json:"genre"`
	Duration uint   `json:"duration"`
	Filepath string `json:"filepath"`
}

type Album struct {
	ID            uint   `json:"id"`
	Title         string `json:"title"`
	Artist        string `json:"artist"`
	TrackIDs      []uint `json:"trackIds"`
	CoverFilepath string `json:"coverFilepath"`
}

type Artist struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	AlbumIDs   []uint `json:"albumIds"`
	TrackIDs   []uint `json:"trackIds"`
	IsVerified bool   `json:"isVerified"`
}

type Playlist struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	OwnerID  uint   `json:"ownerId"`
	IsPublic bool   `json:"isPublic"`
}

type User struct {
	ID           uint        `json:"id"`
	Username     string      `json:"username"`
	PasswordHash string      `json:"passwordHash"`
	PlaylistIDs  []uint      `json:"playlistIds"`
	PlayerState  PlayerState `json:"playerState"`
}

type PlayerState struct {
	TrackID               uint    `json:"trackId"`
	Progress              uint    `json:"progress"` // In seconds
	Device                string  `json:"device"`
	ShuffleEnabled        bool    `json:"shuffleEnabled"`
	RepeatPlaylistEnabled bool    `json:"repeatPlaylistEnabled"`
	RepeatTrackEnabled    bool    `json:"repeatTrackEnabled"`
	Volume                float32 `json:"volume"` // From 0.0 to 1.0
}
