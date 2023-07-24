package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	Gray   = "\033[37m"
)

func main() {
	CL()
}

/* COMMANDS */

func CL() {

	PRINT_STATISTICS()

	fmt.Println(Cyan + "\n<< COMMANDS: add | q >>" + Reset)
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

	INCOME = Question("Income: ")
	CASH = Question("Cash: ")
	PETROL = Question("Petrol: ")
	HOURS = Question("Hours: ")
	Clear_Screen()

	CL()
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

/* CALCULATIONS */
var HOURS float64 = 0
var INCOME float64 = 0
var PETROL float64 = 0
var CASH float64 = 0

func PRINT_STATISTICS() {

	fmt.Println()
	fmt.Println(Cyan + "<<___________ VK BOLT COURIER CALCULATOR v1 ___________>>" + Reset)
	fmt.Println()

	fmt.Println(Cyan + "__________________TAXABLE INCOME_____________________" + Reset)
	INCOME_TAXED := INCOME - CASH
	fmt.Println()
	fmt.Println("-> " + Green + "(INCOME" + Reset + " - " + Green + "CASH)" + Reset + " = " + Green + TWO_DECIMAL_POINTS(INCOME) + Reset + " - " + Green + TWO_DECIMAL_POINTS(CASH) + Reset + " = " + Yellow + TWO_DECIMAL_POINTS(INCOME_TAXED) + Reset + Green + " EUR" + Reset)
	fmt.Println()

	fmt.Println(Cyan + "__________________TAX________________________________" + Reset)
	TAX := INCOME_TAXED * 0.2
	fmt.Println()
	fmt.Println("-> " + Green + "(TAXABLE INCOME" + Reset + " * " + Red + "0.2)" + Reset + " = " + Green + TWO_DECIMAL_POINTS(INCOME_TAXED) + Reset + " * " + Red + "0.2" + Reset + " = " + Red + TWO_DECIMAL_POINTS(TAX) + " EUR" + Reset)
	fmt.Println()

	fmt.Println(Cyan + "__________________INCOME AFTER TAX___________________" + Reset)
	INCOME_AFTER_TAX := INCOME_TAXED - TAX
	fmt.Println()
	fmt.Println("-> " + "(" + Green + "TAXED INCOME" + Reset + " - " + Green + "TAX" + Reset + ")" + " = " + Green + TWO_DECIMAL_POINTS(INCOME_TAXED) + Reset + " - " + Reset + Red + TWO_DECIMAL_POINTS(TAX) + Reset + " = " + Yellow + TWO_DECIMAL_POINTS(INCOME_AFTER_TAX) + Reset + Green + " EUR" + Reset)
	fmt.Println()

	fmt.Println(Cyan + "__________________REVENUE___________________________" + Reset)
	REVENUE := (INCOME_AFTER_TAX - PETROL) + CASH
	fmt.Println()
	fmt.Println("-> " + "(" + Green + "INCOME AFTER TAX" + Reset + " - " + Red + "PETROL" + Reset + ")" + " + " + Green + "CASH" + Reset + " = " + "(" + Green + TWO_DECIMAL_POINTS(INCOME_AFTER_TAX) + Reset + " - " + Red + TWO_DECIMAL_POINTS(PETROL) + Reset + ")" + " + " + Green + TWO_DECIMAL_POINTS(CASH) + Reset + " = " + Yellow + TWO_DECIMAL_POINTS(REVENUE) + Reset + Green + " EUR" + Reset)
	fmt.Println()

	fmt.Println(Cyan + "__________________PER HOUR__________________________" + Reset)
	PER_HOUR := REVENUE / HOURS
	fmt.Println()
	fmt.Println("-> " + Green + "REVENUE" + Reset + "/" + Green + "HOURS" + Reset + " = " + Green + TWO_DECIMAL_POINTS(REVENUE) + Reset + "/" + Purple + TWO_DECIMAL_POINTS(HOURS) + Reset + " = " + Yellow + TWO_DECIMAL_POINTS(PER_HOUR) + Reset + Green + " EUR/H" + Reset)
	fmt.Println()

	fmt.Println(Cyan + "____________________________________________________" + Reset)
	fmt.Println(Red + "PROFIT SUMMARY: " + Reset + Yellow + TWO_DECIMAL_POINTS(REVENUE) + Reset + Green + " EUR" + Reset + " (" + Yellow + TWO_DECIMAL_POINTS(PER_HOUR) + Reset + Green + " EUR/H" + Reset + ")")
}

/* Help */

func Error(err error, location string) {
	if err != nil {
		fmt.Println(" << Function name: ", location+" >>")
		fmt.Println(err.Error())

	}
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

func Convert_CRLF_To_LF(reader *bufio.Reader) string {

	// Read the answer
	input, _ := reader.ReadString('\n')

	// Convert CRLF to LF
	input = strings.Replace(input, "\r\n", "", -1) /* "\r\n" was before.  */

	return input
}
