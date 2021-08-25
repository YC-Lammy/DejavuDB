module src

go 1.16

require (
	github.com/DmitriyVTitov/size v1.1.0
	github.com/fatih/color v1.12.0
	github.com/mattn/go-sqlite3 v1.14.8
	github.com/shirou/gopsutil v3.21.7+incompatible
	github.com/shirou/gopsutil/v3 v3.21.7
	golang.org/x/term v0.0.0-20210615171337-6886f2dfbf5b
	rogchap.com/v8go v0.6.0
	src/tensorflow v0.0.0
	src/sql v0.0.0
	src/datastore v0.0.0
)

replace (
	src/tensorflow => ./tensorflow
	src/sql => ./sql
	src/datastore => ./datastore
)
