import { createContext, useContext, useEffect, useMemo, useState } from "react"

type Theme = "dark" | "light" | "system"

interface ThemeState {
  theme: Theme
  resolvedTheme: "dark" | "light"
  setTheme: (theme: Theme) => void
}

const ThemeContext = createContext<ThemeState>({
  theme: "system",
  resolvedTheme: "light",
  setTheme: () => {},
})

function getSystemTheme(): "dark" | "light" {
  return window.matchMedia("(prefers-color-scheme: dark)").matches ? "dark" : "light"
}

function ThemeProvider({ children, defaultTheme = "system" }: { children: React.ReactNode; defaultTheme?: Theme }) {
  const [theme, setTheme] = useState<Theme>(() => (localStorage.getItem("theme") as Theme) ?? defaultTheme)

  const resolvedTheme = useMemo<"dark" | "light">(() => {
    return theme === "system" ? getSystemTheme() : theme
  }, [theme])

  useEffect(() => {
    const root = document.documentElement
    root.classList.remove("light", "dark")
    root.classList.add(resolvedTheme)
    localStorage.setItem("theme", theme)
  }, [theme, resolvedTheme])

  useEffect(() => {
    const media = window.matchMedia("(prefers-color-scheme: dark)")
    const handleChange = () => {
      if (theme === "system") {
        const systemTheme = media.matches ? "dark" : "light"
        document.documentElement.classList.remove("light", "dark")
        document.documentElement.classList.add(systemTheme)
      }
    }
    media.addEventListener("change", handleChange)
    return () => media.removeEventListener("change", handleChange)
  }, [theme])

  return (
    <ThemeContext.Provider value={{ theme, resolvedTheme, setTheme }}>
      {children}
    </ThemeContext.Provider>
  )
}

function useTheme() {
  return useContext(ThemeContext)
}

export { ThemeProvider, useTheme }