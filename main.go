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
	byword := flag.String("word", "", "get word occurrences instead")
	flag.Parse()

	list := getWordCloud(*dir, *filter, *byword)
	sum := 0

	for _, item := range list {
		fmt.Println(item.word, item.count)
		sum += item.count
	}

	fmt.Println("total", sum)

	fmt.Println("finished", time.Since(start))
}

func getWordCloud(dir string, filter, word string) []entry {
	notletterreg := regexp.MustCompile("[^a-zA-Z0-9_<>]+")
	m := make(map[string]*entry)
	files := getFiles(dir, filter)
	for _, f := range files {
		s := notletterreg.ReplaceAllString(f.contents, " ")
		for _, w := range strings.Fields(s) {
			if word != "" && w != word {
				continue
			}
			key := w
			if word != "" {
				key = f.filename
			}
			if _, ok := m[key]; !ok {
				m[key] = &entry{word: key}
			}
			m[key].files = append(m[key].files, f.filename)
			m[key].count++
		}
	}

	list := []entry{}
	for _, v := range m {
		list = append(list, *v)
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
			f := file{filename: fn}
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
