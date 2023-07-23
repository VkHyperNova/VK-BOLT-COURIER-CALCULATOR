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

var DB_URL = "./week.json"

func main() {

	CHECK_JSON()

	// Open data
	WEEK_DATABASE = OPEN_JSON()

	// Start commandline
	CL()

}

/******************************************************************************/
/* COMMANDS */

func CL() {

	PRINT_STATISTICS()

	fmt.Println("\n\n<< COMMANDS: add | q >>")
	fmt.Print("=> ")

	reader := bufio.NewReader(os.Stdin)

	for {

		command := Convert_CRLF_To_LF(reader)

		switch command {
		case "add":
			Add()
		case "q":
			Quit("clear")
		default:
			Clear_Screen()
			CL()
		}
	}
}

func Add() {

	INCOME_TODAY = Question("Income: ")
	fmt.Println("Set to: ", INCOME_TODAY)

	PETROL_TODAY = Question("Petrol: ")
	fmt.Println("Set to: ", PETROL_TODAY)

	HOURS_TODAY = Question("Hours: ")
	fmt.Println("Set to: ", HOURS_TODAY)

	CASH_TODAY = Question("Cash: ")
	fmt.Println("Set to: ", CASH_TODAY)

	Save()
	CL()
}

func Save() {
	DAY_DATA := CONSTRUCT_DAY_TRANSACTIONS()
	WEEK_DATABASE = append(WEEK_DATABASE, DAY_DATA)
	DAY_DATA_AS_BYTE := Convert_To_Byte(WEEK_DATABASE)
	WRITE_FILE(DAY_DATA_AS_BYTE)
	fmt.Println("<< Updated! >>")
}

func Question(question string) float64 {
start:
	var answer string
	fmt.Print("\n", question)
	fmt.Scanln(&answer)

	if answer == "" {
		answer = "0"
	}

	floatValue, err := strconv.ParseFloat(answer, 64)
	if err != nil {
		fmt.Println("<< Error:", err)
		goto start
	}

	return floatValue
}

/******************************************************************************/
/* CALCULATIONS */
var currentTime = time.Now()
var currentMonth = currentTime.Month().String()

var WEEK_DATABASE []week_day

func CALCULATE_HOURS() float64 {
	var currentWeekHours float64 = 0
	for i := 0; i < len(WEEK_DATABASE); i++ {
		if WEEK_DATABASE[i].MONTH == currentMonth {
			currentWeekHours += WEEK_DATABASE[i].HOURS
		}
	}
	return currentWeekHours
}

func CALCULATE_INCOME() float64 {
	var INCOME float64 = 0
	for i := 0; i < len(WEEK_DATABASE); i++ {
		if WEEK_DATABASE[i].MONTH == currentMonth {
			INCOME += WEEK_DATABASE[i].INCOME
		}
	}
	return INCOME
}

func CALCULATE_PETROL() float64 {
	var PETROL float64 = 0
	for i := 0; i < len(WEEK_DATABASE); i++ {
		if WEEK_DATABASE[i].MONTH == currentMonth {
			PETROL += WEEK_DATABASE[i].PETROL
		}
	}
	return PETROL
}

func CALCULATE_CASH() float64 {
	var CASH float64 = 0
	for i := 0; i < len(WEEK_DATABASE); i++ {
		if WEEK_DATABASE[i].MONTH == currentMonth {
			CASH += WEEK_DATABASE[i].CASH
		}
	}
	return CASH
}

func CALCULATE_REVENUE() float64 {

	INCOME := CALCULATE_INCOME()
	fmt.Println("-----------------", INCOME)
	CASH := CALCULATE_CASH()
	fmt.Println("-----------------", CASH)
	INCOME_AFTER_CASH := INCOME - CASH
	fmt.Println("-----------------", INCOME_AFTER_CASH)
	TAX := 0.2 * INCOME_AFTER_CASH
	fmt.Println("-----------------", TAX)
	REVENUE := INCOME_AFTER_CASH - TAX
	fmt.Println("-----------------", REVENUE)
	PETROL := CALCULATE_PETROL()

	return REVENUE - PETROL
}

func CALCULATE_PER_HOUR() float64 {
	HOURS := CALCULATE_HOURS()
	REVENUE := CALCULATE_REVENUE()
	PER_HOUR := REVENUE / HOURS

	return PER_HOUR
}

func PRINT_STATISTICS() {
	
	fmt.Println("<< VK BOLT COURIER CALCULATOR >>")
	HOURS := CALCULATE_HOURS()
	fmt.Println("-> ", TWO_DECIMAL_POINTS(HOURS), "HOURS per month")
	PETROL := CALCULATE_PETROL()
	fmt.Println("-> PETROL: ", TWO_DECIMAL_POINTS(PETROL), "EURS per month")
	INCOME := CALCULATE_INCOME()
	fmt.Println("-> INCOME: ", TWO_DECIMAL_POINTS(INCOME), "EUR per month")
	CASH := CALCULATE_CASH()
	fmt.Println("-> +", TWO_DECIMAL_POINTS(CASH), "CASH per month")
	REVENUE := CALCULATE_REVENUE()
	fmt.Println("-> +", TWO_DECIMAL_POINTS(REVENUE), "REVENUE per month")
	PER_HOURS := CALCULATE_PER_HOUR()
	fmt.Println("-> ", TWO_DECIMAL_POINTS(PER_HOURS), "EUR/H")

}

/******************************************************************************/

type week_day struct {
	ID       int     `json:"id"`
	DATE     string  `json:"date"`
	INCOME   float64 `json:"income"`
	PETROL   float64 `json:"petrol"`
	REVENUE  float64 `json:"revenue"`
	HOURS    float64 `json:"hours"`
	CASH     float64 `json:"cash"`
	PER_HOUR float64 `json:"per_hour"`
	MONTH    string  `json:"month"`
}

var INCOME_TODAY float64
var PETROL_TODAY float64
var REVENUE_TODAY float64
var HOURS_TODAY float64
var CASH_TODAY float64
var MONTH_TODAY string

func CONSTRUCT_DAY_TRANSACTIONS() week_day {

	return week_day{
		ID:      Get_Unique_ID(WEEK_DATABASE),
		DATE:    Get_Current_Time("15:04 (02.01.2006)"),
		INCOME:  INCOME_TODAY,
		PETROL:  PETROL_TODAY,
		REVENUE: REVENUE_TODAY,
		HOURS:   HOURS_TODAY,
		CASH:    CASH_TODAY,
		MONTH:   GET_MONTH_NAME(),
	}
}

func GET_MONTH_NAME() string {
	currentTime := time.Now()
	return currentTime.Month().String()
}

/******************************************************************************/

func READ_FILE() []byte {
	file, err := ioutil.ReadFile(DB_URL)
	Error(err, "ReadFile")
	return file
}

func WRITE_FILE(dataBytes []byte) {

	var err = ioutil.WriteFile(DB_URL, dataBytes, 0644)
	Error(err, "WriteToFile")
}

func Error(err error, location string) {
	if err != nil {
		fmt.Println(" << Function name: ", location+" >>")
		fmt.Println(err.Error())

	}
}

func Get_Unique_ID(data []week_day) int {

	if len(data) == 0 {
		return 1
	}

	return data[len(data)-1].ID + 1
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

func TWO_DECIMAL_POINTS(number float64) string {
	return fmt.Sprintf("%.2f", number)
}

func CHECK_JSON() {

	if _, err := os.Stat(DB_URL); errors.Is(err, os.ErrNotExist) {
		WRITE_FILE([]byte("[]"))
	}
}

func Convert_CRLF_To_LF(reader *bufio.Reader) string {

	// Read the answer
	input, _ := reader.ReadString('\n')

	// Convert CRLF to LF
	input = strings.Replace(input, "\r\n", "", -1) /* "\r\n" was before.  */

	return input
}

func OPEN_JSON() []week_day {

	file := READ_FILE()
	data := CONVERT_TO_WEEK_DAY(file)

	return data
}

func CONVERT_TO_WEEK_DAY(file []byte) []week_day {

	a := []week_day{}

	err := json.Unmarshal(file, &a)
	Error(err, "CONVERT_TO_WEEK_DAY")

	return a
}
