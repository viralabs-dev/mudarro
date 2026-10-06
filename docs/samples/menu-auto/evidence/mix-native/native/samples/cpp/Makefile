.PHONY: build test fail
build: build/sample
build/sample: main.cpp
	mkdir -p build
	g++ -Wall -Wextra -Werror -pedantic -std=c++17 main.cpp -o build/sample
test: build
	test "$$(./build/sample)" = "42"
	@printf "REAL_CPP_TEST_PASS\n"
fail:
	@printf "OWNED_RECIPE_EXIT_7\n"
	@exit 7
