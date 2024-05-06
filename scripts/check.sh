#!/bin/bash

if [ $(expr $(date +%s) - $(git log origin/$1 --date=format:'%s' | head -10 | grep 'Date:' | awk '{print $2}')) -lt 200 ];then echo true > check_$1.ini; else echo false > check_$1.ini; fi
