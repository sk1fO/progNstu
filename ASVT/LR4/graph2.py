import matplotlib.pyplot as plt
import numpy as np


def plot_performance_comparison():
    """График сравнения производительности PThreads и OpenMP"""
    # Пример данных (замените на реальные результаты)
    threads = [1, 2, 4, 8, 16]
    pthreads_times = [58.9, 31.0, 14.2, 8.2, 6.2]  # В секундах
    openmp_times = [55.5, 29.6, 14.9, 8.6, 7.5]    # В секундах
    
    plt.figure(figsize=(10, 6))
    plt.plot(threads, pthreads_times, 'o-', label='PThreads')
    plt.plot(threads, openmp_times, 's-', label='OpenMP')
    
    plt.title('Сравнение производительности PThreads и OpenMP\n(Матрица 5000×5000)')
    plt.xlabel('Количество потоков')
    plt.ylabel('Время выполнения (сек)')
    plt.xticks(threads)
    plt.grid(True, linestyle='--', alpha=0.7)
    plt.legend()
    plt.savefig('performance_comparison.png', dpi=300, bbox_inches='tight')
    plt.show()

def plot_speedup():
    """График ускорения в зависимости от количества потоков"""
    threads = [1, 2, 4, 8, 16]
    speedup_pthreads = [1.0, 58.9/31, 58.9/14.2, 58.9/8.2, 58.9/6.2]
    speedup_openmp = [1.0, 58.9/29.6, 58.9/14.9, 58.9/8.6, 58.9/7.5]
    
    plt.figure(figsize=(10, 6))
    plt.plot(threads, speedup_pthreads, 'o-', label='PThreads')
    plt.plot(threads, speedup_openmp, 's-', label='OpenMP')
    plt.plot(threads, threads, '--', color='gray', label='Идеальное ускорение')
    
    plt.title('Коэффициент ускорения в зависимости от количества потоков\n(Матрица 5000×5000)')
    plt.xlabel('Количество потоков')
    plt.ylabel('Коэффициент ускорения')
    plt.xticks(threads)
    plt.grid(True, linestyle='--', alpha=0.7)
    plt.legend()
    plt.savefig('speedup_comparison.png', dpi=300, bbox_inches='tight')
    plt.show()

def plot_matrix_size_impact():
    """График зависимости времени выполнения от размера матрицы"""
    sizes = [1000, 2000, 3000, 4000, 5000]
    times_1thread = [0.217, 3.3, 12.0, 28.1, 55.1]     # В секундах
    times_8threads = [0.035, 0.395, 1.6, 3.1, 6.1]     # В секундах
    
    plt.figure(figsize=(10, 6))
    plt.plot(sizes, times_1thread, 'o-', label='1 поток')
    plt.plot(sizes, times_8threads, 's-', label='16 потоков (PThreads)')
    
    plt.title('Зависимость времени выполнения от размера матрицы')
    plt.xlabel('Размер матрицы (N×N)')
    plt.ylabel('Время выполнения (сек)')
    plt.xticks(sizes)
    plt.grid(True, linestyle='--', alpha=0.7)
    plt.legend()
    plt.savefig('matrix_size_impact.png', dpi=300, bbox_inches='tight')
    plt.show()

def plot_optimal_threads():
    """График для определения оптимального числа потоков"""
    threads = [1, 2, 4, 6, 8, 10, 12, 14, 16]
    times = [58.9, 31.0, 14.2, 9.9, 8.2, 6.6, 5.4, 4.9, 6.3]  # В секундах
    
    plt.figure(figsize=(10, 6))
    plt.plot(threads, times, 'o-')
    
    plt.title('Определение оптимального числа потоков\n(Матрица 5000×5000, OpenMP)')
    plt.xlabel('Количество потоков')
    plt.ylabel('Время выполнения (сек)')
    plt.xticks(threads)
    plt.grid(True, linestyle='--', alpha=0.7)
    
    # Пометка оптимального значения
    optimal = min(zip(threads, times), key=lambda x: x[1])
    plt.axvline(x=optimal[0], color='r', linestyle='--', alpha=0.5)
    plt.text(optimal[0]+0.5, optimal[1]+10, f'Оптимум: {optimal[0]} потоков', color='r')
    
    plt.savefig('optimal_threads.png', dpi=300, bbox_inches='tight')
    plt.show()

def main():
    print("Построение графиков для лабораторной работы №4...")
    plot_performance_comparison()
    plot_speedup()
    plot_matrix_size_impact()
    plot_optimal_threads()
    print("Графики сохранены в текущую директорию")

if __name__ == "__main__":
    main()