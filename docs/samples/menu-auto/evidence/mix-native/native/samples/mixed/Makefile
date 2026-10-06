.PHONY: owned-build owned-test owned-fail
owned-build: build/c build/cpp
build/c: main.c
	mkdir -p build
	gcc -Wall -Wextra -Werror -pedantic -std=c11 main.c -o build/c
build/cpp: main.cpp
	mkdir -p build
	g++ -Wall -Wextra -Werror -pedantic -std=c++17 main.cpp -o build/cpp
owned-test: owned-build
	test "$$(./build/c)" = "21"
	test "$$(./build/cpp)" = "42"
	@printf "REAL_MIXED_TEST_PASS\n"
owned-fail:
	@printf "OWNED_RECIPE_EXIT_7\n"
	@exit 7
