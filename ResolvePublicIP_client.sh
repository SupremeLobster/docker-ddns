#!/bin/bash

secret=$1;
user=$2;

res=`curl --location --request GET "lmbfao.ddns.net:8053/resolve?secret=$secret&domain=$user" 2> /dev/null`;
res=`echo $res | cut -f 2 -d "," | cut -f 2 -d ":"`;

res=${res%\"}
res=${res#\"}
echo $res

exit 0;