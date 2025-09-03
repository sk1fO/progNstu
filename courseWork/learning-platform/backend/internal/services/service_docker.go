package services

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
)

type DockerService struct {
	cli *client.Client
}

func NewDockerService() (*DockerService, error) {
	// Указываем версию API явно для совместимости
	cli, err := client.NewClientWithOpts(
		client.WithHostFromEnv(),
		client.WithAPIVersionNegotiation(), // Важно для совместимости версий
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create Docker client: %v", err)
	}

	// Проверяем подключение к Docker
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = cli.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Docker: %v", err)
	}

	return &DockerService{cli: cli}, nil
}

type RunCodeRequest struct {
	Code     string
	Language string
	Timeout  time.Duration
}

type RunCodeResponse struct {
	Output string
	Error  string
}

func (s *DockerService) RunCode(ctx context.Context, req RunCodeRequest) (*RunCodeResponse, error) {
	// Определяем образ в зависимости от языка
	image := s.getImageByLanguage(req.Language)

	// Создаем контейнер
	resp, err := s.cli.ContainerCreate(ctx, &container.Config{
		Image:        image,
		Cmd:          []string{"sh", "-c", fmt.Sprintf("echo '%s' > /app/code && %s", escapeSingleQuotes(req.Code), s.getExecutionCommand(req.Language))},
		Tty:          false,
		AttachStdout: true,
		AttachStderr: true,
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeTmpfs,
				Target: "/app",
			},
		},
		Resources: container.Resources{
			Memory:   100 * 1024 * 1024, // 100MB
			NanoCPUs: 1e9,               // 1 CPU core
		},
	}, nil, nil, "")
	if err != nil {
		return nil, fmt.Errorf("failed to create container: %v", err)
	}
	defer s.removeContainer(context.Background(), resp.ID) // Используем отдельный контекст для очистки

	// Запускаем контейнер
	if err := s.cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return nil, fmt.Errorf("failed to start container: %v", err)
	}

	// Ожидаем завершения с таймаутом
	statusCh, errCh := s.cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("execution timeout: %v", ctx.Err())
	case err := <-errCh:
		if err != nil {
			return nil, fmt.Errorf("error waiting for container: %v", err)
		}
	case <-statusCh:
		// Контейнер завершил работу
	}

	// Получаем логи
	out, err := s.cli.ContainerLogs(ctx, resp.ID, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get container logs: %v", err)
	}
	defer out.Close()

	output, err := decodeDockerOutput(out)
	if err != nil {
		return nil, fmt.Errorf("failed to decode output: %v", err)
	}

	// Проверяем, была ли ошибка выполнения
	if strings.Contains(output, "error") || strings.Contains(output, "Error") {
		return &RunCodeResponse{Output: "", Error: output}, nil
	}

	return &RunCodeResponse{Output: output, Error: ""}, nil
}

func (s *DockerService) getImageByLanguage(language string) string {
	switch strings.ToLower(language) {
	case "python":
		return "python:3.11-slim"
	case "javascript", "js":
		return "node:18-slim"
	case "go", "golang":
		return "golang:1.21-alpine"
	default:
		return "python:3.11-slim"
	}
}

func (s *DockerService) getExecutionCommand(language string) string {
	switch strings.ToLower(language) {
	case "python":
		return "python /app/code"
	case "javascript", "js":
		return "node /app/code"
	case "go", "golang":
		return "go run /app/code"
	default:
		return "python /app/code"
	}
}

func (s *DockerService) removeContainer(ctx context.Context, containerID string) {
	// Не блокируемся на удалении контейнера
	go func() {
		_ = s.cli.ContainerRemove(context.Background(), containerID, container.RemoveOptions{
			Force: true,
		})
	}()
}

func decodeDockerOutput(output io.Reader) (string, error) {
	var result strings.Builder
	scanner := bufio.NewScanner(output)

	for scanner.Scan() {
		line := scanner.Text()
		// Docker добавляет 8-байтовый заголовок к каждой строке, который нужно обрезать
		if len(line) > 8 {
			// Для stdout байт 0x01, для stderr 0x02, первые 4 байта - размер сообщения
			// Простая эвристика: пропускаем первые 8 байт
			result.WriteString(line[8:])
		} else {
			result.WriteString(line)
		}
		result.WriteString("\n")
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return strings.TrimSpace(result.String()), nil
}

func escapeSingleQuotes(input string) string {
	return strings.ReplaceAll(input, "'", "'\"'\"'")
}
