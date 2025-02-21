#ifndef FLIGHT_SCHEDULE_H // Проверка, чтобы избежать повторного включения заголовочного файла
#define FLIGHT_SCHEDULE_H // Определение маркера для предотвращения повторного включения

#include <string> 
#include <unordered_map> 
#include <unordered_set> 
#include <vector> 
#include <sstream> 

// Перечисление типов команд, которые могут быть обработаны
enum class Type {
    CREATE_PLANE, // Команда для создания записи о самолете
    PLANES_FOR_TOWN, // Команда для получения списка самолетов, летящих через город
    TOWNS_FOR_PLANE, // Команда для получения списка городов, через которые летит самолет
    PLANES // Команда для отображения всех самолетов и их маршрутов
};

// Класс для управления расписанием полетов
class FlightSchedule {
public:
    void processCommand(Type type, const std::string& input); // Метод для обработки команд

private:
    void createPlane(std::istringstream& iss); // Метод для создания записи о самолете
    void getPlanesForTown(const std::string& town); // Метод для получения списка самолетов по городу
    void getTownsForPlane(const std::string& plane); // Метод для получения списка городов по самолету
    void displayPlanes(); // Метод для отображения всех самолетов и их маршрутов

    std::unordered_map<std::string, std::vector<std::string>> planeMap; // Хеш-таблица для хранения самолетов и их маршрутов
    std::unordered_map<std::string, std::unordered_set<std::string>> townMap; // Хеш-таблица для хранения городов и самолетов, летящих через них
};

#endif // FLIGHT_SCHEDULE_H // Конец условия для предотвращения повторного включения