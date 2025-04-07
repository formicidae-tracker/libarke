add_library(arke INTERFACE)
target_sources(arke INTERFACE arke.c)

if(NOT ARKE_IS_IMPORTED)
	set(TEST_SRC_FILES arke-private-conversion-test.cpp)

	set(TEST_HDR_FILES)

	add_executable(
		libarke-tests EXCLUDE_FROM_ALL ${TEST_HDR_FILES} ${TEST_SRC_FILES}
	)

	target_link_libraries(libarke-tests GTest::gtest_main arke)

	gtest_discover_tests(libarke-tests)

	add_test(libarke-tests libarke-tests)
	add_dependencies(check libarke-tests)
endif(NOT ARKE_IS_IMPORTED)
