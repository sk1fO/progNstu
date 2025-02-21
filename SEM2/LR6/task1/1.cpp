#include <iostream>
#include <vector>
#include <ctime>
using namespace std;

// Функция для генерации случайных вещественных чисел в диапазоне [0, 100]
double randomDouble() {
    return static_cast<double>(rand()) / RAND_MAX * 100; // Генерация случайного числа от 0 до 100
}
// Функция для инициализации квадратной матрицы
vector<vector<double>> initializeMatrix(int N) {
    vector<vector<double>> matrix(N, vector<double>(N)); // Создание матрицы размером NxN
    for (int i = 0; i < N; ++i) { // Цикл по строкам
        for (int j = 0; j < N; ++j) { // Цикл по столбцам
            matrix[i][j] = randomDouble(); // Заполнение элемента случайным числом
        }
    }
    return matrix; // Возвращение заполненной матрицы
}

// Функция для создания новой матрицы из левой верхней четверти исходной матрицы
vector<vector<double>> createQuarterMatrix(const vector<vector<double>>& matrix) {
    int N = matrix.size(); // Определение размера исходной матрицы
    int quarterSize = N / 2; // Определение размера четверти матрицы
    vector<vector<double>> quarterMatrix(quarterSize, vector<double>(quarterSize)); // Создание четверти матрицы
    for (int i = 0; i < quarterSize; ++i) { // Цикл по строкам четверти
        for (int j = 0; j < quarterSize; ++j) { // Цикл по столбцам четверти
            quarterMatrix[i][j] = matrix[i][j]; // Копирование элемента из исходной матрицы
        }
    }
    return quarterMatrix; // Возвращение четверти матрицы
}

int main() {
    
    srand(time(NULL)); 

    int N;
    cout << "Введите четное число N: "; 
    cin >> N; 

    if (N % 2 != 0) { // Проверка, является ли число N четным
        cout << "Число N должно быть четным." << endl; // Вывод сообщения об ошибке
        return 1; // Завершение программы с кодом ошибки
    }

    vector<vector<double>> matrix = initializeMatrix(N); // Инициализация исходной матрицы
    cout << "Исходная матрица:" << endl; // Вывод сообщения
    for (const auto& row : matrix) { // Цикл по строкам матрицы
        for (double val : row) { // Цикл по элементам строки
            cout << val << " "; // Вывод элемента
        }
        cout << endl;
    }

    vector<vector<double>> quarterMatrix = createQuarterMatrix(matrix); // Создание четверти матрицы
    cout << "Левая верхняя четверть матрицы:" << endl; // Вывод сообщения
    for (const auto& row : quarterMatrix) { // Цикл по строкам четверти матрицы
        for (double val : row) { // Цикл по элементам строки
            cout << val << " "; // Вывод элемента
        }
        cout << endl;
    }

    return 0;
}