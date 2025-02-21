#include <iostream>
#include <vector>
#include <ctime>
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

// Функция для генерации простого числа по методу Поклингтона
int pr_pokling(vector<int> factors) {
    int F = 1; // Инициализация произведения простых чисел
    for (int q : factors){
            F *= q; // Перемножаем простые числа
    }
    int R = F >> 1; // Деление на 2
    R -= R % 2; // Убеждаемся, что R четное
    int n = R * F + 1; // Генерация потенциального простого числа

    return n; // Возвращаем потенциальное простое число
}

// Функция для проверки простоты числа по методу Поклингтона
int is_pr_pokling(int n, int t, vector<int> factors, int& k) {
    srand(time(NULL)); // Инициализация генератора случайных чисел
    vector<int> randoms; // Вектор для хранения случайных чисел
    for (int i = 0; i < t; i++){
            randoms.push_back(rand() % (n - 1) + 1); // Генерация случайных чисел
    }
    for (int a : randoms){
            if (powerMod(a, n - 1, n) != 1){
                    k++; // Увеличиваем счетчик k
                    return 0; // Проверка на простоту
            }
    }

    bool is_zero_flag; // Флаг для проверки
    for (int a : randoms) {
        is_zero_flag = true; // Установка флага
        for (int q : factors) {
            if (powerMod(a, (n - 1) / q, n) == 1) {
                is_zero_flag = false; // Сброс флага
                break;
            }
        }
        if (is_zero_flag){
                return 1; // Возвращаем 1, если число составное
        }
    }
    return 0; // Возвращаем 0, если число простое
}

int main() {

    vector<int> sieve = create_sieve(500); // Создание решета Эратосфена

    vector<int> poklin_primes; // Вектор для хранения простых чисел по методу Поклингтона
    vector<vector<int>> factors(10); // Вектор для хранения разложений
    vector<int> k_values(10, 0); // Вектор для хранения значений k

    cout << "Поклингтон" << endl; // Вывод заголовка
    for (int i = 0; i < 10; i++) {
            cout << i + 1 << "\t";
    }
    cout << endl;
    for (int i = 0; i < 10; i++) {
        factors[i] = {sieve[i + 3]}; // Заполнение разложений
        poklin_primes.push_back(pr_pokling(factors[i])); // Генерация простых чисел по методу Поклингтона
        cout << poklin_primes[i] << "\t"; // Вывод простых чисел
    }
    cout << endl;
    for (int i = 0; i < 10; i++) {
        int result = is_pr_pokling(poklin_primes[i], 10, factors[i], k_values[i]); // Проверка простых чисел по методу Поклингтона
        cout << (result == 1 ? "+" : "-") << "\t"; // Вывод результата проверки
    }
    cout << endl;

    for (int k : k_values) {
        cout << k << "\t"; // Вывод значений k
    }
    cout << endl;

    return 0;
}