package main

import (
	"fmt"
	gc "github.com/gbin/goncurses"
	"os/exec"
)

func w3imgDisplay(wm *gc.Window, path string) {
	_, cols := wm.YX()
	cmd := fmt.Sprintf("echo -e '0;1;%d;%d;400;300;;;;;%s\\n4;\\n3;' | /usr/lib/w3m/w3mimgdisplay", cols*9, 100, path)
	// print(cmd)
	out, _ := exec.Command("bash", "-c", cmd).Output()
	fmt.Printf("%s", out)
}
func main() {
	stdscr, _ := gc.Init()
	defer gc.End()

	// gc.Raw(true)
	// gc.Echo(false)
	gc.Cursor(0)
	stdscr.Keypad(true)
	gc.InitPair(1, gc.C_RED, gc.C_BLACK)

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
	dwin := menuwin.Derived(60, 50, 3, 1)
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
		case gc.KEY_DOWN:
			menu.Driver(gc.REQ_DOWN)
			item := menu.Current(nil)
			w3imgDisplay(pic_view, item.Description())
			stdscr.Refresh()
			menuwin.Refresh()
		case gc.KEY_UP:
			menu.Driver(gc.REQ_UP)
			item := menu.Current(nil)
			w3imgDisplay(pic_view, item.Description())
			stdscr.Refresh()
			menuwin.Refresh()
		case ' ':
			menu.Driver(gc.REQ_TOGGLE)
		case gc.KEY_ENTER, gc.KEY_RETURN:
			var list string
			for _, item := range menu.Items() {
				if item.Value() {
					list += "\"" + item.Name() + "\" "
					println(list)
				}
			}
			stdscr.Move(20, 0)
			stdscr.ClearToEOL()
			stdscr.MovePrint(20, 40, list)
			stdscr.Refresh()
		default:
			menu.Driver(gc.DriverActions[ch])
		}
	}
}
