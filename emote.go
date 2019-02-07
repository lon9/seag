package main

import "strings"

// BaseImageURL is base url of emoticon
const BaseImageURL = "https://steamcommunity.com/economy/emoticon/"

// Emote is object of a emoticon
type Emote struct {
	Name  string    `json:"name"`
	URL   string    `json:"url"`
	Price string    `json:"price"`
	HLS   []float64 `json:"hls"`
	Pos   []int     `json:"pos"`
	Game  string    `json:"game"`
}

// GetImageURL returns image url
func (e *Emote) GetImageURL() string {
	return BaseImageURL + strings.Replace(e.Name, ":", "", -1)
}
