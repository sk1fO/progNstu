package services

import (
	"bufio"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type ValidationService struct {
}

func NewValidationService() *ValidationService {
	return &ValidationService{}
}

type TestResult struct {
	Passed      bool   `json:"passed"`
	Input       string `json:"input"`
	Expected    string `json:"expected"`
	Actual      string `json:"actual"`
	Description string `json:"description"`
}

// TestCase defines the structure for test cases
type TestCase struct {
	Input          string `json:"input"`
	ExpectedOutput string `json:"expected_output"`
	Description    string `json:"description"`
}

func (s *ValidationService) RunTest(code, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return "", fmt.Errorf("ошибка создания Docker клиента: %v", err)
	}

	// Экранируем код и входные данные
	escapedCode := escapeSingleQuotes(code)
	escapedInput := escapeSingleQuotes(input)

	// Создаем контейнер
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "cpp-runner",
		Cmd:   []string{"sh", "-c", fmt.Sprintf("echo '%s' > main.cpp && g++ -std=c++11 main.cpp -o main && echo '%s' | timeout 10s ./main 2>&1", escapedCode, escapedInput)},
		Tty:   false,
	}, &container.HostConfig{
		NetworkMode: "none",
		AutoRemove:  false,
	}, nil, nil, "")
	if err != nil {
		return "", fmt.Errorf("ошибка создания контейнера: %v", err)
	}

	defer func() {
		cli.ContainerRemove(context.Background(), resp.ID, container.RemoveOptions{Force: true})
	}()

	// Запускаем контейнер
	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return "", fmt.Errorf("ошибка запуска контейнера: %v", err)
	}

	// Ждем завершения
	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case <-statusCh:
	case err := <-errCh:
		if err != nil {
			return "", fmt.Errorf("ошибка ожидания контейнера: %v", err)
		}
	case <-ctx.Done():
		return "", fmt.Errorf("таймаут выполнения")
	}

	// Получаем логи
	out, err := cli.ContainerLogs(ctx, resp.ID, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
	})
	if err != nil {
		return "", fmt.Errorf("ошибка получения логов: %v", err)
	}
	defer out.Close()

	// Читаем вывод
	var output strings.Builder
	scanner := bufio.NewScanner(out)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 8 {
			output.WriteString(line[8:] + "\n")
		} else {
			output.WriteString(line + "\n")
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("ошибка чтения логов: %v", err)
	}

	return strings.TrimSpace(output.String()), nil
}

func (s *ValidationService) ValidateSolution(code string, tests []TestCase) ([]TestResult, bool) {
	var results []TestResult
	allPassed := true

	for _, test := range tests {
		actualOutput, err := s.RunTest(code, test.Input)
		if err != nil {
			results = append(results, TestResult{
				Passed:      false,
				Input:       test.Input,
				Expected:    test.ExpectedOutput,
				Actual:      "Ошибка выполнения: " + err.Error(),
				Description: test.Description,
			})
			allPassed = false
			continue
		}

		// Нормализуем вывод для сравнения (убираем лишние пробелы и переносы)
		expectedNormalized := strings.TrimSpace(test.ExpectedOutput)
		actualNormalized := strings.TrimSpace(actualOutput)

		passed := expectedNormalized == actualNormalized
		if !passed {
			allPassed = false
		}

		results = append(results, TestResult{
			Passed:      passed,
			Input:       test.Input,
			Expected:    expectedNormalized,
			Actual:      actualNormalized,
			Description: test.Description,
		})
	}

	return results, allPassed
}

func escapeSingleQuotes(s string) string {
	return strings.ReplaceAll(s, "'", "'\\''")
}
