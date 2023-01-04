package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

type entry struct {
	word  string
	count int
	files []string
}

type file struct {
	filename string
	contents string
}

func main() {
	start := time.Now()
	dir := flag.String("dir", ".", "directory")
	filter := flag.String("filter", "*", "filename filter")
	flag.Parse()

	list := getWordCloud(*dir, *filter)

	for _, item := range list {
		fmt.Println(item.word, item.count)
	}

	fmt.Println("Time", time.Since(start))
}

func getWordCloud(dir string, filter string) []entry {
	letterreg := regexp.MustCompile("[a-z]+")
	notletterreg := regexp.MustCompile("[^a-zA-Z0-9_<>]+")
	m := make(map[string]entry)
	files := getFiles(dir, filter)
	for _, f := range files {
		for _, line := range strings.Split(f.contents, "\n") {
			s := notletterreg.ReplaceAllString(line, " ")
			for _, w := range strings.Fields(s) {
				if letterreg.MatchString(w) {
					if ent, ok := m[w]; ok {
						ent.count++
						ent.files = append(ent.files, f.filename)
						m[w] = ent
					} else {
						ent := entry{word: w, count: 1, files: []string{f.filename}}
						m[w] = ent
					}
				}
			}
		}
	}

	list := []entry{}
	for _, v := range m {
		list = append(list, v)
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].count < list[j].count
	})
	return list
}

func getFiles(dir string, filter string) []file {
	files := []file{}
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		_, fn := filepath.Split(path)
		if strings.Contains(fn, filter) || filter == "*" {
			f := file{filename: path}
			b, err := os.ReadFile(path)
			if err != nil {
				log.Println("can't read file", path, err)
				return err
			}

			f.contents = string(b)
			files = append(files, f)
		}
		return nil
	})
	if err != nil {
		log.Println(err)
	}
	return files
}
