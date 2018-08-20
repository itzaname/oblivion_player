package manager

import (
	"github.com/veandco/go-sdl2/sdl"
)

func (m *Instance) drawFilledBox(x, y, w, h int32) {
	rc := &sdl.Rect{
		X: x,
		Y: y,
		W: w,
		H: h,
	}
	m.Renderer.FillRect(rc)
}

func (m *Instance) drawText(x, y int32, text string, r, g, b uint8, fontname string, size int) error {
	font, err := m.GetFont(fontname, size)
	if err != nil {
		return err
	}

	surface, err := font.RenderUTF8_Blended(text, sdl.Color{r, g, b, 255})
	if err != nil {
		return err
	}
	defer surface.Free()

	texture, err := m.Renderer.CreateTextureFromSurface(surface)
	if err != nil {
		return err
	}
	defer texture.Destroy()

	rc := &sdl.Rect{
		X: x,
		Y: y,
		W: surface.W,
		H: surface.H,
	}

	return m.Renderer.Copy(texture, nil, rc)
}

func (m *Instance) drawTextBG(x, y int32, text string, fg sdl.Color, bg sdl.Color, fontname string, size int) error {
	m.Renderer.SetDrawColor(0, 0, 0, 0)

	font, err := m.GetFont(fontname, size)
	if err != nil {
		return err
	}

	surface, err := font.RenderUTF8_Shaded(text, fg, bg)
	if err != nil {
		return err
	}
	defer surface.Free()

	texture, err := m.Renderer.CreateTextureFromSurface(surface)
	if err != nil {
		return err
	}
	defer texture.Destroy()

	_, _, width, height, err := texture.Query()
	if err != nil {
		return err
	}

	rc := &sdl.Rect{
		X: x,
		Y: y,
		W: width,
		H: height,
	}
	m.Renderer.SetDrawColor(0, 0, 0, 0)

	return m.Renderer.Copy(texture, nil, rc)
}

func (m *Instance) drawTextBGPadded(x, y int32, text string, fg sdl.Color, bg sdl.Color, paddingx, paddingy int32, fontname string, size int) error {
	m.Renderer.SetDrawColor(0, 0, 0, 0)

	font, err := m.GetFont(fontname, size)
	if err != nil {
		return err
	}

	surface, err := font.RenderUTF8_Blended(text, fg)
	if err != nil {
		return err
	}
	defer surface.Free()

	texture, err := m.Renderer.CreateTextureFromSurface(surface)
	if err != nil {
		return err
	}
	defer texture.Destroy()

	_, _, width, height, err := texture.Query()
	if err != nil {
		return err
	}

	rc := &sdl.Rect{
		X: x,
		Y: y,
		W: width,
		H: height,
	}
	m.Renderer.SetDrawColor(bg.R, bg.G, bg.B, bg.A)
	m.drawFilledBox(x-paddingx, y-paddingy, width+paddingx*2, height+paddingy*2)

	return m.Renderer.Copy(texture, nil, rc)
}

func (m *Instance) drawTextCentered(x, y, w, h int32, text string, r, g, b uint8, fontname string, size int) error {
	//

	font, err := m.GetFont(fontname, size)
	if err != nil {
		return err
	}

	surface, err := font.RenderUTF8_Blended(text, sdl.Color{r, g, b, 255})
	if err != nil {
		return err
	}
	defer surface.Free()

	//gl.ReadPixels(0, 0, winw, winh, gl.RGBA, gl.UNSIGNED_BYTE)

	texture, err := m.Renderer.CreateTextureFromSurface(surface)
	if err != nil {
		return err
	}
	defer texture.Destroy()

	_, _, width, height, err := texture.Query()
	if err != nil {
		return err
	}

	rc := &sdl.Rect{
		X: x + w/2 - width/2,
		Y: y + h/2 - height/2,
		W: width,
		H: height,
	}

	return m.Renderer.Copy(texture, nil, rc)
}

func (m *Instance) getTextSize(text string, fontname string, size int) (int32, int32, error) {
	font, err := m.GetFont(fontname, size)
	if err != nil {
		return 0, 0, err
	}

	surface, err := font.RenderUTF8_Blended(text, sdl.Color{255, 255, 255, 255})
	if err != nil {
		return 0, 0, err
	}
	defer surface.Free()

	return surface.W, surface.H, nil
}
