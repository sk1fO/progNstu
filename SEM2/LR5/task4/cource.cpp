#include <iostream>
#include <vector>
#include <algorithm>
using namespace std;

void NEW_STUDENTS(vector<int>& students, int& number) {  //добавление студентов в конец очереди
	cout << "Введите количество студентов: "; //запрос у пользователя количества студентов для добавления и сохранение этого значения в переменной numbers
	
	int numbers;
	cin >> numbers;
	for (int i = number; i < number + numbers; i++) { //добавление numbers новых студентов в вектор, начиная со значения number, и увеличение значения number на numbers
		students.push_back(number + i - 1);
	}
	cout << "Welcome " << numbers << " clever students!" << endl;
	number += numbers; //переменная number используется как счетчик добавленных студентов
}

void SUSPICIOUS(vector<int>& studentsInDanger, int number) { //студент с порядковым номером number_student является крайне подозрительным и входит в топ-лист на отчисление
	cout << "Введите номер студента, который входит в топ-лист на отчисление: ";
	int num;
	cin >> num;
	if (num > number) {  //проверка, есть ли такой студент
		cout << "Такого студента нет" << endl;
	}
	else {
		studentsInDanger.push_back(num);
		cout << "Студент находится в топ-листе на отчисление" << endl;
	}
}

void IMMORTAL(vector<int>& studentsInDanger) {  //студент с порядковым номером number_student является неприкасаемым и из топ-листа на отчисление уходит
	cout << "Введите номер студента, который уходит из топ-листа на отчисление: ";
	int num;
	cin >> num;
	if (count(studentsInDanger.begin(), studentsInDanger.end(), num) == 0) {   //ошибка, если такого студента нет
		cout << "Incorrect" << endl;
	}
	else {
		studentsInDanger.erase(find(studentsInDanger.begin(), studentsInDanger.end(), num)); //удаление из топ-листа
		cout << "Student " << num << " is immortal" << endl;
	}
}

void TOP_LIST(vector<int> studentsInDanger) { //студент с порядковым номером number_student является крайне подозрительным и входит в топ-лист на отчисление
	sort(studentsInDanger.begin(), studentsInDanger.end());
	for (auto i : studentsInDanger) {
		cout << i << endl;
	}
}

void SCOUNT(vector<int> studentsInDanger) {
	cout << "Количество студентов в топ-листе на отчисление: " << studentsInDanger.size() << endl;
}
