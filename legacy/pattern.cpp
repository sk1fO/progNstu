#include <iostream>
#include <string>
#include <cctype> // Для std::isalpha, std::isdigit, std::isupper, std::islower

using namespace std; // Используем пространство имен std

// Класс для валидации пароля
class PasswordValidator {
public:
    // Метод для проверки пароля
    bool isValid(const string& password) const {
        // Проверка длины пароля
        if (password.length() < 8) {
            cout << "Пароль должен содержать минимум 8 символов." << endl;
            return false;
        }

        // Флаги для проверки условий
        bool hasLatinLetter = false; // Есть ли латинские буквы
        bool hasDigit = false; // Есть ли цифры
        bool hasSpecialChar = false; // Есть ли служебные символы
        bool hasUppercase = false; // Есть ли заглавные буквы
        bool hasLowercase = false; // Есть ли строчные буквы
        // Набор служебных символов
        const string specialChars = "?!@*_+-%&";

        // Проходим по каждому символу пароля
        for (char ch : password) {
            if (isalpha(ch)) { // Проверка на латинские буквы
                hasLatinLetter = true;
                if (isupper(ch)) { // Проверка на заглавные буквы
                    hasUppercase = true;
                } else if (islower(ch)) { // Проверка на строчные буквы
                    hasLowercase = true;
                }
            } else if (isdigit(ch)) { // Проверка на цифры
                hasDigit = true;
            } else if (specialChars.find(ch) != string::npos) { // Проверка на служебные символы
                hasSpecialChar = true;
            }
        }

        // Проверяем все условия
        if (!hasLatinLetter) {
            cout << "Пароль должен содержать латинские буквы." << endl;
            return false;
        }
        if (!hasDigit) {
            cout << "Пароль должен содержать хотя бы одну цифру." << endl;
            return false;
        }
        if (!hasSpecialChar) {
            cout << "Пароль должен содержать хотя бы один служебный символ." << endl;
            return false;
        }
        if (!hasUppercase || !hasLowercase) {
            cout << "Пароль должен содержать хотя бы одну заглавную и одну строчную букву." << endl;
            return false;
        }

        return true;
    }
};

// Класс для регистрации пользователя
class User {
private:
    string name; // Имя пользователя
    string password; // Пароль пользователя

public:
    // Конструктор для создания пользователя
    User(const string& name, const string& password) : name(name), password(password) {}

    // Метод для вывода информации о пользователе
    void displayInfo() const {
        cout << "Пользователь зарегистрирован:" << endl;
        cout << "Имя: " << name << endl;
        cout << "Пароль: " << password << endl;
    }
};

int main() {
    string name, password;
    PasswordValidator validator; // Создаем объект для валидации пароля

    // Запрашиваем имя пользователя
    cout << "Введите имя пользователя: ";
    getline(cin, name); 

    // Запрашиваем пароль
    cout << "Введите пароль: ";
    getline(cin, password); 

    // Паттерн "Стратегия": Используем PasswordValidator для проверки пароля
    if (validator.isValid(password)) { // Если пароль валиден
        // Регистрируем пользователя
        User user(name, password); // Создаем объект пользователя
        user.displayInfo(); // Выводим информацию о пользователе
    } else {
        cout << "Регистрация не удалась. Пароль не соответствует требованиям." << endl;
    }

    return 0;
}