package main

import (
	"fmt"
	// "log"
	"os/exec"
	fs "path/filepath"
	"strings"
)

type Files interface {
	getName()
	getPath()
}

type File struct {
	name string
	path string
}

func (f File) getPath() string {
	return f.path

}

func (f File) getName() string {
	return f.name
}

func createFile(path string) File {
	file_split := strings.Split(path, "/")
	name := file_split[len(file_split)-1]
	return File{name: name, path: path}
}

func createFiles(paths []string) []File {

	var files []File
	for _, i := range paths {
		files = append(files, createFile(i))
	}
	return files
}

func getFiles() []File {
	glob_path := get_device_path()
	files := glob_files(glob_path)
	type_files := createFiles(files)
	return type_files
}

func get_device_path() string {
	cmd := "lsblk -l | grep mmcblk0p1 | awk '{ print $7 }'"
	// fmt.Println(cmd)
	// parts := strings.Fields(cmd)
	out, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		fmt.Println("error occured")
		fmt.Printf("%s", err)
	}
	s := strings.TrimSpace(string(out))
	ret := []string{s, "/DCIM/1*CANON/*JPG"}
	return strings.Join(ret, "")
}

func glob_files(path string) []string {
	ret, err := fs.Glob(path)
	if err != nil {
		println(err)
	}
	return ret
}
