#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🚀 Запуск C++ Learning Platform...${NC}"

# Проверяем установлен ли Docker
if ! command -v docker &> /dev/null; then
    echo -e "${RED}❌ Docker не установлен. Пожалуйста, установите Docker сначала.${NC}"
    exit 1
fi

# Проверяем установлен ли Node.js
if ! command -v node &> /dev/null; then
    echo -e "${RED}❌ Node.js не установлен. Пожалуйста, установите Node.js сначала.${NC}"
    exit 1
fi

# Проверяем установлен ли Go
if ! command -v go &> /dev/null; then
    echo -e "${RED}❌ Go не установлен. Пожалуйста, установите Go сначала.${NC}"
    exit 1
fi

# Функция для проверки портов
check_port() {
    if lsof -Pi :$1 -sTCP:LISTEN -t >/dev/null ; then
        return 0
    else
        return 1
    fi
}

# Функция для ожидания порта
wait_for_port() {
    local port=$1
    local timeout=30
    local counter=0
    
    echo -e "${YELLOW}⏳ Ожидаем запуск сервиса на порту $port...${NC}"
    
    while ! check_port $port; do
        if [ $counter -eq $timeout ]; then
            echo -e "${RED}❌ Таймаут ожидания порта $port${NC}"
            return 1
        fi
        sleep 1
        ((counter++))
    done
    return 0
}

# Функция для обработки Ctrl+C
cleanup() {
    echo -e "\n${YELLOW}🛑 Останавливаем сервисы...${NC}"
    
    # Останавливаем фронтенд
    if [ ! -z "$FRONTEND_PID" ]; then
        kill $FRONTEND_PID 2>/dev/null
    fi
    
    # Останавливаем бэкенд
    if [ ! -z "$BACKEND_PID" ]; then
        kill $BACKEND_PID 2>/dev/null
    fi
    
    echo -e "${GREEN}✅ Все сервисы остановлены${NC}"
    exit 0
}

# Устанавливаем обработчик Ctrl+C
trap cleanup INT

# Создаем папки если их нет
mkdir -p data docker

# 1. Собираем Docker образ для C++
echo -e "${BLUE}🐳 Собираем Docker образ для C++...${NC}"
docker build -f backend/docker/Dockerfile.cpp -t cpp-runner . || {
    echo -e "${RED}❌ Ошибка сборки Docker образа${NC}"
    exit 1
}

# 2. Запускаем бэкенд
echo -e "${BLUE}🔧 Запускаем бэкенд (Go)...${NC}"
cd backend

# Устанавливаем зависимости Go
echo -e "${YELLOW}📦 Устанавливаем зависимости Go...${NC}"
go mod download

# Запускаем бэкенд в фоне
go run main.go &
BACKEND_PID=$!
cd ..

# Ждем запуска бэкенда
if wait_for_port 8080; then
    echo -e "${GREEN}✅ Бэкенд запущен на http://localhost:8080${NC}"
else
    echo -e "${RED}❌ Не удалось запустить бэкенд${NC}"
    cleanup
    exit 1
fi

# 3. Запускаем фронтенд
echo -e "${BLUE}⚛️  Запускаем фронтенд (React)...${NC}"
cd frontend

# Проверяем установлены ли зависимости Node.js
if [ ! -d "node_modules" ]; then
    echo -e "${YELLOW}📦 Устанавливаем зависимости Node.js...${NC}"
    npm install || {
        echo -e "${RED}❌ Ошибка установки зависимостей Node.js${NC}"
        cleanup
        exit 1
    }
fi

# Запускаем фронтенд в фоне
npm start &
FRONTEND_PID=$!
cd ..

# Ждем запуска фронтенда
if wait_for_port 3000; then
    echo -e "${GREEN}✅ Фронтенд запущен на http://localhost:3000${NC}"
else
    echo -e "${RED}❌ Не удалось запустить фронтенд${NC}"
    cleanup
    exit 1
fi

echo -e "${GREEN}🎉 Все сервисы успешно запущены!${NC}"
echo -e "${BLUE}📍 Фронтенд: http://localhost:3000${NC}"
echo -e "${BLUE}📍 Бэкенд:   http://localhost:8080${NC}"
echo -e "${BLUE}📍 API Docs: http://localhost:8080 (используйте инструменты разработчика)${NC}"
echo -e "\n${YELLOW}ℹ️  Нажмите Ctrl+C для остановки всех сервисов${NC}"

# Бесконечный цикл чтобы скрипт не завершался
while true; do
    sleep 1
done