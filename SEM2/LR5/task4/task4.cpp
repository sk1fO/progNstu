#include <iostream>
#include <unordered_map>
#include <unordered_set>
#include <string>
#include <vector>

using namespace std;

// Мапа в мапе для хранения друзей
unordered_map<string, unordered_set<string>> friends;

// Функция для добавления друзей
void addFriends(const string& i, const string& j) {
    if (i.empty() || j.empty()) {
        cout << "ОШИБКА: Имена не могут быть пустыми." << endl;
        return;
    }
    friends[i].insert(j);
    friends[j].insert(i);
    cout << "Дружба между " << i << " и " << j << " добавлена." << endl;
}

// Функция для подсчета количества друзей
int countFriends(const string& i) {
    if (i.empty()) {
        cout << "ОШИБКА: Имя не может быть пустым." << endl;
        return -1;
    }
    return friends[i].size();
}

// Функция для проверки, являются ли два пользователя друзьями
string areFriends(const string& i, const string& j) {
    if (i.empty() || j.empty()) {
        return "ОШИБКА: Имена не могут быть пустыми.";
    }
    if (friends[i].count(j)) {
        return "ДА";
    }
    return "НЕТ";
}

int main() {
    while (true) {
        string command;
        cout << "Введите команду (FRIENDS, COUNT, QUESTION, exit): ";
        cin >> command;

        if (command == "exit") {
            cout << "Завершение программы." << endl;
            break;
        } else if (command == "FRIENDS") {
            string i, j;
            cin >> i >> j;
            if (i.empty() || j.empty()) {
                cout << "ОШИБКА: Имена не могут быть пустыми." << endl;
                continue;
            }
            addFriends(i, j);
        } else if (command == "COUNT") {
            string i;
            cin >> i;
            if (i.empty()) {
                cout << "ОШИБКА: Имя не может быть пустым." << endl;
                continue;
            }
            int count = countFriends(i);
            if (count != -1) {
                cout << "Количество друзей у " << i << ": " << count << endl;
            }
        } else if (command == "QUESTION") {
            string i, j;
            cin >> i >> j;
            if (i.empty() || j.empty()) {
                cout << "ОШИБКА: Имена не могут быть пустыми." << endl;
                continue;
            }
            cout << "Являются ли " << i << " и " << j << " друзьями? " << areFriends(i, j) << endl;
        } else {
            cout << "ОШИБКА: Неизвестная команда." << endl;
        }
    }

    return 0;
}