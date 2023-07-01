cmake_minimum_required(VERSION 3.14)
project(ruijie)

set(CMAKE_CXX_STANDARD 14)

set(CURL_LIBRARY "-lcurl")
find_package(CURL REQUIRED)

add_executable(ruijie main.cpp)

include_directories(${CURL_INCLUDE_DIR})
target_link_libraries(ruijie ${CURL_LIBRARIES})