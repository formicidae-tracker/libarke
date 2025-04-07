cmake_policy(SET CMP0048 NEW)
cmake_policy(SET CMP0076 NEW)

project(
	libarke
	LANGUAGES C CXX
	VERSION 0.1
)

cmake_minimum_required(VERSION 3.13..3.23)

set(CMAKE_C_STANDARD 99)
set(CMAKE_CXX_STANDARD 17)

if(NOT CMAKE_CURRENT_SOURCE_DIR STREQUAL CMAKE_SOURCE_DIR)
	set(ARKE_IS_IMPORTED 1)
endif(NOT CMAKE_CURRENT_SOURCE_DIR STREQUAL CMAKE_SOURCE_DIR)

include_directories(${PROJECT_SOURCE_DIR}/include)

if(NOT ARKE_IS_IMPORTED)
	enable_testing()

	set(CMAKE_CTEST_COMMAND ctest -V)
	add_custom_target(
		check COMMAND ${CMAKE_CTEST_COMMAND} ARGS --output-on-failure
	)

	include(FetchContent)

	FetchContent_Declare(
		googletest
		GIT_REPOSITORY https://github.com/google/googletest.git
		GIT_TAG 9bb354fa8325fa31faa1e12627b25ab445e6eed3
	)
	# For Windows: Prevent overriding the parent project's compiler/linker
	# settings
	set(gtest_force_shared_crt
		ON
		CACHE BOOL "" FORCE
	)

	FetchContent_MakeAvailable(googletest)

	include(GoogleTest)

endif(NOT ARKE_IS_IMPORTED)

add_subdirectory(src-c)

target_include_directories(arke INTERFACE ${CMAKE_CURRENT_SOURCE_DIR}/include)
