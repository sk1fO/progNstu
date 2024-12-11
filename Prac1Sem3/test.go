package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

// Define the schema structure
type Schema struct {
	Name       string              `json:"name"`
	TupleLimit int                 `json:"tuples_limit"`
	Structure  map[string][]string `json:"structure"`
}

// Define the table structure
type Table struct {
	Name         string
	Columns      []string
	PrimaryKey   string
	SequenceFile string
	LockFile     string
	Files        []string
	Lock         *sync.Mutex
}

// Define the database structure
type Database struct {
	Name       string
	TupleLimit int
	Tables     map[string]*Table
}

// Read the schema configuration
func readSchema(filePath string) (*Schema, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var schema Schema
	if err := json.Unmarshal(data, &schema); err != nil {
		return nil, err
	}
	return &schema, nil
}

// Create the necessary directories and files based on the schema
func createDatabase(schema *Schema) (*Database, error) {
	db := &Database{
		Name:       schema.Name,
		TupleLimit: schema.TupleLimit,
		Tables:     make(map[string]*Table),
	}

	baseDir := schema.Name
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, err
	}

	for tableName, columns := range schema.Structure {
		tableDir := filepath.Join(baseDir, tableName)
		if err := os.MkdirAll(tableDir, 0755); err != nil {
			return nil, err
		}

		primaryKey := fmt.Sprintf("%s_pk", tableName)
		sequenceFile := filepath.Join(tableDir, fmt.Sprintf("%s_sequence", tableName))
		lockFile := filepath.Join(tableDir, fmt.Sprintf("%s_lock", tableName))

		// Initialize sequence file
		if _, err := os.Stat(sequenceFile); os.IsNotExist(err) {
			if err := ioutil.WriteFile(sequenceFile, []byte("0"), 0644); err != nil {
				return nil, err
			}
		}

		// Initialize lock file
		if _, err := os.Stat(lockFile); os.IsNotExist(err) {
			if err := ioutil.WriteFile(lockFile, []byte("false"), 0644); err != nil {
				return nil, err
			}
		}

		// Initialize first CSV file
		csvFile := filepath.Join(tableDir, "1.csv")
		if _, err := os.Stat(csvFile); os.IsNotExist(err) {
			file, err := os.Create(csvFile)
			if err != nil {
				return nil, err
			}
			writer := csv.NewWriter(file)
			writer.Write(append([]string{primaryKey}, columns...))
			writer.Flush()
			file.Close()
		}

		db.Tables[tableName] = &Table{
			Name:         tableName,
			Columns:      append([]string{primaryKey}, columns...),
			PrimaryKey:   primaryKey,
			SequenceFile: sequenceFile,
			LockFile:     lockFile,
			Files:        []string{csvFile},
			Lock:         &sync.Mutex{},
		}
	}

	return db, nil
}

// Implement the SELECT statement
func selectData(db *Database, query string) ([][]string, error) {
	parts := strings.Split(query, " ")
	if len(parts) < 4 || parts[2] != "FROM" {
		return nil, fmt.Errorf("invalid SELECT command. Usage: SELECT <columns> FROM <tables>")
	}

	selectCols := strings.Split(parts[1], ",")
	fromTables := strings.Split(parts[3], ",")

	var result [][]string

	for _, table := range fromTables {
		tbl, ok := db.Tables[table]
		if !ok {
			return nil, fmt.Errorf("table %s does not exist", table)
		}

		for _, file := range tbl.Files {
			rows, err := readCSV(file)
			if err != nil {
				return nil, err
			}

			for _, row := range rows[1:] { // Skip header
				var selectedRow []string
				for _, col := range selectCols {
					colParts := strings.Split(col, ".")
					if colParts[0] == table {
						idx := -1
						for i, c := range tbl.Columns {
							if c == colParts[1] {
								idx = i
								break
							}
						}
						if idx != -1 {
							selectedRow = append(selectedRow, row[idx])
						}
					}
				}
				result = append(result, selectedRow)
			}
		}
	}

	return result, nil
}

// Implement the WHERE clause
func whereClause(rows [][]string, conditions string) [][]string {
	var header []string = rows[0]
	var result [][]string
	for _, row := range rows {
		if evaluateCondition(header, row, conditions) {
			result = append(result, row)
		}
	}
	return result
}

// Evaluate the WHERE condition
func evaluateCondition(header, row []string, conditions string) bool {
	// Simple evaluation for equality
	parts := strings.Split(conditions, "=")
	col := parts[0]
	value := parts[1]

	colParts := strings.Split(col, ".")
	idx := -1
	for i, c := range header { // from header!!!
		if c == colParts[1] {
			idx = i
			break
		}
	}

	if idx != -1 {
		return row[idx] == value
	}

	return false
}

// Implement the INSERT statement
func insertData(db *Database, table string, values []string) error {
	tbl, ok := db.Tables[table]
	if !ok {
		return fmt.Errorf("table %s does not exist", table)
	}

	tbl.Lock.Lock()
	defer tbl.Lock.Unlock()

	// Get next primary key
	seq, err := getNextSequence(tbl.SequenceFile)
	if err != nil {
		return err
	}

	// Append primary key to values
	values = append([]string{strconv.Itoa(seq)}, values...)

	// Determine which file to write to
	lastFile := tbl.Files[len(tbl.Files)-1]
	rows, err := readCSV(lastFile)
	if err != nil {
		return err
	}

	if len(rows) >= db.TupleLimit {
		// Create a new file
		newFile := filepath.Join(filepath.Dir(lastFile), fmt.Sprintf("%d.csv", len(tbl.Files)+1))
		file, err := os.Create(newFile)
		if err != nil {
			return err
		}
		writer := csv.NewWriter(file)
		writer.Write(tbl.Columns)
		writer.Write(values)
		writer.Flush()
		file.Close()
		tbl.Files = append(tbl.Files, newFile)
	} else {
		// Write to the last file
		file, err := os.OpenFile(lastFile, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		writer := csv.NewWriter(file)
		writer.Write(values)
		writer.Flush()
		file.Close()
	}

	return nil
}

// Get the next sequence number
func getNextSequence(file string) (int, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return 0, err
	}
	seq, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		return 0, err
	}
	seq++
	if err := ioutil.WriteFile(file, []byte(strconv.Itoa(seq)), 0644); err != nil {
		return 0, err
	}
	return seq, nil
}

// Implement the DELETE statement
func deleteData(db *Database, table string, conditions string) error {
	tbl, ok := db.Tables[table]
	if !ok {
		return fmt.Errorf("table %s does not exist", table)
	}

	tbl.Lock.Lock()
	defer tbl.Lock.Unlock()

	for i, file := range tbl.Files {
		rows, err := readCSV(file)
		if err != nil {
			return err
		}

		header := rows[0]
		newRows := [][]string{header} // Keep the header

		for _, row := range rows[1:] {
			if !evaluateCondition(header, row, conditions) {
				newRows = append(newRows, row)
			}
		}

		// If the file becomes empty after deletion, remove it
		if len(newRows) == 1 {
			if err := os.Remove(file); err != nil {
				return err
			}
			// Remove the file from the list
			tbl.Files = append(tbl.Files[:i], tbl.Files[i+1:]...)
		} else {
			// Write the new rows back to the file
			file, err := os.Create(file)
			if err != nil {
				return err
			}
			writer := csv.NewWriter(file)
			for _, row := range newRows {
				if err := writer.Write(row); err != nil {
					file.Close()
					return err
				}
			}
			writer.Flush()
			file.Close()
		}
	}

	return nil
}

// Read CSV file
func readCSV(file string) ([][]string, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return records, nil
}

// Main function to start the database
func main() {
	// Read the schema configuration
	schema, err := readSchema("schema.json")
	if err != nil {
		log.Fatalf("Failed to read schema: %v", err)
	}

	// Create the database
	db, err := createDatabase(schema)
	if err != nil {
		log.Fatalf("Failed to create database: %v", err)
	}

	// Console menu for direct SQL command input
	for {
		fmt.Print("\nEnter an SQL command (or 'exit' to quit): ")

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		command := scanner.Text()

		if command == "exit" {
			fmt.Println("Exiting...")
			return
		}

		parts := strings.Split(command, " ")

		switch parts[0] {
		case "INSERT":
			if len(parts) < 5 || parts[1] != "INTO" || parts[3] != "VALUES" {
				fmt.Println("Invalid INSERT command. Usage: INSERT INTO <table> VALUES (<values>)")
				continue
			}
			table := parts[2]
			values := strings.Trim(parts[4], "()")
			valuesList := strings.Split(values, ",")
			if err := insertData(db, table, valuesList); err != nil {
				fmt.Println("Error inserting data:", err)
			} else {
				fmt.Println("Data inserted successfully.")
			}
		case "SELECT":
			if len(parts) < 3 || parts[2] != "FROM" {
				fmt.Println("Invalid SELECT command. Usage: SELECT <columns> FROM <tables>")
				continue
			}
			query := strings.Join(parts, " ")
			result, err := selectData(db, query)
			if err != nil {
				fmt.Println("Error selecting data:", err)
			} else {
				for _, row := range result {
					fmt.Println(row)
				}
			}
		case "DELETE":
			if len(parts) < 4 || parts[1] != "FROM" {
				fmt.Println("Invalid DELETE command. Usage: DELETE FROM <table> [WHERE <conditions>]")
				continue
			}
			table := parts[2]
			conditions := ""
			if len(parts) > 4 && parts[3] == "WHERE" {
				conditions = strings.Join(parts[4:], " ")
			}
			if err := deleteData(db, table, conditions); err != nil {
				fmt.Println("Error deleting data:", err)
			} else {
				fmt.Println("Data deleted successfully.")
			}
		default:
			fmt.Println("Unknown command. Supported commands: INSERT, SELECT, DELETE, exit")
		}
	}
}
