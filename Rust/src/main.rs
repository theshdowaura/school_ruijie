// main.rs
use clap::{Arg, Command};
use reqwest::blocking::Client;
use reqwest::header::{HeaderMap, HeaderValue, ACCEPT, CONTENT_TYPE, USER_AGENT};
use std::{thread, time::Duration};
use std::str::Split;
use std::process::Command as PingCommand;

const CAPTIVE_SERVER_URL: &str = "http://www.google.cn/generate_204";
const PING_HOST: &str = "180.101.50.188";

fn get_captive_server_response() -> Result<(u16, String), Box<dyn std::error::Error>> {
    let response = reqwest::blocking::get(CAPTIVE_SERVER_URL)?;
    let status_code = response.status().as_u16();
    let body = response.text()?;
    Ok((status_code, body))
}

fn get_login_url_from_html_code(html_code: &str) -> (String, String) {
    let parts: Split<char> = html_code.split('?');
    let login_page_url = parts.clone().next().unwrap_or("");
    let login_url = login_page_url.replace("index.jsp", "InterFace.do?method=login");
    let query_string = parts
        .skip(1)
        .next()
        .unwrap_or_default()
        .replace('&', "%2526")
        .replace('=', "%253D");
    (login_url, query_string)
}

fn login(
    login_url: &str,
    username: &str,
    password: &str,
    query_string: &str,
    services_password: &str,
) -> Result<String, Box<dyn std::error::Error>> {
    let client = Client::new();
    let login_post_data = format!(
        "userId={}&password={}&service=&queryString={}&operatorPwd={}&operatorUserId=&validcode=&passwordEncrypt=false",
        username, password, query_string, services_password
    );

    let mut headers = HeaderMap::new();
    headers.insert(ACCEPT, HeaderValue::from_static("text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8"));
    headers.insert(CONTENT_TYPE, HeaderValue::from_static("application/x-www-form-urlencoded; charset=UTF-8"));
    headers.insert(USER_AGENT, HeaderValue::from_static("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.132 Safari/537.36"));

    let response = client
        .post(login_url)
        .headers(headers)
        .body(login_post_data)
        .send()?;

    let body = response.text()?;
    Ok(body)
}

fn start_ping_check(username: &str, password: &str, services_password: &str) {
    loop {
        let output = PingCommand::new("ping")
            .arg("-c 1")
            .arg(PING_HOST)
            .output();

        match output {
            Ok(result) if !result.status.success() => {
                println!("Ping failed, packet loss detected. Re-authenticating...");
                if let Err(err) = re_authenticate(username, password, services_password) {
                    println!("Re-authentication failed: {}", err);
                } else {
                    println!("Re-authentication successful!");
                }
            }
            Ok(_) => println!("Network is stable."),
            Err(err) => println!("Ping failed: {}", err),
        }
        thread::sleep(Duration::from_secs(10));
    }
}

fn re_authenticate(username: &str, password: &str, services_password: &str) -> Result<(), Box<dyn std::error::Error>> {
    let (status_code, response_body) = get_captive_server_response()?;

    if status_code == 204 {
        println!("You are already online!");
        return Ok(());
    }

    let (login_url, query_string) = get_login_url_from_html_code(&response_body);
    let login_result = login(&login_url, username, password, &query_string, services_password)?;
    println!("{}", login_result);
    Ok(())
}

fn main() {
    let matches = Command::new("Network Login")
        .arg(Arg::new("username").short('u').num_args(1).required(true))
        .arg(Arg::new("password").short('p').num_args(1).required(true))
        .arg(Arg::new("services_password").short('c').num_args(1).required(true))
        .arg(Arg::new("environment").short('e').num_args(1))
        .get_matches();

    let username = matches.get_one::<String>("username").unwrap();
    let password = matches.get_one::<String>("password").unwrap();
    let services_password = matches.get_one::<String>("services_password").unwrap();
    let environment = matches.get_one::<String>("environment").map(String::as_str).unwrap_or("");

    // Initial authentication
    if let Err(err) = re_authenticate(username, password, services_password) {
        println!("Initial authentication failed: {}", err);
        return;
    }

    // Start Ping check if environment is "on"
    if environment == "on" {
        start_ping_check(username, password, services_password);
    }
}
