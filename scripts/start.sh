#!/bin/bash

configFile=$1 #'conf/config_dev.yaml'
apiDir=$2 #'handler|config'
localIP=$3 #'150.109.95.231'

uNames=$(uname -s)
osName=${uNames:0:4}
if [ "$osName" == "Darw" ]; then
  echo "MacOS"
elif [ "$osName" == "Linu" ]; then
  echo "GNU/Linux"
elif [ "$osName" == "MING" ]; then
  echo "Windows, Git-bash"
else
  echo "Unknown OS"
  exit 2
fi

if [ "$osName" != "Darw" ]; then
  gateway=$(ip route get 8.8.8.8 | grep via | cut -d' ' -f3)
else
  gateway=$(route get 8.8.8.8 | grep gateway | cut -d':' -f2)
fi

source /etc/profile



if [ ! -z "$localIP" ]; then
  netp=${gateway%.*}
  localIP=$(ifconfig | grep $netp)
  localIP=${localIP##inet*}
  localIP=$(echo $localIP | cut -d' ' -f2)
fi

#swag init
#swag init --exclude $(ls -d */ | grep -v 'handler|config' | tr '\n' ',')
swag init --exclude $(ls -d */ | grep -v "$apiDir" | tr '\n' ',')

serverPort=$(grep "addr: :" $configFile | cut -d ":" -f3 | cut -d " " -f1)
configPort=$(grep "@host" main.go | cut -d ":" -f2)
if [ "$osName" != "Darw" ]; then
  sed -i "s/.*:${configPort}/\/\/ @host ${localIP}:${serverPort}/g" main.go
  sed -i "s/.*:${configPort}/	Host:             \"${localIP}:${serverPort}/g" docs/docs.go
  sed -i "s/.*:${configPort}/    \"host\": \"${localIP}:${serverPort}/g" docs/swagger.json
  sed -i "s/.*:${configPort}/host: ${localIP}:${serverPort}/g" docs/swagger.yaml
else
  sed -i "" "s/.*:${configPort}/\/\/ @host ${localIP}:${serverPort}/g" main.go
  sed -i "" "s/.*:${configPort}/	Host:             \"${localIP}:${serverPort}/g" docs/docs.go
  sed -i "" "s/.*:${configPort}/    \"host\": \"${localIP}:${serverPort}/g" docs/swagger.json
  sed -i "" "s/.*:${configPort}/host: ${localIP}:${serverPort}/g" docs/swagger.yaml
fi

echo "package base

var ApiDesc = map[string]string{" >base/api_desc.go
basePath=$(grep "basePath" docs/swagger.yaml | cut -d' ' -f2)
grep -e "  /" -e "      description" docs/swagger.yaml >temp_api.log
tag=1
apiName=
apiDesc=
while read line; do
  count=$(echo $line | grep "description:" | wc -l)
  if [ $tag == 3 ] && [ $count == 0 ]; then
    tag=1
  fi
  line=$(echo $line)
  if [ $tag == 1 ]; then
    apiName=$(echo $line | cut -d':' -f1)
    let tag++
  elif [ $tag == 2 ]; then
    apiDesc=$(echo $line | cut -d' ' -f2)
    echo "	\"${basePath}${apiName}\": \"${apiDesc}\"," >>base/api_desc.go
    let tag++
  elif [ $tag == 3 ]; then
    tag=1
  fi
done <temp_api.log
echo "}
" >>base/api_desc.go
rm -rf temp_api.log
echo "=============================================================================="
echo "The Server Visit IP is: $localIP"
echo "The API Doc Address is: $localIP:${serverPort}/swagger/index.html"
echo "The API Base Address is: $localIP:${serverPort}/api/v1/"
echo "=============================================================================="

go run main.go -c $configFile
