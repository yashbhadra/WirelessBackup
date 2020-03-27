package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/jlaffaye/ftp"
)

const (
	curDir    = "C:/Users/yash.bhanushali/BackupTrial-1"
	ipaddress = "192.168.1.101:2221"
	username  = "anonymous"
	password  = "android"
)

type File struct {
	Name string
	Type ftp.EntryType
	Path string
}

func main() {
	fmt.Println("Hello world")
	con := Connection(ipaddress)
	files := make(chan *File)
	quit := make(chan int)
	go GetFiles(con, files, quit)
	var entry *File

	for {
		select {
		case entry = <-files:
			//fmt.Println(entry.Name)
			Store(entry, con)
		case <-quit:
			return
		}
	}

}
func Connection(ipaddress string) *ftp.ServerConn {
	con, err := ftp.Dial(ipaddress, ftp.DialWithTimeout(time.Second))
	if err != nil {
		log.Fatal("Couldnt connect to the server.")
	}
	err = con.Login("anonymous", "android")
	if err != nil {
		log.Fatal("cannot login")
	}
	return con
}
func GetFiles(con *ftp.ServerConn, files chan *File, quit chan int) {

	//fmt.Println(file)
	//log.Fatal("BAs")
	//inodes := make(chan *ftp.Entry)
	var folder string = ""
	RecursiveFetch(con, folder, files)
	quit <- 0

}

/*var en *ftp.Entry
for i := 0; i <= len(inodes); i++ {
	en = <-inodes

}*/

/*resp, err := ioutil.ReadAll(file)
if err != nil {
	fmt.Println(err)
}
println(string(resp))
/*buff := bytes.NewBufferString("Hello World")
err = con.Stor("c.txt", buff)
if err != nil {
	fmt.Println(err)
}*/
func RecursiveFetch(con *ftp.ServerConn, folder string, files chan *File) {
	entries, err := con.List(folder)
	if err != nil {
		log.Fatal(err)
	}

	for _, en := range entries {

		if en.Type == 0 {
			//Code to include all the file into a map or channel or sometthing
			temp := &File{
				Name: en.Name,
				Type: en.Type,
				Path: folder,
			}
			files <- temp

		} else if en.Type == 1 {
			result := filepath.Join(folder, en.Name)
			temp := &File{
				Name: en.Name,
				Type: en.Type,
				Path: result,
			}
			files <- temp
			//fmt.Println(folder)
			RecursiveFetch(con, result, files)

		} else {
			continue
		}
	}
	return
}
func Store(entry *File, con *ftp.ServerConn) {
	if entry.Type == 0 {
		path := filepath.Join(curDir, entry.Path, entry.Name)
		file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(entry.Path + entry.Name)
		result, err := con.Retr(filepath.Join(entry.Path, entry.Name))
		if err != nil {
			log.Fatal("I m here", err)
		}
		_, err = io.Copy(file, result)
		if err != nil {
			log.Fatal("Cannot Copy", err)
		}
		result.Close()

	} else if entry.Type == 1 {

		//fmt.Println(curDir + "/" + entry.Name)

		path := filepath.Join(curDir, entry.Path)
		fmt.Println(path)
		os.Mkdir(path, 0777)

	} else {
		return
	}

}
