#!/bin/bash

/bin/cp -rf main.go main_bak.txt
uNames=`uname -s`
osName=${uNames: 0: 4}
if [ "$osName" == "Darw" ];then
	echo "MacOS X"
	sed -i "" "s/.*docs.*//g" main.go
	sed -i "" "s/.*swaggerFiles.*//g" main.go
	sed -i "" "s/.*ginSwagger.*//g" main.go
elif [ "$osName" == "Linu" ];then
	echo "GNU/Linux"
	sed -i "s/.*docs.*//g" main.go
	sed -i "s/.*swaggerFiles.*//g" main.go
	sed -i "s/.*ginSwagger.*//g" main.go
elif [ "$osName" == "MING" ];then
	echo "Windows, Git-bash"
	sed -i "s/.*docs.*//g" main.go
	sed -i "s/.*swaggerFiles.*//g" main.go
	sed -i "s/.*ginSwagger.*//g" main.go
else
	echo "Unknown OS, Not Support"
	rm -rf main_bak.txt
	exit 1
fi
go build main.go
mkdir -p ventoy/source/manage/
/bin/cp -rf ansible ventoy/source/manage/
/bin/cp -rf cert ventoy/source/manage/
/bin/cp -rf conf ventoy/source/manage/
/bin/cp -rf scheme ventoy/source/manage/
/bin/cp -rf main ventoy/source/manage/
/bin/cp -rf web.log ventoy/source/manage/
rm -rf main.go
mv main_bak.txt main.go
