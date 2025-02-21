#include <iostream>
#include <cmath>
using namespace std;

// Функция 1:
double function1(double x) {
    double y = -(6 * x + 36) / 7;
    return y;
}

// Функция 2:
double function2(double x) {
    double y = pow(x, 3) + 1.5 * pow(x, 2) - 2.5 * x - 3;
    return y;
}

// Функция 3:
double function3(double x) {
    double y = -2 * x + 10;
    return y;
}

// Вычисление и вывод таблицы значений
void compute_table(double X_start, double X_end, double dx) {
    // Вывод заголовка таблицы
    cout << "X\t| Y" << endl;
    cout << "-----------------------" << endl;

    // Итерация по диапазону значений x
    double x = X_start;
    while (x <= X_end){
        double y; // Переменная для хранения значения функции

        // Определение, какую функцию использовать в зависимости от значения x
        if (x <= -2.5){
            y = function1(x);
        }
        else if (x >= -2.5 && x <= 2){
            y = function2(x);
        }
        else if(x >= 2 && x <= 6){
            y = function3(x);
        }

        // Вывод текущих значений x и y
        cout << x << "\t| " << y << endl;

        // Увеличение x на шаг
        x += dx;
    }
}

int main() {
    // Установка диапазона значений x и шага
    double X_start = -6;
    double X_end = 6;
    double dx = 0.5;

    // Вызов функции compute_table для генерации и вывода таблицы
    compute_table(X_start, X_end, dx);

    return 0;
}