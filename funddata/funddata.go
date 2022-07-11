package funddata

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

func GetFund(code string, page int) ([][]string, error) {
	url := "https://fundf10.eastmoney.com/F10DataApi.aspx?type=lsjz&code=%s&page=%d&per=40"
	response, err := http.Get(fmt.Sprintf(url, code, page))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	result := make([][]string, 0)
	pages := make([]string, 0)
	re := string(body)[len(string(body))-35:]
	pagesData := regexp.MustCompile(`[0-9]*`).FindAll([]byte(re), -1)
	for _, v := range pagesData {
		if string(v) != "" {
			pages = append(pages, string(v))
		}
	}
	result = append(result, pages)

	express := `<td[a-z\s=\']*?>([0-9\.%-]+)</td>`
	data := regexp.MustCompile(express).FindAllSubmatch(body, -1)
	//fmt.Printf("%q", data)
	for i := 0; i < len(data); i += 4 {
		//fmt.Printf("%q\n", data[i][1])
		//fmt.Printf("%q\n", data[i+1][1])
		//fmt.Printf("%q\n", data[i+2][1])
		//fmt.Printf("%q\n", data)
		//fmt.Println("..........................................")
		_, err := time.Parse("2006-01-02", string(data[i][1]))
		if err != nil {
			fmt.Println("日期错误：", string(data[i][1]))
			i--
		}
		if strings.HasSuffix(string(data[i+2][1]), "%") {
			result = append(result, []string{
				string(data[i][1]),
				string(data[i+1][1]),
				"0",
				string(data[i+2][1]),
			})
			i--
			continue
		}
		if len(data) <= i+3 || !strings.HasSuffix(string(data[i+3][1]), "%") {
			result = append(result, []string{
				string(data[i][1]),
				string(data[i+1][1]),
				"0",
				"0%",
			})
			i--
		} else {
			result = append(result, []string{
				string(data[i][1]),
				string(data[i+1][1]),
				string(data[i+2][1]),
				string(data[i+3][1]),
			})
		}
	}

	return result, err
}
