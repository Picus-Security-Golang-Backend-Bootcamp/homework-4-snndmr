package infrastructure

import (
	"encoding/csv"
	"github.com/Picus-Security-Golang-Backend-Bootcamp/homework-4-snndmr/internal/domain/author"
	"github.com/Picus-Security-Golang-Backend-Bootcamp/homework-4-snndmr/internal/domain/book"
	"log"
	"os"
	"strconv"
	"sync"
)

func GetBooksFromCSV(path string) chan *book.Book {
	jobs := initializeJobs(path)
	results := make(chan *book.Book)
	waitGroup := new(sync.WaitGroup)

	for w := 0; w <= 2; w++ {
		waitGroup.Add(1)
		go convertToBookStruct(jobs, results, waitGroup)
	}

	go func() {
		waitGroup.Wait()
		close(results)
	}()
	return results
}

func GetAuthorsFromCSV(path string) chan *author.Author {
	jobs := initializeJobs(path)
	results := make(chan *author.Author)
	waitGroup := new(sync.WaitGroup)

	for w := 0; w <= 2; w++ {
		waitGroup.Add(1)
		go convertToAuthorStruct(jobs, results, waitGroup)
	}

	go func() {
		waitGroup.Wait()
		close(results)
	}()
	return results
}

func initializeJobs(path string) chan []string {
	jobs := make(chan []string, 5)

	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("%s", err)
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal("Pending I/O operations canceled.")
		}
	}(file)

	lines, err := csv.NewReader(file).ReadAll()
	if err != nil {
		log.Fatalf("%s", err)
	}

	for _, line := range lines[1:] {
		jobs <- line
	}
	close(jobs)
	return jobs
}

func convertToBookStruct(jobs <-chan []string, results chan<- *book.Book, workGroup *sync.WaitGroup) {
	defer workGroup.Done()
	for job := range jobs {
		pageCount, _ := strconv.Atoi(job[3])
		stockCount, _ := strconv.Atoi(job[4])
		price, _ := strconv.ParseFloat(job[5], 64)
		authorID, _ := strconv.Atoi(job[6])
		results <- book.NewBook(job[0], job[1], job[2], pageCount, stockCount, price, false, uint32(authorID))
	}
}

func convertToAuthorStruct(jobs <-chan []string, results chan<- *author.Author, workGroup *sync.WaitGroup) {
	defer workGroup.Done()
	for job := range jobs {
		results <- author.New(job[0])
	}
}
