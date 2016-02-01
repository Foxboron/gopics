package main

import (
	gc "github.com/gbin/goncurses"
)

func main() {

	config := getConfig()
	wins := createInterface()

	for {

		ch := wins.GetChar()

		switch ch {
		case 'q':
			if wins.tmp != nil {
				print("yay")
				wins.tmp.Delete()
			} else {
				gc.Raw(false)
				gc.Echo(true)
				gc.Cursor(1)
				return
			}

		case gc.KEY_DOWN, 'j':
			wins.Action(gc.REQ_DOWN)
			wins.Refresh()

		case gc.KEY_UP, 'k':
			wins.Action(gc.REQ_UP)
			wins.Refresh()
		case ' ':
			wins.Action(gc.REQ_TOGGLE)
			wins.Action(gc.REQ_DOWN)
			wins.Refresh()

		// Upload
		case 'u':
			pan := wins.createTempWindow("", false, 20, 80)
			str := wins.getInput("Album: ")
			wins.Refresh()
			if len(str) > 0 {

				var list []string
				for _, item := range wins.files.Items() {
					if item.Value() {
						list = append(list, item.Description())
					}
				}
				r := createRsync(config)
				r = r.addFiles(list)
				str = r.Upload(str)
				wins.cmdUpload(str, pan)
			}
		// Select all
		case 'V':
			for _, item := range wins.files.Items() {
				item.SetValue(true)
			}
			wins.Refresh()
		case gc.KEY_ENTER, gc.KEY_RETURN:

		default:
			wins.files.Driver(gc.DriverActions[ch])
		}
	}
}
