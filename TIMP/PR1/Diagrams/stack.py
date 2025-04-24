import matplotlib.pyplot as plt


def draw_layered_architecture():
    fig, ax = plt.subplots(figsize=(8, 6))
    ax.set_xlim(0, 5)
    ax.set_ylim(0, 4)
    ax.axis('off')

    # Слои
    layers = [
        ("Presentation Layer\n(Веб-интерфейс)", 2.5, 3, '#FFAAAA'),
        ("Business Layer\n(Логика обработки)", 2.5, 2, '#AAFFAA'),
        ("Data Layer\n(База данных)", 2.5, 1, '#AAAAFF')
    ]

    for text, x, y, color in layers:
        ax.add_patch(plt.Rectangle((x - 2, y - 0.4), 4, 0.8,
                                   fill=True, color=color))
        ax.text(x, y, text, ha='center', va='center')

    # Стрелки
    ax.annotate("", xy=(2.5, 2.6), xytext=(2.5, 2.4),
                arrowprops=dict(arrowstyle="->", lw=2))
    ax.annotate("", xy=(2.5, 1.6), xytext=(2.5, 1.4),
                arrowprops=dict(arrowstyle="->", lw=2))

    plt.title("Трёхслойная архитектура", pad=20)
    plt.tight_layout()
    plt.savefig('layered_architecture.png', dpi=300)
    plt.close()


draw_layered_architecture()