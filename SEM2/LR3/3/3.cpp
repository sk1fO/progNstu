    #include <iostream>
    #include <vector>
    #include <numeric>
    #include <cmath>
    using namespace std;

    // Структура для хранения результатов аппроксимации
    struct ApproximationResult {
        double a; // Коэффициент a в уравнении y = a*x + b
        double b; // Коэффициент b в уравнении y = a*x + b
    };

    // Структура для хранения результатов корреляции
    struct CorrelationResult {
        double correlation; // Коэффициент корреляции
        double t_value; // Значение t-критерия Стьюдента
    };

    // Функция для моделирования изменения температуры кофе
    vector<double> coffee(double T, double Ts, double r, int time) {
        vector<double> temperatures; // Вектор для хранения температур
        for (int t = 0; t <= time; t++) {
            double temperature = Ts + (T - Ts) * exp(-r * t); // Формула для расчета температуры
            temperatures.push_back(temperature); // Добавление температуры в вектор
        }
        return temperatures; // Возвращение вектора температур
    }

    // Функция для вычисления аппроксимации данных
    ApproximationResult aproks(const vector<double>& x, const vector<double>& y) {
        double x_sum = accumulate(x.begin(), x.end(), 0.0); // Сумма значений x
        double y_sum = accumulate(y.begin(), y.end(), 0.0); // Сумма значений y
        double x2_sum = inner_product(x.begin(), x.end(), x.begin(), 0.0); // Сумма квадратов x
        double xy_sum = inner_product(x.begin(), x.end(), y.begin(), 0.0); // Сумма произведений x и y
        size_t n = x.size(); // Количество элементов в векторах x и y

        double a = (n * xy_sum - x_sum * y_sum) / (n * x2_sum - x_sum * x_sum); // Вычисление коэффициента a
        double b = (y_sum - a * x_sum) / n; // Вычисление коэффициента b

        return {a, b}; // Возвращение результатов аппроксимации
    }

    // Функция для вычисления корреляции данных
    CorrelationResult korrel(const vector<double>& x, const vector<double>& y) {
        double x_mean = accumulate(x.begin(), x.end(), 0.0) / x.size(); // Среднее значение x
        double y_mean = accumulate(y.begin(), y.end(), 0.0) / y.size(); // Среднее значение y
        double sum_xy = 0.0, sum_x2 = 0.0, sum_y2 = 0.0; // Инициализация сумм
        size_t n = x.size(); // Количество элементов в векторах x и y

        for (size_t i = 0; i < n; ++i) {
            double dx = x[i] - x_mean, dy = y[i] - y_mean; // Вычисление отклонений от средних
            sum_xy += dx * dy; // Сумма произведений отклонений
            sum_x2 += dx * dx; // Сумма квадратов отклонений x
            sum_y2 += dy * dy; // Сумма квадратов отклонений y
        }

        double r = sum_xy / (sqrt(sum_x2) * sqrt(sum_y2)); // Вычисление коэффициента корреляции
        double t = r * sqrt(n - 2) / sqrt(1 - r * r); // Вычисление t-значения

        return {r, t}; // Возвращение результатов корреляции
    }

    int main() {
        int T = 88; // Начальная температура кофе
        int Ts = 26; // Температура комнаты
        double r = 0.01; // Коэффициент остывания
        int time = 60; // Временной предел в минутах

        vector<double> temperatures = coffee(T, Ts, r, time); // Получение вектора температур
        vector<double> times(time + 1); // Создание вектора времени
        iota(times.begin(), times.end(), 0); // Заполнение вектора времени значениями от 0 до time

        ApproximationResult approx_result = aproks(times, temperatures); // Аппроксимация данных
        CorrelationResult corr_result = korrel(times, temperatures); // Корреляция данных

        // Вывод результатов аппроксимации
        cout << "Результат апроксимации:" << endl <<"a:" << approx_result.a << ", b:" << approx_result.b << endl;
        // Вывод результатов корреляции
        cout << "Результат корелляции:" << endl << "Корелляция: " << corr_result.correlation << endl << "Значение Т: " << corr_result.t_value << endl;

        // Вывод температуры кофе в зависимости от времени
        for (int i = 0; i < temperatures.size(); i++) {
            cout << "Время-" << times[i] << ":  " << temperatures[i] << " C" << endl;
        }

        return 0; // Завершение программы
    }