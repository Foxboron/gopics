package main

import (
	"fmt"
	gc "github.com/gbin/goncurses"
	"os/exec"
)

func w3imgDisplay(wm *gc.Window, path string) {
	// stdscr, _ := gc.Init()
	// defer gc.End()
	// _, x := stdscr.MaxYX()
	_, cols := wm.YX()
	_, maxcols := wm.MaxYX()
	// getDismensions()
	pixels := 8
	// print(cols)
	// println()
	// print(pixels * cols)
	w := maxcols * 7
	h := float64(w) / 1.5
	cmd := fmt.Sprintf("echo -e '0;1;%d;%d;%d;%d;;;;;%s\\n4;\\n3;' | /usr/lib/w3m/w3mimgdisplay", (cols+3)*pixels, 50, w, int(h), path)
	// print(cmd)
	out, _ := exec.Command("bash", "-c", cmd).Output()
	fmt.Printf("%s", out)
}

func main() {
	stdscr, _ := gc.Init()
	defer gc.End()
	gc.Raw(true)
	gc.Echo(false)
	gc.Cursor(0)
	// stdscr.Keypad(true)

	// build the menu items
	menu_items := getFiles()
	items := make([]*gc.MenuItem, len(menu_items))
	i := 0
	for key, val := range menu_items {
		items[i], _ = gc.NewItem(key, val)
		defer items[i].Free()
		i++
	}

	// create the menu
	menu, _ := gc.NewMenu(items)
	defer menu.Free()
	rows, cols := stdscr.MaxYX()

	menuwin, _ := gc.NewWindow(rows, cols/2, 0, 0)
	menuwin.Keypad(true)

	menu.SetWindow(menuwin)
	dwin := menuwin.Derived(20, 40, 3, 1)
	menu.SubWindow(dwin)
	menu.Mark(" * ")
	menu.Option(gc.O_ONEVALUE, false)

	// Print centered menu title
	y, x := menuwin.MaxYX()
	title := "My Menu"
	menuwin.Box(0, 0)
	menuwin.ColorOn(1)
	menuwin.MovePrint(1, (x/2)-(len(title)/2), title)
	menuwin.ColorOff(1)
	menuwin.MoveAddChar(2, 0, gc.ACS_LTEE)
	menuwin.HLine(2, 1, gc.ACS_HLINE, x-3)
	menuwin.MoveAddChar(2, x-2, gc.ACS_RTEE)

	gc.NewPanel(menuwin)

	pic_view, _ := gc.NewWindow(rows, cols/2, 0, cols/2)
	pic_view.Box(0, 0)
	// pic_view.MovePrintf(y+1, x+1, "TEST")
	gc.NewPanel(pic_view)

	y, x = stdscr.MaxYX()
	stdscr.MovePrint(y-2, 1, "'q' to exit")
	stdscr.Refresh()

	menu.Post()
	defer menu.UnPost()
	menuwin.Refresh()
	pic_view.Refresh()

	for {
		gc.Update()
		ch := menuwin.GetChar()

		switch ch {
		case 'q':
			return

		case gc.KEY_DOWN, 'j':
			menu.Driver(gc.REQ_DOWN)
			item := menu.Current(nil)
			w3imgDisplay(pic_view, item.Description())
			stdscr.Refresh()
			menuwin.Refresh()

		// Upload
		case 'u':

			msg := "Album: "
			row, col := stdscr.MaxYX()
			row, col = (row/2)-1, (col-len(msg))/2

			gc.Raw(false)
			gc.Echo(true)
			gc.Cursor(1)
			input, _ := gc.NewWindow(3, 30, 20, 10)
			input.Box(0, 0)
			input.MovePrint(1, 1, msg)
			panel := gc.NewPanel(input)
			panel.Top()
			str, err := input.GetString(10)
			if err != nil {
				println("error")
			}
			panel.Delete()

			gc.Raw(true)
			gc.Echo(false)
			gc.Cursor(0)
			pic_view.Refresh()
			menuwin.Refresh()
			stdscr.Refresh()

			var list []string
			for _, item := range menu.Items() {
				if item.Value() {
					list = append(list, item.Description())
				}
			}
			r := createRsync()
			r = r.addFiles(list)
			str = r.Upload(str)
			cmdUpload(str)
		// Select all
		case 'V':
			for _, item := range menu.Items() {
				item.Selectable(true)
			}
		case gc.KEY_UP, 'k':
			menu.Driver(gc.REQ_UP)
			item := menu.Current(nil)
			w3imgDisplay(pic_view, item.Description())
			stdscr.Refresh()
			menuwin.Refresh()
		case ' ':
			menu.Driver(gc.REQ_TOGGLE)
			menu.Driver(gc.REQ_DOWN)

		case gc.KEY_ENTER, gc.KEY_RETURN:
		default:
			menu.Driver(gc.DriverActions[ch])
		}
	}
}
