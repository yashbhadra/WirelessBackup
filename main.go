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
	ipaddress = "192.168.1.114:2121"
)

func main() {
	fmt.Println("Hello world")
	con := Connection(ipaddress)
	files := make(chan *ftp.Entry)
	quit := make(chan int)
	go GetFiles(con, files, quit)
	var entry *ftp.Entry

	for {
		select {
		case entry = <-files:
			fmt.Println(entry.Name)
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
	err = con.Login("anonymous", "")
	if err != nil {
		log.Fatal("cannot login")
	}
	return con
}
func GetFiles(con *ftp.ServerConn, files chan *ftp.Entry, quit chan int) {

	//fmt.Println(file)
	//log.Fatal("BAs")
	//inodes := make(chan *ftp.Entry)
	var folder string = "/Downloads/"
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
func RecursiveFetch(con *ftp.ServerConn, folder string, files chan *ftp.Entry) {
	entries, err := con.List(folder)
	if err != nil {
		log.Fatal(err)
	}
	for _, en := range entries {

		if en.Type == 0 {
			//Code to include all the file into a map or channel or sometthing
			files <- en

		} else if en.Type == 1 {
			files <- en
			//result := filepath.Join(folder, en.Name)
			//RecursiveFetch(con, result, files)
			//continue

		} else {
			continue
		}
	}
	return
}
func Store(con ftp.ServerConn, en *ftp.Entry) {
	path := filepath.Join(curDir, en.Name)
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	result, err := con.Retr(en.Name)
	if err != nil {
		log.Fatal(err)
	}
	_, err = io.Copy(file, result)
	if err != nil {
		log.Fatal("Cannot Copyy", err)
	}
	result.Close()

}
