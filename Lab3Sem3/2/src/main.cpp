#include <iostream>
#include <fstream>
#include <sstream>
#include <stdexcept>
#include "hashtable.h"

// Функция для разбора аргументов командной строки
void parseCommandLine(int argc, char** argv, std::string& file, std::string& query) {
    if (argc != 5) {
        throw std::runtime_error("Использование: ./main --file <имя_файла> --query '<команда>'");
    }

    for (int i = 1; i < argc; i += 2) {
        if (std::string(argv[i]) == "--file") {
            file = argv[i + 1];
        } else if (std::string(argv[i]) == "--query") {
            query = argv[i + 1];
        } else {
            throw std::runtime_error("Неизвестная опция: " + std::string(argv[i]));
        }
    }
}

// Функция для выполнения команды
void executeCommand(HashTable& ht, const std::string& command) {
    std::istringstream iss(command);
    std::string cmd, key, value;

    iss >> cmd >> key >> value;

    if (cmd == "HSET") {
        ht.insert(key, value);
    } else if (cmd == "HDEL") {
        ht.remove(key);
    } else if (cmd == "HGET") {
        std::string result = ht.search(key);
        std::cout << "Результат: " << result << std::endl;
    } else if (cmd == "PRINT") {
        ht.printTable();
    } else {
        throw std::runtime_error("Неизвестная команда: " + cmd);
    }
}

// Функция для загрузки данных из файла
void loadFromFile(HashTable& ht, const std::string& filename) {
    try {
        ht.deserializeFromFile(filename);
    } catch (const std::exception& e) {
        ht.serializeToFile(filename);
    }
}

// Функция для сохранения данных в файл
void saveToFile(const HashTable& ht, const std::string& filename) {
    ht.serializeToFile(filename);
}

int main(int argc, char** argv) {
    try {
        std::string file, query;
        parseCommandLine(argc, argv, file, query);

        HashTable ht;

        // Загрузка данных из файла
        loadFromFile(ht, file);

        // Выполнение команды
        executeCommand(ht, query);

        // Сохранение данных в файл
        saveToFile(ht, file);

    } catch (const std::exception& e) {
        std::cerr << "Ошибка: " << e.what() << std::endl;
        return 1;
    }

    return 0;
}