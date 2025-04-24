#include <iostream>
#include <vector>
#include <random>
#include <chrono>
#include <thread>
#include <pthread.h>
#include <omp.h>

// Структура для передачи данных в поток (PThreads)
struct ThreadData {
    const std::vector<std::vector<double>>* A;
    const std::vector<std::vector<double>>* B;
    std::vector<std::vector<double>>* C;
    int start_row;
    int end_row;
};

// Функция для потока (PThreads)
void* matrix_multiply_thread(void* arg) {
    ThreadData* data = static_cast<ThreadData*>(arg);
    const auto& A = *data->A;
    const auto& B = *data->B;
    auto& C = *data->C;
    
    int n = A.size();
    for (int i = data->start_row; i < data->end_row; ++i) {
        for (int k = 0; k < n; ++k) {
            for (int j = 0; j < n; ++j) {
                C[i][j] += A[i][k] * B[k][j];
            }
        }
    }
    return nullptr;
}

// Многопоточное умножение с PThreads
void DGEMM_pthreads(const std::vector<std::vector<double>>& A, 
                   const std::vector<std::vector<double>>& B, 
                   std::vector<std::vector<double>>& C, 
                   int num_threads) {
    int n = A.size();
    std::vector<pthread_t> threads(num_threads);
    std::vector<ThreadData> thread_data(num_threads);
    
    int rows_per_thread = n / num_threads;
    int extra_rows = n % num_threads;
    
    int current_row = 0;
    for (int i = 0; i < num_threads; ++i) {
        thread_data[i].A = &A;
        thread_data[i].B = &B;
        thread_data[i].C = &C;
        thread_data[i].start_row = current_row;
        
        int rows = rows_per_thread + (i < extra_rows ? 1 : 0);
        thread_data[i].end_row = current_row + rows;
        current_row += rows;
        
        pthread_create(&threads[i], nullptr, matrix_multiply_thread, &thread_data[i]);
    }
    
    for (int i = 0; i < num_threads; ++i) {
        pthread_join(threads[i], nullptr);
    }
}

// Многопоточное умножение с OpenMP
void DGEMM_openmp(const std::vector<std::vector<double>>& A, 
                 const std::vector<std::vector<double>>& B, 
                 std::vector<std::vector<double>>& C) {
    int n = A.size();
    #pragma omp parallel for
    for (int i = 0; i < n; ++i) {
        for (int k = 0; k < n; ++k) {
            for (int j = 0; j < n; ++j) {
                C[i][j] += A[i][k] * B[k][j];
            }
        }
    }
}

// Генерация случайной матрицы
std::vector<std::vector<double>> generate_random_matrix(int n) {
    std::random_device rd;
    std::mt19937 gen(rd());
    std::uniform_real_distribution<double> dis(0.0, 1.0);

    std::vector<std::vector<double>> matrix(n, std::vector<double>(n));
    #pragma omp parallel for
    for (int i = 0; i < n; ++i) {
        for (int j = 0; j < n; ++j) {
            matrix[i][j] = dis(gen);
        }
    }
    return matrix;
}

// Тестирование производительности
void test_performance(int n, int num_threads) {
    // Генерация матриц
    auto A = generate_random_matrix(n);
    auto B = generate_random_matrix(n);
    std::vector<std::vector<double>> C(n, std::vector<double>(n, 0.0));
    
    // Тестирование PThreads
    auto start = std::chrono::high_resolution_clock::now();
    DGEMM_pthreads(A, B, C, num_threads);
    auto end = std::chrono::high_resolution_clock::now();
    auto duration = std::chrono::duration_cast<std::chrono::milliseconds>(end - start);
    std::cout << "PThreads (n=" << n << ", threads=" << num_threads << "): " 
              << duration.count() << " ms\n";
    
    // Тестирование OpenMP
    C = std::vector<std::vector<double>>(n, std::vector<double>(n, 0.0));
    omp_set_num_threads(num_threads);
    start = std::chrono::high_resolution_clock::now();
    DGEMM_openmp(A, B, C);
    end = std::chrono::high_resolution_clock::now();
    duration = std::chrono::duration_cast<std::chrono::milliseconds>(end - start);
    std::cout << "OpenMP (n=" << n << ", threads=" << num_threads << "): " 
              << duration.count() << " ms\n";
}

int main(int argc, char* argv[]) {
    if (argc < 3) {
        std::cerr << "Usage: " << argv[0] << " <matrix_size> <num_threads>" << std::endl;
        return 1;
    }
    
    int n = std::stoi(argv[1]);
    int num_threads = std::stoi(argv[2]);
    
    test_performance(n, num_threads);
    
    return 0;
}