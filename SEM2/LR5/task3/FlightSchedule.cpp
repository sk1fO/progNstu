#include "FlightSchedule.h"
#include <iostream>

using namespace std;

void FlightSchedule::processCommand(Type type, const string& input) {
    istringstream iss(input);
    string token;
    switch (type) {
        case Type::CREATE_PLANE:
            createPlane(iss);
            break;
        case Type::PLANES_FOR_TOWN:
            iss >> token;
            getPlanesForTown(token);
            break;
        case Type::TOWNS_FOR_PLANE:
            iss >> token;
            getTownsForPlane(token);
            break;
        case Type::PLANES:
            displayPlanes();
            break;
    }
}

void FlightSchedule::createPlane(istringstream& iss) {
    string planeName;
    string town;
    vector<string> towns;
    iss >> planeName;

    // Проверка на существование самолета
    if (planeMap.count(planeName)) {
        cout << "Ошибка: Самолет с именем " << planeName << " уже существует." << endl;
        return;
    }

    while (iss >> town) {
        towns.push_back(town);
    }

    // Проверка на отсутствие остановок
    if (towns.empty()) {
        cout << "Ошибка: Невозможно создать самолет без остановок." << endl;
        return;
    }

    // Проверка на одну остановку
    if (towns.size() < 2) {
        cout << "Ошибка: Самолет должен иметь как минимум две остановки." << endl;
        return;
    }

    // Проверка на дубликаты городов
    unordered_set<string> uniqueTowns(towns.begin(), towns.end());
    if (uniqueTowns.size() != towns.size()) {
        cout << "Ошибка: В маршруте присутствуют дубликаты городов." << endl;
        return;
    }

    planeMap[planeName] = towns;
    for (const auto& t : towns) {
        townMap[t].insert(planeName);
    }
    cout << "Самолет " << planeName << " успешно создан." << endl;
}


void FlightSchedule::getPlanesForTown(const string& town) {
    if (townMap.count(town)) {
        cout << "Самолеты, летящие через город " << town << ": ";
        for (const auto& plane : townMap[town]) {
            cout << plane << " ";
        }
        cout << endl;
    } else {
        cout << "Город " << town << " не найден в расписании." << endl;
    }
}

void FlightSchedule::getTownsForPlane(const string& plane) {
    if (planeMap.count(plane)) {
        cout << "Самолет " << plane << " летит через города: ";
        for (const auto& town : planeMap[plane]) {
            cout << town << " ";
        }
        cout << endl;
        cout << "Другие самолеты, оставшиеся в этих городах: ";
        unordered_set<string> otherPlanes;
        for (const auto& town : planeMap[plane]) {
            for (const auto& p : townMap[town]) {
                if (p != plane) {
                    otherPlanes.insert(p);
                }
            }
        }
        for (const auto& p : otherPlanes) {
            cout << p << " ";
        }
        cout << endl;
    } else {
        cout << "Самолет " << plane << " не найден в расписании." << endl;
    }
}

void FlightSchedule::displayPlanes() {
    cout << "Все самолеты в расписании:" << endl;
    for (const auto& [plane, towns] : planeMap) {
        cout << plane << ": ";
        for (const auto& town : towns) {
            cout << town << " ";
        }
        cout << endl;
    }
}