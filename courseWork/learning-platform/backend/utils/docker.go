package utils

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func RunCppCode(code string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return "", fmt.Errorf("ошибка создания Docker клиента: %v", err)
	}

	// Создаем контейнер из нашего образа
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "cpp-runner",
		Cmd:   []string{"sh", "-c", "echo '" + escapeSingleQuotes(code) + "' > main.cpp && g++ -std=c++11 main.cpp -o main && (timeout 10s ./main || echo 'Execution timeout') 2>&1"},
		Tty:   false,
	}, &container.HostConfig{
		NetworkMode: "none", // Отключаем сеть для безопасности
		AutoRemove:  false,  // НЕ удаляем автоматически - нужно получить логи
	}, nil, nil, "")
	if err != nil {
		return "", fmt.Errorf("ошибка создания контейнера: %v", err)
	}

	// Запускаем контейнер
	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return "", fmt.Errorf("ошибка запуска контейнера: %v", err)
	}

	// Ждем завершения с таймаутом
	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)

	select {
	case <-statusCh:
		// Контейнер завершился
	case err := <-errCh:
		if err != nil {
			return "", fmt.Errorf("ошибка ожидания контейнера: %v", err)
		}
	case <-ctx.Done():
		return "", fmt.Errorf("таймаут выполнения кода")
	}

	// Получаем логи ПЕРЕД удалением контейнера
	out, err := cli.ContainerLogs(ctx, resp.ID, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     false,
	})
	if err != nil {
		// Пытаемся удалить контейнер даже если не удалось получить логи
		cli.ContainerRemove(ctx, resp.ID, container.RemoveOptions{Force: true})
		return "", fmt.Errorf("ошибка получения логов: %v", err)
	}
	defer out.Close()

	// Читаем и обрабатываем логи
	var output strings.Builder
	scanner := bufio.NewScanner(out)
	for scanner.Scan() {
		line := scanner.Text()
		// Docker добавляет префиксы к логам (8 байт заголовка), убираем их
		if len(line) > 8 {
			output.WriteString(line[8:] + "\n")
		} else {
			output.WriteString(line + "\n")
		}
	}

	if err := scanner.Err(); err != nil {
		cli.ContainerRemove(ctx, resp.ID, container.RemoveOptions{Force: true})
		return "", fmt.Errorf("ошибка чтения логов: %v", err)
	}

	// Удаляем контейнер после получения логов
	err = cli.ContainerRemove(ctx, resp.ID, container.RemoveOptions{Force: true})
	if err != nil {
		log.Printf("Предупреждение: не удалось удалить контейнер %s: %v", resp.ID, err)
	}

	return strings.TrimSpace(output.String()), nil
}

func escapeSingleQuotes(s string) string {
	return strings.ReplaceAll(s, "'", "'\\''")
}

// Альтернативная версия с использованием exec для большей надежности
func RunCppCodeAlternative(code string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return "", fmt.Errorf("ошибка создания Docker клиента: %v", err)
	}

	// Создаем контейнер
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "cpp-runner",
		Tty:   false,
	}, &container.HostConfig{
		NetworkMode: "none",
		AutoRemove:  false,
	}, nil, nil, "")
	if err != nil {
		return "", fmt.Errorf("ошибка создания контейнера: %v", err)
	}

	defer func() {
		// Гарантированно удаляем контейнер при выходе
		cli.ContainerRemove(context.Background(), resp.ID, container.RemoveOptions{Force: true})
	}()

	// Запускаем контейнер
	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return "", fmt.Errorf("ошибка запуска контейнера: %v", err)
	}

	// Создаем файл с кодом
	execConfig := container.ExecOptions{
		Cmd:          []string{"sh", "-c", "echo '" + escapeSingleQuotes(code) + "' > main.cpp"},
		AttachStdout: true,
		AttachStderr: true,
	}

	execID, err := cli.ContainerExecCreate(ctx, resp.ID, execConfig)
	if err != nil {
		return "", fmt.Errorf("ошибка создания exec: %v", err)
	}

	// Выполняем команду записи кода
	attach, err := cli.ContainerExecAttach(ctx, execID.ID, container.ExecAttachOptions{})
	if err != nil {
		return "", fmt.Errorf("ошибка присоединения к exec: %v", err)
	}
	attach.Close()

	// Компилируем код

	execConfig = container.ExecOptions{
		Cmd:          []string{"g++", "-std=c++11", "main.cpp", "-o", "main"},
		AttachStdout: true,
		AttachStderr: true,
	}

	execID, err = cli.ContainerExecCreate(ctx, resp.ID, execConfig)
	if err != nil {
		return "", fmt.Errorf("ошибка создания exec для компиляции: %v", err)
	}

	attach, err = cli.ContainerExecAttach(ctx, execID.ID, container.ExecStartOptions{})
	if err != nil {
		return "", fmt.Errorf("ошибка присоединения к exec компиляции: %v", err)
	}
	defer attach.Close()

	// Читаем вывод компиляции
	compileOutput, _ := io.ReadAll(attach.Reader)
	compileOutputStr := string(compileOutput)

	// Если есть ошибки компиляции, возвращаем их
	if strings.Contains(compileOutputStr, "error:") {
		return compileOutputStr, nil
	}

	// Запускаем скомпилированную программу
	execConfig = container.ExecOptions{
		Cmd:          []string{"timeout", "10s", "./main"},
		AttachStdout: true,
		AttachStderr: true,
	}

	execID, err = cli.ContainerExecCreate(ctx, resp.ID, execConfig)
	if err != nil {
		return "", fmt.Errorf("ошибка создания exec для выполнения: %v", err)
	}

	attach, err = cli.ContainerExecAttach(ctx, execID.ID, container.ExecStartOptions{})
	if err != nil {
		return "", fmt.Errorf("ошибка присоединения к exec выполнения: %v", err)
	}
	defer attach.Close()

	// Читаем вывод выполнения
	output, _ := io.ReadAll(attach.Reader)
	outputStr := string(output)

	return strings.TrimSpace(outputStr), nil
}
