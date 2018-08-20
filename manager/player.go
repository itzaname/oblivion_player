package manager

import (
	"log"

	"fmt"
	"github/itzaname/oblivion_player/oblivion"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

var (
	LastMouseMove  int64 = 0
	ControlsActive       = false
)

func (m *Instance) GetControlsActive() bool {
	if m.Player.Paused() {
		return true
	}

	return (LastMouseMove + 3) > time.Now().Unix()
}

func (m *Instance) renderSettings() {
	menuitemheight := int32(30)
	menuwidth := int32(150)
	menuheight := int32(0)

	// Calc height
	for _, item := range m.Player.Menu {
		menuheight += menuitemheight + menuitemheight*int32(len(item.Items))
	}

	mx, my, _ := sdl.GetMouseState()
	w, h := m.Window.GetSize()
	menux := int32(w) - menuwidth - 10
	menuy := int32(h) - menuheight - 55

	// Render
	renderoffset := int32(0)
	for _, item := range m.Player.Menu {
		itemy := renderoffset
		twidth, theight, err := m.getTextSize(item.Name, "helvetica.ttf", 14)
		if err != nil {
			return
		}

		m.Renderer.SetDrawColor(20, 20, 20, 255)
		m.drawFilledBox(menux, menuy+itemy, menuwidth, menuitemheight)

		renderoffset += menuitemheight + menuitemheight*int32(len(item.Items))
		m.drawTextBG(
			menux+menuwidth/2-twidth/2,
			menuy+itemy+menuitemheight/2-theight/2,
			item.Name,
			sdl.Color{255, 255, 255, 255},
			sdl.Color{20, 20, 20, 255},
			"helvetica.ttf",
			14)

		for i := int32(0); i < int32(len(item.Items)); i++ {
			itemx := menux
			itemy := menuy + itemy + menuitemheight*(i+1)

			if item.Items[i].Selected || isInBounds(itemx, itemy, menuwidth, menuitemheight, int32(mx), int32(my)) {
				m.Renderer.SetDrawColor(60, 60, 60, 255)
			} else {
				m.Renderer.SetDrawColor(40, 40, 40, 255)
			}
			m.drawFilledBox(itemx, itemy, menuwidth, menuitemheight)
			r, g, b, a, err := m.Renderer.GetDrawColor()
			if err != nil {
				return
			}

			twidth, theight, err := m.getTextSize(item.Items[i].Name, "helvetica.ttf", 14)
			if err != nil {
				return
			}

			m.drawTextBG(
				menux+menuwidth/2-twidth/2,
				itemy+menuitemheight/2-theight/2,
				item.Items[i].Name,
				sdl.Color{255, 255, 255, 255},
				sdl.Color{r, g, b, a},
				"helvetica.ttf",
				14)
		}
	}
}

func (m *Instance) renderControls() {
	w, h := m.Window.GetSize()
	mx, my, _ := sdl.GetMouseState()
	//m.Renderer.SetDrawColor(20, 20, 20, 255)
	//m.Renderer.Clear()
	m.Renderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND)
	//m.Renderer.Set(sdl.BLENDMODE_BLEND)

	/// EXIT BUTTON ///
	m.Renderer.SetDrawColor(20, 20, 20, 255) // Holder
	m.drawFilledBox(0, 20, 40, 28)

	texture, err := m.GetTexture("left-arrow")
	if err != nil {
		log.Println(err)
		return
	}

	m.Renderer.Copy(texture, nil, &sdl.Rect{
		X: 10,
		Y: 22,
		W: 24,
		H: 24,
	})

	////////////////////////////////////////////////////////////////

	// Holder
	m.Renderer.SetDrawColor(20, 20, 20, 255)
	m.drawFilledBox(0, int32(h)-40, int32(w), 40)

	// Bars
	m.Renderer.SetDrawColor(61, 61, 61, 255) // Holder
	m.drawFilledBox(0, int32(h)-45, int32(w), 5)

	if m.Player.Paused() {
		texture, err = m.GetTexture("play")
		if err != nil {
			log.Println(err)
			return
		}
	} else {
		texture, err = m.GetTexture("pause")
		if err != nil {
			log.Println(err)
			return
		}
	}

	m.Renderer.Copy(texture, nil, &sdl.Rect{
		X: 10,
		Y: int32(h) - 32,
		W: 24,
		H: 24,
	})

	texture, err = m.GetTexture("arrows-alt")
	if err != nil {
		log.Println(err)
		return
	}

	m.Renderer.Copy(texture, nil, &sdl.Rect{
		X: int32(w) - 34,
		Y: int32(h) - 32,
		W: 24,
		H: 24,
	})

	texture, err = m.GetTexture("cog")
	if err != nil {
		log.Println(err)
		return
	}

	m.Renderer.Copy(texture, nil, &sdl.Rect{
		X: int32(w) - 65,
		Y: int32(h) - 32,
		W: 24,
		H: 24,
	})

	position, err := m.Player.GetPosition()
	if err != nil {
		return
	}

	duration, err := m.Player.GetDuration()
	if err != nil {
		return
	}

	m.Renderer.SetDrawColor(40, 40, 40, 255) // Progress
	m.drawFilledBox(0, int32(h)-45, int32(float64(w)*(position/duration)), 5)

	if isInBounds(0, int32(h)-45, int32(w), 5, int32(mx), int32(my)) {
		m.Renderer.SetDrawColor(200, 200, 200, 255) // Progress
		m.drawFilledBox(0, int32(h)-45, int32(mx), 5)
	}

	m.drawTextBG(45, int32(h)-26, getHHMMSS(position)+"/"+getHHMMSS(duration), sdl.Color{255, 255, 255, 255}, sdl.Color{20, 20, 20, 255}, "helvetica.ttf", 14)
	//m.drawText(45, int32(h)-26, getHHMMSS(position)+"/"+getHHMMSS(duration), 255, 255, 255, "helvetica.ttf", 14)

}

func (m *Instance) logicControls() {
	w, h := m.Window.GetSize()
	mx, my, _ := sdl.GetMouseState()

	mousehand := false
	if ControlsActive {
		// Exit Button
		if isInBounds(0, 20, 40, 28, int32(mx), int32(my)) {
			mousehand = true
		}
		// Fullscreen Button
		if isInBounds(int32(w)-34, int32(h)-32, 24, 24, int32(mx), int32(my)) {
			mousehand = true
		}
		// Play Pause Button
		if isInBounds(10, int32(h)-32, 40, 28, int32(mx), int32(my)) {
			mousehand = true
		}
		// Progress Selection
		if isInBounds(0, int32(h)-45, int32(w), 5, int32(mx), int32(my)) {
			mousehand = true
		}
		// Settings button
		if isInBounds(int32(w)-65, int32(h)-32, 24, 24, int32(mx), int32(my)) {
			mousehand = true
		}
	}

	if mousehand {
		m.SetCursor(sdl.SYSTEM_CURSOR_HAND)
	} else {
		m.SetCursor(sdl.SYSTEM_CURSOR_ARROW)
	}
}

func (m *Instance) logicClickControls(mx, my int32) {
	w, h := m.Window.GetSize()
	//mx, my, state := sdl.GetMouseState()

	// CLICKING PLAYER PAUSE/PLAY
	if ControlsActive && isInBounds(0, 0, int32(w), int32(h)-45, mx, my) {
		m.Player.TogglePause()
	} else if !ControlsActive && isInBounds(0, 0, int32(w), int32(h), mx, my) {
		m.Player.TogglePause()
	}

	// ANYTHING BEYOND IS A CONTROL BUTTON
	if !ControlsActive {
		return
	}

	// EXIT BUTTON
	if isInBounds(0, 20, 40, 28, int32(mx), int32(my)) {
		m.Player.Stop()
		m.Window.SetFullscreen(0)
		sdl.ShowCursor(sdl.ENABLE)
		m.LoadAnime()
		m.State = 1
	}

	// FULLSCREEN BUTTON
	if isInBounds(int32(w)-34, int32(h)-32, 24, 24, int32(mx), int32(my)) {
		if (m.Window.GetFlags() & sdl.WINDOW_FULLSCREEN_DESKTOP) != 0 {
			m.Window.SetFullscreen(0)
		} else {
			m.Window.SetFullscreen(sdl.WINDOW_FULLSCREEN_DESKTOP)
		}
	}

	// PlAY/PAUSE BUTTON
	if isInBounds(10, int32(h)-32, 40, 28, int32(mx), int32(my)) {
		m.Player.TogglePause()
	}

	// PROGRESS SELECTION
	if isInBounds(0, int32(h)-45, int32(w), 5, int32(mx), int32(my)) {
		m.Renderer.SetDrawColor(200, 200, 200, 255) // Progress
		m.drawFilledBox(0, int32(h)-45, int32(mx), 5)
		if duration, err := m.Player.GetDuration(); err == nil {
			m.Player.SetPosition((float64(mx) / float64(w)) * duration)
		}
	}

	// SETTINGS BUTTON
	if isInBounds(int32(w)-65, int32(h)-32, 24, 24, int32(mx), int32(my)) {
		m.Player.MenuOpen = !m.Player.MenuOpen
	}
}

func (m *Instance) renderView() {
	m.Player.Render()
}

func (m *Instance) PlayerView() {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			m.Active = false
			t.Type = 22
			return
		case *sdl.KeyUpEvent:
			if t.Keysym.Sym == sdl.K_SPACE {
				m.Player.TogglePause()
			}
			if t.Keysym.Sym == sdl.K_ESCAPE {
				m.Window.SetFullscreen(0)
			}
			if t.Keysym.Sym == sdl.K_f {
				if (m.Window.GetFlags() & sdl.WINDOW_FULLSCREEN_DESKTOP) != 0 {
					m.Window.SetFullscreen(0)
				} else {
					m.Window.SetFullscreen(sdl.WINDOW_FULLSCREEN_DESKTOP)
				}
			}
			break
		case *sdl.MouseButtonEvent:
			if t.Button == sdl.BUTTON_LEFT && t.State == sdl.RELEASED {
				m.logicClickControls(t.X, t.Y)
			}
		case *sdl.MouseMotionEvent:
			LastMouseMove = time.Now().Unix()
			break
		}
	}

	// Fix for exiting from events
	if m.State != 2 {
		return
	}

	ControlsActive = m.GetControlsActive()

	m.renderView()
	m.logicControls()
	if ControlsActive {
		sdl.ShowCursor(sdl.ENABLE)
		m.renderControls()
		if m.Player.MenuOpen {
			m.renderSettings()
		}
	} else {
		m.Player.MenuOpen = false
		sdl.ShowCursor(sdl.DISABLE)
	}
	m.Renderer.Present()
	m.Player.MPV_GL.ReportFlip(0)

	if m.Player.Complete() {
		m.PlayingAnime.Episode = fmt.Sprintf("%d", m.PlayingEpisode)
		m.PlayAnime(m.PlayingAnime)
		oblivion.ReportEpisodeWatchTime(m.PlayingAnime.ID, m.PlayingAnime.Episode, 0)
		oblivion.SetCurrentWatchingEpisode(m.PlayingAnime.ID, m.PlayingAnime.Episode)
		m.Player.ShowMessage(fmt.Sprintf("Playing %d/%s", m.PlayingEpisode, m.PlayingAnime.Episodes))
		log.Println("Setting Episode")
		m.Player.LastPosition = 0
		return
	}

	go func() {
		position, err := m.Player.GetPosition()
		if err != nil {
			return
		}

		m.Player.LastPosition = position

		if position-m.PlayingTimer > 5 {
			m.PlayingTimer = position
			if err := oblivion.ReportEpisodeWatchTime(m.PlayingAnime.ID, fmt.Sprintf("%d", m.PlayingEpisode), position); err != nil {
				log.Println(err)
			}
		}
	}()
}
