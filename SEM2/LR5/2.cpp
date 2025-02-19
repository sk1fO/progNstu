#include <iostream>
#include <queue>
#include <vector>
#include <string>
#include <sstream>
#include <iomanip>
#include <algorithm> // Добавлено для min_element

using namespace std;

// Структура для представления посетителя
class Visitor {
public:
    string ticketNumber;
    int visitDuration;

    Visitor(string ticketNumber, int visitDuration)
        : ticketNumber(ticketNumber), visitDuration(visitDuration) {}
};

// Класс для электронной очереди
class Queue {
private:
    vector<Visitor> visitors;
    int windowCount;

public:
    Queue(int windows) : windowCount(windows) {}

    // Добавление посетителя в очередь
    void enqueue(int duration) {
        string ticketNumber = generateTicket();
        visitors.emplace_back(ticketNumber, duration);
        cout << ">>> " << ticketNumber << endl;
    }

    // Генерация номера талона
    string generateTicket() {
        static int counter = 0;
        counter++;
        stringstream ss;
        ss << "T" << setw(3) << setfill('0') << counter;
        return ss.str();
    }

    // Распределение посетителей по окнам
    void distribute() {
        vector<int> windowTimes(windowCount, 0);
        vector<vector<string>> windowVisitors(windowCount);

        for (const auto& visitor : visitors) {
            // Находим окно с наименьшим временем ожидания
            int minWindow = min_element(windowTimes.begin(), windowTimes.end()) - windowTimes.begin();
            windowTimes[minWindow] += visitor.visitDuration; // Обновляем время окна
            windowVisitors[minWindow].push_back(visitor.ticketNumber); // Добавляем талон к окну
        }

        // Вывод результата распределения
        for (int i = 0; i < windowCount; i++) {
            cout << ">>> Окно " << (i + 1) << " (" << windowTimes[i] << " минут): ";
            for (size_t j = 0; j < windowVisitors[i].size(); j++) {
                cout << windowVisitors[i][j];
                if (j != windowVisitors[i].size() - 1) {
                    cout << ", "; // Комма для разделения талонов
                }
            }
            cout << endl;
        }
    }
};

int main() {
    int windowCount;
    cout << ">>> Введите кол-во окон" << endl;
    cin >> windowCount;

    Queue queue(windowCount);
    string command;

    while (true) {
        cout << "<<< ";
        cin >> command;

        if (command == "ENQUEUE") {
            int duration;
            cin >> duration;
            queue.enqueue(duration);
        } else if (command == "DISTRIBUTE") {
            queue.distribute();
            break; // Завершение программы после распределения
        } else {
            cout << "<<< Неизвестная команда: " << command << endl; // Уведомление о неизвестной команде
        }
    }
    return 0;
}

