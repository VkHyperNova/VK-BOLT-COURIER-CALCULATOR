package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// ALL TIME STATS
var TOTAL_INCOME float64
var TOTAL_REVENUE float64
var TOTAL_PETROL float64
var TOTAL_HOURS float64

func main() {

	CHECK_FILES()

	// Open data
	WEEK_DATABASE = OPEN_JSON()

	// Start commandline
	CL()

	// delete old week and save info to month.json
	// update stats after entry
	// print current info
	// 20% income
	// cash

}

func Print_week_date() {
	week_start := WEEK_DATABASE[0].DATE
	current_date := len(WEEK_DATABASE) - 1
	week_end := WEEK_DATABASE[current_date].DATE
	fmt.Println(week_start + " - " + week_end)
}

/******************************************************************************/

var INCOME_WEEK float64
var PETROL_WEEK float64
var WEEK_DATABASE []day_transactions
var HOURS_WEEK float64
var TAX float64
var CASH_WEEK float64

func CALCULATE_PETROL_WEEK() {
	for i := 0; i < len(WEEK_DATABASE); i++ {
		PETROL_WEEK += WEEK_DATABASE[i].PETROL
	}
	fmt.Println("-> PETROL: ", PETROL_WEEK * -1, "EUR")
}

func CALCULATE_CASH_WEEK() {
	for i := 0; i < len(WEEK_DATABASE); i++ {
		CASH_WEEK += WEEK_DATABASE[i].CASH
	}
	fmt.Println("--> CASH: +", CASH_WEEK, "EUR")
}

func CALCULATE_INCOME_WEEK() {
	for i := 0; i < len(WEEK_DATABASE); i++ {
		INCOME_WEEK += WEEK_DATABASE[i].INCOME
	}
	fmt.Println("--> INCOME: ", INCOME_WEEK, "EUR")
}

func CALCULATE_REVENUE_AFTER_TAX() {
	INCOME_WEEK = INCOME_WEEK - CASH_WEEK
	TAX = INCOME_WEEK * 0.2
	fmt.Println("-> TAXES TO BE PAYED: ", TAX * -1, "EUR")

	INCOME_WEEK = INCOME_WEEK - TAX
	INCOME_WEEK = INCOME_WEEK + CASH_WEEK
	INCOME_WEEK = INCOME_WEEK - PETROL_WEEK
	fmt.Println("<-- REVENUE AFTER TAX: ", INCOME_WEEK, "EUR -->")
}

func PRINT_STATISTICS() {
	fmt.Println("\n<-------------- BOLT FINANCE CALCULATOR -------------->\n")
	CALCULATE_PETROL_WEEK()
	CALCULATE_INCOME_WEEK()
	CALCULATE_CASH_WEEK()
	CALCULATE_REVENUE_AFTER_TAX()
}

/******************************************************************************/

func CL() {
	// Print
	PRINT_STATISTICS()

	fmt.Println("\n\n<< COMMANDS: add | stats | months | q >>")

	reader := bufio.NewReader(os.Stdin)

	for {

		command := Convert_CRLF_To_LF(reader)

		switch command {
		case "add":
			Add()
		case "months":
			Print_Months()
		case "q":
			Quit("clear")
		}
	}
}

func Add() {
	// Get income and petrol info
	fmt.Print("income today: ")
	fmt.Scanln(&INCOME_TODAY)

	fmt.Print("petrol: ")
	fmt.Scanln(&PETROL_TODAY)

	fmt.Print("hours: ")
	fmt.Scanln(&HOURS_TODAY)

	fmt.Print("cash: ")
	fmt.Scanln(&CASH_TODAY)

	// save to file
	DAY_DATA := CONSTRUCT_DAY_TRANSACTIONS()
	WEEK_DATABASE = append(WEEK_DATABASE, DAY_DATA)
	DAY_DATA_AS_BYTE := Convert_To_Byte(WEEK_DATABASE)
	WRITE_FILE("week.json", DAY_DATA_AS_BYTE)

	//

	fmt.Println("Updated!")
	CL()
}

func Print_Months() {
	// Print months file
}

func CHECK_FILES() {

	// Create week
	if _, err := os.Stat("./week.json"); errors.Is(err, os.ErrNotExist) {
		MAKE_EMPTY_JSON_FILE("week.json")
	}

	// Create month
	if _, err := os.Stat("./months.json"); errors.Is(err, os.ErrNotExist) {
		MAKE_EMPTY_JSON_FILE("months.json")
	}

}

func MAKE_EMPTY_JSON_FILE(name string) {
	WRITE_FILE(name, []byte("[]"))
}

func Convert_CRLF_To_LF(reader *bufio.Reader) string {

	// Read the answer
	input, _ := reader.ReadString('\n')

	// Convert CRLF to LF
	input = strings.Replace(input, "\r\n", "", -1) /* "\r\n" was before.  */

	return input
}

func OPEN_JSON() []day_transactions {

	file := READ_FILE("./week.json")
	data := CONVERT_TO_DATA(file)

	return data
}

func CONVERT_TO_DATA(file []byte) []day_transactions {

	transactions := []day_transactions{}

	err := json.Unmarshal(file, &transactions)
	Error(err, "ConvertToTransactions")

	return transactions
}

/******************************************************************************/

type day_transactions struct {
	ID      int     `json:"id"`
	DATE    string  `json:"date"`
	INCOME  float64 `json:"income"`
	PETROL  float64 `json:"petrol"`
	REVENUE float64 `json:"revenue"`
	HOURS   float64 `json:"hours"`
	TAX     float64 `json:"tax"`
	CASH    float64 `json:"cash"`
}

var INCOME_TODAY float64
var PETROL_TODAY float64
var REVENUE_TODAY float64
var HOURS_TODAY float64
var TAX_TODAY float64
var CASH_TODAY float64

func CONSTRUCT_DAY_TRANSACTIONS() day_transactions {

	var Day = day_transactions{
		ID:      Get_Unique_website_ID(WEEK_DATABASE),
		DATE:    Get_Current_Time("15:04 (02.01.2006)"),
		INCOME:  INCOME_TODAY,
		PETROL:  PETROL_TODAY,
		REVENUE: REVENUE_TODAY,
		HOURS:   HOURS_TODAY,
		TAX:     TAX_TODAY,
		CASH:    CASH_TODAY,
	}

	return Day
}

/******************************************************************************/

type months struct {
	ID      int     `json:"id"`
	MONTH   float64 `json:"month"`
	REVENUE float64 `json:"revenue"`
	PETROL  float64 `json:"petrol"`
	INCOME  float64 `json:"income"`
	HOURS   float64 `json:"hours"`
}

/******************************************************************************/

func READ_FILE(filename string) []byte {
	file, err := ioutil.ReadFile(filename)
	Error(err, "ReadFile")
	return file
}

func WRITE_FILE(filename string, dataBytes []byte) {

	var err = ioutil.WriteFile(filename, dataBytes, 0644)
	Error(err, "WriteToFile")
}

func Error(err error, location string) {
	if err != nil {
		fmt.Println(" << Function name: ", location+" >>")
		fmt.Println(err.Error())

	}
}

func Get_Unique_website_ID(data []day_transactions) int {

	if len(data) == 0 {
		return 1
	}

	return data[len(data)-1].ID + 1
}

func Convert_String_To_Int(a string) int {
	b, _ := strconv.Atoi(a)
	return b
}

func Get_Current_Time(format string) string {
	TimeNow := time.Now()
	FormattedTimeNow := TimeNow.Format(format)
	return FormattedTimeNow
}

func Convert_To_Byte(data interface{}) []byte {
	dataBytes, err := json.MarshalIndent(data, "", "  ")
	Error(err, "Convert_To_Byte")

	return dataBytes
}

func Quit(clear string) {

	if clear == "clear" {
		Clear_Screen()
	}

	os.Exit(0)
}

func Clear_Screen() {

	if runtime.GOOS == "linux" {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}
