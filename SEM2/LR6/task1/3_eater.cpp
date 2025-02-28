#include <iostream>
#include <vector>
#include <ctime>
#include <unistd.h>
using namespace std;

// Функция для инициализации матрицы с устойчивой фигурой и двумя планерами
vector<vector<int>> initializeMatrix(int rows, int cols) {
    vector<vector<int>> matrix(rows, vector<int>(cols, 0)); // Создание матрицы, заполненной нулями

    if (rows >= 11 && cols >= 11) { // Проверка, что матрица достаточно большая
        // блок
        matrix[2][9] = 1;
        matrix[1][9] = 1;
        matrix[2][10] = 1;
        matrix[1][10] = 1;

        // Цветок
        matrix[6][2] = 1;
        matrix[6][3] = 1;
        matrix[5][3] = 1;
        matrix[4][3] = 1;

        matrix[3][4] = 1;
        matrix[4][5] = 1;
        matrix[2][5] = 1;
        matrix[3][6] = 1;

        // планер 
        matrix[1][1] = 1;
        matrix[2][2] = 1;
        matrix[2][3] = 1;
        matrix[3][1] = 1;
        matrix[3][2] = 1;
    }

    return matrix; // Возвращение матрицы
}

// Функция для подсчета соседей клетки
int countNeighbors(const vector<vector<int>>& matrix, int x, int y) {
    int count = 0; // Инициализация счетчика соседей
    int rows = matrix.size(); // Определение количества строк
    int cols = matrix[0].size(); // Определение количества столбцов

    for (int i = -1; i <= 1; ++i) { // Цикл по соседним строкам
        for (int j = -1; j <= 1; ++j) { // Цикл по соседним столбцам
            if (i == 0 && j == 0) continue; // Пропуск самой клетки
            int newX = (x + i + rows) % rows; // Определение новой строки с учетом зацикливания
            int newY = (y + j + cols) % cols; // Определение нового столбца с учетом зацикливания
            count += matrix[newX][newY]; // Увеличение счетчика, если соседняя клетка живая
        }
    }

    return count; // Возвращение количества соседей
}

// Функция для обновления состояния клеток
vector<vector<int>> updateMatrix(const vector<vector<int>>& matrix) {
    int rows = matrix.size(); // Определение количества строк
    int cols = matrix[0].size(); // Определение количества столбцов
    vector<vector<int>> newMatrix(rows, vector<int>(cols)); // Создание новой матрицы для обновленных состояний

    for (int i = 0; i < rows; ++i) { // Цикл по строкам
        for (int j = 0; j < cols; ++j) { // Цикл по столбцам
            int neighbors = countNeighbors(matrix, i, j); // Подсчет соседей для текущей клетки
            if (matrix[i][j] == 1) { // Если текущая клетка живая
                newMatrix[i][j] = (neighbors == 2 || neighbors == 3) ? 1 : 0; // Правило для живой клетки
            } else { // Если текущая клетка мертвая
                newMatrix[i][j] = (neighbors == 3) ? 1 : 0; // Правило для мертвой клетки
            }
        }
    }

    return newMatrix; // Возвращение обновленной матрицы
}

// Функция для вывода матрицы на экран
void printMatrix(const vector<vector<int>>& matrix) {
    for (const auto& row : matrix) { // Цикл по строкам матрицы
        for (int val : row) { // Цикл по элементам строки
            cout << (val ? '1' : ' ') << " "; // Вывод '1' для живой клетки и пробела для мертвой
        }
        cout << endl; // Переход на новую строку
    }
}

int main() {
    srand(time(NULL)); // Инициализация генератора случайных чисел

    int rows, cols;
    cout << "Введите размеры матрицы (строки и столбцы): ";
    cin >> rows >> cols;

    // Инициализация матрицы с устойчивой фигурой и двумя планерами
    vector<vector<int>> matrix = initializeMatrix(rows, cols);

    while (true) {
        printMatrix(matrix); // Вывод текущего состояния матрицы
        matrix = updateMatrix(matrix); // Обновление состояния матрицы
        usleep(500000); // Задержка в 0.5 секунды
        system("clear"); // Очистка экрана 
    }

    return 0;
}