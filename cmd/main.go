package main

/*
	TCC 2025 - Dataset Enrichment Pipeline
	Code written by Marcus  
	GitHub: https://github.com/PyMarcus

	This application is part of a research project aimed at comparing two AI agents 
	to identify "atoms of confusion" — minimal code elements that cause misunderstandings — 
	in the context of Software Engineering.

	The goal is not only to detect these confusing elements but also to evaluate and compare 
	the agents' abilities to recognize such patterns in source code questions and datasets.

	Main functionalities:
	- Load configuration and initialize database connection
	- Retrieve datasets and related questions
	- Dispatch concurrent tasks using a worker pool
	- Query AI agents and handle their responses
	- Log errors and persist structured responses ("Atoms") into the database
*/

func main() {
	execute()
}
