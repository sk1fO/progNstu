#include <iostream>
#include <vector>
#include <thread>
#include <mutex>
#include <semaphore>
#include <barrier>
#include <atomic>
#include <chrono>
#include <condition_variable>

using namespace std;
mutex console_mutex;

// создание семафора счетчика с начальным значением 4 (до четырех потоков смогут одновременно получить доступ к ресурсу, уменьшая счетчик при вызове acquire, и увеличивая при вызове release) 
// если все разрешения уже заняты, потоки будут блокироваться при вызове acquire
counting_semaphore sem(4);

atomic<int> counter(0); // атомарная переменная - это переменная, операции с которой выполняются без возможности прерывания другими потоками

atomic_flag spinlock = ATOMIC_FLAG_INIT; // инициализация атомарного флага для использования в качестве спин-блокировки
atomic_flag spinwait = ATOMIC_FLAG_INIT;

mutex cv_mutex; //  создание объекта мьютекса
condition_variable cv; // создание объекта условной переменной
atomic<bool> is_condition_met = false; // переменная изменяется из разных потоков

int threadCount = 0;
condition_variable barrierCond;

bool is_monitor_locked = false;

atomic<long long> total_execution_time(0); // Общее время выполнения всех потоков
atomic<int> completed_threads(0); // Количество потоков, завершивших выполнение

mutex total_execution_mutex;

void MonitorEnter() { // блокирует монитор, чтобы предотвратить доступ других потоков
    mutex m;
    unique_lock<mutex> lock(m);
    cv.wait(lock, [&]() { return !is_monitor_locked; });// Ожидаем, пока монитор не разблокируется
    is_monitor_locked = true;// Устанавливаем флаг блокировки монитора
}

void MonitorExit() { // разблокирует монитор, позволяя другим потокам получить доступ
    mutex m;
    unique_lock<mutex> lock(m);
    is_monitor_locked = false; // Сбрасываем флаг блокировки монитора
    cv.notify_one();// Уведомление одного потока, ожидающего блокировку монитора
}

char generateRandomChar() { // функция для генерации случайных символов
    return static_cast<char>('!' + rand() % 94); // добавление случайного числа к коду символа '!', чтобы создавать случайные символы из диапазона от '!' до '~'
}

void race(int thread_id, int choice, int raceLength, int numberThreads) { // функция, выполняемая в каждом потоке

    auto start_time = chrono::high_resolution_clock::now();
    barrier barrier(numberThreads);
    mutex m;

    switch (choice) {
    case 1: //поток блокируется с помощью мьютекса, чтобы гарантировать, что только один поток может выполнять вывод на экран
        for (int i = 0; i < raceLength; ++i) {
            m.lock();
            cout << generateRandomChar();
            m.unlock();
            // после вывода символа поток разблокируется, затем приостанавливается на 100 миллисекунд
            this_thread::sleep_for(chrono::milliseconds(100));
        }
        break;

    case 2:
        for (int i = 0; i < raceLength; ++i) {
            sem.acquire();  // уменьшение счетчика семафора на 1, блокировка, если счетчик равен 0
            cout << generateRandomChar();
            sem.release(); // увеличение счетчика после использования ресурса
            this_thread::sleep_for(chrono::milliseconds(100));
        }
        break;

    case 3: // управление доступом к ресурсу с помощью атомарных операций, таких как увеличение и уменьшение значения переменной
        for (int i = 0; i < raceLength; ++i) {
            while (counter.fetch_add(1) >= 1) { // если значение counter оказывается больше или равно 1, происходит отмена увеличения 
                counter.fetch_sub(1);
                this_thread::yield(); // вызов this_thread::yield(), который позволяет другим потокам выполниться
            }
            cout << generateRandomChar();
            counter.fetch_sub(1); // уменьшение счетчика на 1
            this_thread::sleep_for(chrono::milliseconds(100));
        }
        break;

    case 4:  // используется для того, чтобы поток достиг барьера и остановил свое выполнение, ожидая остальных потоков до тех пор, пока все они не достигнут этой точки синхронизации
        for (int i = 0; i < raceLength; ++i) { 
            unique_lock<mutex> lock(m); 
            threadCount++; // Увеличение счетчика потоков
            if (threadCount == numberThreads) { // Если достигнуто заданное количество потоков, то происходит сброс счетчика и уведомление всех потоков о достижении барьера
                threadCount = 0; 
                barrierCond.notify_all(); 
                this_thread::sleep_for(chrono::milliseconds(100)); 
            }
            else {
                barrierCond.wait(lock, [] { return threadCount == 0; }); // Ожидание условия с использованием мьютекса
            }
            cout << generateRandomChar(); 
        }
        break; 

    case 5:
        for (int i = 0; i < raceLength; ++i) {
            // ожидание доступности ресурса
            while (spinlock.test_and_set(memory_order_acquire)) { // проверка текущего значения атомарного флага и установка его в true, если оно было false
                // аргумент memory_order_acquire обеспечивает правильный порядок операций с памятью для синхронизации
                // this_thread::yield(); // передача управления другим потокам
            }
            cout << generateRandomChar(); // использование ресурса
            spinlock.clear(memory_order_release); // снятие флага после использования ресурса
            this_thread::sleep_for(chrono::milliseconds(100));
        }
        break;

    case 6:
        for (int i = 0; i < raceLength; i++) { 
            this_thread::sleep_for(chrono::milliseconds(100)); 

            while (spinwait.test_and_set(memory_order_acquire)) { // Цикл ожидания освобождения ресурса
                this_thread::yield(); 
            }
            cout << generateRandomChar(); 
            spinwait.clear(memory_order_release); // Освобождение ресурса
        }
        break; 
        
    case 7:
        for (int i = 0; i < raceLength; i++) {
            MonitorEnter(); // блокирует монитор, чтобы предотвратить доступ других потоков
            cout << generateRandomChar();
            MonitorExit(); // разблокирует монитор, позволяя другим потокам получить доступ
            std::this_thread::sleep_for(std::chrono::milliseconds(100));
        }
        break;

    default:
        cout << "Некорректный выбор команды" << endl;
        break;
    }
    cout << endl;

    auto end_time = chrono::high_resolution_clock::now(); // получение текущего времени
    auto duration = chrono::duration_cast<chrono::milliseconds>(end_time - start_time);
    //unique_lock<mutex> lock(console_mutex);
    cout << "Поток " << thread_id << ": Время выполнения - " << duration.count() << " миллисекунд." << endl; // вывод информации о времени выполнения потока 

    // Обновление общего времени выполнения и количества завершенных потоков
    {
        lock_guard<mutex> lock(total_execution_mutex);
        total_execution_time += duration.count();
        completed_threads++;
    }

    // Ожидание завершения всех потоков
    if (completed_threads == numberThreads) {
        // Вычисление и вывод среднего времени выполнения
        long long average_execution_time = total_execution_time / numberThreads;
        unique_lock<mutex> lock(console_mutex);
        cout << "Среднее время выполнения всех потоков: " << average_execution_time << " миллисекунд." << endl;
    }
}

int main() {
    int numberThreads;
    int raceLength;

    srand(static_cast<unsigned int>(time(nullptr)));

    cout << "Введите количество потоков: ";
    cin >> numberThreads;
    cout << "Введите длину гонки: ";
    cin >> raceLength;

    int choice;
    do {
        cout << "Возможные примитивы синхронизации: " << endl;
        cout << "1. Mutexes" << endl;
        cout << "2. Semaphore" << endl;
        cout << "3. SemaphoreSlim" << endl;
        cout << "4. Barrier" << endl;
        cout << "5. SpinLock" << endl;
        cout << "6. SpinWait" << endl;
        cout << "7. Monitor" << endl;
        cout << "Выберите примитив синхронизации из предложенных: ";
        cin >> choice;
        if (choice < 1 || choice > 7) {
            cout << "Неверный ввод. Пожалуйста, выберите снова." << endl;
        }
    } while (choice < 1 || choice > 7);

    vector<thread> threads;

    // создание соответствующего количества потоков, каждый из которых вызывает функцию "race" с заданными параметрами
    for (int i = 0; i < numberThreads; ++i) {
        threads.emplace_back(race, i, choice, raceLength, numberThreads);
        // emplace_back - это метод вектора в C++, который создает новый объект прямо в конце вектора
    }

    for (auto& thread : threads) { // цикл "join" ожидает завершения всех созданных потоков перед завершением основной программы
        thread.join();
    }
    cout << endl;
    return 0;
}
