#!/bin/bash

# Сборка проекта с покрытием
make clean
make

# Запуск тестов
./boost_test
./google_test

# Генерация отчета о покрытии
gcovr -r . --html --html-details -o coverage.html

# Открытие отчета
open coverage.html