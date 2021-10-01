package user

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
	Id             uint32
	Gid            uint32 // largest gid
	Group          string
	permission     uint16
	Issue_date     time.Time
	Expiry_time    time.Time
	Domain         string
}

type user_group struct {
	Id    int
	Users map[string]user // name of user
	Name  string
}

type user_token struct {
	name        string
	user_name   string
	user_group  string
	expiry_time time.Time
}

var user_tokens = map[string]*user_token{}

var number_of_users int = 1

var users = map[string]*user{}

var home_dir string

var groups = map[string]*user_group{ // user_map will not be exposed to the out front
	/*
		GID 1–99 are reserved for the system and application use.
		GID 100+ allocated for the user’s group.
		UIDs 1–99 are reserved for other predefined accounts.
		UID 100–999 are reserved by system for administrative and system accounts/groups.
		UID 1000–10000 are occupied by applications account.
	*/
	"adm":     &user_group{Id: 1},    // admin, nearest to root
	"sudo":    &user_group{Id: 27},   // config permission, upgrade and maintainance
	"dev":     &user_group{Id: 30},   // developers, view logs and cofigs
	"analyst": &user_group{Id: 80},   // analystics, no admin permissions
	"user":    &user_group{Id: 100},  // regular user, no additional permissions
	"other":   &user_group{Id: 1000}, // public access, for application account
}

func init() {
	d, _ := os.UserHomeDir()
	home_dir = d

	origin := path.Join(home_dir, "dejavuDB", "users")
	os.Mkdir(origin, os.ModePerm)

	if _, err := os.Stat(path.Join(origin, "root")); os.IsNotExist(err) {

		sauce := make([]byte, 16)
		rand.Read(sauce)
		h := sha256.New()
		h.Write(sauce)
		h.Write([]byte(""))
		root := user{Name: "root", Id: 1, Gid: 1, Group: "adm", Domain: "localhost", Password_sauce: sauce, Password_sum: h.Sum(nil)}
		//groups["adm"].Users["root"] = root
		f, _ := os.Create(path.Join(origin, "root"))
		b, _ := json.Marshal(root)
		f.Write(b)
		if err != nil {
			return
		}
		f.Close()
	}
	arr, _ := ioutil.ReadDir(origin)
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
		//groups[new.Group].Users[new.Name] = new
		fmt.Println(new.Name)
	}

	os.Chdir(origin)
}

func UserExist(username string) (*user, bool) { // return the user and bool

	if v, ok := users[username]; ok {
		return v, true
	}

	return nil, false
}

func Useradd(message string) error { //this function can only be executed on router
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
	if _, ok := UserExist(name); ok {
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
	if _, ok := groups[group]; !ok {
		return errors.New("user group does not exist")
	}
	// generate password hash
	sauce := make([]byte, 16)
	rand.Read(sauce)
	h := sha256.New()
	h.Write(sauce)
	h.Write([]byte(password))
	new := user{Name: name, Id: uint32(id), Issue_date: time.Now(),
		Expiry_time: expire, Password_sum: h.Sum(nil), Password_sauce: sauce}
	groups[group].Users[name] = new

	f, _ := os.Create(path.Join(home_dir, "dejavuDB", "users") + string(os.PathSeparator) + name)
	enc := gob.NewEncoder(f)
	err := enc.Encode(new)
	if err != nil {
		return err
	}
	f.Close()
	return nil
}

func Groupadd(message string) error {
	message = strings.Replace(message, "groupadd ", "", 1)
	splited := strings.Split(message, " ")
	name := splited[len(splited)-1]
	id := 1000 + number_of_users

	groups[name] = &user_group{Id: id, Users: map[string]user{}}

	return nil
}

func Userid(message string) error {

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
