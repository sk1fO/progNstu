#include "FlightSchedule.h" 
#include <iostream> 

using namespace std; 

// Метод для обработки команд в зависимости от их типа
void FlightSchedule::processCommand(Type type, const string& input) {
    istringstream iss(input); // Создание потока строк для обработки входной строки
    string token; // Переменная для хранения токенов (частей строки)
    switch (type) { // Определение действия в зависимости от типа команды
        case Type::CREATE_PLANE: // Если команда для создания самолета
            createPlane(iss); // Вызов метода для создания самолета
            break;
        case Type::PLANES_FOR_TOWN: // Если команда для получения самолетов по городу
            iss >> token; // Чтение города из входной строки
            getPlanesForTown(token); // Вызов метода для получения самолетов по городу
            break;
        case Type::TOWNS_FOR_PLANE: // Если команда для получения городов по самолету
            iss >> token; // Чтение самолета из входной строки
            getTownsForPlane(token); // Вызов метода для получения городов по самолету
            break;
        case Type::PLANES: // Если команда для отображения всех самолетов
            displayPlanes(); // Вызов метода для отображения всех самолетов
            break;
    }
}

// Метод для создания записи о самолете и его маршруте
void FlightSchedule::createPlane(istringstream& iss) {
    string planeName; // Переменная для хранения имени самолета
    string town; // Переменная для хранения названия города
    vector<string> towns; // Вектор для хранения списка городов
    iss >> planeName; // Чтение имени самолета из потока
    while (iss >> town) { // Чтение городов из потока
        towns.push_back(town); // Добавление города в список
    }
    planeMap[planeName] = towns; // Сохранение маршрута самолета в карте
    for (const auto& t : towns) { // Обход списка городов
        townMap[t].insert(planeName); // Добавление самолета в карту городов
    }
}

// Метод для получения списка самолетов, летящих через указанный город
void FlightSchedule::getPlanesForTown(const string& town) {
    if (townMap.count(town)) { // Проверка наличия города в карте
        cout << "Planes flying through the city " << town << ": ";
        for (const auto& plane : townMap[town]) { // Обход списка самолетов
            cout << plane << " ";
        }
        cout << endl;
    } else {
        cout << "Town " << town << " not found in schedule." << endl;
    }
}

// Метод для получения списка городов, через которые летит указанный самолет
void FlightSchedule::getTownsForPlane(const string& plane) {
    if (planeMap.count(plane)) { // Проверка наличия самолета в карте
        cout << "Plane " << plane << " flies through the cities: ";
        for (const auto& town : planeMap[plane]) { // Обход списка городов
            cout << town << " ";
        }
        cout << endl;
        cout << "Other planes remaininng in these cities: ";
        unordered_set<string> otherPlanes; // Множество для хранения других самолетов
        for (const auto& town : planeMap[plane]) { // Обход списка городов
            for (const auto& p : townMap[town]) { // Обход списка самолетов в городе
                if (p != plane) { // Проверка, что самолет не тот, для которого запрашивается информация
                    otherPlanes.insert(p); // Добавление самолета в множество
                }
            }
        }
        for (const auto& p : otherPlanes) { // Обход множества самолетов
            cout << p << " ";
        }
        cout << endl;
    } else {
        cout << "Plane " << plane << " not found in the schedule." << endl;
    }
}

// Метод для отображения всех самолетов и их маршрутов
void FlightSchedule::displayPlanes() {
    cout << "All planes on schedule:" << endl;
    for (const auto& [plane, towns] : planeMap) { // Обход карты самолетов
        cout << plane << ": ";
        for (const auto& town : towns) { // Обход списка городов
            cout << town << " ";
        }
        cout << endl;
    }
}