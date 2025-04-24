#include <iostream>
#include <vector>
#include <random>
#include <chrono>
#include <iomanip>
#include <immintrin.h> // Для AVX инструкций

// Базовая реализация DGEMM
void DGEMM(const std::vector<std::vector<double>>& A, 
           const std::vector<std::vector<double>>& B, 
           std::vector<std::vector<double>>& C) {
    int n = A.size();
    for (int i = 0; i < n; ++i) {
        for (int j = 0; j < n; ++j) {
            for (int k = 0; k < n; ++k) {
                C[i][j] += A[i][k] * B[k][j];
            }
        }
    }
}

// Оптимизация 1: Построчный перебор
void DGEMM_opt_1(const std::vector<std::vector<double>>& A, 
                const std::vector<std::vector<double>>& B, 
                std::vector<std::vector<double>>& C) {
    int n = A.size();
    for (int i = 0; i < n; ++i) {
        for (int k = 0; k < n; ++k) {
            double temp = A[i][k];
            for (int j = 0; j < n; ++j) {
                C[i][j] += temp * B[k][j];
            }
        }
    }
}

// Оптимизация 2: Блочный перебор
void DGEMM_opt_2(const std::vector<std::vector<double>>& A, 
                const std::vector<std::vector<double>>& B, 
                std::vector<std::vector<double>>& C, 
                int block_size) {
    int n = A.size();
    for (int i = 0; i < n; i += block_size) {
        for (int j = 0; j < n; j += block_size) {
            for (int k = 0; k < n; k += block_size) {
                // Обрабатываем блок
                for (int ii = i; ii < std::min(i + block_size, n); ++ii) {
                    for (int kk = k; kk < std::min(k + block_size, n); ++kk) {
                        double temp = A[ii][kk];
                        for (int jj = j; jj < std::min(j + block_size, n); ++jj) {
                            C[ii][jj] += temp * B[kk][jj];
                        }
                    }
                }
            }
        }
    }
}

// Оптимизация 3: Векторизация с использованием AVX
void DGEMM_opt_3(const std::vector<std::vector<double>>& A, 
                const std::vector<std::vector<double>>& B, 
                std::vector<std::vector<double>>& C) {
    int n = A.size();
    for (int i = 0; i < n; ++i) {
        for (int k = 0; k < n; ++k) {
            __m256d a = _mm256_set1_pd(A[i][k]);
            for (int j = 0; j < n; j += 4) {
                if (j + 4 <= n) {
                    __m256d b = _mm256_loadu_pd(&B[k][j]);
                    __m256d c = _mm256_loadu_pd(&C[i][j]);
                    c = _mm256_fmadd_pd(a, b, c);
                    _mm256_storeu_pd(&C[i][j], c);
                } else {
                    // Обработка хвоста, если размер не кратен 4
                    for (int jj = j; jj < n; ++jj) {
                        C[i][jj] += A[i][k] * B[k][jj];
                    }
                }
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
    for (int i = 0; i < n; ++i) {
        for (int j = 0; j < n; ++j) {
            matrix[i][j] = dis(gen);
        }
    }
    return matrix;
}

// Тестирование производительности
void test_performance(int n, int block_size = 64) {
    // Генерация матриц
    auto A = generate_random_matrix(n);
    auto B = generate_random_matrix(n);
    std::vector<std::vector<double>> C(n, std::vector<double>(n, 0.0));
    
    // Тестирование базовой реализации
    auto start = std::chrono::high_resolution_clock::now();
    
    DGEMM(A, B, C);
    auto end = std::chrono::high_resolution_clock::now();
    auto duration = std::chrono::duration_cast<std::chrono::milliseconds>(end - start);
    std::cout << "DGEMM (n=" << n << "): " << duration.count() << " ms\n";
    
    // Тестирование оптимизации 1
    C = std::vector<std::vector<double>>(n, std::vector<double>(n, 0.0));
    start = std::chrono::high_resolution_clock::now();
    DGEMM_opt_1(A, B, C);
    auto end = std::chrono::high_resolution_clock::now();
    auto duration = std::chrono::duration_cast<std::chrono::milliseconds>(end - start);
    std::cout << "DGEMM_opt_1 (n=" << n << "): " << duration.count() << " ms\n";
    
    //Тестирование оптимизации 2
    C = std::vector<std::vector<double>>(n, std::vector<double>(n, 0.0));
    start = std::chrono::high_resolution_clock::now();
    DGEMM_opt_2(A, B, C, block_size);
    auto end = std::chrono::high_resolution_clock::now();
    auto duration = std::chrono::duration_cast<std::chrono::milliseconds>(end - start);
    std::cout << "DGEMM_opt_2 (n=" << n << ", block=" << block_size << "): " 
              << duration.count() << " ms\n";
    
    // Тестирование оптимизации 3 
    C = std::vector<std::vector<double>>(n, std::vector<double>(n, 0.0));
    start = std::chrono::high_resolution_clock::now();
    DGEMM_opt_3(A, B, C);
    auto end = std::chrono::high_resolution_clock::now();
    auto duration = std::chrono::duration_cast<std::chrono::milliseconds>(end - start);
    std::cout << "DGEMM_opt_3 (n=" << n << "): " << duration.count() << " ms\n";
    
}

int main(int argc, char* argv[]) {
    if (argc < 2) {
        std::cerr << "Usage: " << argv[0] << " <matrix_size> [block_size]" << std::endl;
        return 1;
    }
    
    int n = std::stoi(argv[1]);
    int block_size = 64; // Значение по умолчанию
    
    if (argc > 2) {
        block_size = std::stoi(argv[2]);
    }
    
    test_performance(n, block_size);
    
    return 0;
}