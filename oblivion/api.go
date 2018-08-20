package oblivion

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

var AuthCookie *http.Cookie

func Login(user, pass string) error {
	data := url.Values{}
	data.Set("username", user)
	data.Add("password", pass)

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	r, _ := http.NewRequest("POST", "https://oblivion.ws/", strings.NewReader(data.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, err := client.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	for _, cookie := range resp.Cookies() {
		if cookie.Name == "oblivion" {
			AuthCookie = cookie
			ioutil.WriteFile(os.TempDir()+string(os.PathSeparator)+"oblivion"+string(os.PathSeparator)+"authfile", []byte(cookie.Value), os.ModePerm)
			return nil
		}
	}

	return fmt.Errorf("Cookie not found.")
}
