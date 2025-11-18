import { useState } from 'react'

interface Truck {
  id: string
  vin: string
  licensePlate: string
  status: 'Online' | 'Offline' | 'Maintaining'
  lastHeartbeat: string
}

const mockTrucks: Truck[] = [
  { 
    id: 'T-001', 
    vin: 'LFV123...', 
    licensePlate: '沪A88888', 
    status: 'Online', 
    lastHeartbeat: '2025-11-11 06:20:00' 
  },
  { 
    id: 'T-002', 
    vin: 'LFV456...', 
    licensePlate: '沪B77777', 
    status: 'Offline', 
    lastHeartbeat: '2025-11-10 18:00:00' 
  },
  { 
    id: 'T-003', 
    vin: 'LFV789...', 
    licensePlate: '沪C66666', 
    status: 'Maintaining', 
    lastHeartbeat: '2025-11-11 01:15:00' 
  }
]

const StatusIndicator = ({ status }: { status: Truck['status'] }) => {
  const statusConfig = {
    Online: 'bg-green-500 dark:bg-green-400',
    Offline: 'bg-gray-400 dark:bg-gray-500',
    Maintaining: 'bg-yellow-500 dark:bg-yellow-400'
  }

  return (
    <div className="flex items-center gap-2">
      <span className={`w-3 h-3 rounded-full ${statusConfig[status]}`}></span>
      <span className="font-medium text-gray-800 dark:text-gray-200">{status}</span>
    </div>
  )
}

const DashboardOverviewPage = () => {
  const [trucks] = useState<Truck[]>(mockTrucks)

  return (
    <div className="space-y-6">
      {/* Map Placeholder */}
      <div className="bg-gray-300 dark:bg-gray-700 rounded-lg shadow-md h-[500px] flex items-center justify-center border border-gray-200 dark:border-gray-600">
        <p className="text-gray-700 dark:text-gray-300 text-xl font-semibold">
          [地图占位符：实时车队位置将显示在这里]
        </p>
      </div>

      {/* 2-Column Layout */}
      <div className="flex gap-6">
        {/* Left Column - Truck List (70%) */}
        <div className="flex-1" style={{ width: '70%' }}>
          <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6 border border-gray-200 dark:border-gray-700">
            <h2 className="text-xl font-bold text-gray-900 dark:text-gray-100 mb-4">活跃车队</h2>
            
            <div className="overflow-x-auto">
              <table className="w-full">
                <thead>
                  <tr className="border-b-2 border-gray-200 dark:border-gray-600">
                    <th className="text-left py-3 px-4 font-semibold text-gray-700 dark:text-gray-300">车辆识别码</th>
                    <th className="text-left py-3 px-4 font-semibold text-gray-700 dark:text-gray-300">车牌号</th>
                    <th className="text-left py-3 px-4 font-semibold text-gray-700 dark:text-gray-300">状态</th>
                    <th className="text-left py-3 px-4 font-semibold text-gray-700 dark:text-gray-300">最后心跳</th>
                    <th className="text-left py-3 px-4 font-semibold text-gray-700 dark:text-gray-300">操作</th>
                  </tr>
                </thead>
                <tbody>
                  {trucks.map((truck) => (
                    <tr key={truck.id} className="border-b border-gray-100 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-700 transition">
                      <td className="py-4 px-4 text-gray-800 dark:text-gray-300">{truck.vin}</td>
                      <td className="py-4 px-4 text-gray-800 dark:text-gray-300">{truck.licensePlate}</td>
                      <td className="py-4 px-4">
                        <StatusIndicator status={truck.status} />
                      </td>
                      <td className="py-4 px-4 text-gray-600 dark:text-gray-400">{truck.lastHeartbeat}</td>
                      <td className="py-4 px-4">
                        <div className="flex gap-2">
                          <button className="bg-green-600 dark:bg-green-700 text-white px-3 py-1 rounded text-sm font-medium hover:bg-green-700 dark:hover:bg-green-600 transition duration-150">
                            上线
                          </button>
                          <button className="bg-orange-600 dark:bg-orange-700 text-white px-3 py-1 rounded text-sm font-medium hover:bg-orange-700 dark:hover:bg-orange-600 transition duration-150">
                            下线
                          </button>
                        </div>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
        </div>

        {/* Right Column - Instrument Panel (30%) */}
        <div className="w-full" style={{ maxWidth: '30%' }}>
          <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6 border border-gray-200 dark:border-gray-700">
            <h2 className="text-xl font-bold text-gray-900 dark:text-gray-100 mb-6">车辆详情（模拟）</h2>
            
            <div className="space-y-6">
              {/* Speed */}
              <div className="text-center py-4 bg-blue-50 dark:bg-blue-900/30 rounded-lg border border-blue-100 dark:border-blue-800">
                <p className="text-sm text-gray-600 dark:text-gray-400 mb-1">速度</p>
                <p className="text-4xl font-bold text-blue-600 dark:text-blue-400">88 km/h</p>
              </div>

              {/* Battery */}
              <div>
                <div className="flex justify-between items-center mb-2">
                  <p className="text-sm font-medium text-gray-700 dark:text-gray-300">电池电量</p>
                  <p className="text-lg font-bold text-gray-800 dark:text-gray-200">72%</p>
                </div>
                <div className="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-3 overflow-hidden">
                  <div 
                    className="bg-gradient-to-r from-green-400 to-green-600 dark:from-green-500 dark:to-green-700 h-3 rounded-full transition-all duration-300"
                    style={{ width: '72%' }}
                  ></div>
                </div>
              </div>

              {/* GPS */}
              <div className="py-4 px-4 bg-gray-50 dark:bg-gray-700/50 rounded-lg border border-gray-200 dark:border-gray-600">
                <p className="text-sm text-gray-600 dark:text-gray-400 mb-1">GPS 坐标</p>
                <p className="text-lg font-mono text-gray-800 dark:text-gray-200">[30.123, 121.456]</p>
              </div>

              {/* Current Mode */}
              <div className="py-4 px-4 bg-indigo-50 dark:bg-indigo-900/30 rounded-lg border-2 border-indigo-200 dark:border-indigo-700">
                <p className="text-sm text-gray-600 dark:text-gray-400 mb-1">当前模式</p>
                <p className="text-xl font-bold text-indigo-700 dark:text-indigo-400">自动驾驶 (L4)</p>
              </div>

              {/* Additional Stats */}
              <div className="grid grid-cols-2 gap-4 pt-4">
                <div className="text-center py-3 bg-gray-50 dark:bg-gray-700/50 rounded-lg border border-gray-200 dark:border-gray-600">
                  <p className="text-xs text-gray-600 dark:text-gray-400">燃油</p>
                  <p className="text-lg font-bold text-gray-800 dark:text-gray-200">45L</p>
                </div>
                <div className="text-center py-3 bg-gray-50 dark:bg-gray-700/50 rounded-lg border border-gray-200 dark:border-gray-600">
                  <p className="text-xs text-gray-600 dark:text-gray-400">温度</p>
                  <p className="text-lg font-bold text-gray-800 dark:text-gray-200">22°C</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

export default DashboardOverviewPage

