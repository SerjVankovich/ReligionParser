package utils

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html/charset"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func GetPage(url string) *goquery.Document {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	res.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.80 Safari/537.36")
	res.Header.Set("Accept-Language", "ru-RU,ru;q=0.9,en;q=0.8")
	defer res.Body.Close()

	utf8Str, err := charset.NewReader(res.Body, res.Header.Get("Content-Type"))

	if err != nil {
		log.Fatal(err)
	}

	doc, err := goquery.NewDocumentFromReader(utf8Str)

	if err != nil {
		log.Fatal(err)
	}

	return doc
}

func GetReligion(channel chan *ReligionShape) {
	for i := 1; i < 33; i++  {
		fmt.Println("Page:", i)
		go getOnePage("http://xn--b1amnebsh.ru-an.info/%D1%80%D0%B5%D0%BB%D0%B8%D0%B3%D0%B8%D1%8F/" + strconv.Itoa(i), channel)
	}
}

func getOnePage(url string, ch chan *ReligionShape) {
	document := GetPage(url)
	document.Find(".news-view").Each(func(i int, selection *goquery.Selection) {
		channel := make(chan string)
		news := &ReligionShape{}

		//fmt.Println("Title:", selection.Find("h2.h2").Text())
		//fmt.Println("Description:", selection.Find(".text.news-view-text").Text())
		link, _ := selection.Find("a").Attr("href")
		go getFullDescription(link, channel)


		dirtyString := strings.Split(selection.Find(".sb2").Text(), ",")
		date, author := toDateAndAuthor(dirtyString)

		//fmt.Println("Date:", date)
		//fmt.Println("Author:", author)
		//fmt.Println("Views:", convertToNum(strings.Trim(selection.Find("span.prosm-views").Text(), "\t\n ")))
		//fmt.Println("Shares: ", convertToNum(strings.Trim(selection.Find("span.prosm-shares").Text(), "\t\n ")))

		var stars int
		selection.Find("img").Each(func(i int, selection *goquery.Selection) {
			src, _ := selection.Attr("src")
			if src == "http://ru-an.info/Pictures/Icons/star_colored.png" {
				stars++
			}
		})
		//fmt.Println("Stars:", stars)

		description := selection.Find(".text.news-view-text").Text()

		if description == "" {
			description = selection.Find(".text").Text()
		}


		news.Title = selection.Find("h2.h2").Text()
		news.Description = description
		news.FullDescription = <- channel
		news.Author = author
		news.Day = date[0]
		news.Month = date[1]
		news.Year = date[2]
		news.Shows = convertToNum(strings.Trim(selection.Find("span.prosm-views").Text(), "\t\n "))
		news.Repost = convertToNum(strings.Trim(selection.Find("span.prosm-shares").Text(), "\t\n "))
		news.Rating= stars
		fmt.Println(news)

		ch <- news


		fmt.Println()
	})

}

func getFullDescription(url string, channel chan string)  {
	descriptionPage := GetPage(url)

	fullDescription := descriptionPage.Find(".text21.news-content").Text()

	channel <- fullDescription

}

func toDateAndAuthor(str []string) ([]int, string) {
	var dateArr []int
	var author string
	if len(str) != 2 {
		dateStrArr := strings.Split(str[0], " ")
		day, err := strconv.Atoi(dateStrArr[0])
		if err != nil {
			log.Fatal(err)
		}

		year, err := strconv.Atoi(dateStrArr[2])
		if err != nil {
			log.Fatal(err)
		}

		dateArr = []int{day, convertMonth(dateStrArr[1]), year}
	} else {
		dateStrArr := strings.Split(str[1], " ")
		day, err := strconv.Atoi(dateStrArr[1])

		if err != nil {
			log.Fatal(err)
		}

		year, err := strconv.Atoi(dateStrArr[3])
		if err != nil {
			log.Fatal(err)
		}

		dateArr = []int{day, convertMonth(dateStrArr[2]), year}
		author = str[0]
	}

	return dateArr, author
}

func convertMonth(month string) int {
	switch month {
	case "января":
		return 1

	case "февраля":
		return 2

	case "марта":
		return 3

	case "апреля":
		return 4

	case "мая":
		return 5

	case "июня":
		return 6

	case "июля":
		return 7

	case "августа":
		return 8

	case "сентября":
		return 9

	case "октября":
		return 10

	case "ноября":
		return 11

	case "декабря":
		return 12

	default:
		return 0
	}
}

func convertToNum(str string) int  {
	if str == "" {
		return 0
	}

	num, _ := strconv.Atoi(strings.Join(strings.Split(str, " "), ""))

	return num
}

func LineArr(matrix [][]*ReligionShape) []*ReligionShape  {
	var retArr []*ReligionShape
	for _, value := range matrix {
		for _, rel := range value {
			retArr = append(retArr, rel)
		}
	}

	return retArr
}
