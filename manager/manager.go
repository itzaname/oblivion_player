package manager

import (
	"github/itzaname/oblivion_player/player"

	"fmt"

	"github/itzaname/oblivion_player/oblivion"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_image"
	"github.com/veandco/go-sdl2/sdl_ttf"

	"github.com/go-gl/gl/v2.1/gl"
)

const FrameRate = 60

type Instance struct {
	Window   *sdl.Window
	Renderer *sdl.Renderer
	Player   player.Instance
	State    int
	Fonts    map[string]*ttf.Font
	Cursors  map[sdl.SystemCursor]*sdl.Cursor
	Textures map[string]*sdl.Texture
	Watching []oblivion.Anime
	Settings oblivion.UserSettings
	Active   bool
	// Stuff
	PlayingAnime    oblivion.Anime
	PlayingEpisodes []oblivion.Episode
	PlayingEpisode  int
	PlayingTimer    float64
}

func New() (Instance, error) {
	startstate := 0
	if file, err := os.Open(os.TempDir() + string(os.PathSeparator) + "oblivion" + string(os.PathSeparator) + "authfile"); err == nil {
		data, err := ioutil.ReadAll(file)
		if err == nil {
			oblivion.AuthCookie = &http.Cookie{
				Name:  "oblivion",
				Value: string(data),
			}
			log.Println("Logged in with existing cookie")
			startstate = 1
		}
		file.Close()
	}

	settings, _ := oblivion.GetUserSettings()

	sdl.Init(sdl.INIT_VIDEO)
	gl.Init()

	if err := ttf.Init(); err != nil {
		return Instance{}, fmt.Errorf("Failed to initialize TTF: %s\n", err)
	}

	window, err := sdl.CreateWindow("Oblivion", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, 1280, 720, sdl.WINDOW_OPENGL|sdl.WINDOW_SHOWN|sdl.WINDOW_RESIZABLE)
	if err != nil {
		return Instance{}, err
	}

	sdl.SetHintWithPriority("SDL_VIDEO_MINIMIZE_ON_FOCUS_LOSS", "false", sdl.HINT_OVERRIDE)

	_, err = sdl.GL_CreateContext(window)
	if err != nil {
		return Instance{}, err
	}

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		return Instance{}, err
	}

	p, err := player.New(window, renderer)
	if err != nil {
		return Instance{}, err
	}

	inst := Instance{
		Player:   p,
		Window:   window,
		Renderer: renderer,
		State:    startstate,
		Active:   true,
		Settings: settings,
	}

	err = inst.LoadImage("play", "resource/play.png")
	if err != nil {
		return Instance{}, err
	}

	err = inst.LoadImage("pause", "resource/pause.png")
	if err != nil {
		return Instance{}, err
	}

	err = inst.LoadImage("cog", "resource/cog.png")
	if err != nil {
		return Instance{}, err
	}

	err = inst.LoadImage("left-arrow", "resource/left-arrow.png")
	if err != nil {
		return Instance{}, err
	}

	err = inst.LoadImage("arrows-alt", "resource/arrows-alt.png")
	if err != nil {
		return Instance{}, err
	}

	err = inst.LoadImage("chevron-right", "resource/chevron-right.png")
	if err != nil {
		return Instance{}, err
	}

	return inst, err
}

func (m *Instance) LoadAnime() error {
	anime, err := oblivion.GetWatchList()
	if err != nil {
		return fmt.Errorf("WatchList: %s", err.Error())
	}

	for _, a := range anime {
		err = m.LoadImageURL(a.ID, a.Image)
		if err != nil {
			return err
		}
	}

	m.Watching = anime

	return nil
}

func (m *Instance) GetFont(name string, size int) (*ttf.Font, error) {
	if m.Fonts == nil {
		m.Fonts = make(map[string]*ttf.Font)
	}

	value, exists := m.Fonts[fmt.Sprintf("%s-%d", name, size)]
	if exists {
		return value, nil
	}

	font, err := ttf.OpenFont("fonts"+string(os.PathSeparator)+name, size)
	if err != nil {
		return nil, err
	}

	m.Fonts[fmt.Sprintf("%s-%d", name, size)] = font

	return font, nil
}

func (m *Instance) SetCursor(cursor sdl.SystemCursor) {
	if m.Cursors == nil {
		m.Cursors = make(map[sdl.SystemCursor]*sdl.Cursor)
	}

	value, exists := m.Cursors[cursor]
	if exists {
		sdl.SetCursor(value)
		return
	}

	m.Cursors[cursor] = sdl.CreateSystemCursor(cursor)
	sdl.SetCursor(m.Cursors[cursor])
}

func (m *Instance) LoadImageURL(name, url string) error {
	if _, err := os.Stat(os.TempDir() + string(os.PathSeparator) + "oblivion" + string(os.PathSeparator) + name); err == nil {
		return nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(os.TempDir()+string(os.PathSeparator)+"oblivion"+string(os.PathSeparator)+name, data, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func (m *Instance) LoadImage(name, path string) error {
	if m.Textures == nil {
		m.Textures = make(map[string]*sdl.Texture)
	}

	image, err := img.Load(path)
	if err != nil {
		return err
	}
	defer image.Free()

	newtx, err := m.Renderer.CreateTextureFromSurface(image)
	if err != nil {
		return err
	}

	m.Textures[name] = newtx

	return nil
}

func (m *Instance) GetTexture(name string) (*sdl.Texture, error) {
	if m.Textures == nil {
		m.Textures = make(map[string]*sdl.Texture)
	}

	texture, exists := m.Textures[name]
	if exists {
		return texture, nil
	}

	image, err := img.Load(os.TempDir() + string(os.PathSeparator) + "oblivion" + string(os.PathSeparator) + name)
	if err != nil {
		return nil, err
	}
	defer image.Free()

	newtx, err := m.Renderer.CreateTextureFromSurface(image)
	if err != nil {
		return nil, err
	}

	m.Textures[name] = newtx

	return newtx, nil
}

func (m *Instance) Run() {
	for {
		if !m.Active {
			break
		}
		switch m.State {
		case 0:
			m.LoginView()
			break
		case 1:
			m.BrowseView()
			break
		case 2:
			m.PlayerView()
			break
		}
	}
}

func (m *Instance) SelectSubtitle(subtitle interface{}) {
	log.Println(subtitle.(oblivion.Media).Desc)
}

func (m *Instance) PlayAnime(anime oblivion.Anime) error {
	curepisode, err := strconv.ParseInt(anime.Episode, 10, 64)
	if err != nil {
		return err
	}

	episodes, err := oblivion.GetEpisodeList(anime.ID)
	if err != nil {
		return err
	}

	m.PlayingAnime = anime
	m.PlayingEpisodes = episodes
	m.PlayingTimer = 0

	for _, ep := range episodes {
		if ep.Episode == int(curepisode)+1 {
			video := []oblivion.Media{}
			audio := []oblivion.Media{}
			subtitles := []oblivion.Media{}
			for _, media := range ep.Media {
				switch media.Type {
				case "video":
					video = append(video, media)
					break
				case "audio":
					audio = append(audio, media)
					break
				case "subtitle":
					subtitles = append(subtitles, media)
					break
				}
			}
			found := false
			if len(video) > 0 {
				m.Player.LoadVideo(video[0].Url)
			} else {
				m.State = 1
				DisplayBrowseError("Episode not downloaded.", 5)
				return nil
			}
			if len(audio) > 0 {

				for _, aud := range audio {
					if aud.Lang == m.Settings.Language {
						m.Player.LoadAudio(aud.Url)
						found = true
						break
					}
				}
				if !found {
					m.Player.LoadAudio(audio[0].Url)
				}
			}
			if len(subtitles) > 0 {
				if m.Settings.Language != "eng" || !found {
					m.Player.LoadSubs(subtitles[0].Url)
				}
			}

			m.PlayingEpisode = int(curepisode)+1
			m.Player.ClearMenu()
			m.Player.MenuOpen = false
			{
				items := []player.SubItem{}
				for index, sub := range subtitles {
					name := sub.Desc
					if name == "" {
						name = sub.Lang
					}
					if name == "" {
						name = fmt.Sprintf("Subtitle %d", index + 1)
					}
					items = append(items, player.SubItem{
						name,
						sub,
						false,
					})
				}
				m.Player.AddMenu("Subtitles", m.SelectSubtitle, items)
			}

			{
				items := []player.SubItem{}
				for index, aud := range audio {
					name := aud.Desc
					if name == "" {
						name = aud.Lang
					}
					if name == "" {
						name = fmt.Sprintf("Audio %d", index + 1)
					}
					items = append(items, player.SubItem{
						name,
						aud,
						false,
					})
				}
				m.Player.AddMenu("Audio", m.SelectSubtitle, items)
			}

			m.Player.SetPositionForce(oblivion.GetEpisodeWatchTime(m.PlayingAnime.ID, fmt.Sprintf("%d", ep.Episode)))

			m.Player.Play()
			return nil
		}
	}

	m.State = 1
	AnimateShow(0)
	m.Window.SetFullscreen(0)
	sdl.ShowCursor(sdl.ENABLE)
	m.SetCursor(sdl.SYSTEM_CURSOR_ARROW)
	DisplayBrowseError("Episode not downloaded.", 5)

	return fmt.Errorf("Not found")
}
