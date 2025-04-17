package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"github.com/fatih/color"
)

/*
 This script does webscrapping to download the classes from csv file
*/

const STATUS_OK string = "200 OK"

type dataset struct{
	class string 
	atom string 
	snippet string 
	line int 
	statusCode string
}

// readDataset reads the given CSV file and returns a slice of dataset structs.
// It parses the file, extracting each row and converting it into a dataset struct.
// The function assumes the CSV has a specific format that matches the fields of the dataset struct.
// 
// Parameters:
//   - fileName (string): The path to the CSV file to be read.
//	 - baseLink (string): The path to the resource in github.
// Returns:
//   - []dataset: A slice of dataset structs representing the data from the CSV file.
//
// Example usage:
//   datasets := readDataset("data.csv")
func readDataset(fileName, baseLink string) []dataset{
	file, err := os.Open(fileName)

	if err != nil{
		color.Red(fmt.Sprintf("[-] Fail to read csv file: %s", fileName))
		panic(err)
	}

	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ','
	reader.LazyQuotes = true

	records, err := reader.ReadAll()
	records = records[1:] // skip title

	if err != nil{
		color.Red(fmt.Sprintf("[-] Fail to read csv all file: %s", fileName))
		panic(err)
	}

	rows := []dataset{}
	for _, record := range records{

		line, _ := strconv.Atoi(record[3])
		row := dataset{
			class: formatLink(baseLink, strings.ReplaceAll(record[0], ".", "/") + ".java"),
			atom: record[1],
			snippet: record[2],
			line: line,
		}
		rows = append(rows, row)
	}

	return rows
}

func formatLink(baseLink, resource string) string{
	return baseLink + resource
}

func getClassFromGithub(link *dataset) string{
	response, err := http.Get(link.class)
	if err != nil{
		color.Red(fmt.Sprintf("Fail to download link: %s. Status: %s",link.class, response.Status))
		return ""
	}
	defer response.Body.Close()

	if response.Status == STATUS_OK{
		link.statusCode = STATUS_OK
		content, err := io.ReadAll(response.Body)
		if err != nil {
			color.Red(fmt.Sprintf("[-] Fail to read page content: %v", err))
			return ""
		}

		color.Blue("[+] Download Success: " + STATUS_OK)
		code := string(content)
		return code
	}else{
		link.statusCode = response.Status
	}
	return ""
}

func insertIntoDatabase(row dataset){
	classContent := getClassFromGithub(&row)
	log.Println(classContent)
}

func main(){
	basePath := "dataset/CSV/"

	baseLink := make(map[string]string)
	baseLink["aoc"] = "https://raw.githubusercontent.com/redisson/redisson/master/redisson/src/main/java/"
	baseLink["fastutil"] = "https://raw.githubusercontent.com/vigna/fastutil/master/src/"
	baseLink["jimfs"] = "https://raw.githubusercontent.com/google/jimfs/master/jimfs/src/main/java/"
	baseLink["moshi"] = "https://raw.githubusercontent.com/square/moshi/master/moshi/src/main/java/"
	baseLink["ucrop"] = "https://raw.githubusercontent.com/Yalantis/uCrop/develop/ucrop/src/main/java/"

	files, err := os.ReadDir(basePath)

	if err != nil{
		color.Red(fmt.Sprintf("[-] Fail to read dir: %s", basePath))
		panic(err)
	}

	for _, file := range files{
		if !file.IsDir(){
			name := file.Name()
			for _, row := range readDataset(
				path.Join(basePath, name),
				baseLink[strings.Split(name, "-")[1]]){
					insertIntoDatabase(row)
				}
				break
		}
		break
	}
	color.Green("[+] Completed!")
}