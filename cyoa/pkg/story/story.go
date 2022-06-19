package story

import "encoding/json"

type Story map[string]StoryArc

type StoryArc struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []struct {
		Text string `json:"text"`
		Arc  string `json:"arc"`
	} `json:"options"`
}

func ParseStory(text []byte) (Story, error) {
	var s Story
	if err := json.Unmarshal(text, &s); err != nil {
		return nil, err
	}
	return s, nil
}
