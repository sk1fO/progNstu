#include "header.h"
#include "cource.cpp"
#include <iostream>
#include <vector>
using namespace std;
enum class Type {
	NEW_STUDENTS = 1,
	SUSPICIOUS,
	IMMORTAL,
	TOP_LIST,
	SCOUNT,
	EXIT,
};

void menu() {
	cout << "<1> - NEW_STUDENTS" << endl << "<2> - SUSPICIOUS" << endl << "<3> - IMMORTAL" << endl << "<4> - TOP_LIST" << endl << "<5> - SCOUNT" << endl << "<6> - EXIT" << endl;
}

void functions() {
    bool exit = true;
    vector<int> students;
    vector<int> studentsInDanger;
    int number = 1;
    int command;
    menu();
    while (exit) {
        cout << "Введите номер команды: ";
        cin >> command;
        switch (command) { //значение переменной указанной в условии switch сравнивается со значениями, которые следуют за ключевым словом case
        case static_cast<int>(Type::NEW_STUDENTS): // при преобразовании из "size_t" в "int" возможна потеря данных
            NEW_STUDENTS(students,number); //добавление студентов в конец очереди
            break;
        case static_cast<int>(Type::SUSPICIOUS):
            SUSPICIOUS(studentsInDanger, students.size()); //студент с порядковым номером number_student является крайне подозрительным и входит в топ-лист на отчисление
            break;
        case static_cast<int>(Type::IMMORTAL):
            IMMORTAL(studentsInDanger); //студент с порядковым номером number_student является неприкасаемым и из топ-листа на отчисление уходит
            break;
        case static_cast<int>(Type::TOP_LIST):
            TOP_LIST(studentsInDanger); //вывод отсортированного списка студентов, входящих в топ-лист на отчисление
            break;
        case static_cast<int>(Type::SCOUNT):
            SCOUNT(studentsInDanger); //вывод количества студентов, входящих в топ-лист на отчисление
            break;
        case static_cast<int>(Type::EXIT):
            cout << "Выход из программы" << endl;
            exit = false;
        }
    }
}
int main() {
    functions();
}
