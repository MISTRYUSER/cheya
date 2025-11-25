import { Routes, Route, Navigate } from 'react-router-dom'
import LoginPage from './pages/LoginPage'
import RegisterPage from './pages/RegisterPage'
import SaaSLayout from './layouts/SaaSLayout'
import DashboardOverviewPage from './pages/DashboardOverviewPage'
import DashboardMapPage from './pages/DashboardMapPage'
import DashboardListPage from './pages/DashboardListPage'
import Dashboard from './components/Dashboard'

/**
 * 检查用户是否已登录
 */
function isAuthenticated(): boolean {
  const token = localStorage.getItem('access_token')
  const expiresAt = localStorage.getItem('token_expires_at')
  
  if (!token || !expiresAt) {
    return false
  }
  
  // 检查 token 是否过期
  const now = Date.now()
  if (now > parseInt(expiresAt)) {
    // Token 已过期，清除存储
    localStorage.removeItem('access_token')
    localStorage.removeItem('username')
    localStorage.removeItem('token_expires_at')
    return false
  }
  
  return true
}

/**
 * 受保护的路由组件
 */
function ProtectedRoute({ children }: { children: React.ReactNode }) {
  if (!isAuthenticated()) {
    return <Navigate to="/login" replace />
  }
  return <>{children}</>
}

/**
 * 公开路由组件（已登录用户访问时自动跳转）
 */
function PublicRoute({ children }: { children: React.ReactNode }) {
  if (isAuthenticated()) {
    return <Navigate to="/" replace />
  }
  return <>{children}</>
}

function App() {
  return (
    <Routes>
      {/* 主页 - 需要登录 */}
      <Route
        path="/"
        element={
          <ProtectedRoute>
            <Dashboard />
          </ProtectedRoute>
        }
      />
      
      {/* 登录页 - 已登录用户自动跳转 */}
      <Route
        path="/login"
        element={
          <PublicRoute>
            <LoginPage />
          </PublicRoute>
        }
      />
      
      {/* 注册页 - 已登录用户自动跳转 */}
      <Route
        path="/register"
        element={
          <PublicRoute>
            <RegisterPage />
          </PublicRoute>
        }
      />
      
      {/* Dashboard 布局 - 需要登录 */}
      <Route
        path="/dashboard"
        element={
          <ProtectedRoute>
            <SaaSLayout />
          </ProtectedRoute>
        }
      >
        <Route path="overview" element={<DashboardOverviewPage />} />
        <Route path="map" element={<DashboardMapPage />} />
        <Route path="list" element={<DashboardListPage />} />
      </Route>
    </Routes>
  )
}

export default App

