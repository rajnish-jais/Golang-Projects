//package main
//
//import (
//	"encoding/csv"
//	"fmt"
//	"log"
//	"os"
//	"time"
//)
//
//const (
//	incoming = 0
//	outgoing = 1
//)
//
//type Pair struct {
//	values [2]interface{}
//}
//
//func MakePair(k, v interface{}) Pair {
//	return Pair{values: [2]interface{}{k, v}}
//}
//
//func (p Pair) Get(i int) interface{} {
//	return p.values[i]
//}
//
//func main() {
//	file, err := os.Open("True-Caller Test/calling.csv")
//
//	// Checks for the error
//	if err != nil {
//		log.Fatal("Error while reading the file", err)
//	}
//
//	// Closes the file
//	defer file.Close()
//
//	reader := csv.NewReader(file)
//
//	records, err := reader.ReadAll()
//
//	recordMap := make(map[Pair]float64)
//	for _, v := range records {
//		startTime, _ := time.Parse(time.RFC3339, v[2])
//		endTime, _ := time.Parse(time.RFC3339, v[3])
//		duration := endTime.Sub(startTime)
//		recordMap[MakePair(v[0], outgoing)] = duration.Minutes()
//		recordMap[MakePair(v[1], incoming)] = duration.Minutes()
//	}
//
//	PrintCallDuration(recordMap, "Adam")
//
//}
//
//func PrintCallDuration(record map[Pair]float64, name string) {
//	fmt.Printf("Total incoming duration:%v minutes\n", record[MakePair(name, incoming)])
//	fmt.Printf("Total outgoing duration:%v minutes", record[MakePair(name, outgoing)])
//}

package main
