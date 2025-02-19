#include <iostream>
#include <vector>
#include <string>
#include <map>
using namespace std;

string getStopsFromInput() { //получение от пользователя списка остановок для троллейбуса
	string trl_stop;
	cout << "Enter the names of the stops: ";
	getline(cin, trl_stop);
	return trl_stop;
}

void addStopAndRouteToVectors(string& stop, vector<string>& stops, map <string, vector<string>>& route, string& number) { //добавление остановки в вектор stops и связывание ее с номером троллейбуса в карте route 
	stops.push_back(stop); //добавляем остановку в вектор
	route[stop].push_back(number);  //связываем троллейбус и остановки
	route[number].push_back(stop);
}

void CREATE_TRL(vector<string>& trls, vector<string>& stops, map <string, vector<string>>& route) { //создание троллейбуса с именем trl, который проходит через остановки

    string number, trl_stop;
    cout << "Enter the trolleybus number: ";
    cin >> number;
	cin.ignore();   //игнорируем \n
	trls.push_back(number);  //добавляем номер троллейбуса в вектор троллейбусов

	trl_stop = getStopsFromInput();
    string stop;
	for (auto place : trl_stop) {    
		if (place != ' ') {
			stop.push_back(place);
		}
		else {
			addStopAndRouteToVectors(stop, stops, route, number);
			stop.clear();
		}
	}
	addStopAndRouteToVectors(stop, stops, route, number);
	cout << "Trolleybus added" << endl;
}

void TRL_IN_STOP(map <string, vector<string>>& route) { //вывод всех троллейбусов, которые проходят через конкретную остановку
	//используется map, в котором ключами являются названия остановок, а значениями - векторы с названиями троллейбусов, проходящих через эти остановки

	string stop;
	cout << "Enter the name of the stop: ";
	cin >> stop;
	cout << "The following trolleybuses pass through the stop: ";

	for (auto place : route[stop]) { //проход по элементам вектора, соответствующего указанной остановке, и вывод на экран названия этих троллейбусов
		cout << place << " ";
	}
	cout << endl;
}

void STOPS_IN_TRL(map <string, vector<string>>& route) { //вывод всех остановок, которые проезжает определенный троллейбус, для каждой остановки прописать, какие троллейбусы идут через эту остановку (не включая trl)

	string number_trl;
	cout << "Enter the trolleybus number: ";
	cin >> number_trl;

	cout << "This trolleybus passes through the following stops: "; //остановки, через которые проезжает троллейбус

	for (auto place : route[number_trl]) {
		cout << place << endl;

		cout << "Other trolleybuses also pass through this stop: "; //другие троллейбусы, проезжающие через эту остановку
		for (auto numbers : route[place]) {
			if (numbers != number_trl) {
				cout << numbers << " ";
			}
		}
		cout << endl;
	}
	cout << endl;
}

void TRLS(vector<string>& trls, map <string, vector<string>>& route) { //отображение всех троллейбусов с указанием остановок

	cout << "Trolleybus timetable with stops: " << endl; //составление расписания рейсов
	for (auto place : trls) {
		cout << "Trolleybus number " << place << " drives through stops: ";
		for (auto trll : route[place]) {
			cout << trll << " ";
		}
	}
	cout << endl;
}
