ORIGIN_DIR=`pwd`;
cd $(dirname $(readlink -f $0));
cd src;

case $1 in 

"run")
    go run *.go
    ;;

"build")
    go buid *.go
    ;;
esac

cd .. ;

cd $ORIGIN_DIR