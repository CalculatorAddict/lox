build:
	clang -std=c17 -Wall -Wextra -g -fsanitize=address,undefined *.c -o main

run:
	clang -std=c17 *.c -g -o main && ./main

clean:
	rm -f main
	rm -rf *.dSYM

.PHONY: build clean run