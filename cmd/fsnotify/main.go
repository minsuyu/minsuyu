package main

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
)

const InClose fsnotify.Op = 0x18

// const InCloseWrite fsnotify.Op = 0x8
const InCloseWrite fsnotify.Op = 128

func main() {
	watcher, err := fsnotify.NewWatcher()

	if err != nil {
		fmt.Println("create watcher ", err.Error())
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				//fmt.Printf("[EVENT] %s\n", event)
				if event.Op&InCloseWrite == InCloseWrite {
					fmt.Printf("[Close_Write] 0x8 감지: %s\n", event.Name)
				} else if event.Op&InClose == InClose {
					fmt.Printf("[Close] 0x18 감지: %s\n", event.Name)
				} else {
					fmt.Printf("[EVENT] %s\n", event)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				fmt.Println("watcher error ", err.Error())
			}
		}
	}()

	err = watcher.Add("/home/minsu/")
	if err != nil {
		fmt.Println("watch add err: ", err.Error())
	}
	fmt.Println("wating watch")
	<-make(chan bool)
}
