import { useState, useEffect } from 'react'

type Theme = 'light' | 'dark'

const ThemeToggle = () => {
  const [theme, setTheme] = useState<Theme>('light')

  // 初始化主题
  useEffect(() => {
    // 1. 优先从 localStorage 读取
    const savedTheme = localStorage.getItem('theme') as Theme | null
    
    // 2. 如果没有保存的主题，检查系统偏好
    const systemPrefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches
    
    const initialTheme = savedTheme || (systemPrefersDark ? 'dark' : 'light')
    
    setTheme(initialTheme)
    applyTheme(initialTheme)
  }, [])

  // 应用主题到 HTML 标签
  const applyTheme = (newTheme: Theme) => {
    const html = document.documentElement
    
    if (newTheme === 'dark') {
      html.classList.add('dark')
    } else {
      html.classList.remove('dark')
    }
  }

  // 切换主题
  const toggleTheme = () => {
    const newTheme: Theme = theme === 'light' ? 'dark' : 'light'
    
    setTheme(newTheme)
    applyTheme(newTheme)
    localStorage.setItem('theme', newTheme)
  }

  return (
    <button
      onClick={toggleTheme}
      className="p-2 rounded-lg bg-gray-100 hover:bg-gray-200 dark:bg-gray-700 dark:hover:bg-gray-600 transition-colors duration-200"
      aria-label="切换主题"
      title={theme === 'light' ? '切换到暗黑模式' : '切换到亮色模式'}
    >
      {theme === 'light' ? (
        // 月亮图标（亮色模式时显示）
        <svg 
          className="w-5 h-5 text-gray-700" 
          fill="none" 
          stroke="currentColor" 
          viewBox="0 0 24 24"
        >
          <path 
            strokeLinecap="round" 
            strokeLinejoin="round" 
            strokeWidth={2} 
            d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z" 
          />
        </svg>
      ) : (
        // 太阳图标（暗黑模式时显示）
        <svg 
          className="w-5 h-5 text-yellow-400" 
          fill="none" 
          stroke="currentColor" 
          viewBox="0 0 24 24"
        >
          <path 
            strokeLinecap="round" 
            strokeLinejoin="round" 
            strokeWidth={2} 
            d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z" 
          />
        </svg>
      )}
    </button>
  )
}

export default ThemeToggle



















