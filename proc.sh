#!/bin/bash
# test_process_analysis.sh

echo "=== Анализ процессов Go тестов ==="

# Запускаем тесты и отслеживаем процессы
echo "1. Запуск всех тестов:"
ps aux | grep "go test" | grep -v grep || echo "Нет процессов go test"

echo "2. Запускаем тесты и мониторим процессы:"
go test -v ./... &
TEST_PID=$!

sleep 1
echo "3. Процессы во время выполнения:"
pstree -p $TEST_PID 2>/dev/null || ps -o pid,ppid,command -p $TEST_PID

wait $TEST_PID
echo "4. Тесты завершены"

