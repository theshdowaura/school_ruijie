#!/bin/bash

#If received logout parameter, send a logout request to eportal server
if [ "${1}" = "logout" ]; then
  userIndex=`curl -s -A "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.91 Safari/537.36" -I http://10.8.2.2/eportal/redirectortosuccess.jsp | grep -o 'userIndex=.*'` #Fetch user index for logout request
  logoutResult=`curl -s -A "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.91 Safari/537.36" -d "${userIndex}" http://10.8.2.2/eportal/InterFace.do?method=logout`
  echo $logoutResult
  exit 0
fi

#If received parameters is less than 3, print usage
if [ "${#}" -lt "3" ]; then
  echo "Usage: ./ruijie_jmu.sh service username password"
  echo "Service parameter can be \"chinamobile\", \"chinanet\" and \"chinaunicom\". If service parameter do not set as these value, it will use campus network as default internet service provider."
  echo "Example: ./ruijie_jmu.sh chinanet 201620000000 123456"
  echo "if you want to logout, use: ./ruijie_jmu.sh logout"
  exit 1
fi

#Exit the script when is already online, use www.google.cn/generate_204 to check the online status
captiveReturnCode=`curl -s -I -m 10 -o /dev/null -s -w %{http_code} http://www.google.cn/generate_204`
if [ "${captiveReturnCode}" = "204" ]; then
  echo "You are already online!"
  exit 0
fi

#If not online, begin Ruijie Auth

#Get Ruijie login page URL
loginPageURL=`curl -s "http://www.google.cn/generate_204" | awk -F \' '{print $2}'`

chinamobile="%25E7%25A7%25BB%25E5%258A%25A8%25E5%25AE%25BD%25E5%25B8%25A6%25E6%258E%25A5%25E5%2585%25A5"
chinanet="%25E7%2594%25B5%25E4%25BF%25A1%25E5%25AE%25BD%25E5%25B8%25A6%25E6%258E%25A5%25E5%2585%25A5"
chinaunicom="%25E8%2581%2594%25E9%2580%259A%25E5%25AE%25BD%25E5%25B8%25A6%25E6%258E%25A5%25E5%2585%25A5"
campus="%25E6%2595%2599%25E8%2582%25B2%25E7%25BD%2591%25E6%258E%25A5%25E5%2585%25A5"

service=""

if [ "${1}" = "chinamobile" ]; then
  echo "Use ChinaMobile as internet service provider."
  service="${chinamobile}"
fi

if [ "${1}" = "chinanet" ]; then
  echo "Use ChinaNet as internet service provider."
  service="${chinanet}"
fi

if [ "${1}" = "chinaunicom" ]; then
  echo "Use ChinaUnicom as internet service provider."
  service="${chinaunicom}"
fi

if [ -z "${service}" ]; then
  echo "Use Campus Network internet service provider."
  service="${campus}"
fi

#Structure loginURL
loginURL=`echo ${loginPageURL} | awk -F \? '{print $1}'`
loginURL="${loginURL/index.jsp/InterFace.do?method=login}"

#Structure quertString
queryString=`echo ${loginPageURL} | awk -F \? '{print $2}'`
queryString="${queryString//&/%2526}"
queryString="${queryString//=/%253D}"

#Send Ruijie eportal auth request and output result
if [ -n "${loginURL}" ]; then
  authResult=`curl -s -A "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.91 Safari/537.36" -e "${loginPageURL}" -b "EPORTAL_COOKIE_USERNAME=; EPORTAL_COOKIE_PASSWORD=; EPORTAL_COOKIE_SERVER=; EPORTAL_COOKIE_SERVER_NAME=; EPORTAL_AUTO_LAND=; EPORTAL_USER_GROUP=; EPORTAL_COOKIE_OPERATORPWD=;" -d "userId=${2}&password=${3}&service=${service}&queryString=${queryString}&operatorPwd=&operatorUserId=&validcode=&passwordEncrypt=false" -H "Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8" -H "Content-Type: application/x-www-form-urlencoded; charset=UTF-8" "${loginURL}"`
  echo $authResult
fi