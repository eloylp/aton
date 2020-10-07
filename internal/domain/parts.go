package domain

import (
	"time"
)

type VideoStream struct {
	ID             string    `json:"id"`
	URL            string    `json:"url"`
	Status         string    `json:"status"`
	LastConnection time.Time `json:"lastConnection"`
}

type Face struct {
	ID      string    `json:"id"`
	Name    string    `json:"name"`
	Samples []*Sample `json:"samples"`
}

type Sample struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Data []byte `json:"data"`
}

type Match struct {
	ID          string       `json:"id"`
	Face        *Face        `json:"faceId"`
	VideoStream *VideoStream `json:"videoStream"`
	MatchTime   time.Time    `json:"matchTime"`
}
