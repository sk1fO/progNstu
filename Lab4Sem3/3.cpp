#include <iostream>
#include <thread>
#include <mutex>
#include <atomic>
#include <chrono>
#include <semaphore.h>

using namespace std;

mutex console_mutex;

atomic<int> sharedResource(0); // общий ресурс для чтения и записи
atomic<bool> exitRequested(false); // флаг для запроса выхода

sem_t semaphore;

void Writer(int id) {
    while (!exitRequested) {
        sem_wait(&semaphore); // ожидание доступа к ресурсу
        {
            lock_guard<mutex> lock(console_mutex); // захват мьютекса для вывода в консоль
            cout << "Писатель " << id << " пишет: " << ++sharedResource << endl;
        }
        sem_post(&semaphore); // освобождение ресурса
        this_thread::sleep_for(chrono::milliseconds(100));    // задержка для имитации работы
    }
}

void Reader(int id) {
    while (!exitRequested) {
        sem_wait(&semaphore);
        {
            lock_guard<mutex> lock(console_mutex);
            cout << "Читатель " << id << " читает: " << sharedResource << endl;
        }
        sem_post(&semaphore);
        this_thread::sleep_for(chrono::milliseconds(50));
    }
}

int main() {

    sem_init(&semaphore, 0, 1); //0 - указывает, что семафор используется для потоков
    //1 - устанавливает начальное значение семафора в 1, что означает, что изначально он разрешает доступ (семафор открыт)

    thread readers[20]; // создание массивов потоков для читателей и писателей
    thread writers[20];

    for (int i = 0; i < 20; ++i) { // запуск потоков
        readers[i] = thread(Reader, i);
        writers[i] = thread(Writer, i);
    }

    this_thread::sleep_for(chrono::milliseconds(1000));
    exitRequested = true; // установка флага выхода

    for (auto& reader : readers) { // ожидание завершения потоков
        reader.join();
    }

    for (auto& writer : writers) {
        writer.join();
    }

    cout << "Главный поток завершил работу" << endl;
    sem_destroy(&semaphore); // уничтожение семафора
    return 0;
}
