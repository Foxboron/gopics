package main

import (
	"bufio"
	"fmt"
	gc "github.com/gbin/goncurses"
	"os/exec"
	"strings"
	"time"
)

type Rsyncs interface {
	addFiles()
	Upload()
}

type Rsync struct {
	rsync      string "/usr/bin/rsync"
	flags      string "-HPavzc"
	path       string "./www/pics/albums"
	hook       string "./bin/rsync-pics"
	includes   []string
	excludes   []string
	ssh        string "trinity"
	folder     string ""
	local_path string
}

func createRsync() Rsync {
	return Rsync{
		rsync:      "/usr/bin/rsync",
		flags:      "-HPavzc",
		path:       "./www/pics/albums",
		hook:       "./bin/rsync-pics",
		includes:   []string{"*/"},
		excludes:   []string{"*"},
		ssh:        "trinity",
		local_path: getDevicePath()}
}

func createFlags(name string, flags []string) string {
	ret := ""
	for _, i := range flags {
		flag_ := strings.Split(i, "/")
		flag := flag_[len(flag_)-1]
		ret += fmt.Sprintf("--%s=\"%s\" ", name, flag)
	}
	return ret
}

func (r Rsync) addFiles(files []string) Rsync {
	r.includes = append(r.includes, files...)
	return r
}

func (r Rsync) Upload(folder string) string {
	shell := ""
	includes := createFlags("include", r.includes)
	excludes := createFlags("exclude", r.excludes)
	shell = r.hook
	shell = fmt.Sprintf("--rsync-path=\"mkdir -p %s/%s && %s\"", r.path, folder, r.hook)
	cmd := "%s %s %s %s %s %s %s:%s/%s"
	ret := fmt.Sprintf(cmd, r.rsync, r.flags, shell, includes, excludes, r.local_path, r.ssh, r.path, folder)
	return ret
}

func (w Interface) cmdUpload(msg string, pan *gc.Panel) {
	cmd := exec.Command("bash", "-c", msg)

	stdout, _ := cmd.StdoutPipe()
	cmd.Start()
	rd := bufio.NewReader(stdout)
	size := 9
	i := 0
	for {
		byt, _, _ := rd.ReadLine()
		str := string(byt)

		if str == "" {
			break
		}

		pan.Window().Scroll(1)

		if i > 13 {
			break
		}
		if len(str) >= 55 {
			str = str[:55]
		}
		pan.Window().Box(0, 0)
		pan.Window().MovePrintln(size, 1, string(str))
		pan.Window().Refresh()
		w.Refresh()
		i++
		if str == "" {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
	pan.Delete()
}
