package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"syscall"
	"text/tabwriter"

	"golang.org/x/term"
)

func Client_dial(router_addr string, mycfg []byte) error {

	defer wg.Done()

	conn, err := net.Dial("tcp", router_addr)

	if err != nil {
		log.Println(err)
		return err
	}

	defer conn.Close()

	send(conn, mycfg) // send my config to router, router reads and decides

	connbuff := bufio.NewReader(conn)

	_, err = recieve(connbuff) // ignore config

	if err != nil {
		log.Println(err)
		return err
	}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter username: ")
	scanner.Scan()
	user := scanner.Text()

	fmt.Print("Enter password: ")

	bytepw, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		os.Exit(1)
	}
	pass := string(bytepw)
	fmt.Printf("\nYou've entered: %q\n\n", pass)

	for {

		fmt.Print(user + "$ ")
		// Scans a line from Stdin(Console)
		scanner.Scan()
		// Holds the string that scanned
		text := scanner.Text()

		err = command_syntax_checker(text)

		if err == nil {
			if len(text) != 0 {

				send(conn, []byte("CLIENT "+text))

				result, err := recieve(connbuff) // wait for message to recieve

				if err != nil {
					log.Println(err)

				} else {

					command_result_output(text, result)
				}
			}
		} else {
			fmt.Println(err)
		}

	}
}

func command_syntax_checker(text string) error {
	splited := strings.Split(text, " ")
	if !(contains(list_of_commands, splited[0])) {
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
	default:
		fmt.Println(result)
	}
}
