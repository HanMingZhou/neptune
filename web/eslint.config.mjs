import js from '@eslint/js'
import eslintConfigPrettier from 'eslint-config-prettier'
import pluginVue from 'eslint-plugin-vue'
import globals from 'globals'
import tseslint from 'typescript-eslint'
import vueParser from 'vue-eslint-parser'

const sharedGlobals = {
  ...globals.browser,
  ...globals.node,
  ...globals.es2024
}

const tsRecommendedConfigs = tseslint.configs.recommended.map((config) => ({
  ...config,
  files: ['**/*.ts']
}))

export default [
  {
    ignores: ['coverage/**', 'dist/**', 'node_modules/**', 'src/generated/**']
  },
  {
    files: ['**/*.{ts,vue,mjs}'],
    languageOptions: {
      ecmaVersion: 'latest',
      sourceType: 'module',
      globals: sharedGlobals
    }
  },
  {
    files: ['**/*.mjs'],
    ...js.configs.recommended
  },
  ...tsRecommendedConfigs,
  ...pluginVue.configs['flat/essential'],
  {
    files: ['**/*.vue'],
    languageOptions: {
      parser: vueParser,
      parserOptions: {
        ecmaVersion: 'latest',
        sourceType: 'module',
        parser: tseslint.parser,
        extraFileExtensions: ['.vue']
      },
      globals: sharedGlobals
    }
  },
  {
    files: ['src/**/*.test.ts'],
    languageOptions: {
      globals: {
        ...sharedGlobals,
        ...globals.vitest
      }
    }
  },
  {
    files: ['**/*.{ts,vue,mjs}'],
    plugins: {
      '@typescript-eslint': tseslint.plugin
    },
    rules: {
      'no-console': 'off',
      'vue/no-mutating-props': 'off',
      'vue/multi-word-component-names': 'off',
      '@typescript-eslint/no-explicit-any': 'off',
      '@typescript-eslint/no-unused-vars': [
        'error',
        {
          argsIgnorePattern: '^_',
          caughtErrorsIgnorePattern: '^_',
          varsIgnorePattern: '^_'
        }
      ]
    }
  },
  eslintConfigPrettier
]
