package main

import (
	sqlite "demo/database"
	fund "demo/funddata"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

func getData(code string, page int) [][]string {
	for i := 0; i < 10; i++ {
		result, err := fund.GetFund(code, page)
		if err != nil {
			fmt.Println(err)
			time.Sleep(time.Second * time.Duration(rand.Intn(10)+10))
			continue
		} else {
			return result
		}

	}
	return nil
}

func insertData(data [][]string, dbLength int) []sqlite.QueryTable {
	var item []sqlite.QueryTable
	for i := 1; i < len(data); i++ {
		netValue, err := strconv.ParseFloat(data[i][1], 64)
		growthRate, err := strconv.ParseFloat(data[i][3][:len(data[i][3])-1], 64)
		if err != nil {
			fmt.Println(err)
		}
		item = append(item, sqlite.QueryTable{
			Id:         dbLength + i,
			EquityTime: data[i][0],
			NetValue:   netValue,
			GrowthRate: growthRate,
		})
	}

	return item
}

func work(code string) {
	tableName := "A" + code
	db := sqlite.NewConnPool()
	sqlite.CreateTable(db, tableName)
	ret := sqlite.QueryData(db, tableName, 1)
	var dbLength int
	if len(ret) == 0 {
		dbLength = 0
	} else {
		dbLength = ret[0].Id + 1
	}
	fmt.Println(ret, dbLength)

	//code := "513050"
	data := getData(code, 1)
	pages, err := strconv.Atoi(data[0][1])
	if err != nil {
		fmt.Println(err)
	}

	for i := 1; i < pages; i++ {
		data = getData(code, i)
		//fmt.Println(data)
		fmt.Println("get data finished.")
		item := insertData(data, dbLength)
		fmt.Println("type data finished.")
		//fmt.Println(item)
		sqlite.InsertData(db, tableName, item)
		fmt.Println("insert data finished.")
		fmt.Printf("fetch %d page finished all %d pages, it is sleep rand 10.\n", i, pages)
		time.Sleep(time.Duration(rand.Intn(10)+10) * time.Second)

		ret := sqlite.QueryData(db, tableName, 1)
		dbLength = ret[0].Id
	}
}

func deleteDatabase(tableName string) {
	db := sqlite.NewConnPool()
	ret := sqlite.QueryData(db, tableName, 1)
	var dbLength int
	if len(ret) == 0 {
		dbLength = 0
	} else {
		dbLength = ret[0].Id + 1
	}

	for i := 0; i < dbLength; i++ {
		sqlite.RemoveData(db, tableName, i)
	}
}

func main() {
	code := "502010"
	work(code)
	//deleteDatabase("A513050")
}
