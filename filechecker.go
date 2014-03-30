package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type FileInfo struct {
	size  int64
	hash  string
	paths []string
}

func Hash(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}

	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	h := md5.New()
	buf := make([]byte, 1024)
	for {
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if n == 0 {
			break
		}
		h.Write(buf[:n])
	}
	sum := h.Sum(nil)
	hashstr := fmt.Sprintf("%x", sum)
	//fmt.Printf("%s\n", hashstr)
	return hashstr, nil
}

func visit(path string, f os.FileInfo, result map[string]*FileInfo, sizeinfo map[int64]*FileInfo, err error) error {
	//fmt.Printf("Visited: %s : ", path)
	if f.IsDir() {
		//fmt.Print("skipping dir\n")
		return nil
	}

	// check the size of the file. If it's a new size then skip it (no dupes yet)
	size := f.Size()
	fi, ok := sizeinfo[size]
	if !ok {
		sizeinfo[f.Size()] = &FileInfo{f.Size(), "", []string{path}}
		return nil
	}

	//fmt.Print(".")
	// duplicate size
	// if the FileInfo from sizeinfo doesn't have a hash, then
	// that means it's not in the result set
	if fi.hash == "" {
		fi.hash, err = Hash(fi.paths[0])
		if err != nil {
			return err
		}

		result[fi.hash] = fi
	}

	// hash the file at this path and add it to the set.
	hashstr, err := Hash(path)
	if err != nil {
		return err
	}

	if fi, ok := result[hashstr]; !ok {
		result[hashstr] = &FileInfo{f.Size(), hashstr, []string{path}}
	} else {
		fmt.Print("O")
		fi.paths = append(fi.paths, path)
	}

	return nil
}

func visitor(sizemap map[int64]*FileInfo, result map[string]*FileInfo) func(string, os.FileInfo, error) error {
	return func(path string, f os.FileInfo, err error) error {
		return visit(path, f, result, sizemap, err)
	}
}

func walk(root string, result map[string]*FileInfo) {

	sizemap := make(map[int64]*FileInfo)

	err := filepath.Walk(root, visitor(sizemap, result))
	if err != nil {
		fmt.Printf("Trouble walking the file system: %s\n", err)
	}
}

func main() {
	flag.Parse()

	result := make(map[string]*FileInfo)

	for _, root := range flag.Args() {
		walk(root, result)
	}
	fmt.Println("")
	for hash, fi := range result {
		if len(fi.paths) == 1 {
			continue
		}

		fmt.Printf("%s\n", hash)
		for i := range fi.paths {
			fmt.Printf("\t%s\n", fi.paths[i])
		}
	}
}
