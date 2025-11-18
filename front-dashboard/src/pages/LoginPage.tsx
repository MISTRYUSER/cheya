import { Link } from 'react-router-dom'
import ThemeToggle from '../components/ThemeToggle'

const LoginPage = () => {
  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 dark:from-gray-900 dark:to-blue-900 flex items-center justify-center px-4 relative">
      {/* 主题切换按钮 - 固定在右上角 */}
      <div className="absolute top-6 right-6">
        <ThemeToggle />
      </div>

      {/* 登录卡片 - 毛玻璃效果 */}
      <div className="max-w-md w-full bg-white/80 dark:bg-gray-800/80 backdrop-blur-lg rounded-lg shadow-lg border border-gray-200 dark:border-gray-700 p-8 mx-auto mt-20">
        <h1 className="text-3xl font-bold text-gray-800 dark:text-gray-100 mb-6 text-center">
          登录 - CheYa 平台
        </h1>
        
        <form className="space-y-6">
          <div>
            <label htmlFor="username" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              用户名
            </label>
            <input
              type="text"
              id="username"
              name="username"
              className="w-full px-4 py-3 border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 rounded-lg focus:ring-2 focus:ring-blue-500 dark:focus:ring-blue-400 focus:border-transparent transition duration-200 placeholder-gray-400 dark:placeholder-gray-500"
              placeholder="请输入用户名"
            />
          </div>

          <div>
            <label htmlFor="password" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              密码
            </label>
            <input
              type="password"
              id="password"
              name="password"
              className="w-full px-4 py-3 border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 rounded-lg focus:ring-2 focus:ring-blue-500 dark:focus:ring-blue-400 focus:border-transparent transition duration-200 placeholder-gray-400 dark:placeholder-gray-500"
              placeholder="请输入密码"
            />
          </div>

          <button
            type="submit"
            className="w-full bg-blue-600 dark:bg-blue-700 text-white py-3 rounded-lg font-semibold hover:bg-blue-700 dark:hover:bg-blue-600 transition duration-200 shadow-md hover:shadow-lg"
          >
            登录
          </button>
        </form>

        <p className="mt-6 text-center text-sm text-gray-600 dark:text-gray-400">
          还没有账号？{' '}
          <Link to="/register" className="text-blue-600 dark:text-blue-400 hover:text-blue-700 dark:hover:text-blue-300 font-semibold hover:underline">
            立即注册
          </Link>
        </p>
      </div>
    </div>
  )
}

export default LoginPage


