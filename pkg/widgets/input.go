package widgets

import "github.com/awesome-gocui/gocui"

type Input struct {
	name      string
	x, y      int
	w         int
	maxLength int
}

func NewInput(name string, maxX, maxY, maxLength int, ctx ctx) (*Input, error) {
	x1 := maxX/2 - 50
	x2 := maxX/2 + 50
	y1 := maxY/2 - 1

	if err := ctx.SetKeyBinding("remote", gocui.KeyEnter, gocui.ModNone, func(gui *gocui.Gui, view *gocui.View) error {
		if len(view.BufferLines()) > 0 {
			if image, _ := view.Line(0); image != "" {
				ctx.ScanImage(image)
			}
		}
		gui.Mouse = true
		gui.Cursor = false
		return gui.DeleteView("remote")
	}); err != nil {
		return nil, err
	}

	if err := ctx.SetKeyBinding("remote", gocui.KeyEsc, gocui.ModNone, func(gui *gocui.Gui, view *gocui.View) error {
		gui.Mouse = true
		gui.Cursor = false
		return gui.DeleteView("remote")
	}); err != nil {
		return nil, err
	}

	return &Input{name: name, x: x1, y: y1, w: x2, maxLength: maxLength}, nil
}

func (i *Input) Layout(g *gocui.Gui) error {
	v, err := g.SetView(i.name, i.x, i.y, i.w, i.y+2, 0)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = " Enter remote image name "
		v.Editor = i
		v.Editable = true
		g.Cursor = true
		v.TitleColor = gocui.ColorGreen
		v.FrameColor = gocui.ColorGreen
	}
	return nil
}

func (i *Input) Edit(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	cx, _ := v.Cursor()
	ox, _ := v.Origin()
	limit := ox+cx+1 > i.maxLength
	switch {
	case ch != 0 && mod == 0 && !limit:
		v.EditWrite(ch)
	case key == gocui.KeySpace && !limit:
		v.EditWrite(' ')
	case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
		v.EditDelete(true)
	}
}
