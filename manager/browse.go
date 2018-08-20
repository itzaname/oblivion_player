package manager

import (
	"log"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	PosterW int32 = 225
	PosterH int32 = 315
	Margin  int32 = 20
)

var (
	Animating                    = false
	AnimateTarget        int32   = 0
	AnimateStart         int32   = 1
	AnimateColorTarget   int32   = 0
	AnimateColorStart    int32   = 0
	AnimateStartTime     float64 = 0
	AnimateSceneComplete         = false
	ErrorString                  = ""
	ErrorEnd             int64   = 0
	ScrollOffset         int32   = 0
	ScrollBar                    = false
	ScrollHeight         int32   = 0
	ScrollDrag                   = false
	LoadAnime                    = true
)

func getTotalHeight(items, collums int32) int32 {
	return 10 + (PosterH+Margin)*((items-1)/collums) + PosterH
}

func AnimateHide(w int32) {
	Animating = true
	AnimateStart = 0
	AnimateTarget = int32(w)
	AnimateColorTarget = 255
	AnimateColorStart = 0
	AnimateStartTime = float64(time.Now().UnixNano()) / float64(time.Second)
	AnimateSceneComplete = true
}

func AnimateShow(w int32) {
	Animating = true
	AnimateStart = 0
	AnimateTarget = 0
	AnimateColorTarget = 0
	AnimateColorStart = 0
	AnimateStartTime = float64(time.Now().UnixNano()) / float64(time.Second)
	AnimateSceneComplete = false
}

func DisplayBrowseError(text string, length int) {
	ErrorString = text
	ErrorEnd = time.Now().Unix() + int64(length)
}

func (m *Instance) renderWatching() {
	w, h := m.Window.GetSize()
	mx, my, _ := sdl.GetMouseState()

	collums := int32(w) / (PosterW + Margin)
	if collums <= 0 {
		collums = 1
	}

	xoffset := animate(1, AnimateStartTime, AnimateStart, AnimateTarget)
	aoffset := animate(1, AnimateStartTime, AnimateColorStart, AnimateColorTarget)
	if Animating && xoffset == AnimateTarget {
		Animating = false
		if AnimateSceneComplete {
			m.State = 2
		}
		return
	}

	theight := getTotalHeight(int32(len(m.Watching)), collums)
	if theight > int32(h) {
		ScrollHeight = theight + Margin
		ScrollBar = true
		w -= 15
	} else {
		ScrollBar = false
		ScrollOffset = 0
	}

	width := int32(w) / collums

	m.Renderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND)
	m.Renderer.SetDrawColor(43, 46, 53, 255)
	m.Renderer.Clear()

	if Animating {
		m.Renderer.SetDrawColor(0, 0, 0, uint8(aoffset))
		m.drawFilledBox(0, 0, int32(w), int32(h))
	}

	m.Renderer.SetDrawColor(0, 0, 0, 150)
	for i := int32(0); i < int32(len(m.Watching)); i++ {
		texture, err := m.GetTexture(m.Watching[int(i)].ID)
		if err != nil {
			continue
		}

		m.Renderer.Copy(texture, nil, &sdl.Rect{
			X: xoffset + (i%collums)*width + ((width - PosterW) / 2),
			Y: -ScrollOffset + 10 + (PosterH+Margin)*(i/collums),
			W: PosterW,
			H: PosterH,
		})

		if isInBounds(xoffset+(i%collums)*width+((width-PosterW)/2), -ScrollOffset+10+(PosterH+Margin)*(i/collums), PosterW, PosterH, int32(mx), int32(my)) {
			m.drawFilledBox(xoffset+(i%collums)*width+((width-PosterW)/2), -ScrollOffset+10+(PosterH+Margin)*(i/collums), PosterW, PosterH)
		}

		m.drawTextBGPadded(
			xoffset+(i%collums)*width+((width-PosterW)/2)+5,
			-ScrollOffset+10+(PosterH+Margin)*(i/collums)+5,
			m.Watching[int(i)].Episode+"/"+m.Watching[int(i)].Episodes,
			sdl.Color{240, 240, 240, 255},
			sdl.Color{0, 0, 0, 220},
			4,
			2,
			"helvetica.ttf",
			10)

	}

	if ScrollBar {
		m.Renderer.SetDrawColor(40, 40, 40, 255)
		m.drawFilledBox(int32(w), 0, 15, int32(h))
		m.Renderer.SetDrawColor(60, 60, 60, 255)
		scrollsize := int32(h) * int32(h) / theight
		m.drawFilledBox(int32(w), int32(float64(h)*(float64(ScrollOffset)/float64(theight))), 15, scrollsize)
	}

	if time.Now().Unix() < ErrorEnd {
		errorw := int32(500)
		errorh := int32(30)
		m.Renderer.SetDrawColor(60, 60, 60, 255)
		m.drawFilledBox(int32(w)/2-errorw/2, 20, errorw, errorh)
		m.drawTextCentered(int32(w)/2-errorw/2, 20, errorw, errorh, ErrorString, 255, 255, 255, "helvetica.ttf", 18)
	}

	m.Renderer.Present()
	sdl.Delay(1000 / FrameRate)
}

func (m *Instance) logicWatching() {

	if Animating {
		return
	}
	w, h := m.Window.GetSize()
	if AnimateTarget != 0 {
		AnimateShow(int32(w))
	}
	mx, my, state := sdl.GetMouseState()
	if ScrollHeight > 0 {
		if (state&sdl.Button(sdl.BUTTON_LEFT)) != 0 && ScrollBar && isInBounds(int32(w-15), 0, 15, int32(h), int32(mx), int32(my)) {
			ScrollDrag = true
		} else if state&sdl.Button(sdl.BUTTON_LEFT) == 0 {
			ScrollDrag = false
		}

		if ScrollDrag {
			scrollsize := int32(h) * int32(h) / ScrollHeight
			ScrollOffset = int32(float64(ScrollHeight) * (float64(int32(my)-scrollsize/2) / float64(h)))
		}

		if ScrollOffset < 0 {
			ScrollOffset = 0
		}
		if ScrollOffset > ScrollHeight-int32(h) {
			ScrollOffset = ScrollHeight - int32(h)
		}
	}

	handcursor := false

	collums := int32(w) / (PosterW + Margin)
	if collums <= 0 {
		collums = 1
	}

	width := int32(w) / collums

	for i := int32(0); i < int32(len(m.Watching)); i++ {
		if isInBounds((i%collums)*width+((width-PosterW)/2), -ScrollOffset+10+(PosterH+Margin)*(i/collums), PosterW, PosterH, int32(mx), int32(my)) {
			handcursor = true

		}
	}
	if handcursor {
		m.SetCursor(sdl.SYSTEM_CURSOR_HAND)
	} else {
		m.SetCursor(sdl.SYSTEM_CURSOR_ARROW)
	}
}

func (m *Instance) logicClickWatching() {
	w, _ := m.Window.GetSize()
	if AnimateTarget != 0 {
		AnimateShow(int32(w))
	}
	mx, my, _ := sdl.GetMouseState()

	collums := int32(w) / (PosterW + Margin)
	if collums <= 0 {
		collums = 1
	}

	width := int32(w) / collums

	for i := int32(0); i < int32(len(m.Watching)); i++ {
		if isInBounds((i%collums)*width+((width-PosterW)/2), -ScrollOffset+10+(PosterH+Margin)*(i/collums), PosterW, PosterH, int32(mx), int32(my)) {
			go m.PlayAnime(m.Watching[i])
			AnimateHide(int32(w))

		}
	}
}

func (m *Instance) BrowseView() {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			m.Active = false
			return
		case *sdl.MouseWheelEvent:
			if ScrollBar {
				ScrollOffset += t.Y * -53
			}
		case *sdl.MouseButtonEvent:
			if t.State == sdl.PRESSED {
				m.logicClickWatching()
			}
			break
		}
	}
	if LoadAnime {
		LoadAnime = false
		go func() {
			log.Println(m.LoadAnime())
		}()
	}
	m.logicWatching()
	m.renderWatching()
}
