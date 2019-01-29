add_library(arke STATIC arke-src-all.c)


set(TEST_SRC_FILES main-check.cpp arke-private-conversion-test.cpp)

set(TEST_HDR_FILES)

add_executable(libarke-tests EXCLUDE_FROM_ALL ${TEST_SRC_FILES} ${TEST_HDR_FILES})

target_link_libraries(libarke-tests gtest Threads::Threads)

add_test(libarke-tests libarke-tests)
add_dependencies(check libarke-tests)
