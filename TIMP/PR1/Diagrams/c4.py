import matplotlib.pyplot as plt
from matplotlib.patches import Rectangle, Circle, Arrow

def draw_c4_diagram():
    fig, ax = plt.subplots(figsize=(10, 8))
    ax.set_xlim(0, 10)
    ax.set_ylim(0, 10)
    ax.axis('off')
    
    # Внешние акторы
    ax.add_patch(Circle((1, 8), 0.5, fill=True, color='lightblue'))
    ax.text(1, 8, "Оператор", ha='center', va='center')
    
    ax.add_patch(Circle((1, 5), 0.5, fill=True, color='lightgreen'))
    ax.text(1, 5, "Патруль", ha='center', va='center')
    
    # Система
    ax.add_patch(Rectangle((3, 3), 6, 5, fill=False, lw=2))
    ax.text(6, 8.2, "Система мониторинга", ha='center', fontsize=12)
    
    # Компоненты системы
    ax.add_patch(Rectangle((4, 6), 2, 1.5, fill=True, color='#FFDDDD'))
    ax.text(5, 6.75, "Веб-интерфейс\n(React)", ha='center')
    
    ax.add_patch(Rectangle((4, 4), 2, 1.5, fill=True, color='#DDFFDD'))
    ax.text(5, 4.75, "Бэкенд-сервис\n(Go)", ha='center')
    
    ax.add_patch(Rectangle((7, 4), 2, 1.5, fill=True, color='#DDDDFF'))
    ax.text(8, 4.75, "База данных\n(PostgreSQL)", ha='center')
    
    # Связи
    ax.annotate("", xy=(1.5, 8), xytext=(3, 6.75),
                arrowprops=dict(arrowstyle="->"))
    ax.annotate("", xy=(5, 5.5), xytext=(5, 4.75),
                arrowprops=dict(arrowstyle="->"))
    ax.annotate("", xy=(7, 4.75), xytext=(5, 4.75),
                arrowprops=dict(arrowstyle="->"))
    ax.annotate("", xy=(5, 4), xytext=(1.5, 5),
                arrowprops=dict(arrowstyle="->"))
    
    plt.title("C4-диаграмма: Контейнеры системы", pad=20)
    plt.tight_layout()
    plt.savefig('c4_diagram.png', dpi=300)
    plt.close()

draw_c4_diagram()