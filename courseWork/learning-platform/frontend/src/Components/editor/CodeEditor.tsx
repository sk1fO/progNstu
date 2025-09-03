import React, { useState } from 'react';
import Editor from '@monaco-editor/react';
import { assignmentsAPI } from '../../services/api';
import { Assignment, TestResult } from '../../types';

interface CodeEditorProps {
  assignment: Assignment;
}

const CodeEditor: React.FC<CodeEditorProps> = ({ assignment }) => {
  const [code, setCode] = useState(assignment.starter_code);
  const [output, setOutput] = useState<string>('');
  const [testResults, setTestResults] = useState<TestResult[]>([]);
  const [isLoading, setIsLoading] = useState(false);

  const handleRunCode = async () => {
    setIsLoading(true);
    setOutput('');
    setTestResults([]);
    
    try {
      const response = await assignmentsAPI.runCode(assignment.id, code);
      const { all_passed, test_results, output: resultOutput } = response.data;
      
      setOutput(resultOutput || '');
      setTestResults(test_results || []);
      
      if (all_passed) {
        setOutput(prev => prev + '\n\n✅ Все тесты пройдены успешно!');
      } else {
        setOutput(prev => prev + '\n\n❌ Некоторые тесты не пройдены');
      }
    } catch (error: any) {
      setOutput(`Ошибка: ${error.response?.data?.error || error.message}`);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div style={{ display: 'flex', flexDirection: 'column', height: '100%' }}>
      <div style={{ marginBottom: '20px' }}>
        <h2>{assignment.title}</h2>
        <p>{assignment.description}</p>
      </div>

      <div style={{ flex: 1, minHeight: '300px', marginBottom: '20px' }}>
        <Editor
          height="100%"
          defaultLanguage={assignment.language || 'python'}
          value={code}
          onChange={(value) => setCode(value || '')}
          options={{
            minimap: { enabled: false },
            fontSize: 14,
            automaticLayout: true,
            scrollBeyondLastLine: false,
          }}
        />
      </div>

      <div style={{ marginBottom: '20px' }}>
        <button
          onClick={handleRunCode}
          disabled={isLoading}
          style={{
            padding: '10px 20px',
            fontSize: '16px',
            backgroundColor: isLoading ? '#6c757d' : '#007bff',
            color: 'white',
            border: 'none',
            borderRadius: '4px',
            cursor: isLoading ? 'not-allowed' : 'pointer'
          }}
        >
          {isLoading ? 'Запуск...' : 'Запустить код'}
        </button>
      </div>

      {output && (
        <div style={{ 
          backgroundColor: '#f8f9fa', 
          padding: '15px', 
          borderRadius: '4px',
          marginBottom: '20px',
          maxHeight: '200px',
          overflow: 'auto'
        }}>
          <h4>Результат выполнения:</h4>
          <pre style={{ margin: 0, whiteSpace: 'pre-wrap' }}>{output}</pre>
        </div>
      )}

      {testResults.length > 0 && (
        <div style={{ 
          backgroundColor: '#f8f9fa', 
          padding: '15px', 
          borderRadius: '4px',
          maxHeight: '300px',
          overflow: 'auto'
        }}>
          <h4>Результаты тестов:</h4>
          {testResults.map((result, index) => (
            <div
              key={index}
              style={{
                padding: '10px',
                margin: '5px 0',
                backgroundColor: result.passed ? '#d4edda' : '#f8d7da',
                border: `1px solid ${result.passed ? '#c3e6cb' : '#f5c6cb'}`,
                borderRadius: '4px'
              }}
            >
              <div><strong>Тест {index + 1}:</strong> {result.passed ? '✅ Пройден' : '❌ Не пройден'}</div>
              {result.input && <div><strong>Ввод:</strong> {result.input}</div>}
              <div><strong>Ожидаемый вывод:</strong> {result.expected}</div>
              <div><strong>Полученный вывод:</strong> {result.actual}</div>
              {result.error_message && (
                <div><strong>Ошибка:</strong> {result.error_message}</div>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

export default CodeEditor;