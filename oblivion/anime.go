package oblivion

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Anime struct {
	Image    string `json:"image"`
	Episodes string `json:"episodes"`
	Status   string `json:"status"`
	Class    string `json:"class"`
	Title    string `json:"title"`
	ID       string `json:"id"`
	Episode  string `json:"episode"`
	Url      string `json:"url"`
	Titles   string `json:"titles"`
}

type Media struct {
	Desc string `json:"desc"`
	Lang string `json:"lang"`
	Type string `json:"type"`
	Url  string `json:"url"`
}

type Episode struct {
	Episode int     `json:"episode"`
	Media   []Media `json:"media"`
	Summary string  `json:"summary"`
	Title   string  `json:"title"`
}

func GetWatchList() ([]Anime, error) {
	animeList := []Anime{}
	if AuthCookie == nil {
		return animeList, fmt.Errorf("Not logged in")
	}

	client := &http.Client{}
	r, _ := http.NewRequest("GET", "http://oblivion.ws/api/anime/list/watching", nil)
	r.AddCookie(AuthCookie)

	resp, err := client.Do(r)
	if err != nil {
		return animeList, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return animeList, fmt.Errorf("Bad status %d", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return animeList, err
	}

	err = json.Unmarshal(data, &animeList)
	if err != nil {
		return animeList, err
	}

	return animeList, nil
}

func GetEpisodeList(anime string) ([]Episode, error) {
	animeList := []Episode{}
	if AuthCookie == nil {
		return animeList, fmt.Errorf("Not logged in")
	}

	client := &http.Client{}
	r, _ := http.NewRequest("GET", fmt.Sprintf("http://oblivion.ws/api/anime/%s/episode/list", anime), nil)
	r.AddCookie(AuthCookie)

	resp, err := client.Do(r)
	if err != nil {
		return animeList, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return animeList, fmt.Errorf("Bad status %d", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return animeList, err
	}

	err = json.Unmarshal(data, &animeList)
	if err != nil {
		return animeList, err
	}

	return animeList, nil
}
