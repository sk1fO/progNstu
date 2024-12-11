package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// Set represents a set of elements
type Set struct {
	elements map[string]bool
}

// NewSet creates a new set
func NewSet() *Set {
	return &Set{
		elements: make(map[string]bool),
	}
}

// Add adds an element to the set
func (s *Set) Add(element string) {
	s.elements[element] = true
}

// Delete removes an element from the set
func (s *Set) Delete(element string) {
	delete(s.elements, element)
}

// Contains checks if the set contains an element
func (s *Set) Contains(element string) bool {
	_, exists := s.elements[element]
	return exists
}

// Command represents a command in the JSON file
type Command struct {
	Action  string `json:"action"`
	Element string `json:"element"`
}

// ProcessFile reads a JSON file and processes commands
func ProcessFile(filePath string, query string) error {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	var commands []Command
	if err := json.Unmarshal(file, &commands); err != nil {
		return err
	}

	set := NewSet()

	for _, cmd := range commands {
		switch cmd.Action {
		case "SETADD":
			set.Add(cmd.Element)
		case "SETDEL":
			set.Delete(cmd.Element)
		case "SET_AT":
			result := set.Contains(cmd.Element)
			fmt.Printf("SET_AT %s: %v\n", cmd.Element, result)
		default:
			fmt.Printf("Unknown command: %s\n", cmd.Action)
		}
	}

	// Process the query
	parts := strings.SplitN(query, " ", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid query format")
	}

	action, element := parts[0], parts[1]

	switch action {
	case "SETADD":
		set.Add(element)
		commands = append(commands, Command{Action: "SETADD", Element: element})
		fmt.Printf("Added %s to the set\n", element)
	case "SETDEL":
		set.Delete(element)
		commands = append(commands, Command{Action: "SETDEL", Element: element})
		fmt.Printf("Deleted %s from the set\n", element)
	case "SET_AT":
		result := set.Contains(element)
		fmt.Printf("SET_AT %s: %v\n", element, result)
	default:
		return fmt.Errorf("unknown command: %s", action)
	}

	// Save changes back to the file
	fileData, err := json.Marshal(commands)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(filePath, fileData, 0644); err != nil {
		return err
	}

	return nil
}

func main() {
	if len(os.Args) != 5 || os.Args[1] != "--file" || os.Args[3] != "--query" {
		fmt.Println("Usage: ./set_program --file <path_to_json_file> --query <query>")
		return
	}

	filePath := os.Args[2]
	query := os.Args[4]

	if err := ProcessFile(filePath, query); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
