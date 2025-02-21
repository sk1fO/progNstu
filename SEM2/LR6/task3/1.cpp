#include <iostream>
#include <vector>
#include <iomanip>
#include <cmath>
using namespace std;

// Функция для выбора главного элемента
void pivot(vector<vector<double>>& A, vector<double>& B, int k) {
    int n = A.size();
    int maxRow = k;
    // Поиск строки с максимальным элементом в столбце k
    for (int i = k + 1; i < n; ++i) {
        if (abs(A[i][k]) > abs(A[maxRow][k])) {
            maxRow = i;
        }
    }
    // Если главный элемент не на диагонали, меняем строки местами
    if (maxRow != k) {
        swap(A[k], A[maxRow]);
        swap(B[k], B[maxRow]);
    }
}

// Функция для прямого хода метода Гаусса
void forwardGauss(vector<vector<double>>& A, vector<double>& B) {
    int n = A.size();
    // Прямой ход метода Гаусса
    for (int k = 0; k < n; ++k) {
        pivot(A, B, k); // Выбор главного элемента
        for (int i = k + 1; i < n; ++i) {
            double factor = A[i][k] / A[k][k]; // Коэффициент для исключения
            for (int j = k; j < n; ++j) {
                A[i][j] -= factor * A[k][j]; // Исключение элемента
            }
            B[i] -= factor * B[k]; // Обновление вектора свободных членов
        }
    }
}

// Функция для обратного хода метода Гаусса
vector<double> backwardGauss(vector<vector<double>>& A, vector<double>& B) {
    int n = A.size();
    vector<double> X(n);
    // Обратный ход метода Гаусса
    for (int i = n - 1; i >= 0; --i) {
        X[i] = B[i];
        for (int j = i + 1; j < n; ++j) {
            X[i] -= A[i][j] * X[j]; // Вычитание уже найденных значений
        }
        X[i] /= A[i][i]; // Деление на диагональный элемент
    }
    return X;
}

// Функция для решения СЛАУ методом Гаусса
vector<double> solveSystemGauss(vector<vector<double>>& A, vector<double>& B) {
    forwardGauss(A, B); // Прямой ход
    return backwardGauss(A, B); // Обратный ход
}

// Функция для вычисления нормы вектора
double vectorNorm(const vector<double>& V) {
    double norm = 0.0;
    // Сумма квадратов элементов вектора
    for (double v : V) {
        norm += v * v;
    }
    return sqrt(norm); // Возвращаем квадратный корень
}

// Функция для вычисления поэлементной разности двух векторов
vector<double> vectorDifference(const vector<double>& V1, const vector<double>& V2) {
    vector<double> result(V1.size());
    // Вычитание элементов векторов
    for (size_t i = 0; i < V1.size(); ++i) {
        result[i] = V1[i] - V2[i];
    }
    return result;
}

// Функция для решения СЛАУ методом простой итерации
vector<double> solveSystemIterative(const vector<vector<double>>& A, const vector<double>& B, double epsilon, int& iterations) {
    int n = A.size();
    vector<double> X(n, 0.0); // Начальное приближение
    vector<double> X_new(n);
    double error = epsilon + 1; // Начальная ошибка

    iterations = 0;
    // Итерационный процесс
    while (error > epsilon) {
        for (int i = 0; i < n; ++i) {
            X_new[i] = B[i];
            for (int j = 0; j < n; ++j) {
                if (i != j) {
                    X_new[i] -= A[i][j] * X[j]; // Исключение других переменных
                }
            }
            X_new[i] /= A[i][i]; // Деление на диагональный элемент
        }
        error = vectorNorm(vectorDifference(X_new, X)); // Вычисление ошибки
        X = X_new; // Обновление текущего приближения
        iterations++; // Увеличение счетчика итераций
    }
    return X;
}

// Функция для вывода таблицы
void printTable(int iteration, const vector<double>& X, double error) {
    cout << iteration << "\t"; // Вывод номера итерации
    for (auto x : X) {
        cout << setprecision(4) << x << "\t"; // Вывод значений переменных
    }
    cout << setprecision(4) << error << "\t" << endl; // Вывод ошибки
}

int main() {
    double M = -0.88, N = 0.1, P = 0.91;
    vector<vector<double>> A = {
        {M, -0.04, 0.21, -18},
        {0.25, -1.23, N, -0.09},
        {-0.21, N, 0.8, -0.13},
        {0.15, -1.31, 0.06, P}
    };
    vector<double> B = {-1.24, P, 2.56, M};

// Решение методом Гаусса
    vector<double> X_gauss = solveSystemGauss(A, B);
    cout << "Решение системы уравнений методом Гаусса:" << endl;
    for (int i = 0; i < X_gauss.size(); ++i) {
        cout << "x" << i + 1 << " = " << X_gauss[i] << endl;
    }

    // Решение методом простой итерации
    int iterations = 0;
    double epsilon = 1e-3;
    vector<double> X_iterative = solveSystemIterative(A, B, epsilon, iterations);
    cout << endl << "Решение системы уравнений методом простых итераций:" << endl;
    for (int i = 0; i < X_iterative.size(); ++i) {
        cout << "x" << i + 1 << " = " << X_iterative[i] << endl;
    }
    cout << "Количество итераций: " << iterations << endl;

    // Вывод таблицы для метода простой итерации
    cout << "N\t" << "X1\t" << "X2\t" << "X3\t" << "X4\t" << "E" << endl;
    vector<double> X_prev(4, 0.0);
    for (int i = 0; i < iterations; ++i) {
        vector<double> X_new(4);
        for (int j = 0; j < 4; ++j) {
            X_new[j] = B[j];
            for (int k = 0; k < 4; ++k) {
                if (j != k) {
                    X_new[j] -= A[j][k] * X_prev[k];
                }
            }
            X_new[j] /= A[j][j];
        }
        double error = vectorNorm(vectorDifference(X_new, X_prev));
        printTable(i + 1, X_new, error);
        X_prev = X_new;
    }

    return 0;
}