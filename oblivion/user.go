package oblivion

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type UserSettings struct {
	Language string
	Admin    bool
	AutoPlay bool
	AutoNext bool
}

func GetUserSettings() (UserSettings, error) {
	settings := UserSettings{}
	if AuthCookie == nil {
		return settings, fmt.Errorf("Not logged in")
	}

	client := &http.Client{}
	r, _ := http.NewRequest("GET", "http://oblivion.ws/api/settings", nil)
	r.AddCookie(AuthCookie)

	resp, err := client.Do(r)
	if err != nil {
		return settings, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return settings, fmt.Errorf("Bad status %d", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return settings, err
	}

	err = json.Unmarshal(data, &settings)
	if err != nil {
		return settings, err
	}

	return settings, nil
}

func ReportEpisodeWatchTime(anime, epsisode string, time float64) error {
	if AuthCookie == nil {
		return fmt.Errorf("Not logged in")
	}

	data := url.Values{}
	data.Set("time", fmt.Sprintf("%f", time))

	client := &http.Client{}
	r, _ := http.NewRequest("POST", fmt.Sprintf("https://oblivion.ws/api/time/report/%s/%s", anime, epsisode), strings.NewReader(data.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	r.AddCookie(AuthCookie)

	resp, err := client.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("Bad status %d", resp.StatusCode)
	}

	return nil
}

func GetEpisodeWatchTime(anime, epsisode string) float64 {
	if AuthCookie == nil {
		return 0
	}

	client := &http.Client{}
	r, _ := http.NewRequest("GET", fmt.Sprintf("http://oblivion.ws/api/time/get/%s/%s", anime, epsisode), nil)
	r.AddCookie(AuthCookie)

	resp, err := client.Do(r)
	if err != nil {
		return 0
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return 0
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0
	}

	time, err := strconv.ParseFloat(string(data), 64)
	if err != nil {
		return 0
	}

	return time
}

func SetCurrentWatchingEpisode(anime, epsisode string) error {
	if AuthCookie == nil {
		return fmt.Errorf("Cookie not set")
	}

	client := &http.Client{}
	r, _ := http.NewRequest("GET", fmt.Sprintf("https://oblivion.ws/api/anime/set/%s/%s", anime, epsisode), nil)
	r.AddCookie(AuthCookie)

	resp, err := client.Do(r)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("Bad status %d", resp.StatusCode)
	}

	return nil
}