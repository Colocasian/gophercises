package main

import (
	"encoding/csv"
	"io"
	"log"

	"github.com/Colocasian/gophercises/quiz/pkg/quiz"
)

func readCSV(r io.Reader) *quiz.Quizzer {
	reader := csv.NewReader(r)
	reader.FieldsPerRecord = 2

	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("invalid CSV file: %v", err)
	}

	booklet := make([]quiz.QA, len(records))
	for i := range records {
		booklet[i] = quiz.QA{Q: records[i][0], A: records[i][1]}
	}
	return quiz.NewQuizzer(booklet)
}
