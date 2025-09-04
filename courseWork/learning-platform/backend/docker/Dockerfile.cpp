FROM ubuntu:22.04

# Устанавливаем необходимые пакеты
RUN apt-get update && apt-get install -y \
    g++ \
    build-essential \
    coreutils \
    && rm -rf /var/lib/apt/lists/*

# Создаем рабочую директорию
WORKDIR /app

# Устанавливаем ограничения для безопасности
RUN useradd -m appuser && chown -R appuser:appuser /app
USER appuser

# Устанавливаем переменные окружения для безопасности
ENV PATH="/usr/bin:${PATH}"
ENV LD_LIBRARY_PATH=""

CMD ["/bin/bash"]