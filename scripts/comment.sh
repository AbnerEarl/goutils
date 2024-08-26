#!/bin/bash

# To run this script, the following parameters are required:
# $modelPath: modelPath, the orm model path. default: model
# $packageName: packageName, gen file package name. default: base
# $fileName: fileName,gen file name. default: comment.go

modelPath=$1
packageName=$2
fileName=$3

if [ -z $modelPath ]; then
  modelPath="model"
fi

if [ -z $packageName ]; then
  packageName="base"
fi

if [ -z $fileName ]; then
  fileName="comment.go"
fi

files=($(ls $modelPath | grep -v init.go))
mkdir -p $(dirname "$fileName")
echo "package $packageName

var Comment = map[string]map[string]interface{}{" >$fileName
for i in $(seq ${#files[@]}); do
  names=($(cat model/${files[i - 1]} | grep -e "comment:" -e "BaseModel$" -A1 | grep "comment:" -B1 | awk -F'json:' '{print $2}' | awk -F '\"' '{print $2}' | sed 's/^$/###/g'))
  comments=($(cat model/${files[i - 1]} | grep -e "comment:" -e "BaseModel$" -A1 | grep "comment:" -B1 | awk -F'comment:' '{print $2}' | awk -F "\'" '{print $2}' | sed 's/^$/###/g'))
  if [ ${#names[@]} -lt 1 ]; then
    continue
  fi
  modelNames=($(cat model/${files[i - 1]} | grep -e "comment:" -e "BaseModel$" | grep -e "comment:" B2 | grep " struct " | awk '{print $2}' ))
  k=0
  tableName=$(cat model/${files[i - 1]} | grep -e "${modelNames[$k]}.*TableName" -A2 | grep "return " | awk '{print $NF}')
  tableName=$(echo ${tableName} | awk '{print $1}')
  echo "	${tableName}: map[string]interface{}{" >>$fileName
  for j in $(seq ${#names[@]}); do
    if [ $j -gt 0 ] && [ ${names[j - 1]} == "###" ] && [ ${names[j]} == "###" ]; then
        continue
    fi
    if [ ${names[j - 1]} == "###" ] && [ $j -gt 1 ]; then
      echo "	}," >>$fileName
      let k++
      tableName=$(cat model/${files[i - 1]} | grep -e "${modelNames[$k]}.*TableName" -A2 | grep "return " | awk '{print $NF}')
      echo "	${tableName}: map[string]interface{}{" >>$fileName
      continue
    fi
    if [[ ${names[j - 1]} == "###" ]];then
      continue
    fi
    echo "		\"${names[j - 1]}\": \"${comments[j - 1]}\"," >>$fileName
  done
  echo "	}," >>$fileName
done

echo "}
" >>$fileName
