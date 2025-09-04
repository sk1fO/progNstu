FROM gcc:latest

# Устанавливаем необходимые инструменты
RUN apt-get update && apt-get install -y \
    python3 \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Этот образ будет использоваться для запуска C++ кода
# Базовый образ GCC уже содержит все необходимое для компиляции