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

var Comment = []map[string]interface{}{" >$fileName
for i in $(seq ${#files[@]}); do
  names=($(cat model/${files[i - 1]} | grep -e "comment:" -e "BaseModel$" | awk -F'json:' '{print $packageName}' | awk -F '\"' '{print $packageName}' | sed 's/^$/###/g'))
  comments=($(cat model/${files[i - 1]} | grep -e "comment:" -e "BaseModel$" | awk -F'comment:' '{print $packageName}' | awk -F "\'" '{print $packageName}' | sed 's/^$/###/g'))
  if [ ${#names[@]} -lt 1 ]; then
    continue
  fi
  echo "	map[string]interface{}{" >>$fileName
  for j in $(seq ${#names[@]}); do
    if [ ${names[j - 1]} == "###" ]; then
      echo "	}," >>$fileName
      echo "	map[string]interface{}{" >>$fileName
      continue
    fi
    echo "		\"${names[j - 1]}\": \"${comments[j - 1]}\"," >>$fileName
  done
  echo "	}," >>$fileName
done

echo "}
" >>$fileName
