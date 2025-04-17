/*
 This script performs web scraping and database insertion for Java class files.

 It reads one or more CSV files from the local dataset directory. Each CSV contains metadata about
 Java classes such as class path, atom, code snippet, and line number. For each entry, the script:

 1. Constructs the full URL to the Java class file using a predefined GitHub base path.
 2. Downloads the class file content via HTTP.
 3. Inserts the downloaded content and metadata into a PostgreSQL database.

 Environment variables for database connection are loaded from a `.env` file using the `godotenv` package.
*/

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

const STATUS_OK string = "200 OK"

type dataset struct {
	class      string
	atom       string
	snippet    string
	line       string
	githubLink string
	statusCode string
}

type dBSettings struct {
	host     string
	port     string
	user     string
	password string
	database string
}

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
	link.githubLink = link.class
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

	sql := "INSERT INTO dataset (class, atom, snippet, line, status_code, github_link) VALUES ($1, $2, $3, $4, $5, $6)"

	_, err := pool.Exec(context.Background(), sql, classContent, row.atom, row.snippet, row.line, row.statusCode, row.githubLink)

	if err != nil {
		color.Red(fmt.Sprintf("[-] Error to insert %v into database: %v", row, err))
		return
	}

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
