#!/bin/bash

moduleName=$1
go test -v -timeout 20m -covermode=count -coverprofile=report/test_report.out -run="^Test" -coverpkg=$(go list ./... | grep -ve "/test" -ve "/docs" | grep "$moduleName/" | tr '\n' ',') ./...
result=$?
bash recover.sh
if [ $result == 0 ]; then
  go tool cover -html=report/test_report.out
else
  echo "==========================================================================================================="
  echo "The test case has some error, please check and fix it."
  echo "==========================================================================================================="
fi
