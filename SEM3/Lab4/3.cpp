#include <iostream>
#include <thread>
#include <vector>
#include <mutex>
#define MAX 20 // Максимальное число процессов и ресурсов

using namespace std;

// Класс реализующий алгоритм Банкера для предотвращения дедлока
class BankersAlgorithm
{
private:
    int al[MAX][MAX], m[MAX][MAX], n[MAX][MAX], avail[MAX]; // Матрицы для хранения информации о ресурсах
    int nop, nor, k, result[MAX], pnum, work[MAX], finish[MAX]; // Переменные для работы алгоритма
    mutex mtx; // Мьютекс для синхронизации потоков

public:
    BankersAlgorithm(); // Конструктор класса
    void input(); // Функция для ввода данных пользователем
    void method(); // Основной метод алгоритма Банкера
    int search(int); // Метод поиска процесса, который может быть завершен
    void display(); // Функция для отображения результатов
    void check_process(int i); // Метод для проверки возможности завершения процесса
};

// Конструктор класса
BankersAlgorithm::BankersAlgorithm()
{
    k = 0; // Начальное значение количества завершенных процессов
    for (int i = 0; i < MAX; i++) // Инициализация всех матриц и массивов нулями
    {
        for (int j = 0; j < MAX; j++)
        {
            al[i][j] = 0;
            m[i][j] = 0;
            n[i][j] = 0;
        }
        avail[i] = 0;
        result[i] = 0;
        finish[i] = 0;
    }
}

// Функция для ввода данных пользователем
void BankersAlgorithm::input()
{
    int i, j;
    cout << "Введите количество процессов: ";
    cin >> nop;
    cout << "Введите количество ресурсов: ";
    cin >> nor;
    cout << "Введите выделенные ресурсы для каждого процесса: " << endl;
    for (i = 0; i < nop; i++)
    {
        cout << "\nПроцесс " << i;
        for (j = 0; j < nor; j++)
        {
            cout << "\nРесурс " << j << ": ";
            cin >> al[i][j];
        }
    }
    cout << "Введите максимальные ресурсы, необходимые для каждого процесса: " << endl;
    for (i = 0; i < nop; i++)
    {
        cout << "\nПроцесс " << i;
        for (j = 0; j < nor; j++)
        {
            cout << "\nРесурс " << j << ": ";
            cin >> m[i][j];
            n[i][j] = m[i][j] - al[i][j];
        }
    }
    cout << "Введите текущие доступные ресурсы в системе: ";
    for (j = 0; j < nor; j++)
    {
        cout << "Ресурс " << j << ": ";
        cin >> avail[j];
        work[j] = -1;
    }
    for (i = 0; i < nop; i++)
        finish[i] = 0;
}

// Метод для проверки возможности завершения процесса
void BankersAlgorithm::check_process(int i)
{
    if (finish[i] == 0) // Если процесс еще не завершен
    {
        pnum = search(i); // Проверяем, можно ли завершить процесс
        if (pnum != -1)
        {
            mtx.lock(); // Заблокируем мьютекс для синхронизации
            result[k++] = i; // Добавляем процесс в список завершенных
            finish[i] = 1; // Отметим процесс как завершенный
            for (int j = 0; j < nor; j++)
            {
                avail[j] = avail[j] + al[i][j]; // Обновляем доступные ресурсы
            }
            mtx.unlock(); // Разблокируем мьютекс
        }
    }
}

// Основной метод алгоритма Банкера
void BankersAlgorithm::method()
{
    vector<thread> threads; // Вектор для хранения потоков
    int flag;

    while (true)
    {
        flag = 0;
        for (int i = 0; i < nop; i++)
        {
            threads.emplace_back(&BankersAlgorithm::check_process, this, i); // Создаем потоки для проверки процессов
        }

        for (auto &th : threads)
        {
            th.join(); // Ждем завершения всех потоков
        }
        threads.clear();

        for (int j = 0; j < nor; j++)
        {
            if (avail[j] != work[j]) // Проверяем, изменились ли доступные ресурсы
                flag = 1;
        }

        for (int j = 0; j < nor; j++)
        {
            work[j] = avail[j]; // Обновляем рабочую копию доступных ресурсов
        }

        if (flag == 0)
            break;
    }
}

// Метод поиска процесса, который может быть завершен
int BankersAlgorithm::search(int i)
{
    for (int j = 0; j < nor; j++)
    {
        if (n[i][j] > avail[j]) // Если потребность процесса больше доступных ресурсов
            return -1;
    }
    return 0;
}

// Функция для отображения результатов
void BankersAlgorithm::display()
{
    int i, j;
    cout << endl << "ВЫХОДНЫЕ ДАННЫЕ:";
    cout << endl << "=================";
    cout << endl << "ПРОЦЕСС\t\tВЫДЕЛЕНО\tМАКСИМУМ\tПОТРЕБНОСТЬ";
    for (i = 0; i < nop; i++)
    {
        cout << "\nP" << i + 1 << "\t\t";
        for (j = 0; j < nor; j++)
        {
            cout << al[i][j] << "  "; // Выделенные ресурсы
        }
        cout << "\t";
        for (j = 0; j < nor; j++)
        {
            cout << m[i][j] << "  "; // Максимальные ресурсы
        }
        cout << "\t";
        for (j = 0; j < nor; j++)
        {
            cout << n[i][j] << "  "; // Требуемые ресурсы
        }
    }
    cout << "\nПоследовательность безопасных процессов: \n";
    for (i = 0; i < k; i++)
    {
        int temp = result[i] + 1;
        cout << "P" << temp << " "; // Последовательность безопасных процессов
    }
    cout << "\nПоследовательность небезопасных процессов: \n";
    int flg = 0;
    for (i = 0; i < nop; i++)
    {
        if (finish[i] == 0)
        {
            flg = 1; // Есть небезопасные процессы
        }
        cout << "P" << i << " ";
    }
    cout << endl << "РЕЗУЛЬТАТ:";
    cout << endl << "===========";
    if (flg == 1)
        cout << endl << "Система не находится в безопасном состоянии, возможен дедлок!"; 
    else
        cout << endl << "Система находится в безопасном состоянии, дедлок не произойдет!"; 
}

int main()
{
    cout << "АЛГОРИТМ БАНКИРА ДЛЯ ПРЕДОТВРАЩЕНИЯ ДЕДЛОКА" << endl;
    BankersAlgorithm B;
    B.input();
    B.method();
    B.display();

    return 0;
}