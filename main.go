package main

import (
	gc "github.com/gbin/goncurses"
)

// func writeFiles(w *gc.Window, f []File) {
// 	y, x := w.YX()
// 	maxy, _ := w.MaxYX()
// 	y++
// 	for _, i := range f {
// 		w.MovePrintf(y, x+1, i.getName())
// 		if y+2 == maxy {
// 			break
// 		}
// 		y++
// 	}

// }

func writeFiles(f []File) []*gc.MenuItem {
	ret := make([]*gc.MenuItem, len(f))
	for i, val := range f {
		ret[i], _ = gc.NewItem(val.getName(), "")
	}
	return ret
}

func initWindows(panels map[string]*gc.Panel) map[string]*gc.Panel {
	stdscr, _ := gc.Init()
	defer gc.End()
	rows, cols := stdscr.MaxYX()
	// y, x := stdscr.YX()

	f := getFiles()
	items := writeFiles(f)

	file_list, _ := gc.NewMenu(items)
	// menuwin, _ := gc.NewWindow(rows, cols/2, 0, 0)
	// menuwin.Keypad(true)

	// menuwin.Box(0, 0)
	// file_list.SetWindow(menuwin)
	// dwin := menuwin.Derived(10, 10, 20, 20)
	// file_list.SubWindow(dwin)
	// file_list.Format(5, 1)
	// file_list.Mark(" * ")
	file_list.Post()
	defer file_list.UnPost()
	file_list.Window().Refresh()

	// dwin := menuwin.Derived(6, 38, 3, 1)
	// file_list.SubWindow(dwin)
	// file_list.Mark(" * ")

	file_pan := gc.NewPanel(file_list.Window())
	panels["files"] = file_pan

	pic_view, _ := gc.NewWindow(rows, cols/2, 0, cols/2)
	pic_view.Box(0, 0)
	// pic_view.MovePrintf(y+1, x+1, "TEST")
	pic_pan := gc.NewPanel(pic_view)
	panels["picture"] = pic_pan
	return panels

}

func Update() {
	gc.UpdatePanels()
	gc.Update()
}
func drawEverything(panels_inn map[string]*gc.Panel) {
	// var panels = initWindows(panels_inn)
	// Update()
	// panels["picture"].Window().Erase()
	// panels["files"].Window().Erase()

	// panels["picture"].Window().NoutRefresh()
	// panels["files"].Window().NoutRefresh()
	// f := getFiles()
	// file_list := panels["files"].Window()
	// writeFiles(f)
}

func main() {
	stdscr, _ := gc.Init()
	defer gc.End()
	gc.CBreak(true)
	gc.Echo(true)

	panels_init := make(map[string]*gc.Panel)

	// drawEverything(stdscr)
	initWindows(panels_init)
	for {
		// stdscr.Erase()
		// drawEverything(panels)
		Update()

		switch stdscr.GetChar() {
		case 'q':
			return
		}
	}
}
