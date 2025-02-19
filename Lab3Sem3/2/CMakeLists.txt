cmake_minimum_required(VERSION 3.10)
project(HashTableProject)

# Устанавливаем пути к библиотекам
set(BOOST_ROOT "/usr/include/boost")
set(GTEST_ROOT "/usr/include/gtest")

# Добавляем поддержку C++17
set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED True)

# Добавляем исходные файлы
add_executable(main src/main.cpp src/hashtable.h)

# Добавляем тесты
enable_testing()

# Boost Test
add_executable(boost_test tests/boost_test.cpp src/hashtable.h)
target_link_libraries(boost_test Boost::unit_test_framework)

# Google Test
add_executable(google_test tests/google_test.cpp src/hashtable.h)
target_link_libraries(google_test gtest gtest_main)

# Benchmark
add_executable(benchmark tests/benchmark.cpp src/hashtable.h)
target_link_libraries(benchmark benchmark_main)

# Настройка покрытия кода
if(CMAKE_CXX_COMPILER_ID STREQUAL "GNU" OR CMAKE_CXX_COMPILER_ID MATCHES "Clang")
    set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} --coverage")
endif()

# Добавляем тесты
add_test(NAME boost_test COMMAND boost_test)
add_test(NAME google_test COMMAND google_test)
add_test(NAME benchmark COMMAND benchmark)