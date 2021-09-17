package main

import (
	"crypto/sha256"
	"encoding/gob"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	json "github.com/goccy/go-json"
)

/* user.go is database administration user related
   users in this registration are for internal purpose only,
   application authorization is strongly unrecomanded for security concerns
*/

type user struct {
	Name           string
	Password_sum   []byte
	Password_sauce []byte // some random bytes to secure the password sum
	Id             int
	Group          string
	Issue_date     time.Time
	Expiry_time    time.Time
	Domain         string
	//token          *user_token
}

type user_group struct {
	Id    int
	Users map[string]*user // name of user
}

type user_token struct {
	user_name   string
	user_group  string
	expiry_time time.Time
}

var user_tokens = map[string]*user_token{}

var number_of_users int = 1

var users = map[string]*user{}

var user_map = map[string]*user_group{ // user_map will not be exposed to the out front
	/*
		GID 1–99 are reserved for the system and application use.
		GID 100+ allocated for the user’s group.
		UIDs 1–99 are reserved for other predefined accounts.
		UID 100–999 are reserved by system for administrative and system accounts/groups.
		UID 1000–10000 are occupied by applications account.
	*/
	"adm":      &user_group{Id: 1, Users: map[string]*user{}},    // admin, nearest to root
	"sudo":     &user_group{Id: 27, Users: map[string]*user{}},   // config permission, upgrade and maintainance
	"dev":      &user_group{Id: 30, Users: map[string]*user{}},   // developers, view logs and cofigs
	"analysts": &user_group{Id: 80, Users: map[string]*user{}},   // analystics, no admin permissions
	"user":     &user_group{Id: 100, Users: map[string]*user{}},  // regular user, no additional permissions
	"other":    &user_group{Id: 1000, Users: map[string]*user{}}, // public access, for application account
}

func init() {

	origin := path.Join(home_dir, "dejavuDB")
	os.Chdir(origin)
	os.Mkdir("users", os.ModePerm)
	os.Chdir("users")
	/*os.Mkdir("adm", os.ModePerm)
	os.Mkdir("sudo", os.ModePerm)
	os.Mkdir("dev", os.ModePerm)
	os.Mkdir("analysts", os.ModePerm)
	os.Mkdir("user", os.ModePerm)
	os.Mkdir("other", os.ModePerm)
	*/

	if _, err := os.Stat("root"); os.IsNotExist(err) {

		sauce := make([]byte, 16)
		rand.Read(sauce)
		h := sha256.New()
		h.Write(sauce)
		h.Write([]byte(""))
		root := user{Name: "root", Id: 1, Group: "adm", Domain: "localhost", Password_sauce: sauce, Password_sum: h.Sum(nil)}
		user_map["adm"].Users["root"] = &root
		f, _ := os.Create("root")
		b, _ := json.Marshal(root)
		f.Write(b)
		if err != nil {
			return
		}
		f.Close()
	}
	arr, _ := ioutil.ReadDir(path.Join(origin, "users"))
	for _, v := range arr {
		var new = user{}
		f, err := os.Open(v.Name())
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
		v, _ := ioutil.ReadAll(f)
		err = json.Unmarshal(v, &new)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
		users[new.Name] = &new
		if new.Group == "" {
			panic("group not exist")
		}
		user_map[new.Group].Users[new.Name] = &new
	}

	os.Chdir(origin)
}

func userExist(username string) (*user, bool) { // return the user and bool
	for _, v := range user_map {
		if v, ok := v.Users[username]; ok {
			return v, true
		}
	}
	return nil, false
}

func userLogin(username string, password string) error {
	user, _ := userExist(username)
	if user == nil {
		return errors.New("invalid")
	}
	h := sha256.New()
	h.Write(user.Password_sauce)
	h.Write([]byte(password))
	if string(h.Sum(nil)) == string(user.Password_sum) {
		return nil
	} else {
		return errors.New("invalid")
	}
}

func useradd(message string) error { //this function can only be executed on router
	message = strings.Replace(message, "useradd", "", 1)
	splited := strings.Split(message, " ")
	if len(splited) < 1 {
		return errors.New("not enough arguments")
	}
	fmt.Println(splited)
	name := splited[len(splited)-1]
	group := "user"
	id := 1000 + number_of_users
	password := ""
	var expire time.Time = time.Now().AddDate(100, 0, 0)
	if _, ok := userExist(name); ok {
		return errors.New("username exist")
	}
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
			case "-p": // password
				password = splited[i+1]
			}
		}
	}
	// check if user group exist
	if _, ok := user_map[group]; !ok {
		return errors.New("user group does not exist")
	}
	// generate password hash
	sauce := make([]byte, 16)
	rand.Read(sauce)
	h := sha256.New()
	h.Write(sauce)
	h.Write([]byte(password))
	new := user{Name: name, Id: id, Issue_date: time.Now(),
		Expiry_time: expire, Password_sum: h.Sum(nil), Password_sauce: sauce}
	user_map[group].Users[name] = &new

	f, _ := os.Create(path.Join(home_dir, "dejavuDB", "users") + string(os.PathSeparator) + name)
	enc := gob.NewEncoder(f)
	err := enc.Encode(new)
	if err != nil {
		return err
	}
	f.Close()
	command_result = "sucess"
	return nil

}

func groupadd(message string) error {
	message = strings.Replace(message, "groupadd ", "", 1)
	splited := strings.Split(message, " ")
	name := splited[len(splited)-1]
	id := 1000 + number_of_users

	user_map[name] = &user_group{Id: id, Users: map[string]*user{}}

	return nil
}

func userid(message string) error {

	return nil
}

func chmod(command string) error {
	return nil
}

func chown(command string) error {
	return nil
}

func chgrp(command string) error {
	return nil
}
