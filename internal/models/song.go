package models

import "time"

type SongRequest struct {
	ID            string    `bson:"id" json:"id"`
	RoomID        string    `bson:"room_id" json:"room_id"`
	YouTubeID     string    `bson:"youtube_video_id" json:"youtube_video_id"`
	Title         string    `bson:"title" json:"title"`
	Thumbnail     string    `bson:"thumbnail" json:"thumbnail"`
	URL           string    `bson:"youtube_url" json:"youtube_url"`
	SubmitterName string    `bson:"submitter_name" json:"submitter_name"`
	SubmitterType string    `bson:"submitter_type" json:"submitter_type"`
	Votes         int       `bson:"votes" json:"votes"`
	VotedBy       []string  `bson:"voted_by" json:"voted_by"`
	Status        string    `bson:"status" json:"status"`
	CreatedAt     time.Time `bson:"created_at" json:"created_at"`
}

type SongRequestCreate struct {
	YouTubeID     string `json:"youtube_video_id"`
	Title         string `json:"title"`
	Thumbnail     string `json:"thumbnail"`
	URL           string `json:"youtube_url"`
	SubmitterName string `json:"submitter_name"`
	SubmitterType string `json:"submitter_type"`
}

type VoteRequest struct {
	SessionID string `json:"session_id"`
}

type StatusUpdate struct {
	Status string `json:"status"`
}

