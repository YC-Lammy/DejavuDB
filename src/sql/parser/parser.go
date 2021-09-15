package parser

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var operaters = []string{
	"SELECT", "UPDATE", "DELETE", "INSERT", "INSERTINTO", "VALUES",

	"FROM", "DISTINCT",
	"WHERE", "HAVING", "ORDERBY", "TOP", "GROUPBY",
	"GO",
	"CREATE", "DROP",
	"AND",
	"CASE", "WHEN", "THEN", "ELSE", "END",

	"DECLARE", "SET", "DESC", "DESCRIBE",
}

func Parse(script string) (err error) {

	var script_strings = []string{}

	//
	// parse all strings
	// replace all strings by "%{index}%"
	//

	var string_opened = false
	var last_string_token byte
	var string_buffer = []byte{}

	for i := range script {
		switch token := script[i]; string(token) {
		case `'`, `"`, "`":
			if string_opened {
				if token == last_string_token {

					string_opened = false // close string

					str_token := " %" + strconv.Itoa(len(script_strings)) + "% "
					script = script[:i-(1+len(string_buffer))] + str_token + script[i+1:]
					script_strings = append(script_strings, string(string_buffer))

					string_buffer = []byte{}
				}
			} else {
				string_opened = true
				last_string_token = token
			}
		default:
			string_buffer = append(string_buffer, token)
		}
	}
	if string_opened {
		return errors.New("Invalid bracketing of name " + string(string_buffer))
	}

	//
	// end parsing strings
	//

	script = strings.ToUpper(script)
	script = strings.TrimFunc(script, func(r rune) bool {
		return !unicode.IsGraphic(r)
	})

	script = strings.Replace(script, "SELECT*FROM", "SELECT * FROM", -1)
	script = strings.Replace(script, "SELECT*", "SELECT *", -1)
	script = strings.Replace(script, "*FROM", "* FROM", -1)
	script = strings.Replace(script, "INSET INTO", "INSERTINTO", -1)
	script = strings.Replace(script, "ORDER BY", "ORDERBY", -1)
	script = strings.Replace(script, "GROUP BY", "GROUPBY", -1)
	script = strings.Replace(script, ",", " , ", -1)

	splited := strings.Split(script, " ")

	buf := []byte{}

	opened_condition := false
	condition := []string{}

	for _, w := range []byte(script) {
		if opened_condition {
			switch v := condition[len(condition)-1]; v {

			}
		}
		switch w {
		case ' ':
			switch string(buf) {
			case "SELECT":
			case "INSERTINTO":
			case "UPDATE":
			case "DELETE":
			case "CREATE":
			case "ALTER":
			case "DROP":
			case "DECLARE":
			case "SET":
			}
			buf = []byte{}
		default:
			buf = append(buf, w)
		}
	}

	err = nil

	return nil
}
