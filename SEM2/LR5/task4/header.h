#pragma once //чтобы заголовочный исходный файл при компиляции подключался строго один раз
#include <iostream>
#include <vector>
using namespace std;

void NEW_STUDENTS(vector<int>& students, int& number);
void SUSPICIOUS(vector<int>& studentsInDanger, int number);
void IMMORTAL(vector<int>& studentsInDanger);
void TOP_LIST(vector<int> studentsInDanger);
void SCOUNT(vector <int> studentsInDanger);
