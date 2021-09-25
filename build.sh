#!/bin/bash
ORIGIN_DIR=`pwd`;

for command in node #
do
if ! command -v $command &> /dev/null
then
    echo $command could not be found
    exit 1
fi
done

command=''
outputdir=''

while test $# -gt 0; do
  case "$1" in
    -h|--help)
      echo " "
      echo "$package [options] application [arguments]"
      echo " "
      echo "options:"
      echo "-h, --help                show brief help"
      echo "-a, --action=ACTION       specify an action to use"
      echo "-o, --output-dir=DIR      specify a directory to store output in"
      exit 0
      ;;
    -a)
      shift
      if test $# -gt 0; then
        command=$1
      else
        echo "no process specified"
        exit 1
      fi
      shift
      ;;
    --action*)
      export PROCESS=`echo $1 | sed -e 's/^[^=]*=//g'`
      shift
      ;;
    -o)
      shift
      if test $# -gt 0; then
        outputdir=$1
      else
        echo "no output dir specified"
        exit 1
      fi
      shift
      ;;
    --output-dir*)
      export OUTPUT=`echo $1 | sed -e 's/^[^=]*=//g'`
      shift
      ;;
    run)
       command=run
      ;;
    build)
        command=build
    ;;
  esac
done


cd $(dirname $(readlink -f $0));
cd src;

case $command in

"run")
    go run *.go
    ;;

"build")
    go build *.go
    ;;
*)
    echo unknown command $command
;;
esac

cd .. ;

cd $ORIGIN_DIR
