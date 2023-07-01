#include <iostream>
#include <curl/curl.h>

using namespace std;

const char *captiveServerUrl = "http://www.google.cn/generate_204";

void printUsage() {
    cout << "Usage: ./ruijie service_name username password" << endl;
    cout
            << R"(Service parameter can be "chinamobile", "chinanet" and "chinaunicom". If service parameter do not set as these value, it will use campus network as default internet service provider.)"
            << endl;
    cout << "Example: ./ruijie chinanet 201620000000 123456" << endl;
    cout << "if you want to logout, use: ./ruijie logout" << endl;
}

size_t noPrintWriteCallback(char *ptr, size_t size, size_t nmemb, void *userdata) {
    return size * nmemb;
}

size_t writeCallback(char *ptr, size_t size, size_t nmemb, void *userdata) {
    ((string *) userdata)->append((char *) ptr, size * nmemb);
    return size * nmemb;
}

long getResponseStatusCode(const char *url) {
    CURL *curl;
    CURLcode res;
    long responseCode = -1;
    curl = curl_easy_init();
    if (curl) {
        curl_easy_setopt(curl, CURLOPT_URL, url);
        curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, noPrintWriteCallback);
        res = curl_easy_perform(curl);
        if (res != CURLE_OK) {
            cout << "Can not get response content of " << url << ", exit." << endl;
            exit(EXIT_FAILURE);
        }
        curl_easy_getinfo(curl, CURLINFO_RESPONSE_CODE, &responseCode);
        curl_easy_cleanup(curl);
    }
    return responseCode;
}

string getResponseContent(const char *url) {
    CURL *curl;
    CURLcode res;
    string responseContent;
    curl = curl_easy_init();
    if (curl) {
        curl_easy_setopt(curl, CURLOPT_URL, url);
        curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, writeCallback);
        curl_easy_setopt(curl, CURLOPT_WRITEDATA, &responseContent);
        res = curl_easy_perform(curl);
        if (res != CURLE_OK) {
            cout << "Can not get response content of " << url << ", exit." << endl;
            exit(EXIT_FAILURE);
        }
        curl_easy_cleanup(curl);
    }
    return responseContent;
}

string getServiceStringByServiceName(char *serviceName) {
    string serviceNameString = serviceName;
    if (serviceNameString == "chinamobile") {
        return "%25E7%25A7%25BB%25E5%258A%25A8%25E5%25AE%25BD%25E5%25B8%25A6%25E6%258E%25A5%25E5%2585%25A5";
    } else if (serviceNameString == "chinanet") {
        return "%25E7%2594%25B5%25E4%25BF%25A1%25E5%25AE%25BD%25E5%25B8%25A6%25E6%258E%25A5%25E5%2585%25A5";
    } else if (serviceNameString == "chinaunicom") {
        return "%25E8%2581%2594%25E9%2580%259A%25E5%25AE%25BD%25E5%25B8%25A6%25E6%258E%25A5%25E5%2585%25A5";
    } else {
        return "%25E6%2595%2599%25E8%2582%25B2%25E7%25BD%2591%25E6%258E%25A5%25E5%2585%25A5";
    }
}

string getLoginPageUrlFromHtmlCode(const string &htmlCode) {
    auto leftPosition = htmlCode.find_first_of('\'');
    auto rightPosition = htmlCode.find_last_of('\'');
    return htmlCode.substr(leftPosition + 1, rightPosition - leftPosition - 1);
}

string getLoginUrlFromLoginPageUrl(const string &loginPageUrl) {
    string baseUrl = loginPageUrl.substr(0, loginPageUrl.find_first_of('?'));
    baseUrl.replace(baseUrl.find("index.jsp"), strlen("index.jsp"), "InterFace.do?method=login");
    return baseUrl;
}

string replaceAll(string str, const string &from, const string &to) {
    size_t start_pos = 0;
    while ((start_pos = str.find(from, start_pos)) != string::npos) {
        str.replace(start_pos, from.length(), to);
        start_pos += to.length();
    }
    return str;
}

string getQueryStringFromHtmlCode(string &htmlCode) {
    auto leftPosition = htmlCode.find_first_of('?');
    auto rightPosition = htmlCode.find_last_of('\'');
    string queryString = htmlCode.substr(leftPosition + 1, rightPosition - leftPosition - 1);
    queryString = replaceAll(queryString, "&", "%2526");
    queryString = replaceAll(queryString, "=", "%253D");
    return queryString;
}

string login(const char *loginUrl, const char *username, const char *password, const char *serviceString,
             const char *queryString) {
    CURL *curl;
    CURLcode res;
    string loginResult;
    struct curl_slist *chunk = nullptr;
    chunk = curl_slist_append(chunk,
                              "Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8");
    chunk = curl_slist_append(chunk, "Content-Type: application/x-www-form-urlencoded; charset=UTF-8");
    chunk = curl_slist_append(chunk, "Connection: keep-alive");
    chunk = curl_slist_append(chunk,
                              "User-Agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.132 Safari/537.36");
    // 101 is the length of "userId=&password=&service=&queryString=&operatorPwd=&operatorUserId=&validcode=&passwordEncrypt=false"
    char *postFields = (char *) malloc(sizeof(char) *
                                       (101 + strlen(username) + strlen(password) + strlen(serviceString) +
                                        strlen(queryString) + 1));
    sprintf(postFields,
            "userId=%s&password=%s&service=%s&queryString=%s&operatorPwd=&operatorUserId=&validcode=&passwordEncrypt=false",
            username, password, serviceString, queryString);
    curl = curl_easy_init();
    if (curl) {
        curl_easy_setopt(curl, CURLOPT_URL, loginUrl);
        curl_easy_setopt(curl, CURLOPT_POST, 1);
        curl_easy_setopt(curl, CURLOPT_POSTFIELDS, postFields);
        curl_easy_setopt(curl, CURLOPT_HTTPHEADER, chunk);
        curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, writeCallback);
        curl_easy_setopt(curl, CURLOPT_WRITEDATA, &loginResult);
        res = curl_easy_perform(curl);
        if (res != CURLE_OK) {
            cout << "Can not send login request, exit." << endl;
            exit(EXIT_FAILURE);
        }
        curl_easy_cleanup(curl);
    }
    delete postFields;
    return loginResult;
}

string getUserIndex() {
    char *redirectUrl = nullptr;
    CURL *curl;
    CURLcode res;
    curl = curl_easy_init();
    if (curl) {
        curl_easy_setopt(curl, CURLOPT_URL, "http://10.8.2.2/eportal/redirectortosuccess.jsp");
        curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, noPrintWriteCallback);
        res = curl_easy_perform(curl);
        if (res != CURLE_OK) {
            cout << "Can not send get user index request, exit." << endl;
            exit(EXIT_FAILURE);
        }
        curl_easy_getinfo(curl, CURLINFO_REDIRECT_URL, &redirectUrl);
        curl_easy_cleanup(curl);
    }
    string redirectUrlString = redirectUrl;
    int tempIndex = redirectUrlString.find("userIndex=");
    redirectUrlString = redirectUrlString.substr(tempIndex + strlen("userIndex="),
                                                 redirectUrlString.length() - tempIndex);
    return redirectUrlString;
}

string logout(const string &userIndex) {
    CURL *curl;
    CURLcode res;
    string logoutResult;
    string postFieldsString = "userIndex=";
    postFieldsString += userIndex;
    char *postFields = (char *) malloc(sizeof(char) * (postFieldsString.length() + 1));
    strcpy(postFields, postFieldsString.c_str());
    curl = curl_easy_init();
    if (curl) {
        curl_easy_setopt(curl, CURLOPT_URL, "http://10.8.2.2/eportal/InterFace.do?method=logout");
        curl_easy_setopt(curl, CURLOPT_POST, 1);
        curl_easy_setopt(curl, CURLOPT_POSTFIELDS, postFields);
        curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, writeCallback);
        curl_easy_setopt(curl, CURLOPT_WRITEDATA, &logoutResult);
        res = curl_easy_perform(curl);
        if (res != CURLE_OK) {
            cout << "Can not send logout request, exit." << endl;
            exit(EXIT_FAILURE);
        }
        curl_easy_cleanup(curl);
    }
    delete postFields;
    return logoutResult;
}

int main(int argc, char **argv) {
    if (argc < 3) {
        if (argc == 2 && strcmp(argv[1], "logout") == 0) {
            string userIndex = getUserIndex();
            string logoutResult = logout(userIndex);
            cout << logoutResult << endl;
        } else {
            // Show usage when argc < 3 and no logout
            printUsage();
        }
        return 0;
    }
    // Check network status
    long captiveServerStatusCode = getResponseStatusCode(captiveServerUrl);
    if (captiveServerStatusCode == -1) {
        cout << "Can not initialize curl" << endl;
        return EXIT_FAILURE;
    } else if (captiveServerStatusCode == 204) {
        // Exit when user is already online
        cout << "You are already online!" << endl;
        return EXIT_SUCCESS;
    }
    // Start ruijie login
    char *serviceName = argv[1];
    char *username = argv[2];
    char *password = argv[3];
    string captiveServerResponseContent = getResponseContent(captiveServerUrl);
    string loginPageUrl = getLoginPageUrlFromHtmlCode(captiveServerResponseContent);
    string queryString = getQueryStringFromHtmlCode(captiveServerResponseContent);
    string loginUrl = getLoginUrlFromLoginPageUrl(loginPageUrl);
    string serviceString = getServiceStringByServiceName(serviceName);
    string loginResult = login(loginUrl.c_str(), username, password, serviceString.c_str(), queryString.c_str());
    cout << loginResult << endl;
    return EXIT_SUCCESS;
}