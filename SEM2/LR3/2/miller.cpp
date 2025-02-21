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

// Функция для генерации простого числа по методу Миллера
int generate_prime_miller(vector<int> factors) {
    int m = 1; // Инициализация произведения простых чисел
    for (int q : factors) {
            m *= q; // Перемножаем простые числа
    }
    int n = 2 * m + 1; // Генерация потенциального простого числа
    return n; // Возвращаем потенциальное простое число
}

// Функция для проверки простоты числа по методу Миллера
int test_prime_miller(int n, int rounds, vector<int> factors, int& k) {
    factors.push_back(2); // Добавляем 2 в список простых чисел
    srand(time(NULL)); // Инициализация генератора случайных чисел
    vector<int> randoms; // Вектор для хранения случайных чисел
    for (int i = 0; i < rounds; i++) {
            randoms.push_back(rand() % (n - 1) + 1); // Генерация случайных чисел
    }
    for (int a : randoms) {
            if (powerMod(a, n - 1, n) != 1){
                    k++; // Увеличиваем счетчик k
                    return 0; // Проверка на простоту
            }
    }
    bool is_one_flag; // Флаг для проверки
    for (int q : factors) {
        is_one_flag = true; // Установка флага
        for (int a : randoms) {
            if (powerMod(a, (n - 1) / q, n) != 1) {
                is_one_flag = false; // Сброс флага
                break;
            }
        }
        if (is_one_flag){
                return 0; 
        }
    }
    return 1; // составное
}

int main() {

    vector<int> sieve = create_sieve(500); // Создание решета Эратосфена

    vector<int> miller_primes; // Вектор для хранения простых чисел по методу Миллера
    vector<vector<int>> factors(10); // Вектор для хранения разложений
    vector<int> k_values(10, 0); // Вектор для хранения значений k

    cout << "Миллер" << endl; // Вывод заголовка
    for (int i = 0; i < 10; i++) {
            cout << i + 1 << "\t";
    }
    cout << endl;
    for (int i = 0; i < 10; i++) {
        factors[i] = {sieve[i + 1], sieve[i + 2]}; // Заполнение разложений
        miller_primes.push_back(generate_prime_miller(factors[i])); // Генерация простых чисел по методу Миллера
        cout << miller_primes[i] << "\t"; // Вывод простых чисел
    }
    cout << endl;
    for (int i = 0; i < 10; i++) {
        int result = test_prime_miller(miller_primes[i], 10, factors[i], k_values[i]); // Проверка простых чисел по методу Миллера
        cout << (result == 1 ? "+" : "-") << "\t"; // Вывод результата проверки
    }
    cout << endl;

    for (int k : k_values) {
        cout << k << "\t"; // Вывод значений k
    }
    cout << endl;

    return 0;
}