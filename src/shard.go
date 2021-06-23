package router

import (
	"fmt"
	"time"

	"github.com/DmitriyVTitov/size"
)

var shardData map[string]interface{}
var type_map = map[string]interface{}{}

func NewShard() {
	s := time.Now()
	fmt.Printf("%s : Shard Server Start", s)

	users := map[string]interface{}{
		/*
			GID 1–99 are reserved for the system and application use.
			GID 100+ allocated for the user’s group.
			UIDs 1–99 are reserved for other predefined accounts.
			UID 100–999 are reserved by system for administrative and system accounts/groups.
			UID 1000–10000 are occupied by applications account.
			UID 10000+ are used for user accounts.
		*/
		"adm":     map[string]interface{}{"id": 1},    // admin, nearest to root
		"sudo":    map[string]interface{}{"id": 27},   // config permission, upgrade and maintainance
		"dev":     map[string]interface{}{"id": 30},   // developers, view logs and cofigs
		"monitor": map[string]interface{}{"id": 80},   // analystics, no admin permissions
		"user":    map[string]interface{}{"id": 100},  // regular user, no additional permissions
		"public":  map[string]interface{}{"id": 1000}} // public access, no authorization needed

	shardData["groups"] = users
	println(size.Of(shardData))

}
