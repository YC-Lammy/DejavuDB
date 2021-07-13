package main

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

type user struct {
	name        string
	password    []byte
	id          int
	issue_date  time.Time
	expiry_time time.Time
}

type user_group struct {
	id    int
	users map[string]user // name of user
}

var number_of_users int = 0

var user_map = map[string]user_group{ // user_map will not be exposed to the out front
	/*
		GID 1–99 are reserved for the system and application use.
		GID 100+ allocated for the user’s group.
		UIDs 1–99 are reserved for other predefined accounts.
		UID 100–999 are reserved by system for administrative and system accounts/groups.
		UID 1000–10000 are occupied by applications account.
	*/
	"adm":      user_group{id: 1, users: map[string]user{}},    // admin, nearest to root
	"sudo":     user_group{id: 27, users: map[string]user{}},   // config permission, upgrade and maintainance
	"dev":      user_group{id: 30, users: map[string]user{}},   // developers, view logs and cofigs
	"analysts": user_group{id: 80, users: map[string]user{}},   // analystics, no admin permissions
	"user":     user_group{id: 100, users: map[string]user{}},  // regular user, no additional permissions
	"public":   user_group{id: 1000, users: map[string]user{}}, // public access, for application account
}

func useradd(message string) error { //this function can only be executed on router
	message = strings.Replace(message, "useradd ", "", 1)
	splited := strings.Split(message, " ")
	name := splited[len(splited)-1]
	group := "user"
	id := 1000 + number_of_users
	var expire time.Time = time.Now().AddDate(100, 0, 0)
	if len(splited) > 1 {
		for i, v := range splited {
			switch v {
			case "-G":
				group = splited[i+1]
			case "-u":
				a, _ := strconv.ParseInt(splited[i+1], 10, 64)
				id = int(a)

			case "-e": //expiry date YYYY-MM-DD
				datestr := strings.Split(splited[i+1], "-")
				var date []int
				for _, v := range datestr {
					a, _ := strconv.ParseInt(v, 10, 64)
					date = append(date, int(a))
				}

				expire = time.Date(date[0], time.Month(date[1]), date[2], 0, 0, 0, 0, time.Local)
			}
		}

	}
	if _, ok := user_map[group]; !ok {
		return errors.New("user group does not exist")
	}
	user_map[group].users[name] = user{name: name, id: id, issue_date: time.Now(),
		expiry_time: expire}
	return nil
}

func groupadd(message string) error {
	message = strings.Replace(message, "groupadd ", "", 1)
	splited := strings.Split(message, " ")
	name := splited[len(splited)-1]
	id := 1000 + number_of_users

	user_map[name] = user_group{id: id, users: map[string]user{}}

	return nil
}

func chmod(command string) error {
	return nil
}
