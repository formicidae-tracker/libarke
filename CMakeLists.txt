cmake_policy(SET CMP0048 NEW)

project(libarke LANGUAGES C CXX
	            VERSION 0.1)

cmake_minimum_required(VERSION 3.2)


set(CMAKE_C_STANDARD 99)
set(CMAKE_CXX_STANDARD 11)


include_directories(${PROJECT_SOURCE_DIR}/include)

enable_testing()

set(CMAKE_CTEST_COMMAND ctest -V)
add_custom_target(check COMMAND  ${CMAKE_CTEST_COMMAND})

find_path(GTEST_SRC_ROOT googletest/src/gtest-all.cc
	PATHS /usr/src/googletest)

if(NOT GTEST_SRC_ROOT )
	message(FATAL_ERROR "Could not find gtest")
endif(NOT GTEST_SRC_ROOT )

find_package(Threads REQUIRED)

include_directories(${GTEST_SRC_ROOT}/googletest)
set(GTEST_SRC_FILE ${GTEST_SRC_ROOT}/googletest/src/gtest-all.cc)

add_library(gtest ${GTEST_SRC_FILE})

add_subdirectory(src-c)
