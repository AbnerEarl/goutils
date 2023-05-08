#!/bin/bash

configFile=$1
serverPort=$2

localIP=$(grep host $configFile | awk '{print $NF}')
if [ -z $serverPort ]; then
  serverPort=$(grep "addr: :" $configFile | cut -d ":" -f3 | cut -d " " -f1)
fi

configPort=$(grep "@host" main.go | cut -d ":" -f2)

sed -i "s/.*:${configPort}/\/\/ @host ${localIP}:${serverPort}/g" main.go
sed -i "s/.*:${configPort}/	Host:             \"${localIP}:${serverPort}/g" docs/docs.go
sed -i "s/.*:${configPort}/    \"host\": \"${localIP}:${serverPort}/g" docs/swagger.json
sed -i "s/.*:${configPort}/host: ${localIP}:${serverPort}/g" docs/swagger.yaml

echo "=============================================================================="
echo "The Server Visit IP is: $localIP"
echo "The API Doc Address is: $localIP:${serverPort}/swagger/index.html"
echo "The API Base Address is: $localIP:${serverPort}/api/v1/"
echo "=============================================================================="

mode=$(grep "runmode" $configFile | cut -d' ' -f2)
if [ $mode == "release" ]; then
  echo "prod mode"
else
  bash mock.sh
fi
bash comment.sh
go run main.go -c $configFile
