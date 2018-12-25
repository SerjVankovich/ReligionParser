package main

import (
	"encoding/csv"
	"log"
	"os"
	"parsing/utils"
	"strconv"
)

func main()  {

	channel := make(chan *utils.ReligionShape)

	var bigArr []*utils.ReligionShape

	utils.GetReligion(channel)

	for i := 0; i < 1412; i++{
		bigArr = append(bigArr, <- channel)
	}



	records := [][]string {
		{"title", "author", "description", "fulldescription", "month", "day", "year", "shows", "reposts", "rating"},
	}

	for _, religion := range bigArr {

		records = append(records, []string{religion.Title, religion.Author, religion.Description,
		religion.FullDescription, strconv.Itoa(religion.Month), strconv.Itoa(religion.Day), strconv.Itoa(religion.Year),
		strconv.Itoa(religion.Shows), strconv.Itoa(religion.Repost), strconv.Itoa(religion.Rating)})
	}

	file, err := os.Create("data.csv")
	if err != nil {
		log.Fatal(err)
	}

	writer := csv.NewWriter(file)

	writer.WriteAll(records)
}
