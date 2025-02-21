#include <iostream>
#include <map>
using namespace std;

// Структура для хранения информации о ячейке
struct Cell {
    string product; // Название товара
    int quantity; // Количество товара
};

// Класс для управления складом
class Warehouse {
private:
    map<string, Cell> cells; // Контейнер для хранения ячеек склада
    int totalCapacity; // Общая вместимость склада
    int currentLoad; // Текущая загрузка склада

public:
    Warehouse() : totalCapacity(1920), currentLoad(0) {
        // Инициализация ячеек
        char zones[] = {'A', 'B'}; // Зоны склада
        for (char zone : zones) { // Перебор зон
            for (int shelf = 1; shelf <= 6; ++shelf) { // Перебор полок
                for (int section = 1; section <= 4; ++section) { // Перебор секций
                    for (int row = 1; row <= 4; ++row) { // Перебор рядов
                        string address = zone + to_string(shelf) + to_string(section) + to_string(row); // Формирование адреса ячейки
                        cells[address] = {"Пусто", 0}; // Инициализация ячейки как пустой
                    }
                }
            }
        }
    }

    void addProduct(const string& product, int quantity, const string& address) {
        if (cells.find(address) == cells.end()) { // Проверка наличия ячейки по адресу
            cout << "Адрес ячейки не найден." << endl;
            return;
        }
        Cell& cell = cells[address]; // Получение ссылки на ячейку
        if (cell.product != "Пусто" && cell.product != product) { // Проверка занятости ячейки
            cout << "Ячейка уже занята другим товаром." << endl;
            return;
        }
        if (cell.quantity + quantity > 10) { // Проверка вместимости ячейки
            cout << "Недостаточно места в ячейке." << endl;
            return;
        }
        cell.product = product; // Установка товара в ячейку
        cell.quantity += quantity; // Увеличение количества товара
        currentLoad += quantity; // Увеличение текущей загрузки склада
        cout << "Товар добавлен в ячейку " << address << "." << endl;
    }

    void removeProduct(const string& product, int quantity, const string& address) {
        if (cells.find(address) == cells.end()) { // Проверка наличия ячейки по адресу
            cout << "Адрес ячейки не найден." << endl;
            return;
        }
        Cell& cell = cells[address]; // Получение ссылки на ячейку
        if (cell.product != product) { // Проверка соответствия товара
            cout << "Ячейка содержит другой товар." << endl;
            return;
        }
        if (cell.quantity < quantity) { // Проверка наличия товара в ячейке
            cout << "Недостаточно товара в ячейке для удаления." << endl;
            return;
        }
        cell.quantity -= quantity; // Уменьшение количества товара
        if (cell.quantity == 0) { // Проверка на пустую ячейку
            cell.product = "Пусто";
        }
        currentLoad -= quantity; // Уменьшение текущей загрузки склада
        cout << "Товар удален из ячейки " << address << "." << endl;
    }

    void info() {
        cout << "Общая загрузка склада: " << (currentLoad * 100 / totalCapacity) << "%" << endl; // Вывод загрузки склада
        for (const auto& [address, cell] : cells) { // Перебор ячеек
            if (cell.product != "Пусто") { // Проверка на пустую ячейку
                cout << "Ячейка " << address << ": " << cell.product << " (" << cell.quantity << " единиц)" << endl; // Вывод информации о ячейке
            }
        }
        cout << "Пустые ячейки: ";
        for (const auto& [address, cell] : cells) { // Перебор ячеек
            if (cell.product == "Пусто") { // Проверка на пустую ячейку
                cout << address << " "; // Вывод адреса пустой ячейки
            }
        }
        cout << endl;
    }
};

int main() {

    Warehouse warehouse; // Создание объекта склада
    string command, product, address; // Объявление переменных для команд, товара и адреса
    int quantity; // Объявление переменной для количества товара

    while (true) { // цикл для ввода команд
        cout << "Введите команду (add, remove, info, end): ";
        cin >> command; // Ввод команды

        if (command == "end") { // Проверка на команду завершения
            break;
        }

        if (command == "add") { // Проверка на команду добавления товара
            cout << "Введите название товара, количество и адрес ячейки: ";
            cin >> product >> quantity >> address; // Ввод данных о товаре
            warehouse.addProduct(product, quantity, address); // Добавление товара
        } else if (command == "remove") { // Проверка на команду удаления товара
            cout << "Введите название товара, количество и адрес ячейки: ";
            cin >> product >> quantity >> address; // Ввод данных о товаре
            warehouse.removeProduct(product, quantity, address); // Удаление товара
        } else if (command == "info") { // Проверка на команду вывода информации
            warehouse.info(); // Вывод информации о складе
        } else {
            cout << "Неизвестная команда." << endl; // Вывод сообщения о неизвестной команде
        }
    }

    return 0;
}