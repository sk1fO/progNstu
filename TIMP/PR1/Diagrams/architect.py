import matplotlib.pyplot as plt


def draw_system_architecture():
    fig, ax = plt.subplots(figsize=(12, 8))
    ax.set_xlim(0, 12)
    ax.set_ylim(0, 10)
    ax.axis('off')

    # Пользователи
    ax.add_patch(plt.Circle((1, 8), 0.5, color='lightblue'))
    ax.text(1, 8, "Операторы", ha='center', va='center')

    # Компоненты
    components = [
        ("Nginx\nБалансировщик", 4, 8, '#FFDDDD'),
        ("Сервер 1", 4, 6, '#DDFFDD'),
        ("Сервер 2", 7, 6, '#DDFFDD'),
        ("Сервер 3", 10, 6, '#DDFFDD'),
        ("PostgreSQL\nPrimary", 7, 3, '#DDDDFF'),
        ("PostgreSQL\nReplica", 10, 3, '#DDDDFF'),
        ("Prometheus\nМониторинг", 4, 3, '#FFFFDD')
    ]

    for text, x, y, color in components:
        ax.add_patch(plt.Rectangle((x - 1.5, y - 0.8), 3, 1.6, color=color))
        ax.text(x, y, text, ha='center', va='center')

    # Связи
    connections = [
        (1.5, 8, 2.5, 8),  # Операторы -> Nginx
        (4, 7.2, 4, 6.8),  # Nginx -> Сервер 1
        (4, 6.8, 7, 6.8),  # Сервер 1 -> Сервер 2
        (7, 6.8, 10, 6.8),  # Сервер 2 -> Сервер 3
        (4, 5.2, 4, 3.8),  # Сервер 1 -> Мониторинг
        (7, 5.2, 7, 3.8),  # Сервер 2 -> Primary DB
        (10, 5.2, 10, 3.8),  # Сервер 3 -> Primary DB
        (7, 2.2, 10, 2.2)  # Primary -> Replica
    ]

    for x1, y1, x2, y2 in connections:
        ax.plot([x1, x2], [y1, y2], 'k-', lw=1)
        if y1 == y2:  # Горизонтальные линии
            ax.plot((x2 - 0.2, x2, x2 - 0.2), (y2 - 0.1, y2, y2 + 0.1), 'k-', lw=1)
        else:  # Вертикальные линии
            ax.plot((x2 - 0.1, x2, x2 + 0.1), (y2 + 0.2, y2, y2 + 0.2), 'k-', lw=1)

    plt.title("Системная архитектура", pad=20)
    plt.tight_layout()
    plt.savefig('system_architecture.png', dpi=300)
    plt.close()


draw_system_architecture()