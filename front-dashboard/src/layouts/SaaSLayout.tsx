import { Outlet, NavLink } from 'react-router-dom'
import ThemeToggle from '../components/ThemeToggle'

const SaaSLayout = () => {
  return (
    <div className="flex h-screen bg-gray-100 dark:bg-gray-900">
      {/* Sidebar - 毛玻璃效果 */}
      <aside className="w-64 bg-gray-800/90 dark:bg-gray-900/80 backdrop-blur-lg h-screen flex flex-col border-r border-gray-700 dark:border-gray-800">
        {/* Logo */}
        <div className="p-6 border-b border-gray-700 dark:border-gray-800">
          <h1 className="text-2xl font-bold text-white">CheYa Logo</h1>
          <p className="text-sm text-gray-400 dark:text-gray-500 mt-1">车辆监控平台</p>
        </div>

        {/* Navigation */}
        <nav className="flex-1 px-4 py-6 space-y-2">
          <NavLink
            to="/dashboard/overview"
            className={({ isActive }) =>
              `flex items-center gap-3 px-4 py-3 rounded-lg transition-colors duration-200 ${
                isActive
                  ? 'bg-blue-600 dark:bg-blue-500 text-white'
                  : 'text-gray-300 dark:text-gray-400 hover:bg-gray-700 dark:hover:bg-gray-800 hover:text-white'
              }`
            }
          >
            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
            </svg>
            <span className="font-medium">概览 (Overview)</span>
          </NavLink>

          <NavLink
            to="/dashboard/map"
            className={({ isActive }) =>
              `flex items-center gap-3 px-4 py-3 rounded-lg transition-colors duration-200 ${
                isActive
                  ? 'bg-blue-600 dark:bg-blue-500 text-white'
                  : 'text-gray-300 dark:text-gray-400 hover:bg-gray-700 dark:hover:bg-gray-800 hover:text-white'
              }`
            }
          >
            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 20l-5.447-2.724A1 1 0 013 16.382V5.618a1 1 0 011.447-.894L9 7m0 13l6-3m-6 3V7m6 10l4.553 2.276A1 1 0 0021 18.382V7.618a1 1 0 00-.553-.894L15 4m0 13V4m0 0L9 7" />
            </svg>
            <span className="font-medium">实时地图 (Map)</span>
          </NavLink>

          <NavLink
            to="/dashboard/list"
            className={({ isActive }) =>
              `flex items-center gap-3 px-4 py-3 rounded-lg transition-colors duration-200 ${
                isActive
                  ? 'bg-blue-600 dark:bg-blue-500 text-white'
                  : 'text-gray-300 dark:text-gray-400 hover:bg-gray-700 dark:hover:bg-gray-800 hover:text-white'
              }`
            }
          >
            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01" />
            </svg>
            <span className="font-medium">车辆列表 (List)</span>
          </NavLink>
        </nav>

        {/* Footer */}
        <div className="p-4 border-t border-gray-700 dark:border-gray-800">
          <p className="text-xs text-gray-500 dark:text-gray-600 text-center">© 2025 CheYa Platform</p>
        </div>
      </aside>

      {/* Content Area */}
      <div className="flex-1 flex flex-col overflow-hidden">
        {/* Header - 毛玻璃效果 + 浮动 */}
        <header className="sticky top-0 z-10 bg-white/70 dark:bg-gray-800/70 backdrop-blur-lg shadow-md px-6 py-4 flex items-center justify-between border-b border-gray-200 dark:border-gray-700">
          <h2 className="text-xl font-semibold text-gray-800 dark:text-gray-100">
            CheYa - 自动驾驶卡车监控平台
          </h2>
          <div className="flex items-center gap-4">
            {/* 主题切换按钮 */}
            <ThemeToggle />
            
            {/* 登出按钮 */}
            <button className="bg-red-600 dark:bg-red-700 text-white px-6 py-2 rounded-lg font-semibold hover:bg-red-700 dark:hover:bg-red-800 transition duration-200 shadow-md">
              登出
            </button>
          </div>
        </header>

        {/* Main Content - This is where child routes render */}
        <main className="flex-1 overflow-y-auto bg-gray-100 dark:bg-gray-900 p-6">
          <Outlet />
        </main>
      </div>
    </div>
  )
}

export default SaaSLayout

