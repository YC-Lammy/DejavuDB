package dejavuDB

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/fatih/color"

	"src/lazy"
)

type Database struct {
	connection net.Conn
	connbuff   *bufio.Reader
	Admin      *admin // every connection can login once
	Closed     bool
}

func (db *Database) Send(message string) error {
	_, err := send(db.connection, []byte(message))
	return err
}

func (db *Database) Recieve() (string, error) {
	return recieve(db.connbuff)
}

func (db *Database) Close() {
	db.connection.Close()
}

func (db *Database) Set(location string, value interface{}, datatype ...string) error {
	dtype := ""
	strvalue := ""
	switch v := value.(type) {
	case int:
		dtype = "int"
		strvalue = strconv.Itoa(int(v))
	case int8:
		dtype = "int"
		strvalue = strconv.Itoa(int(v))
	case int16:
		dtype = "int"
		strvalue = strconv.Itoa(int(v))
	case int32:
		dtype = "int"
		strvalue = strconv.Itoa(int(v))
	case int64:
		dtype = "int"
		strvalue = strconv.Itoa(int(v))

	case string:
		dtype = "str"
		strvalue = v

	case float64:
		dtype = "float64"
		strvalue = fmt.Sprintf("%v", v)
	case float32:
		dtype = "float64"
		strvalue = fmt.Sprintf("%v", v)

	case bool:
		dtype = "bool"
		strvalue = strconv.FormatBool(v)

	case []byte:
		dtype = "[]byte"
		strvalue = string(v)

	case []string:
		dtype = "[]string"
		strvalue = strings.Replace(strings.Join(v, ","), " ", "", -1)

	case []int:
		dtype = "[]int"
		strvalue = strings.Replace(fmt.Sprintf("%v", v), " ", ",", -1)
	case []int64:
		dtype = "[]int"
		strvalue = strings.Replace(fmt.Sprintf("%v", v), " ", ",", -1)
	case []int32:
		dtype = "[]int"
		strvalue = strings.Replace(fmt.Sprintf("%v", v), " ", ",", -1)
	case []int16:
		dtype = "[]int"
		strvalue = strings.Replace(fmt.Sprintf("%v", v), " ", ",", -1)
	case []int8:
		dtype = "[]int"
		strvalue = strings.Replace(fmt.Sprintf("%v", v), " ", ",", -1)

	case []float64:
		dtype = "[]float64"
		strvalue = strings.Replace(fmt.Sprintf("%v", v), " ", ",", -1)
	case []float32:
		dtype = "[]float64"
		strvalue = strings.Replace(fmt.Sprintf("%v", v), " ", ",", -1)

	case []bool:
		dtype = "[]bool"
		strvalue = strings.Replace(fmt.Sprintf("%v", v), " ", ",", -1)

	case [][]byte:
		dtype = "[][]byte"
		r := []string{}
		for _, a := range v {
			r = append(r, string(a))
		}
		strvalue = strings.Replace(strings.Join(r, ","), " ", "", -1)
	}

	a, _ := json.Marshal(value)
	if len(a) > 2 {
		if a[0] == '{' && a[len(a)-1] == '}' {
			if strings.Contains(string(a), `":`) {
				dtype = "json"
				strvalue = string(a)
			}
		}
	}

	err := db.Send("Set " + location + " " + strvalue + " " + dtype)
	if err != nil {
		return err
	}
	r, err := db.Recieve()
	if err != nil {
		return err
	}
	if r != "sucess" {
		return errors.New(r)
	}
	return nil
}

func Connect(router_addr string) (*Database, error) {

	conn, err := net.Dial("tcp", router_addr)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	send(conn, []byte(`{"role":"client"}`)) // send my config to router, router reads and decides

	connbuff := bufio.NewReader(conn)

	_, err = recieve(connbuff) // ignore config

	if err != nil {
		return nil, err
	}

	return &Database{connection: conn, connbuff: connbuff, Admin: &admin{connection: conn, connbuf: connbuff, token: ""}}, nil
}

func command_syntax_checker(text string) error {
	splited := strings.Split(text, " ")
	if !(lazy.Contains(list_of_commands, splited[0])) {
		return errors.New(splited[0] + ": command not found")
	}
	return nil
}

func command_result_output(command string, result string) {
	splited := strings.Split(command, " ")
	switch splited[0] {
	case "df", "free":
		table := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
		fmt.Fprint(table, result)
		table.Flush()
	case "dstat":
		color.HiBlue("You did not select any stats, using -cdngy by default.\n")
		header_lengths := []int{8, 13, 14, 16}
		column_lens := []int{0, 0, 0, 0, 0, 0, 0, 0}
		lines := strings.Split(result, ";")
		columns := [][]string{}
		for _, v := range lines {
			if v == "" {
				continue
			}
			tab := strings.Split(v, ",")
			columns = append(columns, tab)
			length := []int{len(tab[0]), len(tab[1]), len(tab[2]), len(tab[3]), len(tab[4]), len(tab[5]), len(tab[6]), len(tab[7])}
			for i, v := range length {
				if v > column_lens[i] {
					column_lens[i] = v
				}
			}
		}
		if column_lens[2]+column_lens[3]+column_lens[4] > 14 {
			header_lengths[2] = column_lens[2] + column_lens[3] + column_lens[4]
		}
		if column_lens[5]+column_lens[6]+column_lens[7] > 16 {
			header_lengths[3] = column_lens[5] + column_lens[6] + column_lens[7]
		}

		color.Blue("%8s   %13s    %"+strconv.Itoa(header_lengths[2])+"s     %"+strconv.Itoa(header_lengths[3])+"s   \n",
			"--Role--", "--cpu-usage--", "--disk-usage--", "--memory-usage--")

		a := strconv.Itoa(column_lens[2])
		b := strconv.Itoa(column_lens[3])
		c := strconv.Itoa(column_lens[4])
		d := strconv.Itoa(column_lens[5])
		e := strconv.Itoa(column_lens[6])
		f := strconv.Itoa(column_lens[7])
		underline := color.New(color.FgHiBlue, color.Underline, color.Bold)
		underline.Print(fmt.Sprintf("%-8s | %13s | %"+a+"s %"+b+"s %"+c+"s | %"+d+"s %"+e+"s %"+f+"s |\n",
			"role", "usage%", "used",
			"avaliable", "use%", "used",
			"avaliable", "use%"))

		for _, v := range columns[1:] {

			fmt.Printf("%-8s | %13s | %"+a+"s %"+b+"s %"+c+"s | %"+d+"s %"+e+"s %"+f+"s |\n",
				v[0], v[1], v[2], v[3], v[4], v[5], v[6], v[7])
		}
	default:
		fmt.Println(result)
	}
}
