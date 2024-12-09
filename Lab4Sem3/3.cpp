#include <iostream>
#include <thread>
#include <vector>
#include <mutex>
#define MAX 20

using namespace std;

class BankersAlgorithm
{
private:
    int al[MAX][MAX], m[MAX][MAX], n[MAX][MAX], avail[MAX];
    int nop, nor, k, result[MAX], pnum, work[MAX], finish[MAX];
    mutex mtx;

public:
    BankersAlgorithm();
    void input();
    void method();
    int search(int);
    void display();
    void check_process(int i);
};

BankersAlgorithm::BankersAlgorithm()
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
}

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

void BankersAlgorithm::check_process(int i)
{
    if (finish[i] == 0)
    {
        pnum = search(i);
        if (pnum != -1)
        {
            mtx.lock();
            result[k++] = i;
            finish[i] = 1;
            for (int j = 0; j < nor; j++)
            {
                avail[j] = avail[j] + al[i][j];
            }
            mtx.unlock();
        }
    }
}

void BankersAlgorithm::method()
{
    vector<thread> threads;
    int flag;

    while (true)
    {
        flag = 0;
        for (int i = 0; i < nop; i++)
        {
            threads.emplace_back(&BankersAlgorithm::check_process, this, i);
        }

        for (auto &th : threads)
        {
            th.join();
        }
        threads.clear();

        for (int j = 0; j < nor; j++)
        {
            if (avail[j] != work[j])
                flag = 1;
        }

        for (int j = 0; j < nor; j++)
        {
            work[j] = avail[j];
        }

        if (flag == 0)
            break;
    }
}

int BankersAlgorithm::search(int i)
{
    for (int j = 0; j < nor; j++)
    {
        if (n[i][j] > avail[j])
            return -1;
    }
    return 0;
}

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
            cout << al[i][j] << "  ";
        }
        cout << "\t";
        for (j = 0; j < nor; j++)
        {
            cout << m[i][j] << "  ";
        }
        cout << "\t";
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