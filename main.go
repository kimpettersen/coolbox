// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !plan9

package main

import (
	"log"
	"os"

	"path/filepath"

	"fmt"

	"github.com/fsnotify/fsnotify"
)

func addWatcherToFolder(watcher *fsnotify.Watcher, folder string) {
	fileInfo, err := os.Stat(folder)
	if err != nil {
		log.Fatalf("Could not find the information of the file %s", folder)
	}
	if fileInfo.IsDir() {
		watcher.Add(folder)
		f, err := os.Open(folder)
		if err != nil {
			fmt.Println(err)
		}
		defer f.Close()
		names, err := f.Readdirnames(-1)
		if err != nil {
			fmt.Println(err)
		}
		for _, name := range names {
			addWatcherToFolder(watcher, filepath.Join(folder, name))
		}
	}
}

func main() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)
				} else if event.Op&fsnotify.Create == fsnotify.Create {
					log.Println("created file:", event.Name)
					addWatcherToFolder(watcher, event.Name)
				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add("/tmp/foo")
	if err != nil {
		log.Fatal(err)
	}
	<-done
}
