package main

import (
	"evaluation/internal/infsysresults"
	"evaluation/internal/infsystem"
	"evaluation/internal/metrics"
	"evaluation/internal/qrels"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	qr := flag.String("qrels", "qrels.txt", "Ruta del fichero con los juicios de relevancia")
	results := flag.String("results", "results.txt", "Ruta del fichero con los resultados del sistema de infomación")
	output := flag.String("output", "output.txt", "Ruta del fichero donde se guardará la evaluación del sistema de información")
	flag.Parse()
	rels, err := qrels.ParseQrelsFile(openFile(*qr))
	if err != nil {
		log.Fatalf("error parsing qrels file: %v", err)
	}
	relevant := qrels.CreateMap(rels)

	res, err := infsysresults.ParseResults(openFile(*results))
	if err != nil {
		log.Fatalf("error parsing results file: %v", err)
	}
	resMap := infsysresults.CreateMap(res)
	collection := infsystem.InfSystem{
		Relevances: relevant,
		ISResults:  resMap,
	}
	summary := metrics.CreateSummary(collection)
	fmt.Println(summary.String())
	err = writeToFile(*output, []byte(summary.String()))
	if err != nil {
		log.Fatalf("error writing file: %v", err)
	}
}

func openFile(filename string) *os.File {
	f, err := os.Open(filename)
	if err != nil {
		log.Printf("error opening %s", filename)
		os.Exit(1)
	}
	return f
}

func writeToFile(filename string, data []byte) error {
	return ioutil.WriteFile(filename, data, 0644)
}
