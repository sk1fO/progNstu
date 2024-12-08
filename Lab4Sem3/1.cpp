#include <iostream>
#include <vector>
#include <thread>
#include <mutex>
#include <condition_variable>
#include <semaphore>
#include <barrier>
#include <atomic>
#include <random>
#include <chrono>

// Функция для гонки с использованием различных примитивов синхронизации
void race_with_mutex(int thread_id, int race_length, std::mutex& mtx, std::vector<char>& track) {
    for (int i = 0; i < race_length; ++i) {
        std::lock_guard<std::mutex> lock(mtx);
        track.push_back('T' + thread_id);
    }
}

void race_with_semaphore(int thread_id, int race_length, std::counting_semaphore<1>& sem, std::vector<char>& track) {
    for (int i = 0; i < race_length; ++i) {
        sem.acquire();
        track.push_back('T' + thread_id);
        sem.release();
    }
}

void race_with_semaphore_slim(int thread_id, int race_length, std::binary_semaphore& sem, std::vector<char>& track) {
    for (int i = 0; i < race_length; ++i) {
        sem.acquire();
        track.push_back('T' + thread_id);
        sem.release();
    }
}

void race_with_barrier(int thread_id, int race_length, std::barrier& b, std::vector<char>& track) {
    for (int i = 0; i < race_length; ++i) {
        track.push_back('T' + thread_id);
        b.arrive_and_wait();
    }
}

void race_with_spinlock(int thread_id, int race_length, std::atomic<bool>& lock, std::vector<char>& track) {
    for (int i = 0; i < race_length; ++i) {
        while (lock.exchange(true, std::memory_order_acquire)) {
            std::this_thread::yield();
        }
        track.push_back('T' + thread_id);
        lock.store(false, std::memory_order_release);
    }
}

void race_with_spinwait(int thread_id, int race_length, std::atomic<bool>& lock, std::vector<char>& track) {
    for (int i = 0; i < race_length; ++i) {
        while (lock.exchange(true, std::memory_order_acquire)) {
            std::this_thread::yield();
        }
        track.push_back('T' + thread_id);
        lock.store(false, std::memory_order_release);
    }
}

void race_with_monitor(int thread_id, int race_length, std::mutex& mtx, std::condition_variable& cv, std::vector<char>& track, std::atomic<bool>& ready) {
    for (int i = 0; i < race_length; ++i) {
        std::unique_lock<std::mutex> lock(mtx);
        cv.wait(lock, [&ready] { return ready.load(); });
        track.push_back('T' + thread_id);
        ready.store(false);
        cv.notify_all();
    }
}

void run_race(int num_threads, int race_length, int sync_type) {
    std::vector<std::thread> threads;
    std::vector<char> track;
    std::mutex mtx;
    std::condition_variable cv;
    std::atomic<bool> ready(false);

    switch (sync_type) {
        case 1: { // Mutex
            for (int i = 0; i < num_threads; ++i) {
                threads.emplace_back(race_with_mutex, i, race_length, std::ref(mtx), std::ref(track));
            }
            break;
        }
        case 2: { // Semaphore
            std::counting_semaphore<1> sem(1);
            for (int i = 0; i < num_threads; ++i) {
                threads.emplace_back(race_with_semaphore, i, race_length, std::ref(sem), std::ref(track));
            }
            break;
        }
        case 3: { // SemaphoreSlim
            std::binary_semaphore sem(1);
            for (int i = 0; i < num_threads; ++i) {
                threads.emplace_back(race_with_semaphore_slim, i, race_length, std::ref(sem), std::ref(track));
            }
            break;
        }
        case 4: { // Barrier
            std::barrier b(num_threads);
            for (int i = 0; i < num_threads; ++i) {
                threads.emplace_back(race_with_barrier, i, race_length, std::ref(b), std::ref(track));
            }
            break;
        }
        case 5: { // SpinLock
            std::atomic<bool> lock(false);
            for (int i = 0; i < num_threads; ++i) {
                threads.emplace_back(race_with_spinlock, i, race_length, std::ref(lock), std::ref(track));
            }
            break;
        }
        case 6: { // SpinWait
            std::atomic<bool> lock(false);
            for (int i = 0; i < num_threads; ++i) {
                threads.emplace_back(race_with_spinwait, i, race_length, std::ref(lock), std::ref(track));
            }
            break;
        }
        case 7: { // Monitor
            for (int i = 0; i < num_threads; ++i) {
                threads.emplace_back(race_with_monitor, i, race_length, std::ref(mtx), std::ref(cv), std::ref(track), std::ref(ready));
            }
            break;
        }
    }

    for (auto& t : threads) {
        t.join();
    }

    std::cout << "Race results: ";
    for (char c : track) {
        std::cout << c << " ";
    }
    std::cout << std::endl;
}

int main() {
    int num_threads, race_length, sync_type;

    std::cout << "Enter the number of threads: ";
    std::cin >> num_threads;

    std::cout << "Enter the length of the race: ";
    std::cin >> race_length;

    std::cout << "Choose synchronization primitive:\n"
              << "1. Mutex\n"
              << "2. Semaphore\n"
              << "3. SemaphoreSlim\n"
              << "4. Barrier\n"
              << "5. SpinLock\n"
              << "6. SpinWait\n"
              << "7. Monitor\n"
              << "Your choice: ";
    std::cin >> sync_type;

    run_race(num_threads, race_length, sync_type);

    return 0;
}