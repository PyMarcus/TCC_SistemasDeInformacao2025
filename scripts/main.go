package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/fatih/color"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

/*
 This script does webscrapping to download the classes from csv file
*/

const STATUS_OK string = "200 OK"

type dataset struct {
	class      string
	atom       string
	snippet    string
	line       string
	statusCode string
}

// readDataset reads the given CSV file and returns a slice of dataset structs.
// It parses the file, extracting each row and converting it into a dataset struct.
// The function assumes the CSV has a specific format that matches the fields of the dataset struct.
//
// Parameters:
//   - fileName (string): The path to the CSV file to be read.
//   - baseLink (string): The path to the resource in github.
//
// Returns:
//   - []dataset: A slice of dataset structs representing the data from the CSV file.
//
// Example usage:
//
//	datasets := readDataset("data.csv")
func readDataset(fileName, baseLink string) []dataset {
	file, err := os.Open(fileName)

	if err != nil {
		color.Red(fmt.Sprintf("[-] Fail to read csv file: %s", fileName))
		panic(err)
	}

	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ','
	reader.LazyQuotes = true

	records, err := reader.ReadAll()
	records = records[1:] // skip title

	if err != nil {
		color.Red(fmt.Sprintf("[-] Fail to read csv all file: %s", fileName))
		panic(err)
	}

	rows := []dataset{}
	for _, record := range records {
		row := dataset{
			class:   formatLink(baseLink, strings.ReplaceAll(record[0], ".", "/")+".java"),
			atom:    record[1],
			snippet: record[2],
			line:    record[3],
		}
		rows = append(rows, row)
	}

	return rows
}

func formatLink(baseLink, resource string) string {
	return baseLink + resource
}

func getClassFromGithub(link *dataset) string {
	response, err := http.Get(link.class)
	if err != nil {
		color.Red(fmt.Sprintf("[-] Error: %v to %v", err, link))
		color.Red(fmt.Sprintf("[-]Fail to download link: %s. Status: %s", link.class, response.Status))
		return ""
	}
	defer response.Body.Close()

	if response.Status == STATUS_OK {
		link.statusCode = STATUS_OK
		content, err := io.ReadAll(response.Body)
		if err != nil {
			color.Red(fmt.Sprintf("[-] Fail to read page content: %v", err))
			return ""
		}

		color.Blue("[+] Download Success: " + STATUS_OK)
		code := string(content)
		return code
	} else {
		color.Yellow(fmt.Sprintf("[-]%s to %s", response.Status, link.class))
		link.statusCode = response.Status
	}
	return ""
}

func insertIntoDatabase(row dataset, pool *pgxpool.Pool) {
	classContent := getClassFromGithub(&row)

	sql := "INSERT INTO dataset (class, atom, snippet, line, status_code) VALUES ($1, $2, $3, $4, $5)"

	_, err := pool.Exec(context.Background(), sql, classContent, row.atom, row.snippet, row.line, row.statusCode)

	if err != nil {
		color.Red(fmt.Sprintf("[-] Error to insert %v into database: %v", row, err))
		return
	}

}

type dBSettings struct {
	host     string
	port     string
	user     string
	password string
	database string
}

func init() {
	err := godotenv.Load()

	if err != nil {
		panic("Env not set " + err.Error())
	}
}

func newDBSettings() *dBSettings {
	return &dBSettings{
		host:     os.Getenv("DATABASE_HOST"),
		database: os.Getenv("DATABASE_NAME"),
		password: os.Getenv("DATABASE_PASSWORD"),
		user:     os.Getenv("DATABASE_USER"),
		port:     os.Getenv("DATABASE_PORT"),
	}
}

func main() {

	dbSettings := newDBSettings()
	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbSettings.user, dbSettings.password, dbSettings.host, dbSettings.port, dbSettings.database)

	pool, err := pgxpool.New(context.Background(), dbUrl)

	if err != nil {
		color.Red("[!] Fail to connect with database: " + err.Error())
		return
	}

	defer pool.Close()

	basePath := "dataset/CSV/"

	baseLink := make(map[string]string)
	baseLink["aoc"] = "https://raw.githubusercontent.com/redisson/redisson/master/redisson/src/main/java/"
	baseLink["fastutil"] = "https://raw.githubusercontent.com/vigna/fastutil/master/src/"
	baseLink["jimfs"] = "https://raw.githubusercontent.com/google/jimfs/master/jimfs/src/main/java/"
	baseLink["moshi"] = "https://raw.githubusercontent.com/square/moshi/master/moshi/src/main/java/"
	baseLink["ucrop"] = "https://raw.githubusercontent.com/Yalantis/uCrop/develop/ucrop/src/main/java/"

	files, err := os.ReadDir(basePath)

	if err != nil {
		color.Red(fmt.Sprintf("[-] Fail to read dir: %s", basePath))
		panic(err)
	}

	for _, file := range files {
		if !file.IsDir() {
			name := file.Name()
			baseName := strings.Split(name, "-")[1]
			for _, row := range readDataset(
				path.Join(basePath, name),
				baseLink[strings.Replace(baseName, ".csv", "", 1)]) {
				insertIntoDatabase(row, pool)
			}
		}
	}
	color.Green("[+] Completed!")
}
