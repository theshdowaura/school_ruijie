# !/usr/bin/env python
# -*-coding:utf-8-*-
import requests

USERNAME = '在这里填写你的用户名'
PASSWORD = '在这里填写你的密码'
SERVICE = r'在这里填写你的互联网接入商，需要对接入商的中文进行两次UrlEncode。如果无需选择互联网接入商，则此处留空。'
CAPTIVE_SERVER = r'http://www.google.cn/generate_204'


def get_captive_server_response():
    return requests.get(CAPTIVE_SERVER)


def login(response):
    response_text = response.text
    login_page_url = response_text.split('\'')[1]
    login_url = login_page_url.split('?')[0].replace('index.jsp', 'InterFace.do?method=login')
    query_string = login_page_url.split('?')[1]
    query_string = query_string.replace('&', '%2526')
    query_string = query_string.replace('=', '%253D')
    headers = {
        'Accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8',
        'Content-Type': 'application/x-www-form-urlencoded; charset=UTF-8',
        'Connection': 'keep-alive',
        'User-Agent': 'Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.132 Safari/537.36'
    }
    login_post_data = 'userId={}&password={}&service={}&queryString={}&operatorPwd=&operatorUserId=&validcode=&passwordEncrypt=false'.format(
        USERNAME, PASSWORD, SERVICE, query_string)
    login_result = requests.post(
        url=login_url,
        data=login_post_data,
        headers=headers
    )
    print(login_result.content.decode('utf-8'))


if __name__ == '__main__':
    captive_server_response = get_captive_server_response()
    if captive_server_response.status_code != 204:
        # Login when user is offline
        login(captive_server_response)
    else:
        # Exit script when user is online
        print('You are already online.')
