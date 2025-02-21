#include <iostream>
#include <vector>
#include <string>
#include <thread>
#include <chrono>
#include <mutex>
#include <random>
#include <iomanip>
#include <sstream>

using namespace std;

// Структура для хранения информации об абонементе
struct Subscription {
    string name;
    chrono::system_clock::time_point purchase_date;
    chrono::system_clock::time_point expiration_date;
};

// Функция для преобразования строки в формате YYYY-MM-DD в time_point
chrono::system_clock::time_point string_to_time_point(const string& date_str) {
    tm tm = {};
    istringstream ss(date_str);
    ss >> get_time(&tm, "%Y-%m-%d");
    if (ss.fail()) {
        throw runtime_error("Invalid date format. Please use YYYY-MM-DD.");
    }
    return chrono::system_clock::from_time_t(mktime(&tm));
}

// Функция для генерации случайной даты в заданном диапазоне
chrono::system_clock::time_point random_date(chrono::system_clock::time_point start, chrono::system_clock::time_point end) {
    static random_device rd;
    static mt19937 gen(rd());
    uniform_int_distribution<> dist(0, chrono::duration_cast<chrono::seconds>(end - start).count());
    return start + chrono::seconds(dist(gen));
}

// Функция для генерации случайных данных
vector<Subscription> generate_data(int size, chrono::system_clock::time_point d1, chrono::system_clock::time_point d2) {
    vector<Subscription> data(size);
    for (int i = 0; i < size; ++i) {
        data[i].name = "Name" + to_string(i);
        data[i].purchase_date = random_date(d1, d2);
        data[i].expiration_date = data[i].purchase_date + chrono::hours(24 * 365); // Абонемент действует год
    }
    return data;
}

// Функция для проверки, истек ли срок действия абонемента
bool is_expired(const Subscription& sub, chrono::system_clock::time_point current_date) {
    return sub.expiration_date < current_date;
}

// Функция обработки данных без использования многопоточности
void process_data_sequentially(const vector<Subscription>& data, chrono::system_clock::time_point d1, chrono::system_clock::time_point d2, chrono::system_clock::time_point current_date, vector<Subscription>& expired_subscriptions) {
    for (const auto& sub : data) {
        if (sub.purchase_date >= d1 && sub.purchase_date <= d2 && is_expired(sub, current_date)) {
            expired_subscriptions.push_back(sub);
        }
    }
}

// Функция обработки данных с использованием многопоточности
void process_data_multithreaded(const vector<Subscription>& data, chrono::system_clock::time_point d1, chrono::system_clock::time_point d2, chrono::system_clock::time_point current_date, int num_threads, vector<Subscription>& expired_subscriptions) {
    vector<thread> threads;
    mutex mutex;

    auto process_chunk = [&](int start, int end) {
        vector<Subscription> local_expired;
        for (int i = start; i < end; ++i) {
            const auto& sub = data[i];
            if (sub.purchase_date >= d1 && sub.purchase_date <= d2 && is_expired(sub, current_date)) {
                local_expired.push_back(sub);
            }
        }
        lock_guard<std::mutex> lock(mutex);
        expired_subscriptions.insert(expired_subscriptions.end(), local_expired.begin(), local_expired.end());
    };

    int chunk_size = data.size() / num_threads;
    for (int i = 0; i < num_threads; ++i) {
        int start = i * chunk_size;
        int end = (i == num_threads - 1) ? data.size() : start + chunk_size;
        threads.emplace_back(process_chunk, start, end);
    }

    for (auto& t : threads) {
        t.join();
    }
}

int main() {
    int array_size, num_threads;
    string d1_str, d2_str;

    cout << "Введите размер массива: ";
    cin >> array_size;

    cout << "Введите число потоков: ";
    cin >> num_threads;

    cout << "Введите дату начала (YYYY-MM-DD): ";
    cin >> d1_str;

    cout << "Введите дату окончания (YYYY-MM-DD): ";
    cin >> d2_str;

    try {
        chrono::system_clock::time_point d1 = string_to_time_point(d1_str);
        chrono::system_clock::time_point d2 = string_to_time_point(d2_str);
        chrono::system_clock::time_point current_date = chrono::system_clock::now();

        vector<Subscription> data = generate_data(array_size, d1, d2);

        // Обработка данных последовательно
        auto start = chrono::high_resolution_clock::now();
        vector<Subscription> expired_subscriptions_sequential;
        process_data_sequentially(data, d1, d2, current_date, expired_subscriptions_sequential);
        auto end = chrono::high_resolution_clock::now();
        chrono::duration<double> duration_sequential = end - start;
        cout << "Время выполниения (однопоточный режим): " << duration_sequential.count() << " секунд\n";

        // Обработка данных параллельно
        start = chrono::high_resolution_clock::now();
        vector<Subscription> expired_subscriptions_multithreaded;
        process_data_multithreaded(data, d1, d2, current_date, num_threads, expired_subscriptions_multithreaded);
        end = chrono::high_resolution_clock::now();
        chrono::duration<double> duration_multithreaded = end - start;
        cout << "Время выполниения (многопоточный режим): " << duration_multithreaded.count() << " секунд\n";

        // Вывод результатов
        cout << "Истекших абонементов (однопоточный режим): " << expired_subscriptions_sequential.size() << "\n";
        cout << "Истекших абонементов (многопоточный режим): " << expired_subscriptions_multithreaded.size() << "\n";
    } catch (const exception& e) {
        cerr << "Error: " << e.what() << endl;
    }

    return 0;
}