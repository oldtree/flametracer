#!/bin/shell

name=flametracer

function build() {
    go build -o $name main.go 
    sc=$?
    if [ $sc -ne 0 ]; then
        echo "build error"
        exit $sc
    else
        timestamp=time
        echo "build ok, timestamp : $time"
    fi
}

function package() {
    build
    tar zcvf $name.tar.gz $name config.json
}

function start() {
    nohup ./$name > /dev/null & 
    pid=`ps -A |grep "$name"| awk '{print $1}'`
    echo "process start : $pid "
}

function stop() {
    pid=`ps -A |grep "$name"| awk '{print $1}'`
    kill -s 9 $pid
}

action=$1
case $action in
    "build" )
        build
        ;;
    "package" )
        package
        ;;
    "start" )
        start
        ;;
    "stop" )
        stop
        ;;
esac