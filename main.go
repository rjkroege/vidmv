package main

import (
	"log"
	"path/filepath"
	"os"
	"strings"
	"io"
)


var totalbytes int64
var totalfiles int64


func init() {
	filelist = make([]string, 0, 100)
}

var filelist []string


func eachfile(path string, info os.FileInfo, err error) error {
//	log.Printf("visiting File: path='%s'\n", path)
	
	// test the file for having the right type
	// 1. a .mov
	// 2. in a Original Media directory

	if filepath.Ext(path) != ".mov" {
		return nil
	}

//	log.Println("path is a mov", path)

	// explode path up into parts

	if strings.Contains(path, "/Original Media/") {
			// log.Printf("Adding File: path='%s'\n", path)
			totalbytes += info.Size()
			totalfiles += 1
			filelist = append(filelist, path)
	}
	return nil
}

func main() {
//	log.Println("hello")

	// 1. Walk file
	var rootpath string

	if len(os.Args) > 1 {
		// Could conceivably use flags here.
		rootpath = os.Args[1]
	}  else {

	rp, err := os.Getwd()
	if err != nil {
		log.Fatal("Can't get current dir:", err)
	}
		rootpath = rp
	}
	
//	log.Println("rootpath", rootpath)

	totalbytes = 0
	totalfiles = 0

	if err :=	filepath.Walk(rootpath, eachfile); err != nil {
		log.Fatal("bah! can't walk", err)
	}
	
	log.Printf("finished walking %d files, %d bytes\n", totalfiles, totalbytes)
//	for _, v := range filelist {
//		log.Println("	", v)
//	}


	// 2. copy files.
	statchan := make(chan int, 2)
	statchan <- 1

	for _, v := range filelist {
		// target path
		tp := maketargetpath(v,rootpath)
		// log.Println("	", v, "->", tp)
		
		<-statchan
		go copyfile(v, tp, statchan)
	}

	<-statchan
}

func maketargetpath(op, rootpath string) string {
	np := strings.TrimPrefix(op, rootpath)
	np = strings.Map( func(r rune) rune {
		switch r {
		case ' ':
			return '-'
		case '/':
			return '_'
		case '#':
			return ','
		
		} 
		return r
	}, np)
	return np
}

// copyfile copies a file from one place to another
func copyfile(op, np string, statchan chan int) {
	log.Println("starting copyfile", op, "->", np)
	defer func() { statchan <- 1 }()
	rdf, err := os.Open(op)
	if err != err {
		log.Println("copyfile failed to open", op, "because", err)
		return
	}
	defer rdf.Close()

	wrf, err := os.Create(np)
	if err != err {
		log.Println("copyfile failed to create", np, "because", err)
		return
	}
	defer wrf.Close()
	
	if _, err := io.Copy(wrf, rdf); err != nil {
		log.Println("copyfile: failed to copy", op, "to", np, "because", err)
	}

}