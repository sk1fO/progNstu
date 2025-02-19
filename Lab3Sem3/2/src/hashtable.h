#ifndef HASHTABLE_H
#define HASHTABLE_H

#include <string>
#include <iostream>
#include <fstream>
#include <cstring>

// Узел односвязного списка
struct Node {
    std::string key;
    std::string value;
    Node* next;

    Node(const std::string& k, const std::string& v) : key(k), value(v), next(nullptr) {}
};

// Класс хеш-таблицы
class HashTable {
private:
    int TABLE_SIZE = 10; // Начальный размер таблицы
    Node** table; // Массив указателей на узлы
    int count; // Количество элементов в таблице

    // Хеш-функция
    int hashFunction(const std::string& key) const {
        int hash = 0;
        for (char ch : key) {
            hash += ch;
        }
        return hash % TABLE_SIZE;
    }

    // Увеличение размера таблицы и рехэширование
    void resize() {
        int newSize = TABLE_SIZE * 2;
        Node** newTable = new Node*[newSize];
        for (int i = 0; i < newSize; ++i) {
            newTable[i] = nullptr;
        }

        for (int i = 0; i < TABLE_SIZE; ++i) {
            Node* current = table[i];
            while (current != nullptr) {
                Node* next = current->next;
                int newIndex = hashFunction(current->key);
                current->next = newTable[newIndex];
                newTable[newIndex] = current;
                current = next;
            }
        }

        delete[] table;
        table = newTable;
        TABLE_SIZE = newSize;
    }

public:
    // Конструктор
    HashTable() : count(0) {
        table = new Node*[TABLE_SIZE];
        for (int i = 0; i < TABLE_SIZE; ++i) {
            table[i] = nullptr;
        }
    }

    // Деструктор
    ~HashTable() {
        for (int i = 0; i < TABLE_SIZE; ++i) {
            Node* current = table[i];
            while (current != nullptr) {
                Node* temp = current;
                current = current->next;
                delete temp;
            }
        }
        delete[] table;
    }

    // Вставка элемента
    void insert(const std::string& key, const std::string& value) {
        if ((double)count / TABLE_SIZE > 0.65) {
            resize();
        }
        int index = hashFunction(key);
        Node* current = table[index];

        // Проверяем, есть ли уже такой ключ
        while (current != nullptr) {
            if (current->key == key) {
                current->value = value; // Обновляем значение
                return;
            }
            current = current->next;
        }

        // Если ключа нет, добавляем новый узел
        Node* newNode = new Node(key, value);
        newNode->next = table[index];
        table[index] = newNode;
        count++;
    }

    // Поиск элемента
    std::string search(const std::string& key) const {
        int index = hashFunction(key);
        Node* current = table[index];
        while (current != nullptr) {
            if (current->key == key) {
                return current->value;
            }
            current = current->next;
        }
        return "not exists";
    }

    // Удаление элемента
    void remove(const std::string& key) {
        int index = hashFunction(key);
        Node* current = table[index];
        Node* prev = nullptr;
        while (current != nullptr) {
            if (current->key == key) {
                if (prev == nullptr) {
                    table[index] = current->next;
                } else {
                    prev->next = current->next;
                }
                delete current;
                count--;
                return;
            }
            prev = current;
            current = current->next;
        }
    }

    // Вывод таблицы
    void printTable() const {
        for (int i = 0; i < TABLE_SIZE; ++i) {
            Node* current = table[i];
            while (current != nullptr) {
                std::cout << "[" << current->key << ":" << current->value << "] ";
                current = current->next;
            }
        }
        std::cout << std::endl;
    }

    // Перегрузка оператора []
    std::string& operator[](const std::string& key) {
        int index = hashFunction(key);
        Node* current = table[index];
        while (current != nullptr) {
            if (current->key == key) {
                return current->value;
            }
            current = current->next;
        }
        // Если ключ не найден, вставляем новый узел
        insert(key, "");
        return table[hashFunction(key)]->value;
    }

    // Сериализация в бинарный формат
    void serializeToFile(const std::string& filename) const {
        std::ofstream file(filename, std::ios::binary);
        if (!file.is_open()) {
            throw std::runtime_error("Не удалось открыть файл: " + filename);
        }

        // Записываем количество элементов
        file.write(reinterpret_cast<const char*>(&count), sizeof(count));

        // Записываем каждый элемент
        for (int i = 0; i < TABLE_SIZE; ++i) {
            Node* current = table[i];
            while (current != nullptr) {
                size_t keySize = current->key.size();
                size_t valueSize = current->value.size();

                file.write(reinterpret_cast<const char*>(&keySize), sizeof(keySize));
                file.write(current->key.c_str(), keySize);

                file.write(reinterpret_cast<const char*>(&valueSize), sizeof(valueSize));
                file.write(current->value.c_str(), valueSize);

                current = current->next;
            }
        }
    }

    // Десериализация из бинарного формата
    void deserializeFromFile(const std::string& filename) {
        std::ifstream file(filename, std::ios::binary);
        if (!file.is_open()) {
            throw std::runtime_error("Не удалось открыть файл: " + filename);
        }

        // Очищаем существующие данные
        for (int i = 0; i < TABLE_SIZE; ++i) {
            Node* current = table[i];
            while (current != nullptr) {
                Node* temp = current;
                current = current->next;
                delete temp;
            }
            table[i] = nullptr;
        }
        count = 0;

        // Читаем количество элементов
        file.read(reinterpret_cast<char*>(&count), sizeof(count));

        // Читаем каждый элемент
        for (int i = 0; i < count; ++i) {
            size_t keySize, valueSize;
            file.read(reinterpret_cast<char*>(&keySize), sizeof(keySize));
            if (!file) {
                throw std::runtime_error("Не удалось прочитать размер ключа");
            }

            std::string key(keySize, '\0');
            file.read(&key[0], keySize);
            if (!file) {
                throw std::runtime_error("Не удалось прочитать ключ");
            }

            file.read(reinterpret_cast<char*>(&valueSize), sizeof(valueSize));
            if (!file) {
                throw std::runtime_error("Не удалось прочитать размер значения");
            }

            std::string value(valueSize, '\0');
            file.read(&value[0], valueSize);
            if (!file) {
                throw std::runtime_error("Не удалось прочитать значение");
            }

            insert(key, value);
        }
    }
};

#endif // HASHTABLE_H