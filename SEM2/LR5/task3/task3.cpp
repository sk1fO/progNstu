#include "header.h"
#include "source.cpp"
#include <iostream>
#include <vector>
#include <string>
using namespace std;

enum class Type {
    CREATE_TRL = 1,
    TRL_IN_STOP,
    STOPS_IN_TRL,
    TRLS,
    EXIT
};

void menu() {
    cout << "<1> - CREATE_TRL" << endl << "<2> - TRL_IN_STOP" << endl << "<3> - STOPS_IN_TRL" << endl << "<4> - TRLS" << endl << "<5> - EXIT" << endl;
}

void functions() {
    bool exit = true;
    vector<string> trls;    //вектор троллейбусов
    vector<string> stops;   //вектор остановок
    map <string, vector<string>> route; //контейнер для маршрутов
    int command;
    menu();
    while (exit) {
        cout << "Enter the command number: ";
        cin >> command;
        switch (command) {
        case static_cast<int>(Type::CREATE_TRL):
            CREATE_TRL(trls, stops, route); //создание троллейбуса с именем trl, который проходит через остановки
            break;
        case static_cast<int>(Type::TRL_IN_STOP):
            TRL_IN_STOP(route); //вывод всех троллейбусов, которые проходят через конкретную остановку
            break;
         case static_cast<int>(Type::STOPS_IN_TRL):
            STOPS_IN_TRL(route); //вывод всех остановок, которые проезжает троллейбус с именем trl. Для каждой остановки прописать, какие троллейбусы идут через эту остановку (не включая trl)
            break;
        case static_cast<int>(Type::TRLS):
            TRLS(trls, route); //отображение всех троллейбусов с указанием остановок
            break;
        case static_cast<int>(Type::EXIT):
            cout << "Exit the program" << endl;
            exit = false;
        }
    }
}

int main() {
    functions();
}
