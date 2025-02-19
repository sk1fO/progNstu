#include <iostream>
#include <map>
#include <vector>
#include <string>
#include <algorithm>

struct Cell {
    std::string product;
    int quantity;
};

std::map<std::string, Cell> warehouse;
int totalCapacity = 50400; // Вместимость для варианта 8

bool isValidAddress(const std::string& address) {
    // Проверка длины адреса (6 символов)
    if (address.length() != 6) {
        return false;
    }

    // Проверка зоны хранения (A, Б, В)
    char zone = address[0];
    if (zone != 'A' && zone != 'Б' && zone != 'В') {
        return false;
    }

    // Проверка номера стеллажа (01-14)
    std::string shelfStr = address.substr(1, 2);
    try {
        int shelf = std::stoi(shelfStr);
        if (shelf < 1 || shelf > 14) {
            return false;
        }
    } catch (const std::invalid_argument&) {
        return false; // Некорректное число
    }

    // Проверка номера вертикальной секции (01-06)
    std::string sectionStr = address.substr(3, 2);
    try {
        int section = std::stoi(sectionStr);
        if (section < 1 || section > 6) {
            return false;
        }
    } catch (const std::invalid_argument&) {
        return false; // Некорректное число
    }

    // Проверка номера полки (01-20)
    std::string levelStr = address.substr(5, 2);
    try {
        int level = std::stoi(levelStr);
        if (level < 1 || level > 20) {
            return false;
        }
    } catch (const std::invalid_argument&) {
        return false; // Некорректное число
    }

    return true;
}

void addProduct(const std::string& product, int quantity, const std::string& address) {
    if (!isValidAddress(address)) {
        std::cout << "Ошибка: некорректный адрес ячейки " << address << std::endl;
        return;
    }

    if (warehouse.find(address) != warehouse.end()) {
        if (warehouse[address].quantity + quantity > 10) {
            std::cout << "Ошибка: превышена вместимость ячейки " << address << std::endl;
            return;
        }
        warehouse[address].quantity += quantity;
    } else {
        if (quantity > 10) {
            std::cout << "Ошибка: превышена вместимость ячейки " << address << std::endl;
            return;
        }
        warehouse[address] = {product, quantity};
    }
    std::cout << "Товар добавлен в ячейку " << address << std::endl;
}

void removeProduct(const std::string& product, int quantity, const std::string& address) {
    if (!isValidAddress(address)) {
        std::cout << "Ошибка: некорректный адрес ячейки " << address << std::endl;
        return;
    }

    if (warehouse.find(address) == warehouse.end() || warehouse[address].product != product) {
        std::cout << "Ошибка: товар не найден в ячейке " << address << std::endl;
        return;
    }
    if (warehouse[address].quantity < quantity) {
        std::cout << "Ошибка: недостаточно товара в ячейке " << address << std::endl;
        return;
    }
    warehouse[address].quantity -= quantity;
    if (warehouse[address].quantity == 0) {
        warehouse.erase(address);
    }
    std::cout << "Товар удален из ячейки " << address << std::endl;
}

void printWarehouseInfo() {
    int totalUsed = 0;
    std::map<char, int> zoneUsage;
    std::vector<std::string> emptyCells;

    for (int z = 0; z < 3; ++z) {
        char zone = 'A' + z;
        zoneUsage[zone] = 0;
        for (int s = 1; s <= 14; ++s) {
            for (int v = 1; v <= 6; ++v) {
                for (int p = 1; p <= 20; ++p) {
                    std::string address = std::string(1, zone) + (s < 10 ? "0" : "") + std::to_string(s) + (v < 10 ? "0" : "") + std::to_string(v) + (p < 10 ? "0" : "") + std::to_string(p);
                    if (warehouse.find(address) != warehouse.end()) {
                        totalUsed += warehouse[address].quantity;
                        zoneUsage[zone] += warehouse[address].quantity;
                    } else {
                        emptyCells.push_back(address);
                    }
                }
            }
        }
    }

    std::cout << "Общая загруженность склада: " << (totalUsed * 100 / totalCapacity) << "%" << std::endl;
    for (auto& zone : zoneUsage) {
        std::cout << "Загруженность зоны " << zone.first << ": " << (zone.second * 100 / (totalCapacity / 3)) << "%" << std::endl;
    }

    std::cout << "Содержимое ячеек:" << std::endl;
    for (auto& cell : warehouse) {
        std::cout << cell.first << ": " << cell.second.product << " - " << cell.second.quantity << " единиц" << std::endl;
    }

    std::cout << "Пустые ячейки:" << std::endl;
    for (auto& cell : emptyCells) {
        std::cout << cell << std::endl;
    }
}

int main() {
    std::string command;
    while (true) {
        std::cout << "Введите команду: ";
        std::getline(std::cin, command);
        if (command.substr(0, 3) == "ADD") {
            std::string product;
            int quantity;
            std::string address;
            size_t pos = command.find(' ', 4);
            product = command.substr(4, pos - 4);
            size_t pos2 = command.find(' ', pos + 1);
            quantity = std::stoi(command.substr(pos + 1, pos2 - pos - 1));
            address = command.substr(pos2 + 1);
            addProduct(product, quantity, address);
        } else if (command.substr(0, 6) == "REMOVE") {
            std::string product;
            int quantity;
            std::string address;
            size_t pos = command.find(' ', 7);
            product = command.substr(7, pos - 7);
            size_t pos2 = command.find(' ', pos + 1);
            quantity = std::stoi(command.substr(pos + 1, pos2 - pos - 1));
            address = command.substr(pos2 + 1);
            removeProduct(product, quantity, address);
        } else if (command == "INFO") {
            printWarehouseInfo();
        } else if (command == "EXIT") {
            break;
        } else {
            std::cout << "Неизвестная команда" << std::endl;
        }
    }
    return 0;
}