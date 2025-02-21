package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// readFile читает содержимое файла и возвращает его в виде структуры данных
func readFile(t *testing.T, file string) map[string]interface{} {
	data, err := ioutil.ReadFile(file)
	assert.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	assert.NoError(t, err)

	return result
}

// resetFlags сбрасывает флаги перед вызовом main
func resetFlags() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
}

// TestMain_ArrayOperations тестирует операции с массивом
func TestMain_ArrayOperations(t *testing.T) {
	// Создаем временный файл для хранения данных
	tmpFile, err := ioutil.TempFile("", "test_*.json")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	// Команды для выполнения
	commands := []string{
		"APUSH test_array 1",
		"APUSH test_array 2",
		"AGET test_array 0",
		"ALEN test_array",
		"ADEL test_array 0",
		"PRINT test_array",
	}

	// Выполняем команды
	for _, cmd := range commands {
		resetFlags() // Сбрасываем флаги перед каждым вызовом main
		os.Args = []string{"main", "--file", tmpFile.Name(), "--query", cmd}
		main()
	}

	// Проверяем содержимое файла
	data := readFile(t, tmpFile.Name())
	assert.Contains(t, data["test_array"], "2")
	assert.NotContains(t, data["test_array"], "1")
}

// TestMain_ListOperations тестирует операции со списком
func TestMain_ListOperations(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "test_*.json")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	commands := []string{
		"LPUSHHEAD test_list 1",
		"LPUSHTAIL test_list 2",
		"LDELHEAD test_list",
		"PRINT test_list",
	}

	for _, cmd := range commands {
		resetFlags()
		os.Args = []string{"main", "--file", tmpFile.Name(), "--query", cmd}
		main()
	}

	data := readFile(t, tmpFile.Name())
	assert.Contains(t, data["test_list"], "2")
	assert.NotContains(t, data["test_list"], "1")
}

// TestMain_QueueOperations тестирует операции с очередью
func TestMain_QueueOperations(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "test_*.json")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	commands := []string{
		"QPUSH test_queue 1",
		"QPUSH test_queue 2",
		"QPOP test_queue",
		"PRINT test_queue",
	}

	for _, cmd := range commands {
		resetFlags()
		os.Args = []string{"main", "--file", tmpFile.Name(), "--query", cmd}
		main()
	}

	data := readFile(t, tmpFile.Name())
	assert.Contains(t, data["test_queue"], "2")
	assert.NotContains(t, data["test_queue"], "1")
}

// TestMain_StackOperations тестирует операции со стеком
func TestMain_StackOperations(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "test_*.json")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	commands := []string{
		"SPUSH test_stack 1",
		"SPUSH test_stack 2",
		"SPOP test_stack",
		"PRINT test_stack",
	}

	for _, cmd := range commands {
		resetFlags()
		os.Args = []string{"main", "--file", tmpFile.Name(), "--query", cmd}
		main()
	}

	data := readFile(t, tmpFile.Name())
	assert.Contains(t, data["test_stack"], "1")
	assert.NotContains(t, data["test_stack"], "2")
}

// TestMain_HashTableOperations тестирует операции с хеш-таблицей
func TestMain_HashTableOperations(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "test_*.json")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	commands := []string{
		"HSET test_hash key1 value1",
		"HSET test_hash key2 value2",
		"HGET test_hash key1",
		"HDEL test_hash key1",
		"PRINT test_hash",
	}

	for _, cmd := range commands {
		resetFlags()
		os.Args = []string{"main", "--file", tmpFile.Name(), "--query", cmd}
		main()
	}

	data := readFile(t, tmpFile.Name())
	assert.Contains(t, data["test_hash"], "value2")
	assert.NotContains(t, data["test_hash"], "value1")
}

// TestMain_CBTreeOperations тестирует операции с полным бинарным деревом
func TestMain_CBTreeOperations(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "test_*.json")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	commands := []string{
		"TINSERT test_cbtree 1",
		"TINSERT test_cbtree 2",
		"TFIND test_cbtree 1",
		"TISCOMPLETE test_cbtree",
		"PRINT test_cbtree",
	}

	for _, cmd := range commands {
		resetFlags()
		os.Args = []string{"main", "--file", tmpFile.Name(), "--query", cmd}
		main()
	}

	data := readFile(t, tmpFile.Name())
	assert.Contains(t, data["test_cbtree"], "1")
	assert.Contains(t, data["test_cbtree"], "2")
}

// TestMain_ErrorHandling тестирует обработку ошибок
func TestMain_ErrorHandling(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "test_*.json")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	// Некорректная команда
	os.Args = []string{"main", "--file", tmpFile.Name(), "--query", "INVALID_COMMAND"}
	assert.Panics(t, main) // Ожидаем панику

	// Попытка удаления из пустого массива
	os.Args = []string{"main", "--file", tmpFile.Name(), "--query", "ADEL test_array 0"}
	assert.Panics(t, main) // Ожидаем панику

	// Попытка доступа к несуществующему индексу
	os.Args = []string{"main", "--file", tmpFile.Name(), "--query", "AGET test_array 100"}
	assert.Panics(t, main) // Ожидаем панику
}

// TestMain_ArrayPushIndex тестирует команду APUSHINDEX
func TestMain_ArrayPushIndex(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "test_*.json")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	// Добавляем элемент по индексу
	os.Args = []string{"main", "--file", tmpFile.Name(), "--query", "APUSHINDEX test_array 0 1"}
	main()

	// Проверяем результат
	data := readFile(t, tmpFile.Name())
	assert.Contains(t, data["test_array"], "1")
}

// TestMain_BoundaryCases тестирует граничные случаи
func TestMain_BoundaryCases(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "test_*.json")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	// Добавляем элемент в пустой массив
	os.Args = []string{"main", "--file", tmpFile.Name(), "--query", "APUSH test_array 1"}
	main()

	// Удаляем последний элемент
	os.Args = []string{"main", "--file", tmpFile.Name(), "--query", "ADEL test_array 0"}
	main()

	// Проверяем, что массив пуст
	data := readFile(t, tmpFile.Name())
	assert.Empty(t, data["test_array"])

	// Поиск несуществующего элемента
	os.Args = []string{"main", "--file", tmpFile.Name(), "--query", "TFIND test_cbtree 100"}
	main()
}

// TestMain_SingleLinkedListOperations тестирует операции с односвязным списком
func TestMain_SingleLinkedListOperations(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "test_*.json")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	// Добавляем элементы в голову и хвост
	commands := []string{
		"SLPUSHHEAD test_slist 1",
		"SLPUSHTAIL test_slist 2",
		"SLDELHEAD test_slist",
		"PRINT test_slist",
	}

	for _, cmd := range commands {
		resetFlags()
		os.Args = []string{"main", "--file", tmpFile.Name(), "--query", cmd}
		main()
	}

	// Проверяем результат
	data := readFile(t, tmpFile.Name())
	assert.Contains(t, data["test_slist"], "2")
	assert.NotContains(t, data["test_slist"], "1")
}

// TestMain_SaveAndLoad тестирует сохранение и загрузку данных
// func TestMain_SaveAndLoad(t *testing.T) {
// 	tmpFile, err := ioutil.TempFile("", "test_*.json")
// 	assert.NoError(t, err)
// 	defer os.Remove(tmpFile.Name())

// 	// Добавляем данные
// 	os.Args = []string{"main", "--file", tmpFile.Name(), "--query", "APUSH test_array 1"}
// 	main()

// 	// Перезагружаем данные
// 	os.Args = []string{"main", "--file", tmpFile.Name(), "--query", "PRINT test_array"}
// 	main()

// 	// Проверяем результат
// 	data := readFile(t, tmpFile.Name())
// 	assert.Contains(t, data["test_array"], "1")
// }

// TestMain_PrintCommands тестирует команду PRINT для всех структур данных
func TestMain_PrintCommands(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "test_*.json")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	// Добавляем данные
	commands := []string{
		"APUSH test_array 1",
		"LPUSHHEAD test_list 1",
		"QPUSH test_queue 1",
		"SPUSH test_stack 1",
		"HSET test_hash key1 value1",
		"TINSERT test_cbtree 1",
	}

	for _, cmd := range commands {
		resetFlags()
		os.Args = []string{"main", "--file", tmpFile.Name(), "--query", cmd}
		main()
	}

	// Печатаем все структуры данных
	printCommands := []string{
		"PRINT test_array",
		"PRINT test_list",
		"PRINT test_queue",
		"PRINT test_stack",
		"PRINT test_hash",
		"PRINT test_cbtree",
	}

	for _, cmd := range printCommands {
		resetFlags()
		os.Args = []string{"main", "--file", tmpFile.Name(), "--query", cmd}
		main()
	}
}

// TestMain_CBTreeCommands тестирует команды для полного бинарного дерева
func TestMain_CBTreeCommands(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "test_*.json")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	// Добавляем данные
	commands := []string{
		"TINSERT test_cbtree 1",
		"TINSERT test_cbtree 2",
		"TFIND test_cbtree 1",
		"TISCOMPLETE test_cbtree",
	}

	for _, cmd := range commands {
		resetFlags()
		os.Args = []string{"main", "--file", tmpFile.Name(), "--query", cmd}
		main()
	}
}
