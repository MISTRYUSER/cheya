import { useState, FormEvent } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import ThemeToggle from '../components/ThemeToggle'

interface RegisterResponse {
  code: number
  message: string
  user_id?: number
  username?: string
  created_at?: string
  error?: string
}

const RegisterPage = () => {
  const navigate = useNavigate()
  const [formData, setFormData] = useState({
    username: '',
    password: '',
    confirmPassword: '',
    email: '',
  })
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [success, setSuccess] = useState(false)

  const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    setError('')

    // 基础验证
    if (!formData.username || !formData.password || !formData.confirmPassword) {
      setError('请填写所有必填字段')
      return
    }

    if (formData.password !== formData.confirmPassword) {
      setError('两次输入的密码不一致')
      return
    }

    if (formData.password.length < 6) {
      setError('密码长度至少为6位')
      return
    }

    // 邮箱格式验证（如果填写了邮箱）
    if (formData.email) {
      const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
      if (!emailRegex.test(formData.email)) {
        setError('邮箱格式不正确')
        return
      }
    }

    setLoading(true)

    try {
      const payload: {
        username: string
        password: string
        confirm_password: string
        email?: string
      } = {
        username: formData.username,
        password: formData.password,
        confirm_password: formData.confirmPassword,
      }

      if (formData.email) {
        payload.email = formData.email
      }

      const response = await fetch('/api/v1/auth/register', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(payload),
      })

      const data: RegisterResponse = await response.json()

      if (!response.ok || data.code !== 201) {
        // 优先使用 message，其次 error
        const errorMsg = data.message || data.error || '注册失败，请重试'
        setError(errorMsg)
        return
      }

      // 注册成功，显示成功消息并跳转到登录页面
      setSuccess(true)
      setTimeout(() => {
        navigate('/login')
      }, 1500)
    } catch (err) {
      setError('网络错误，请检查连接后重试')
      console.error('注册错误:', err)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-green-50 to-blue-100 dark:from-gray-900 dark:to-green-900 flex items-center justify-center px-4 relative">
      {/* 主题切换按钮 - 固定在右上角 */}
      <div className="absolute top-6 right-6">
        <ThemeToggle />
      </div>

      {/* 注册卡片 - 毛玻璃效果 */}
      <div className="max-w-md w-full bg-white/80 dark:bg-gray-800/80 backdrop-blur-lg rounded-lg shadow-lg border border-gray-200 dark:border-gray-700 p-8 mx-auto mt-20">
        <h1 className="text-3xl font-bold text-gray-800 dark:text-gray-100 mb-6 text-center">
          创建新账号
        </h1>

        {/* 成功提示 */}
        {success && (
          <div className="mb-4 p-3 bg-green-100 dark:bg-green-900/30 border border-green-400 dark:border-green-700 text-green-700 dark:text-green-400 rounded-lg text-sm flex items-center gap-2">
            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
            </svg>
            注册成功！正在跳转到登录页面...
          </div>
        )}

        {/* 错误提示 */}
        {error && !success && (
          <div className="mb-4 p-3 bg-red-100 dark:bg-red-900/30 border border-red-400 dark:border-red-700 text-red-700 dark:text-red-400 rounded-lg text-sm flex items-center gap-2">
            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit} className="space-y-6">
          <div>
            <label htmlFor="username" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              用户名
            </label>
            <input
              type="text"
              id="username"
              name="username"
              value={formData.username}
              onChange={(e) => setFormData({ ...formData, username: e.target.value })}
              className="w-full px-4 py-3 border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 rounded-lg focus:ring-2 focus:ring-green-500 dark:focus:ring-green-400 focus:border-transparent transition duration-200 placeholder-gray-400 dark:placeholder-gray-500"
              placeholder="请选择用户名"
              required
              disabled={loading}
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
              value={formData.password}
              onChange={(e) => setFormData({ ...formData, password: e.target.value })}
              className="w-full px-4 py-3 border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 rounded-lg focus:ring-2 focus:ring-green-500 dark:focus:ring-green-400 focus:border-transparent transition duration-200 placeholder-gray-400 dark:placeholder-gray-500"
              placeholder="创建密码"
              required
              minLength={6}
              disabled={loading}
            />
          </div>

          <div>
            <label htmlFor="confirmPassword" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              确认密码
            </label>
            <input
              type="password"
              id="confirmPassword"
              name="confirmPassword"
              value={formData.confirmPassword}
              onChange={(e) => setFormData({ ...formData, confirmPassword: e.target.value })}
              className="w-full px-4 py-3 border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 rounded-lg focus:ring-2 focus:ring-green-500 dark:focus:ring-green-400 focus:border-transparent transition duration-200 placeholder-gray-400 dark:placeholder-gray-500"
              placeholder="再次输入密码"
              required
              disabled={loading}
            />
          </div>

          <div>
            <label htmlFor="email" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              邮箱（可选）
            </label>
            <input
              type="email"
              id="email"
              name="email"
              value={formData.email}
              onChange={(e) => setFormData({ ...formData, email: e.target.value })}
              className="w-full px-4 py-3 border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 rounded-lg focus:ring-2 focus:ring-green-500 dark:focus:ring-green-400 focus:border-transparent transition duration-200 placeholder-gray-400 dark:placeholder-gray-500"
              placeholder="请输入邮箱（可选）"
              disabled={loading}
            />
          </div>

          <button
            type="submit"
            disabled={loading}
            className="w-full bg-green-600 dark:bg-green-700 text-white py-3 rounded-lg font-semibold hover:bg-green-700 dark:hover:bg-green-600 transition duration-200 shadow-md hover:shadow-lg disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {loading ? '注册中...' : '注册'}
          </button>
        </form>

        <p className="mt-6 text-center text-sm text-gray-600 dark:text-gray-400">
          已有账号？{' '}
          <Link to="/login" className="text-green-600 dark:text-green-400 hover:text-green-700 dark:hover:text-green-300 font-semibold hover:underline">
            立即登录
          </Link>
        </p>
      </div>
    </div>
  )
}

export default RegisterPage


