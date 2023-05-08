#!/bin/bash

# In the service package, create a file and define a variable for mark need to mock method and set return values, such as:
# ========================================================================================================================
# service/mock_method.go

# package service

# var methods = map[interface{}]interface{}{
#   SetDefaultLink:      []interface{}{nil},
#   DelDefaultLink:      []interface{}{nil},
#   SetDefaultRouteLink: []interface{}{nil},
#   DelDefaultRouteLink: []interface{}{nil},
#   RegProxyAccount:     []interface{}{nil, nil},
#  }
# ========================================================================================================================

mkdir -p ../service_bak/
/bin/cp -rf service/* ../service_bak/
uNames=$(uname -s)
osName=$(echo $uNames | cut -c 1-4)
methods=($(grep ": " service/mock_method.go | awk -F':' '{print $1}'))
for i in $(seq ${#methods[@]}); do
  file=$(grep -ri "func ${methods[i - 1]}(" service | awk -F':' '{print $1}')
  name=$(grep -ri "func ${methods[i - 1]}(" service | awk -F':' '{print $NF}')
  bakName=$(echo $name | sed "s/func ${methods[i - 1]}/func ${methods[i - 1]}Bak/g")
  response=$(grep "${methods[i - 1]}" service/mock_method.go | awk -F'{' '{print $NF}')
  response=${response%??}
  if [ "$osName" == "Darw" ]; then
    #    echo "MacOS X"
    response="    return ${response}"
    ranges=($(grep -rn "func " service | grep -n -A 1 "func ${methods[i - 1]}" | awk -F':' '{print $(NF-1)}'))
    start=$(expr ${ranges[0]} + 1)
    end=$(expr ${ranges[1]} - 1)

    #  sed -i "" "${start},${end}d" $file
    sed -i "" -e "${ranges[0]}s/func ${methods[i - 1]}/func ${methods[i - 1]}Bak/g" $file
    #  sed -i "" "${ranges[0]}s/{/{\'$'\n/g" $file
    sed -i "" -e "${ranges[0]}i\\"$'\n'"}" $file
    sed -i "" -e "${ranges[0]}i\\"$'\n'"${response}" $file
    sed -i "" -e "${ranges[0]}i\\"$'\n'"${name}" $file
  elif [ "$osName" == "Linu" ]; then
    #    echo "GNU/Linux"
    response="    return ${response}"
    ranges=($(grep -rn "func " service | grep -n -A 1 "func ${methods[i - 1]}" | awk -F':' '{print $(NF-1)}'))
    start=$(expr ${ranges[0]} + 1)
    end=$(expr ${ranges[1]} - 1)

    sed -i "${ranges[0]}s/func ${methods[i - 1]}/func ${methods[i - 1]}Bak/g" $file
    sed -i "${ranges[0]}i\\}" $file
    sed -i "${ranges[0]}i\\${response}" $file
    sed -i "${ranges[0]}i\\${name}" $file
  elif [ "$osName" == "MING" ]; then
    echo "Windows, Git-bash"

  else
    echo "Unknown OS, Not Support"
    rm -rf ../service_bak
    exit 1
  fi
done
