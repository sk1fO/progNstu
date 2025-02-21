#include <iostream>
#include <vector>
#include <cmath>
using namespace std;

// Функция решета Эратосфена для нахождения всех простых чисел до max_num
vector<int> create_sieve(int max_num) {
    vector<bool> is_prime(max_num + 1, true);
    is_prime[0] = is_prime[1] = false;
    for (int p = 2; p * p <= max_num; ++p) {
        if (is_prime[p]) {
            for (int i = p * p; i <= max_num; i += p) {
                is_prime[i] = false;
            }
        }
    }
    vector<int> primes;
    for (int p = 2; p <= max_num; ++p) {
        if (is_prime[p]) {
            primes.push_back(p);
        }
    }
    return primes;
}

// Функция возведения в степень по модулю
int powerMod(int base, int degree, int modul) {
    int result = 1;
    for (int i = 1; i <= degree; i++) {
        result *= base; // Возведение числа base в степень i
        result %= modul; // Получение остатка от деления на modul
    }
    return result;
}


// Функция для генерации простого числа по алгоритму ГОСТ
int GOST(int q, int t) {
    int N, u, p, step = 1; // Инициализация переменных
    double E; // Инициализация переменной для ошибки
    auto degree = [](int num, int degr) {
        double res = 1;
        while (degr > 0) {
                res *= num;
                degr--;
                }
                return res;
                }; // Лямбда-функция для возведения в степень

    while (true) {
        switch (step) {
        case 1: // Шаг 1
            E = 0; // Инициализация ошибки
            N = ceil(degree(2, t - 1) / q) + ceil(degree(2, t - 1) * E / q); // Вычисление N
            if (N % 2 == 1) {
                N++; // Убеждаемся, что N четное
            }

            u = 0;   // Шаг 2

        case 3:     // Шаг 3
            p = (N + u) * q + 1; // Вычисление потенциального простого числа

            if (p > degree(2, t)) {  // Шаг 4
                step = 1; // Возвращаемся к шагу 1
                break;
            }

            if (powerMod(2, p - 1, p) == 1 && powerMod(2, N + u, p) != 1) {
                    return p; // Проверка на простоту
            }
            u += 2; // Увеличение u на 2
            step = 3; // Возвращаемся к шагу 3
        }
    }
}

int main() {

    vector<int> sieve = create_sieve(500); // Создание решета Эратосфена

    cout << "ГОСТ" << endl; // Вывод заголовка
    for(int i = 3; i < 13; i++){
        cout << GOST(3, i) << "\t";
    }
    cout << endl;
    return 0;
}