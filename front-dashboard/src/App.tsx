import { Routes, Route, Navigate } from 'react-router-dom'
import LoginPage from './pages/LoginPage'
import RegisterPage from './pages/RegisterPage'
import SaaSLayout from './layouts/SaaSLayout'
import DashboardOverviewPage from './pages/DashboardOverviewPage'
import DashboardMapPage from './pages/DashboardMapPage'
import DashboardListPage from './pages/DashboardListPage'

function App() {
  return (
    <Routes>
      <Route path="/" element={<Navigate to="/dashboard/overview" replace />} />
      <Route path="/login" element={<LoginPage />} />
      <Route path="/register" element={<RegisterPage />} />
      
      {/* Nested Routes with SaaS Layout */}
      <Route path="/dashboard" element={<SaaSLayout />}>
        <Route path="overview" element={<DashboardOverviewPage />} />
        <Route path="map" element={<DashboardMapPage />} />
        <Route path="list" element={<DashboardListPage />} />
      </Route>
    </Routes>
  )
}

export default App

