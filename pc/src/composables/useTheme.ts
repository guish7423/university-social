import { ref, watch } from "vue"

type Theme = "dark" | "light"

const STORAGE_KEY = "university-social-theme"
const currentTheme = ref<Theme>(loadTheme())

function loadTheme(): Theme {
  const stored = localStorage.getItem(STORAGE_KEY)
  if (stored === "dark" || stored === "light") return stored
  return "dark"
}

function applyTheme(theme: Theme) {
  if (theme === "light") {
    document.documentElement.setAttribute("data-theme", "light")
  } else {
    document.documentElement.removeAttribute("data-theme")
  }
}

export function useTheme() {
  function toggle() {
    currentTheme.value = currentTheme.value === "dark" ? "light" : "dark"
  }

  function setTheme(theme: Theme) {
    currentTheme.value = theme
  }

  // Apply on load (if already light from previous session)
  if (currentTheme.value === "light") {
    applyTheme("light")
  }

  // Watch for changes
  watch(currentTheme, (theme) => {
    localStorage.setItem(STORAGE_KEY, theme)
    applyTheme(theme)
  })

  return {
    theme: currentTheme,
    isDark: currentTheme,
    toggle,
    setTheme,
  }
}
