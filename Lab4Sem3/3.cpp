#include <iostream>
#include <vector>
#include <thread>
#include <mutex>
#include <condition_variable>
#include <random>
#include <chrono>

// Глобальные переменные
int NUM_PROCESSES;
int NUM_RESOURCES;
std::vector<int> available;
std::vector<std::vector<int>> max;
std::vector<std::vector<int>> allocation;
std::vector<std::vector<int>> need;

std::mutex mtx;
std::condition_variable cv;
bool safe_state = true;

// Функция для инициализации данных
void initialize() {
    // Запрос количества процессов и ресурсов
    std::cout << "Введите количество процессов: ";
    std::cin >> NUM_PROCESSES;
    std::cout << "Введите количество ресурсов: ";
    std::cin >> NUM_RESOURCES;

    // Инициализация векторов
    available.resize(NUM_RESOURCES, 0);
    max.resize(NUM_PROCESSES, std::vector<int>(NUM_RESOURCES, 0));
    allocation.resize(NUM_PROCESSES, std::vector<int>(NUM_RESOURCES, 0));
    need.resize(NUM_PROCESSES, std::vector<int>(NUM_RESOURCES, 0));

    // Автозаполнение случайными данными
    std::random_device rd;
    std::mt19937 gen(rd());
    std::uniform_int_distribution<> dis(1, 10);

    // Заполнение доступных ресурсов
    for (int i = 0; i < NUM_RESOURCES; ++i) {
        available[i] = dis(gen);
    }

    // Заполнение максимальных требований
    for (int i = 0; i < NUM_PROCESSES; ++i) {
        for (int j = 0; j < NUM_RESOURCES; ++j) {
            max[i][j] = dis(gen);
        }
    }

    // Заполнение текущих выделенных ресурсов
    for (int i = 0; i < NUM_PROCESSES; ++i) {
        for (int j = 0; j < NUM_RESOURCES; ++j) {
            allocation[i][j] = dis(gen) % (max[i][j] + 1);
        }
    }

    // Вычисление нужных ресурсов
    for (int i = 0; i < NUM_PROCESSES; ++i) {
        for (int j = 0; j < NUM_RESOURCES; ++j) {
            need[i][j] = max[i][j] - allocation[i][j];
        }
    }

    // Вывод инициализированных данных
    std::cout << "Доступно ресурсов: ";
    for (int i = 0; i < NUM_RESOURCES; ++i) {
        std::cout << available[i] << " ";
    }
    std::cout << std::endl;

    std::cout << "Максимально ресурсов: " << std::endl;
    for (int i = 0; i < NUM_PROCESSES; ++i) {
        for (int j = 0; j < NUM_RESOURCES; ++j) {
            std::cout << max[i][j] << " ";
        }
        std::cout << std::endl;
    }

    std::cout << "Распределение: " << std::endl;
    for (int i = 0; i < NUM_PROCESSES; ++i) {
        for (int j = 0; j < NUM_RESOURCES; ++j) {
            std::cout << allocation[i][j] << " ";
        }
        std::cout << std::endl;
    }

    std::cout << "Потребность: " << std::endl;
    for (int i = 0; i < NUM_PROCESSES; ++i) {
        for (int j = 0; j < NUM_RESOURCES; ++j) {
            std::cout << need[i][j] << " ";
        }
        std::cout << std::endl;
    }
}

// Функция для проверки безопасности
bool is_safe_state() {
    std::vector<bool> finish(NUM_PROCESSES, false);
    std::vector<int> work = available;
    int count = 0;

    while (count < NUM_PROCESSES) {
        bool found = false;
        for (int i = 0; i < NUM_PROCESSES; ++i) {
            if (!finish[i]) {
                bool can_allocate = true;
                for (int j = 0; j < NUM_RESOURCES; ++j) {
                    if (need[i][j] > work[j]) {
                        can_allocate = false;
                        break;
                    }
                }

                if (can_allocate) {
                    for (int j = 0; j < NUM_RESOURCES; ++j) {
                        work[j] += allocation[i][j];
                    }
                    finish[i] = true;
                    found = true;
                    count++;
                }
            }
        }

        if (!found) {
            return false;
        }
    }

    return true;
}

// Функция для запроса ресурсов
void request_resources(int process_id, const std::vector<int>& request) {
    std::unique_lock<std::mutex> lock(mtx);

    bool can_allocate = true;
    for (int i = 0; i < NUM_RESOURCES; ++i) {
        if (request[i] > need[process_id][i]) {
            std::cout << "Процесс " << process_id << " превысило максимально возможное требование." << std::endl;
            can_allocate = false;
            break;
        }
    }

    if (can_allocate) {
        for (int i = 0; i < NUM_RESOURCES; ++i) {
            if (request[i] > available[i]) {
                std::cout << "Процесс " << process_id << " необходимо подождать, ресурсов не хватает." << std::endl;
                can_allocate = false;
                break;
            }
        }
    }

    if (can_allocate) {
        for (int i = 0; i < NUM_RESOURCES; ++i) {
            available[i] -= request[i];
            allocation[process_id][i] += request[i];
            need[process_id][i] -= request[i];
        }

        if (is_safe_state()) {
            std::cout << "Процесс " << process_id << " ресурсы выделены." << std::endl;
        } else {
            std::cout << "Процесс " << process_id << " необходимо подождать, система не в безопасном состоянии." << std::endl;
            for (int i = 0; i < NUM_RESOURCES; ++i) {
                available[i] += request[i];
                allocation[process_id][i] -= request[i];
                need[process_id][i] += request[i];
            }
        }
    }

    cv.notify_all();
}

// Функция для каждого процесса
void process_function(int process_id) {
    std::vector<int> request(NUM_RESOURCES, 0);
    std::random_device rd;
    std::mt19937 gen(rd());
    std::uniform_int_distribution<> dis(0, 5);

    for (int i = 0; i < NUM_RESOURCES; ++i) {
        request[i] = dis(gen);
    }

    request_resources(process_id, request);
}

int main() {
    initialize();

    std::vector<std::thread> threads;
    for (int i = 0; i < NUM_PROCESSES; ++i) {
        threads.push_back(std::thread(process_function, i));
    }

    for (auto& t : threads) {
        t.join();
    }

    return 0;
}