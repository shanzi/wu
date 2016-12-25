package runner

import (
	"github.com/fsnotify/fsnotify"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"time"
)

func watch(path string, abort <-chan struct{}) (<-chan string, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	out := make(chan string)
	go func() {
		defer close(out)
		defer watcher.Close()
		for {
			select {
			case <-abort:
				// Abort watching
				err := watcher.Close()
				if err != nil {
					log.Fatalln("Failed to stop watch")
				}
				return
			case fp := <-watcher.Events:
				if fp.Op == fsnotify.Create {
					info, err := os.Stat(fp.Name)
					if err == nil && info.IsDir() {
						// Add newly created sub directories to watch list
						log.Printf("Add newly diectory ( %s )\n", fp.Name)
						watcher.Add(fp.Name)
					}
				}

				if fp.Op&fsnotify.Write == fsnotify.Write || fp.Op == fsnotify.Remove || fp.Op == fsnotify.Rename {
					out <- fp.Name
				}

			case err := <-watcher.Errors:
				log.Println("Watch Error:", err)
			}
		}
	}()

	// Start watch
	{
		var paths []string
		currpath, _ := os.Getwd()

		readAppDirectories(currpath, &paths)

		log.Println("Start watching...")

		for _, dir := range paths {
			watcher.Add(dir)
			log.Printf("Directory( %s )\n", dir)
		}
	}

	return out, nil
}

func match(in <-chan string, patterns []string) <-chan string {
	out := make(chan string)

	go func() {
		defer close(out)
		for fp := range in {
			info, err := os.Stat(fp)

			if os.IsNotExist(err) {
				log.Printf("Dictory (%s) have been removed\n", fp)
				log.Println("here=======here")
				continue
			}

			if os.IsNotExist(err) || !info.IsDir() {
				_, fn := filepath.Split(fp)
				for _, p := range patterns {
					if ok, _ := filepath.Match(p, fn); ok {
						out <- fp
					}
				}
			}
		}
	}()

	return out
}

// gather delays further operations for a while and gather
// all changes happened in this period
func gather(first string, changes <-chan string, delay time.Duration) []string {
	files := make(map[string]bool)
	files[first] = true
	after := time.After(delay)
loop:
	for {
		select {
		case fp := <-changes:
			files[fp] = true
		case <-after:
			// After the delay, return collected filenames
			break loop
		}
	}

	ret := []string{}
	for k := range files {
		ret = append(ret, k)
	}

	sort.Strings(ret)
	return ret
}

func readAppDirectories(directory string, paths *[]string) {
	fileInfos, err := ioutil.ReadDir(directory)

	if err != nil {
		return
	}

	haveDir := false
	for _, fileinfo := range fileInfos {
		if fileinfo.IsDir() == true && fileinfo.Name() != "." && fileinfo.Name() != ".git" {
			readAppDirectories(directory+"/"+fileinfo.Name(), paths)
			continue
		}

		if haveDir {
			continue
		}

		if filepath.Ext(fileinfo.Name()) == ".go" {
			*paths = append(*paths, directory)
			haveDir = true
		}
	}
}
