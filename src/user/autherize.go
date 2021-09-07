package user

import (
	"time"
	"../lazy"
	"../settings"
)

type token struct{
	permission uint16
	expiry time.Time
}

var tokens = map[string]*token{}

func Login(username, password string) bool{
	name := path.Join(os.UserHomeDir(),"dejavuDB", "users")+ os.PathSeparator + username
	
	if _, err := os.Stat(name); os.IsNotExist(err) {
		return false
	}
	var new = user{}
	f, err := os.Open(name)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	v, _ := ioutil.ReadAll(f)
	v,_:= lazy.DecryptAES([]byte(settings.AES_key),v)
	err = json.Unmarshal(v, &new)

	h := sha256.New()
	h.Write(new.Password_sauce)
	h.Write([]byte(password))
	if string(h.Sum(nil)) == string(new.Password_sum) {
		tokens[lazy.RandString(16)] = &token{ // create a new token that expires after an hour
			permission:new.permission
			expiry: time.Now().Add(time.Hours)
		}
		return true
	}
	return false
}

func CheckToken(token string) (*token,error){
	if v, ok:= tokens[token];ok{
		if time.Now().Sub(v.expiry) > 0 { // the token has expired
			delete(tokens, token)
			return nil, errors.New("token expired")
		}
		return v, nil // authorize passed
	}
	return nil, errors.New("incorrect token")
}

func PermissionByToken(token string) (uint16, error) {
	
	if v, ok:= tokens[token];ok{
		if time.Now().Sub(v.expiry) > 0 { // the token has expired
			delete(tokens, token)
			return 0, errors.New("token expired")
		}
		return v.permission, nil // authorize passed
	}
	return 0, errors.New("incorrect token")
}