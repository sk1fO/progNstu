package utils

import (
	"bufio"
	"context"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func RunCppCode(code string) (string, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return "", err
	}

	// Создаем контейнер
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "gcc", // Используем официальный образ GCC
		Cmd:   []string{"sh", "-c", "echo '" + code + "' > main.cpp && g++ main.cpp -o main && ./main"},
	}, nil, nil, nil, "")
	if err != nil {
		return "", err
	}

	// Запускаем контейнер
	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return "", err
	}

	// Ждем завершения
	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return "", err
		}
	case <-statusCh:
	}

	// Получаем логи (Исправлено: используем container.LogsOptions)
	out, err := cli.ContainerLogs(ctx, resp.ID, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
	})
	if err != nil {
		return "", err
	}
	defer out.Close()

	// Читаем логи
	var output strings.Builder
	scanner := bufio.NewScanner(out)
	for scanner.Scan() {
		output.WriteString(scanner.Text() + "\n")
	}

	// Удаляем контейнер (Исправлено: используем container.RemoveOptions)
	err = cli.ContainerRemove(ctx, resp.ID, container.RemoveOptions{
		Force: true, // Принудительное удаление, если контейнер еще работает
	})
	if err != nil {
		return "", err
	}

	return output.String(), nil
}
