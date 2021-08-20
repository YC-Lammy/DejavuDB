cd src

case $1 in 

"run")
    go run *.go
    ;;

"build")
    go buid *.go
    ;;
esac

cd ..