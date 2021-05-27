#!/bin/bash

while true; do
    secret=$1;
    user=$2;
    addr=`curl ifconfig.me 2> /dev/null`;

    res=`curl --location --request GET "lmbfao.ddns.net:8053/update?secret=$secret&domain=$user&addr=$addr" 2> /dev/null`;

    sleep 5s;
done

exit 0;