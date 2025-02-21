CC = g++
CFLAGS = -std=c++17 -Wall -Wextra -g --coverage
LDFLAGS = -lgtest -lgtest_main -pthread
INCLUDES = -I/usr/include/
LIBS = -L/usr/lib

all: main boost_test google_test

main: src/main.o
	$(CC) $(CFLAGS) $^ -o $@ $(LDFLAGS) $(LIBS)

boost_test: tests/boost_test.o
	$(CC) $(CFLAGS) $^ -o $@ $(LDFLAGS) $(LIBS) -lboost_unit_test_framework

google_test: tests/google_test.o
	$(CC) $(CFLAGS) $^ -o $@ $(LDFLAGS) $(LIBS)

src/main.o: src/main.cpp
	$(CC) $(CFLAGS) -c $< -o $@

tests/boost_test.o: tests/boost_test.cpp
	$(CC) $(CFLAGS) -c $< -o $@

tests/google_test.o: tests/google_test.cpp
	$(CC) $(CFLAGS) -c $< -o $@

clean:
	rm -f main boost_test google_test benchmark *.o

test: all
	./boost_test
	./google_test

coverage: all
	./boost_test
	./google_test
	gcovr -r . --html --html-details -o coverage.html
	open coverage.html

.PHONY: clean test coverage