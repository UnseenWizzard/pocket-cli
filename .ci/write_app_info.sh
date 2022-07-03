#!/bin/sh

APP_INFO="pkg/util/appinfo.go"

if [ "$#" -ne 2 ]; then
    echo "Missing required arguments"
    echo "run as: ./write_app_info.sh APPLICATION_ID VERSION"
fi


touch $APP_INFO

echo "package util

const PocketAppId = \"$1\"
const Version = \"$2\"" > $APP_INFO