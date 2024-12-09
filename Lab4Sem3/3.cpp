#include <iostream>
#include <thread>
#include <vector>
#include <mutex>
#define MAX 20
using namespace std;

class Bankers
{
private:
    int al[MAX][MAX], m[MAX][MAX], n[MAX][MAX], avail[MAX]; // Матрицы текущих, максимума, нужды и доступные ресурсы
    int nop, nor, k, result[MAX], pnum, work[MAX], finish[MAX]; // Число процессов, ресурсов, индекс результата, результат, номер процесса, рабочий массив, завершенные процессы
    mutex mtx; // Мьютекс для синхронизации доступа к общим ресурсам

public:
    Bankers(); // Конструктор класса
    void input(); // Ввод данных
    void method(); // Основной метод алгоритма банкира
    int search(int); // Поиск подходящего процесса
    void display(); // Вывод результатов
    void process(int i); // Обработка процесса в отдельном потоке
};

Bankers::Bankers()
{
    k = 0;
    for (int i = 0; i < MAX; i++)
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
    // Инициализация всех значений нулями
}

void Bankers::input()
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
    cout << "Введите максимальные ресурсы, которые нужны каждому процессу: " << endl;
    for (i = 0; i < nop; i++)
    {
        cout << "\nПроцесс " << i;
        for (j = 0; j < nor; j++)
        {
            cout << "\nРесурс " << j << ": ";
            cin >> m[i][j];
            n[i][j] = m[i][j] - al[i][j]; // Вычисление нужных ресурсов
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
        finish[i] = 0; // Все процессы изначально не завершены
}

void Bankers::process(int i)
{
    unique_lock<mutex> lock(mtx); // Блокировка мьютекса для синхронизации
    if (finish[i] == 0)
    {
        pnum = search(i);
        if (pnum != -1)
        {
            result[k++] = i; // Добавление процесса в результат
            finish[i] = 1; // Пометка процесса как завершенного
            for (int j = 0; j < nor; j++)
            {
                avail[j] += al[i][j]; // Обновление доступных ресурсов
            }
        }
    }
    lock.unlock(); // Разблокировка мьютекса
}

int Bankers::search(int i)
{
    for (int j = 0; j < nor; j++)
        if (n[i][j] > avail[j])
            return -1; // Если нужные ресурсы превышают доступные, процесс не может быть выполнен
    return 0;
}

void Bankers::method()
{
    vector<thread> threads;
    for (int i = 0; i < nop; i++)
    {
        threads.emplace_back(&Bankers::process, this, i); // Создание потока для каждого процесса
    }

    for (auto &t : threads)
    {
        t.join(); // Ожидание завершения всех потоков
    }

    // Проверка, завершены ли все процессы
    bool flag = true;
    for (int i = 0; i < nop; i++)
    {
        if (finish[i] == 0)
        {
            flag = false;
            break;
        }
    }

    if (!flag)
    {
        cout << "Система не находится в безопасном состоянии, возможен deadlock!" << endl;
    }
    else
    {
        cout << "Система находится в безопасном состоянии, deadlock невозможен!" << endl;
    }
}

void Bankers::display()
{
    int i, j;
    cout << endl << "ВЫВОД:";
    cout << endl << "========";
    cout << endl << "ПРОЦЕСС\t     ВЫДЕЛЕННЫЕ\t     МАКСИМАЛЬНЫЕ\tНУЖНЫЕ";
    for (i = 0; i < nop; i++)
    {
        cout << "\nP" << i + 1 << "\t     ";
        for (j = 0; j < nor; j++)
        {
            cout << al[i][j] << "  ";
        }
        cout << "\t     ";
        for (j = 0; j < nor; j++)
        {
            cout << m[i][j] << "  ";
        }
        cout << "\t     ";
        for (j = 0; j < nor; j++)
        {
            cout << n[i][j] << "  ";
        }
    }
    cout << "\nПоследовательность безопасных процессов: \n";
    for (i = 0; i < k; i++)
    {
        int temp = result[i] + 1;
        cout << "P" << temp << " ";
    }
    cout << "\nПоследовательность небезопасных процессов: \n";
    int flg = 0;
    for (i = 0; i < nop; i++)
    {
        if (finish[i] == 0)
        {
            flg = 1;
        }
        cout << "P" << i << " ";
    }
    cout << endl << "РЕЗУЛЬТАТ:";
    cout << endl << "=======";
    if (flg == 1)
        cout << endl << "Система не находится в безопасном состоянии, возможен deadlock!";
    else
        cout << endl << "Система находится в безопасном состоянии, deadlock невозможен!";
}

int main()
{
    cout << " АЛГОРИТМ БАНКИРА ДЛЯ ПРЕДОТВРАЩЕНИЯ ДЕДЛОКА " << endl;
    Bankers B;
    B.input();
    B.method();
    B.display();
    cout << endl;
}