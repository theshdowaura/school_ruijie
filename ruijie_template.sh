#!/bin/sh

#If received parameters is less than 2, print usage
if [ "${#}" -lt "2" ]; then
  echo "Usage: ./ruijie_template.sh username password"
  echo "Example: ./ruijie_template.sh 201620000000 123456 987654"
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

#Structure loginURL
loginURL=`echo ${loginPageURL} | awk -F \? '{print $1}'`
loginURL="${loginURL/index.jsp/InterFace.do?method=login}"
#如果学校的认证需要选择运营商，需要填写此处的service，并删除注释符号
#service="%25E4%25B8%25AD%25E5%259B%25BD%25E7%25A7%25BB%25E5%258A%25A8ChinaMobile"
queryString="wlanuserip=822116b5f1fc86e82bef2ec112e2ca3b&wlanacname=c4f2fd6200d97669e67e88409950b214&ssid=&nasip=9a0225c89437df46244894fce5813368&snmpagentip=&mac=eb0ebea18508cc99aa180f8bc55dedbe&t=wireless-v2&url=2c0328164651e2b4f13b933ddf36628bea622dedcc302b30&apmac=&nasid=c4f2fd6200d97669e67e88409950b214&vid=2a971c39c0ea89a4&port=9bdeca5efa34f87c&nasportid=f5eb983692924fa26e6431fe9df4835f3b700d8620f11a71e53b66e60965df877578077260b1e309"
queryString="${queryString//&/%2526}"
queryString="${queryString//=/%253D}"

#Send Ruijie eportal auth request and output result
if [ -n "${loginURL}" ]; then
  authResult=`curl -s -A "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.91 Safari/537.36" -e "${loginPageURL}" -b "EPORTAL_COOKIE_USERNAME=; EPORTAL_COOKIE_PASSWORD=; EPORTAL_COOKIE_SERVER=; EPORTAL_COOKIE_SERVER_NAME=; EPORTAL_AUTO_LAND=; EPORTAL_USER_GROUP=; EPORTAL_COOKIE_OPERATORPWD=;" -d "userId=${1}&password=${2}&service=${service}&queryString=${queryString}&operatorPwd=${3}&operatorUserId=&validcode=&passwordEncrypt=false" -H "Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8" -H "Content-Type: application/x-www-form-urlencoded; charset=UTF-8" "${loginURL}"`
  echo $authResult
fi
