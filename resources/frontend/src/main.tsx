import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { FluentProvider, webLightTheme, webDarkTheme } from '@fluentui/react-components'
import './index.css'
import App from './App.tsx'
import { ThemeProvider, useTheme } from './components/ThemeProvider.tsx'
import { AuthProvider } from './components/AuthProvider.tsx'
import { FluentToaster } from './lib/toast.tsx'

function ThemedApp() {
  const { resolvedTheme } = useTheme()

  return (
    <FluentProvider theme={resolvedTheme === 'dark' ? webDarkTheme : webLightTheme}>
      <App />
      <FluentToaster />
    </FluentProvider>
  )
}

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <ThemeProvider defaultTheme="system">
      <AuthProvider>
        <ThemedApp />
      </AuthProvider>
    </ThemeProvider>
  </StrictMode>,
)