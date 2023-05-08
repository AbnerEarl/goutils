#!/bin/bash

count=$(ps -ef | grep -E "start.sh|swag.sh|test.sh|auto.sh" | wc -l)
if [ $count -gt 1 ]; then
  /bin/cp -rf ../service_bak/* service/
  rm -rf ../service_bak/*
fi
