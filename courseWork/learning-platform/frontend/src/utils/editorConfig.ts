import { loader } from '@monaco-editor/react';

// Конфигурация Monaco Editor
export const configureMonaco = () => {
  loader.config({
    paths: {
      vs: 'https://cdn.jsdelivr.net/npm/monaco-editor@0.52.2/min/vs'
    }
  });
};

// Тема редактора
export const editorTheme = {
  base: 'vs',
  inherit: true,
  rules: [],
  colors: {
    'editor.background': '#f8f9fa'
  }
};