package manager

import (
	"log"

	"github/itzaname/oblivion_player/oblivion"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	BoxWidth    = int32(400)
	BoxHeight   = int32(195)
	InputWidth  = int32(360)
	InputHeight = int32(35)
)

var (
	UsernameInput  = false
	PasswordInput  = false
	TextInput      = false
	UsernameBuffer = ""
	PasswordBuffer = ""
	CursorState    = false
	LoggedIn       = false
	Failure        = false
	Waiting        = false
)

func (m *Instance) renderLogin() {
	w, h := m.Window.GetSize()

	loginx := int32(w)/2 - BoxWidth/2
	loginy := int32(h)/2 - BoxHeight/2

	m.Renderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND)
	m.Renderer.SetDrawColor(43, 46, 53, 255)
	m.Renderer.Clear()

	// Login Box
	m.Renderer.SetDrawColor(37, 40, 48, 255)
	m.drawFilledBox(loginx, loginy, BoxWidth, BoxHeight)

	// Error
	if Failure {
		m.Renderer.SetDrawColor(200, 125, 125, 255)
		m.drawFilledBox(loginx, loginy-30, BoxWidth, 30)
		m.drawTextCentered(loginx, loginy-30, BoxWidth, 30, "Login failed.", 20, 20, 20, "helvetica.ttf", 18)
	}

	// Login Text
	m.drawTextCentered(loginx, loginy+15, BoxWidth, 22, "MAL Login", 221, 221, 221, "helvetica.ttf", 32)

	// Inputs
	if UsernameInput {
		m.Renderer.SetDrawColor(220, 220, 220, 255)
	} else {
		m.Renderer.SetDrawColor(200, 200, 200, 255)
	}
	m.drawFilledBox(loginx+20, loginy+50, InputWidth, InputHeight)
	if PasswordInput {
		m.Renderer.SetDrawColor(220, 220, 220, 255)
	} else {
		m.Renderer.SetDrawColor(200, 200, 200, 255)
	}
	m.drawFilledBox(loginx+20, loginy+50+InputHeight+10, InputWidth, InputHeight)

	m.drawTextCentered(loginx+20, loginy+50, InputWidth, InputHeight, UsernameBuffer, 40, 40, 40, "helvetica.ttf", 18)
	passwordtext := ""
	for i := 0; i < len(PasswordBuffer); i++ {
		passwordtext += "*"
	}
	m.drawTextCentered(loginx+20, loginy+50+InputHeight+10, InputWidth, InputHeight, passwordtext, 40, 40, 40, "helvetica.ttf", 18)

	// Login Button
	m.Renderer.SetDrawColor(69, 74, 88, 255)
	m.drawFilledBox(loginx+20, loginy+50+InputHeight*2+20, InputWidth, InputHeight)
	m.drawTextCentered(loginx+20, loginy+50+InputHeight*2+20, InputWidth, InputHeight, "Login", 200, 200, 200, "helvetica.ttf", 18)

	m.Renderer.Present()
	sdl.Delay(1000 / FrameRate)
}

func (m *Instance) logicLogin() {
	w, h := m.Window.GetSize()
	mx, my, state := sdl.GetMouseState()

	loginx := int32(w)/2 - BoxWidth/2
	loginy := int32(h)/2 - BoxHeight/2
	if (state & sdl.Button(sdl.BUTTON_LEFT)) != 0 {
		if isInBounds(loginx+20, loginy+50, InputWidth, InputHeight, int32(mx), int32(my)) {
			if !TextInput {
				TextInput = true
				sdl.StartTextInput()
			}
			UsernameInput = true
			PasswordInput = false
		}

		if isInBounds(loginx+20, loginy+50+InputHeight+10, InputWidth, InputHeight, int32(mx), int32(my)) {
			if !TextInput {
				TextInput = true
				sdl.StartTextInput()
			}
			UsernameInput = false
			PasswordInput = true
		}

		if isInBounds(loginx+20, loginy+50+InputHeight*2+20, InputWidth, InputHeight, int32(mx), int32(my)) && !Waiting {
			Waiting = true
			go func() {
				if err := oblivion.Login(UsernameBuffer, PasswordBuffer); err != nil {
					log.Println(err)
					Failure = true
					Waiting = false
					return
				}
				settings, _ := oblivion.GetUserSettings()
				m.Settings = settings

				log.Println("Logged in")
				sdl.Do(func() {
					sdl.StopTextInput()
					m.SetCursor(sdl.SYSTEM_CURSOR_ARROW)
				})
				m.State = 1
				Waiting = false
			}()
		}
	}

	if isInBounds(loginx+20, loginy+50+InputHeight*2+20, InputWidth, InputHeight, int32(mx), int32(my)) {
		m.SetCursor(sdl.SYSTEM_CURSOR_HAND)
		CursorState = true
	} else if CursorState {
		m.SetCursor(sdl.SYSTEM_CURSOR_ARROW)
	}
}

func (m *Instance) LoginView() {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			m.Active = false
			return
		case *sdl.KeyUpEvent:
			if TextInput {
				if t.Keysym.Sym == sdl.K_BACKSPACE {
					if UsernameInput && len(UsernameBuffer) > 0 {
						UsernameBuffer = UsernameBuffer[:len(UsernameBuffer)-1]
					}
					if PasswordInput && len(PasswordBuffer) > 0 {
						PasswordBuffer = PasswordBuffer[:len(PasswordBuffer)-1]
					}
				}
				//log.Println(sdl.GetModState())
				if t.Keysym.Sym == sdl.K_v && (sdl.GetModState()&sdl.KMOD_CTRL) != 0 {
					text, err := sdl.GetClipboardText()
					if err != nil {
						break
					}
					if UsernameInput {
						UsernameBuffer += text
					}
					if PasswordInput {
						PasswordBuffer += text
					}
				}
				break
			}
			break
		case *sdl.TextInputEvent:
			if UsernameInput {
				UsernameBuffer += string(t.Text[0])
			}
			if PasswordInput {
				PasswordBuffer += string(t.Text[0])
			}
			break
		case *sdl.TextEditingEvent:
			log.Println(t.Text)
			log.Println("meme")
			break
		}
	}

	if LoggedIn {
		return
	}
	m.logicLogin()
	m.renderLogin()
}
