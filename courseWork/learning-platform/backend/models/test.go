package models

type TestCase struct {
	Input          string `json:"input"`
	ExpectedOutput string `json:"expected_output"`
	Description    string `json:"description"`
}

type Task struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Difficulty  string     `json:"difficulty"`
	Template    string     `json:"template,omitempty"`
	Tests       []TestCase `json:"tests,omitempty"`
}
