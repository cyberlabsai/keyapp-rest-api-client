#!/bin/bash

if [[ $URL == "" ]];then
    echo "You must export the URL environment variable"
    exit 1
fi

if [[ $2 == "" ]];then
    echo "Usage: $(basename $0) access_token key_name"
    exit 1
else
    ACCESS_TOKEN=$1
    KEY_NAME=$2
fi

curl -x POST \
    "$URL/appkeys" \
    -H "Authorization: bearer $ACCESS_TOKEN" \
    --data "{\"name\": \"$KEY_NAME\"}"
