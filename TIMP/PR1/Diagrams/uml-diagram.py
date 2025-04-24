import matplotlib.pyplot as plt
from matplotlib.patches import Rectangle


def draw_uml_class():
    fig, ax = plt.subplots(figsize=(8, 6))
    ax.set_xlim(0, 10)
    ax.set_ylim(0, 8)
    ax.axis('off')

    # Классы
    classes = [
        ("ConsoleUI", 2, 6, ['+ Run()']),
        ("AlertService", 5, 6, ['+ HandleAlert()']),
        ("MockDatabase", 8, 6, ['+ SaveAlert()'])
    ]

    for name, x, y, methods in classes:
        ax.add_patch(Rectangle((x - 1.5, y - 1), 3, 1.5, fill=True, color='#EEEEEE'))
        ax.text(x, y + 0.3, name, ha='center', va='center', fontweight='bold')

        for i, method in enumerate(methods):
            ax.text(x, y - 0.2 - i * 0.3, method, ha='center', va='center', fontsize=10)

    # Связи
    ax.annotate("", xy=(3.5, 6), xytext=(1.5, 6),
                arrowprops=dict(arrowstyle="->", linestyle='dashed'))
    ax.annotate("", xy=(6.5, 6), xytext=(4.5, 6),
                arrowprops=dict(arrowstyle="->", linestyle='dashed'))

    plt.title("UML: Диаграмма классов", pad=20)
    plt.tight_layout()
    plt.savefig('uml_class.png', dpi=300)
    plt.close()


draw_uml_class()