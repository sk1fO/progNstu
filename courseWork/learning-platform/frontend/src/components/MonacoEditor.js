import React from 'react';
import Editor from '@monaco-editor/react';

function MonacoEditor({ code, onChange, readOnly }) {
  const handleEditorChange = (value) => {
    if (onChange) {
      onChange(value);
    }
  };

  return (
    <div style={editorContainerStyle}>
      <Editor
        height="400px"
        defaultLanguage="cpp"
        value={code}
        onChange={handleEditorChange}
        options={{
          readOnly: readOnly,
          minimap: { enabled: false },
          fontSize: 14,
          lineNumbers: 'on',
          roundedSelection: false,
          scrollBeyondLastLine: false,
          automaticLayout: true,
          theme: 'vs-dark',
          wordWrap: 'on',
          padding: { top: 10 },
        }}
        theme="vs-dark"
      />
    </div>
  );
}

const editorContainerStyle = {
  border: '2px solid #e9ecef',
  borderRadius: '8px',
  overflow: 'hidden',
  marginBottom: '1.5rem'
};

export default MonacoEditor;