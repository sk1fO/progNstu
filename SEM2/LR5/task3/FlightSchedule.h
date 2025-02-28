#ifndef FLIGHT_SCHEDULE_H
#define FLIGHT_SCHEDULE_H

#include <string>
#include <unordered_map>
#include <unordered_set>
#include <vector>
#include <sstream>

enum class Type {
    CREATE_PLANE,
    PLANES_FOR_TOWN,
    TOWNS_FOR_PLANE,
    PLANES
};

class FlightSchedule {
public:
    void processCommand(Type type, const std::string& input);

private:
    void createPlane(std::istringstream& iss);
    void getPlanesForTown(const std::string& town);
    void getTownsForPlane(const std::string& plane);
    void displayPlanes();

    std::unordered_map<std::string, std::vector<std::string>> planeMap;
    std::unordered_map<std::string, std::unordered_set<std::string>> townMap;
};

#endif