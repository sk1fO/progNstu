import React, { useState } from 'react';
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter';
import { atomDark } from 'react-syntax-highlighter/dist/esm/styles/prism';

function CodeEditor({ initialCode, onCodeChange, readOnly = false }) {
  const [code, setCode] = useState(initialCode || '');

  const handleCodeChange = (e) => {
    const newCode = e.target.value;
    setCode(newCode);
    if (onCodeChange) {
      onCodeChange(newCode);
    }
  };

  if (readOnly) {
    return (
      <div style={readOnlyStyle}>
        <SyntaxHighlighter
          language="cpp"
          style={atomDark}
          customStyle={{
            borderRadius: '8px',
            padding: '1rem',
            fontSize: '14px',
            minHeight: '300px',
            margin: 0
          }}
        >
          {code}
        </SyntaxHighlighter>
      </div>
    );
  }

  return (
    <div style={editorContainerStyle}>
      {/* Textarea - прозрачный текст, но видимый курсор */}
      <textarea
        value={code}
        onChange={handleCodeChange}
        placeholder="Введите ваш код на C++"
        style={textareaStyle}
        rows={15}
        spellCheck="false"
      />
      
      {/* Подсветка синтаксиса - абсолютное позиционирование */}
      <div style={highlightContainerStyle}>
        <pre style={preStyle}>
          <code>
            <SyntaxHighlighter
              language="cpp"
              style={atomDark}
              customStyle={{
                background: 'transparent',
                padding: '1rem',
                margin: 0,
                fontSize: '14px'
              }}
              PreTag="div"
              CodeTag="div"
            >
              {code}
            </SyntaxHighlighter>
          </code>
        </pre>
      </div>
    </div>
  );
}

const editorContainerStyle = {
  position: 'relative',
  width: '100%',
  marginBottom: '1.5rem',
  borderRadius: '8px',
  overflow: 'hidden',
  minHeight: '300px',
  backgroundColor: '#2d3748'
};

const readOnlyStyle = {
  width: '100%',
  marginBottom: '1.5rem',
  borderRadius: '8px',
  overflow: 'hidden'
};

const textareaStyle = {
  position: 'relative',
  zIndex: 2,
  width: '100%',
  height: '100%',
  padding: '1rem',
  border: 'none',
  borderRadius: '8px',
  fontSize: '14px',
  lineHeight: '1.5',
  fontFamily: '"Monaco", "Menlo", "Ubuntu Mono", monospace',
  resize: 'vertical',
  minHeight: '300px',
  backgroundColor: 'transparent',
  color: 'transparent',
  caretColor: '#fff',
  outline: 'none'
};

const highlightContainerStyle = {
  position: 'absolute',
  top: 0,
  left: 0,
  right: 0,
  bottom: 0,
  zIndex: 1,
  pointerEvents: 'none',
  overflow: 'auto'
};

const preStyle = {
  margin: 0,
  padding: '1rem',
  fontSize: '14px',
  lineHeight: '1.5',
  fontFamily: '"Monaco", "Menlo", "Ubuntu Mono", monospace'
};

export default CodeEditor;