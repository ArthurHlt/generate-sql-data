package main

import (
	"os"
	"fmt"
	"encoding/base64"
	"github.com/dustin/go-humanize"
)

var file string
var size uint64
var writtenData uint64 = uint64(0)
var currentChunk uint64 = uint64(0)

func main() {
	var err error
	args := os.Args
	if len(args) < 3 {
		fatal("Usage: " + args[0] + " [file size] [file name] (e.g.: generate-sql-data 1mb fakedata.sql)")
		os.Exit(1)
	}
	file = args[2]
	size, err = humanize.ParseBytes(args[1])
	fatalIf("invalid size", err)

	random, err := os.Open("/dev/urandom")
	fatalIf("error opening random file", err)
	defer random.Close()

	fileToWrite, err := os.Create(file)
	fatalIf("error creating file", err)
	defer fileToWrite.Close()

	generateFile(fileToWrite, random)

}
func generateFile(fileToWrite *os.File, random *os.File) {
	initializeFile(fileToWrite)

	for ; writtenData < size; {
		if (!shouldWriteInsert()) {
			writeData(fileToWrite, random)
		} else {
			writeInsertInto(fileToWrite, random)
		}
	}
	fileToWrite.Write([]byte("('" + generateData(random) + "');"))
}
func initializeFile(fileToWrite *os.File) {
	data := []byte(`DROP TABLE IF EXISTS fake_table;
CREATE TABLE fake_table
(
    value VARCHAR(8)
);
INSERT INTO fake_table VALUES`)
	writtenData += uint64(len(data))
	currentChunk += uint64(len(data))
	fileToWrite.Write(data)

}
func shouldWriteInsert() bool {
	if currentChunk > 900 * 1024 {
		currentChunk = 0
		return true
	}
	return false
}
func writeInsertInto(fileToWrite *os.File, random *os.File) {
	data := []byte("('" + generateData(random) + "'); INSERT INTO fake_table VALUES")
	writtenData += uint64(len(data))
	currentChunk += uint64(len(data))
	fileToWrite.Write(data)
}
func writeData(fileToWrite *os.File, random *os.File) {
	data := []byte("('" + generateData(random) + "'), ")
	writtenData += uint64(len(data))
	currentChunk += uint64(len(data))
	fileToWrite.Write(data)
}
func generateData(random *os.File) string {
	data := make([]byte, 5)
	_, err := random.Read(data)
	fatalIf("error get random value", err)
	return base64.StdEncoding.EncodeToString(data)
}
func fatalIf(doing string, err error) {
	if err != nil {
		fatal(doing + ": " + err.Error())
	}
}
func fatal(message string) {
	fmt.Fprintln(os.Stdout, message)
	os.Exit(1)
}