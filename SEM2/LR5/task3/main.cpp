#include <iostream> 
#include <sstream> 
#include "FlightSchedule.h" 

using namespace std; 

int main() { // Точка входа в программу
    FlightSchedule schedule; // Создание объекта класса FlightSchedule для управления расписанием
    string input; // Переменная для хранения входной строки
    while (getline(cin, input)) { // Чтение строк из стандартного ввода до EOF
        istringstream iss(input); // Создание потока строк для обработки входной строки
        string commandStr; // Переменная для хранения команды
        iss >> commandStr; // Чтение команды из потока строк
        if (commandStr == "create_plane") { // Проверка, является ли команда командой создания самолета
            schedule.processCommand(Type::CREATE_PLANE, input.substr(input.find(' ') + 1)); // Обработка команды создания самолета
        } else if (commandStr == "planes_for_town") { // Проверка, является ли команда командой получения самолетов по городу
            schedule.processCommand(Type::PLANES_FOR_TOWN, input.substr(input.find(' ') + 1)); // Обработка команды получения самолетов по городу
        } else if (commandStr == "towns_for_plane") { // Проверка, является ли команда командой получения городов по самолету
            schedule.processCommand(Type::TOWNS_FOR_PLANE, input.substr(input.find(' ') + 1)); // Обработка команды получения городов по самолету
        } else if (commandStr == "planes") { // Проверка, является ли команда командой отображения всех самолетов
            schedule.processCommand(Type::PLANES, ""); // Обработка команды отображения всех самолетов
        } else { // Если команда не распознана
            cout << "Error: " << commandStr << endl; // Вывод сообщения об ошибке
        }
    }
    return 0; // Возвращение кода завершения программы
}