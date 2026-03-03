#!/bin/bash

echo "Запуск бэкенда (Go)..."
cd backend
go run main.go &
BACKEND_PID=$!

echo "Запуск фронтенда (React)..."
cd ../frontend/pass-control
npm run dev &
FRONTEND_PID=$!

echo "Сервер и клиент запущены. Нажмите Ctrl+C для остановки."
wait $BACKEND_PID $FRONTEND_PID