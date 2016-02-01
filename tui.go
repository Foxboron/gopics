package main

import (
	"fmt"
	gc "github.com/gbin/goncurses"
	"os/exec"
)

type Interface struct {
	files   *gc.Menu
	picture *gc.Window
	stdscr  *gc.Window
	tmp     *gc.Panel
}

type Interfaces interface {
	createMenu()
	Refresh()
	w3imgDisplay()
	GetChar()
	Redraw()
}

func (w Interface) Action(act int) {
	w.files.Driver(act)
}

func (w Interface) Redraw() {
	// TODO: Fix
	w.files.Window().Erase()
	w.picture.Erase()
	w.Refresh()
}

func (w Interface) GetChar() gc.Key {
	return w.files.Window().GetChar()
}

func (w Interface) Refresh() {
	w.w3imgDisplay()
	w.files.Window().Refresh()
	w.picture.Refresh()
	w.stdscr.Refresh()
}

func setDefaultOptions() {
	gc.Raw(true)
	gc.Echo(false)
	gc.Cursor(0)
}

func createPictureView(stdscr *gc.Window) *gc.Window {

	rows, cols := stdscr.MaxYX()
	pic_view, _ := gc.NewWindow(rows, cols/2, 0, cols/2)
	pic_view.Box(0, 0)
	// pic_view.MovePrintf(y+1, x+1, "TEST")
	gc.NewPanel(pic_view)

	y, _ := stdscr.MaxYX()
	stdscr.MovePrint(y-2, 1, "'q' to exit")
	stdscr.Refresh()
	pic_view.Refresh()
	return pic_view

}
func createFileList(stdscr *gc.Window) *gc.Menu {

	// build the menu items
	menu_items := getFiles()
	items := make([]*gc.MenuItem, len(menu_items))
	i := 0
	for key, val := range menu_items {
		items[i], _ = gc.NewItem(key, val)
		// defer items[i].Free()
		i++
	}
	menu, _ := gc.NewMenu(items)
	defer menu.Free()
	rows, cols := stdscr.MaxYX()

	// menuwin, _ := gc.NewWindow(rows, cols/2, 3, 1)
	menuwin, _ := gc.NewWindow(rows, cols/2, 0, 0)
	menuwin.Keypad(true)

	menu.SetWindow(menuwin)
	dwin := menuwin.Derived(29, 50, 3, 1)
	menu.SubWindow(dwin)
	menu.Option(gc.O_SHOWDESC, true)
	menu.Mark(" * ")
	menu.Option(gc.O_ONEVALUE, false)

	// Print centered menu title
	_, x := menuwin.MaxYX()
	title := "My Menu"
	menuwin.Box(0, 0)
	menuwin.ColorOn(1)
	menuwin.MovePrint(1, (x/2)-(len(title)/2), title)
	menuwin.ColorOff(1)
	menuwin.MoveAddChar(2, 0, gc.ACS_LTEE)
	menuwin.HLine(2, 1, gc.ACS_HLINE, x-3)
	menuwin.MoveAddChar(2, x-2, gc.ACS_RTEE)

	gc.NewPanel(menuwin)

	menu.Post()
	//defer menu.UnPost()
	menu.Window().Refresh()
	return menu
}

func (win Interface) w3imgDisplay() {
	// stdscr, _ := gc.Init()
	// defer gc.End()
	// _, x := stdscr.MaxYX()
	wm := win.picture
	_, cols := wm.YX()
	_, maxcols := wm.MaxYX()
	// getDismensions()
	pixels := 8
	// print(cols)
	// println()
	// print(pixels * cols)
	w := maxcols * 7
	h := float64(w) / 1.5

	path := win.files.Current(nil).Description()
	cmd := fmt.Sprintf("echo -e '0;1;%d;%d;%d;%d;;;;;%s\\n4;\\n3;' | /usr/lib/w3m/w3mimgdisplay", (cols+3)*pixels, 50, w, int(h), path)
	// print(cmd)
	out, _ := exec.Command("bash", "-c", cmd).Output()
	fmt.Printf("%s", out)
}

func createInterface() Interface {

	stdscr, _ := gc.Init()
	stdscr.Refresh()

	setDefaultOptions()
	menu := createFileList(stdscr)
	picture := createPictureView(stdscr)
	in := Interface{files: menu, picture: picture, stdscr: stdscr}
	in.Refresh()
	return in
}

func (w Interface) getInput(msg string) string {
	input, _ := gc.NewWindow(3, 30, 20, 10)
	input.Box(0, 0)
	input.MovePrint(1, 1, msg)
	panel := gc.NewPanel(input)
	panel.Top()

	gc.Raw(false)
	gc.Echo(true)
	gc.Cursor(1)
	str, err := input.GetString(10)
	if err != nil {
		println("Input error")
	}

	panel.Delete()
	setDefaultOptions()
	return str
}

func (win Interface) createTempWindow(text string, writeable bool, h int, w int) *gc.Panel {
	// gc.Raw(false)
	// gc.Echo(true)
	// gc.Cursor(1)
	_, cols := win.files.Window().MaxYX()

	input, _ := gc.NewWindow(11, cols, 19, 0)
	input.Box(0, 0)
	panel := gc.NewPanel(input)
	panel.Top()
	if len(text) > 0 {
		input.MovePrint(1, 1, text)
	}
	if writeable {
		gc.Raw(false)
		gc.Echo(true)
		gc.Cursor(1)
		_, err := input.GetString(10)
		if err != nil {
			println("Input error")
		}
		setDefaultOptions()
		input.Refresh()
		input.Erase()
		panel.Hide()
		input.Refresh()
		win.Refresh()
	} else {
		input.Refresh()
		panel.Window().ScrollOk(true)
		win.tmp = panel
	}

	return win.tmp
}
