package player

import (
	"fmt"
	"strconv"

	"github.com/YouROK/go-mpv/mpv"
	"github.com/veandco/go-sdl2/sdl"
)

type Instance struct {
	Window       *sdl.Window
	Renderer     *sdl.Renderer
	MPV          *mpv.Mpv
	MPV_GL       *mpv.MpvGL
	Menu         []Menu
	MenuOpen     bool
	LastPosition float64
}

func New(window *sdl.Window, renderer *sdl.Renderer) (Instance, error) {
	m := mpv.Create()
	m.SetOptionString("vo", "opengl-cb")

	err := m.Initialize()
	if err != nil {
		return Instance{}, err
	}

	mgl := m.GetSubApiGL()
	if mgl == nil {
		return Instance{}, err
	}

	err = mgl.InitGL()
	if mgl == nil {
		return Instance{}, err
	}

	return Instance{
		Window:   window,
		Renderer: renderer,
		MPV:      m,
		MPV_GL:   mgl,
	}, nil
}

func (p *Instance) ResetPlayer() error {
	err := p.MPV_GL.UninitGL()
	if err != nil {
		return err
	}
	p.MPV.TerminateDestroy()

	m := mpv.Create()
	m.SetOptionString("vo", "opengl-cb")

	err = m.Initialize()
	if err != nil {
		return err
	}

	mgl := m.GetSubApiGL()
	if mgl == nil {
		return err
	}

	err = mgl.InitGL()
	if mgl == nil {
		return err
	}

	p.MPV_GL = mgl
	p.MPV = m

	return nil
}

func (p *Instance) Render() {
	w, h := p.Window.GetSize()
	p.MPV_GL.Draw(0, w, h)
}

func (p *Instance) LoadVideo(video string) {
	for p.MPV.Command([]string{"loadfile", video}) != nil {
	}
}

func (p *Instance) LoadAudio(audio string) {
	for p.MPV.Command([]string{"audio-add", audio, "cached"}) != nil {
	}
}

func (p *Instance) LoadSubs(subs string) {
	for p.MPV.Command([]string{"sub-add", subs, "cached"}) != nil {
	}
}

func (p *Instance) TogglePause() error {
	return p.MPV.Command([]string{"cycle", "pause"})
}

func (p *Instance) GetDuration() (float64, error) {
	return strconv.ParseFloat(p.MPV.GetPropertyString("duration"), 64)
}

func (p *Instance) GetPosition() (float64, error) {
	return strconv.ParseFloat(p.MPV.GetPropertyString("time-pos"), 64)
}

func (p *Instance) SetPosition(time float64) error {
	return p.MPV.SetPropertyString("time-pos", fmt.Sprintf("%f", time))
}

func (p *Instance) SetPositionForce(time float64) {
	for p.MPV.SetPropertyString("time-pos", fmt.Sprintf("%f", time)) != nil {
	}
}

func (p *Instance) ShowMessage(text string) {
	p.MPV.Command([]string{"show-text", text, "5000"})
}

func (p *Instance) Stop() {
	p.MPV.CommandString("stop")
}

func (p *Instance) Play() {
	p.MPV.SetProperty("pause", mpv.FORMAT_FLAG, false)
}

func (p *Instance) Pause() {
	p.MPV.SetProperty("pause", mpv.FORMAT_FLAG, true)
}

func (p *Instance) Paused() bool {
	paused, err := p.MPV.GetProperty("pause", mpv.FORMAT_FLAG)
	if err != nil {
		return false
	}

	return paused.(bool)
}

func (p *Instance) AddMenu(name string, Callback func(interface{}), items []SubItem) {
	p.Menu = append(p.Menu, Menu{items, name, Callback, false})
}

func (p *Instance) ClearMenu() {
	p.Menu = []Menu{}
}

func (p *Instance) Complete() bool {
	if p.LastPosition < 5 {
		return false
	}
	active, err := p.MPV.GetProperty("idle-active", mpv.FORMAT_FLAG)
	if err != nil {
		return false
	}

	return active.(bool)
}
