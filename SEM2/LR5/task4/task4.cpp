#include <iostream>
#include <map>
#include <vector>
#include <algorithm>
using namespace std;

int main() {
    int N;
    cin >> N; // Считываем количество операций

    map<int, vector<string>> schedule; // Создаем ассоциативный массив для хранения расписания
    int currentMonth = 1; // Начинаем с первого месяца

    for (int i = 0; i < N; ++i) {
        string command;
        getline(cin, command); // Считываем команду полностью

        if (command.substr(0, 5) == "class") { // Если команда начинается с "CLASS"
            int day = stoi(command.substr(6, command.find(' ', 6) - 6)); // Извлекаем день
            string subject = command.substr(command.find(' ', 6) + 1); // Извлекаем предмет
            if (find(schedule[day].begin(), schedule[day].end(), subject) == schedule[day].end()) {
                // Проверяем, существует ли уже такое занятие
                schedule[day].push_back(subject); // Добавляем предмет в расписание на указанный день
            }
        } else if (command == "next") { // Если команда "NEXT"
            map<int, vector<string>> newSchedule; // Создаем новое расписание для следующего месяца
            int lastDay = 0;
            for (const auto& entry : schedule) { // Находим последний день текущего месяца
                if (entry.first > lastDay) {
                    lastDay = entry.first;
                }
            }
            for (const auto& entry : schedule) { // Переносим занятия в новое расписание
                if (entry.first > lastDay) {
                    newSchedule[lastDay].insert(newSchedule[lastDay].end(), entry.second.begin(), entry.second.end());
                } else {
                    newSchedule[entry.first] = entry.second;
                }
            }
            schedule = newSchedule; // Обновляем текущее расписание
            currentMonth++; // Переходим к следующему месяцу
        } else if (command.substr(0, 4) == "view") { // Если команда начинается с "VIEW"
            int day = stoi(command.substr(5)); // Извлекаем день
            if (schedule.find(day) != schedule.end() && !schedule[day].empty()) { // Если есть занятия в этот день
                cout << "В " << day << " день " << schedule[day].size() << " занятий в университете: ";
                for (size_t j = 0; j < schedule[day].size(); ++j) { // Выводим все занятия
                    if (j > 0) cout << ", ";
                    cout << schedule[day][j];
                }
                cout << endl;
            } else { // Если занятий нет
                cout << "В " << day << " день мы свободны!" << endl;
            }
        }
    }

    return 0; // Завершение программы
}