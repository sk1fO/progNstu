import matplotlib.pyplot as plt
from matplotlib.patches import Circle, Rectangle, Arrow


def draw_event_driven():
    fig, ax = plt.subplots(figsize=(10, 6))
    ax.set_xlim(0, 10)
    ax.set_ylim(0, 6)
    ax.axis('off')

    # Компоненты
    components = [
        ("Кнопка 1", 1, 5, '#FFCCCC'),
        ("Кнопка 2", 1, 4, '#FFCCCC'),
        ("Кнопка 3", 1, 3, '#FFCCCC'),
        ("Брокер\n(Kafka/SQS)", 4, 4, '#CCFFCC'),
        ("Обработчик\nсобытий", 7, 5, '#CCCCFF'),
        ("База данных", 7, 3, '#FFFFCC'),
        ("Уведомления", 9, 4, '#FFCCFF')
    ]

    for text, x, y, color in components:
        if "Кнопка" in text:
            ax.add_patch(Circle((x, y), 0.5, color=color))
        else:
            ax.add_patch(Rectangle((x - 1, y - 0.5), 2, 1, color=color))
        ax.text(x, y, text, ha='center', va='center')

    # Связи
    arrows = [
        (1.5, 5, 3, 4), (1.5, 4, 3, 4), (1.5, 3, 3, 4),
        (5, 4, 6, 5), (5, 4, 6, 3),
        (8, 5, 8, 3), (8, 5, 9.5, 4)
    ]

    for x1, y1, x2, y2 in arrows:
        ax.annotate("", xy=(x2, y2), xytext=(x1, y1),
                    arrowprops=dict(arrowstyle="->", lw=1))

    plt.title("Event-Driven Архитектура", pad=20)
    plt.tight_layout()
    plt.savefig('event_driven.png', dpi=300)
    plt.close()


draw_event_driven()