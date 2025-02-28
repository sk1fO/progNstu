#include <iostream>
#include <sstream>
#include "FlightSchedule.h"

using namespace std;

void displayMenu() {
    cout << "Система управления расписанием полетов" << endl;
    cout << "1. Создать самолет" << endl;
    cout << "2. Получить самолеты для города" << endl;
    cout << "3. Получить города для самолета" << endl;
    cout << "4. Показать все самолеты" << endl;
    cout << "5. Выход" << endl;
    cout << "Введите ваш выбор: ";
}

int main() {
    FlightSchedule schedule;
    int choice;
    string input;

    while (true) {
        displayMenu();
        cin >> choice;
        cin.ignore(); // Игнорируем оставшийся символ новой строки

        switch (choice) {
            case 1: {
                cout << "Введите название самолета и города через пробел (например, Самолет1 Город1 Город2 Город3): ";
                getline(cin, input);
                schedule.processCommand(Type::CREATE_PLANE, input);
                break;
            }
            case 2: {
                cout << "Введите название города: ";
                getline(cin, input);
                schedule.processCommand(Type::PLANES_FOR_TOWN, input);
                break;
            }
            case 3: {
                cout << "Введите название самолета: ";
                getline(cin, input);
                schedule.processCommand(Type::TOWNS_FOR_PLANE, input);
                break;
            }
            case 4: {
                schedule.processCommand(Type::PLANES, "");
                break;
            }
            case 5: {
                cout << "Выход из программы..." << endl;
                return 0;
            }
            default: {
                cout << "Неверный выбор. Пожалуйста, попробуйте снова." << endl;
                break;
            }
        }
    }

    return 0;
}