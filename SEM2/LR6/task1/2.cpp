#include <iostream>
#include <vector>
#include <ctime>

using namespace std;

// Функция для генерации случайных целых чисел в диапазоне [1000, 5000]
int randomInt() {
    return rand() % 4000 + 1000; // Генерация случайного числа от 1000 до 5000
}

// Функция для инициализации матрицы
vector<vector<int>> initializeMatrix(int M, int N) {
    vector<vector<int>> matrix(M, vector<int>(N)); // Создание матрицы размером MxN
    for (int i = 0; i < M; ++i) { // Цикл по строкам
        for (int j = 0; j < N; ++j) { // Цикл по столбцам
            matrix[i][j] = randomInt(); // Заполнение элемента случайным числом
        }
    }
    return matrix; // Возвращение заполненной матрицы
}

// Функция для подсчета суммы цифр числа
int sumOfDigits(int number) {
    int sum = 0; // Инициализация суммы цифр
    while (number > 0) { // Пока число больше нуля
        sum += number % 10; // Добавление последней цифры к сумме
        number /= 10; // Удаление последней цифры из числа
    }
    return sum; // Возвращение суммы цифр
}

// Функция для поиска строки с наименьшей суммой цифр элементов
int findRowWithMinDigitSum(const vector<vector<int>>& matrix) {
    int M = matrix.size(); // Определение количества строк
    int minSum = 1000000; // Инициализация минимальной суммы как максимально возможное целое число
    int minRowIndex = -1; // Инициализация индекса строки с минимальной суммой

    for (int i = 0; i < M; ++i) { // Цикл по строкам
        int currentSum = 0; // Инициализация текущей суммы цифр
        for (int j = 0; j < matrix[i].size(); ++j) { // Цикл по элементам строки
            currentSum += sumOfDigits(matrix[i][j]); // Суммирование цифр элемента
        }
        if (currentSum < minSum) { // Если текущая сумма меньше минимальной
            minSum = currentSum; // Обновление минимальной суммы
            minRowIndex = i; // Обновление индекса строки с минимальной суммой
        }
    }

    return minRowIndex; // Возвращение индекса строки с минимальной суммой цифр
}

int main() {
    setlocale(LC_ALL, "rus");
    srand(time(NULL));

    int M, N;
    cout << "Введите размеры матрицы M и N: ";
    cin >> M >> N; 

    vector<vector<int>> matrix = initializeMatrix(M, N); // Инициализация матрицы
    cout << "Исходная матрица:" << endl; // Вывод сообщения
    for (const auto& row : matrix) { // Цикл по строкам матрицы
        for (int val : row) { // Цикл по элементам строки
            cout << val << " "; // Вывод элемента
        }
        cout << endl; // Переход на новую строку
    }

    int minRowIndex = findRowWithMinDigitSum(matrix); // Поиск строки с наименьшей суммой цифр элементов
    cout << "Строка с наименьшей суммой цифр элементов: " << minRowIndex + 1 << endl; // Вывод результата

    return 0;
}