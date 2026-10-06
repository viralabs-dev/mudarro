.PHONY: build test fail
build: build/sample
build/sample: main.c
	mkdir -p build
	gcc -Wall -Wextra -Werror -pedantic -std=c11 main.c -o build/sample
test: build
	test "$$(./build/sample)" = "21"
	@printf "REAL_C_TEST_PASS\n"
fail:
	@printf "OWNED_RECIPE_EXIT_7\n"
	@exit 7
