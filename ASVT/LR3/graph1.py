import matplotlib.pyplot as plt
import numpy as np

def plot_time_comparison():
    """График сравнения времени выполнения разных реализаций"""
    # Пример данных (замените на реальные результаты)
    sizes = [1000, 2000, 3000, 4000, 5000]
    base_times = [1.56, 17.8, 61.1, 178.2, 382.2]      # В секундах
    opt1_times = [0.256, 2.7, 11.1, 24.4, 46.7]        # DGEMM_opt_1
    opt2_times = [0.303, 1.3, 4.1, 9.7, 19.2]          # DGEMM_opt_2 (block=64)
    opt3_times = [0.217, 3.4, 13.3, 31.3, 60.9]        # DGEMM_opt_3
    
    plt.figure(figsize=(10, 6))
    plt.plot(sizes, base_times, 'o-', label='Базовая реализация')
    plt.plot(sizes, opt1_times, 's-', label='Оптимизация 1 (построчная)')
    plt.plot(sizes, opt2_times, 'D-', label='Оптимизация 2 (блочная)')
    plt.plot(sizes, opt3_times, '^-', label='Оптимизация 3 (векторизация)')
    
    plt.title('Сравнение времени выполнения разных реализаций DGEMM')
    plt.xlabel('Размер матрицы (N×N)')
    plt.ylabel('Время выполнения (сек)')
    plt.xticks(sizes)
    plt.grid(True, linestyle='--', alpha=0.7)
    plt.legend()
    plt.savefig('time_comparison.png', dpi=300, bbox_inches='tight')
    plt.show()

def plot_speedup():
    """График ускорения оптимизированных версий"""
    sizes = [1000, 2000, 3000, 4000, 5000]
    speedup_opt1 = [1.56/0.26, 17.8/2.7, 61.1/11.1, 178.2/24.4, 382.2/46.7]
    speedup_opt2 = [1.56/0.30, 17.8/1.3, 61.1/4.10, 178.2/9.70, 382.2/19.2]
    speedup_opt3 = [1.56/0.22, 17.8/3.4, 61.1/13.3, 178.2/31.3, 382.2/60.9]
    
    plt.figure(figsize=(10, 6))
    plt.plot(sizes, speedup_opt1, 's-', label='Оптимизация 1 (построчная)')
    plt.plot(sizes, speedup_opt2, 'D-', label='Оптимизация 2 (блочная)')
    plt.plot(sizes, speedup_opt3, '^-', label='Оптимизация 3 (векторизация)')
    
    plt.title('Ускорение оптимизированных версий относительно базовой')
    plt.xlabel('Размер матрицы (N×N)')
    plt.ylabel('Коэффициент ускорения')
    plt.xticks(sizes)
    plt.grid(True, linestyle='--', alpha=0.7)
    plt.legend()
    plt.savefig('speedup.png', dpi=300, bbox_inches='tight')
    plt.show()

def plot_block_size_impact():
    """График зависимости времени от размера блока (для DGEMM_opt_2)"""
    block_sizes = [16, 32, 64, 128, 256]
    times_1000 = [0.309, 0.272, 0.294, 0.231, 0.23]    # Для матрицы 1000×1000
    times_2000 = [2.3, 1.5, 1.2, 0.968, 0.89]       # Для матрицы 2000×2000
    
    plt.figure(figsize=(10, 6))
    plt.plot(block_sizes, times_1000, 'o-', label='Матрица 1000×1000')
    plt.plot(block_sizes, times_2000, 's-', label='Матрица 2000×2000')
    
    plt.title('Зависимость времени выполнения от размера блока\n(DGEMM_opt_2)')
    plt.xlabel('Размер блока')
    plt.ylabel('Время выполнения (сек)')
    plt.xticks(block_sizes)
    plt.grid(True, linestyle='--', alpha=0.7)
    plt.legend()
    
    # Пометка оптимальных значений
    opt_1000 = block_sizes[np.argmin(times_1000)]
    opt_2000 = block_sizes[np.argmin(times_2000)]
    plt.axvline(x=opt_1000, color='b', linestyle='--', alpha=0.3)
    plt.axvline(x=opt_2000, color='g', linestyle='--', alpha=0.3)
    plt.text(opt_1000+5, min(times_1000)+0.1, f'Оптимум: {opt_1000}', color='b')
    plt.text(opt_2000+5, min(times_2000)+0.3, f'Оптимум: {opt_2000}', color='g')
    
    plt.savefig('block_size_impact.png', dpi=300, bbox_inches='tight')
    plt.show()

def plot_cache_misses():
    """График количества промахов кэша"""
    implementations = ['Базовая', 'Опт.1', 'Опт.2', 'Опт.3']
    cache_misses = [10.4e8, 1.23e8, 1.48e8, 0.54e8]  # Примерные значения
    
    plt.figure(figsize=(10, 6))
    bars = plt.bar(implementations, cache_misses, color=['#1f77b4', '#ff7f0e', '#2ca02c', '#d62728'])
    
    plt.title('Cache-misses\n(Матрица 2000×2000)')
    plt.ylabel('Cache-misses')
    plt.grid(True, axis='y', linestyle='--', alpha=0.7)
    
    # Добавление значений над столбцами
    for bar in bars:
        height = bar.get_height()
        plt.text(bar.get_x() + bar.get_width()/2., height,
                 f'{height/1e6:.1f}M',
                 ha='center', va='bottom')
    
    plt.savefig('cache_misses.png', dpi=300, bbox_inches='tight')
    plt.show()

def main():
    print("Построение графиков для лабораторной работы №3...")
    plot_time_comparison()
    plot_speedup()
    plot_block_size_impact()
    plot_cache_misses()
    print("Графики сохранены в текущую директорию")

if __name__ == "__main__":
    main()