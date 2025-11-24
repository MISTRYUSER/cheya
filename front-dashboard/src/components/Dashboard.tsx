/**
 * Dashboard 组件
 * 实时显示车辆列表和遥测数据
 */

import { useEffect, useState, useRef } from 'react';
import { fetchVehicleList, Vehicle, VehicleStatus, controlVehicle, ControlAction } from '../api/vehicle';

/** WebSocket 遥测数据接口 */
interface TelemetryData {
  vehicle_id: string;
  timestamp: number;
  latitude: number;
  longitude: number;
  speed: number;
}

/** 扩展车辆数据（包含实时遥测） */
interface VehicleWithTelemetry extends Vehicle {
  latitude?: number;
  longitude?: number;
  speed?: number;
  lastUpdate?: number;
}

/**
 * 获取车辆状态的显示文本和样式
 */
function getStatusDisplay(status: VehicleStatus): { text: string; className: string } {
  switch (status) {
    case VehicleStatus.ONLINE:
      return { text: '在线', className: 'bg-green-100 text-green-800' };
    case VehicleStatus.OFFLINE:
      return { text: '离线', className: 'bg-gray-100 text-gray-800' };
    default:
      return { text: '未知', className: 'bg-yellow-100 text-yellow-800' };
  }
}

/**
 * Dashboard 组件
 */
export default function Dashboard() {
  const [vehicles, setVehicles] = useState<VehicleWithTelemetry[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [wsConnected, setWsConnected] = useState(false);
  const wsRef = useRef<WebSocket | null>(null);
  const [notification, setNotification] = useState<{ message: string; type: 'success' | 'error' } | null>(null);

  // 显示通知（3秒后自动消失）
  const showNotification = (message: string, type: 'success' | 'error') => {
    setNotification({ message, type });
    setTimeout(() => setNotification(null), 3000);
  };

  // 车辆控制处理函数
  const handleControlVehicle = async (vin: string, action: ControlAction) => {
    try {
      await controlVehicle(vin, action);
      const actionText = action === 'STOP' ? '紧急停车' : '恢复运行';
      showNotification(`✅ ${actionText}指令已发送！`, 'success');
    } catch (err) {
      const actionText = action === 'STOP' ? '紧急停车' : '恢复运行';
      showNotification(`❌ ${actionText}指令发送失败`, 'error');
    }
  };

  // 获取车辆列表
  useEffect(() => {
    async function loadVehicles() {
      try {
        setLoading(true);
        setError(null);
        const response = await fetchVehicleList(1, 100);
        
        if (response.code === 200) {
          const vehicleList = response.data.items || [];
          setVehicles(vehicleList);
        } else {
          setError('获取车辆列表失败');
        }
      } catch (err) {
        setError(err instanceof Error ? err.message : '未知错误');
      } finally {
        setLoading(false);
      }
    }

    loadVehicles();
  }, []);

  // 定时刷新以检测过期数据
  useEffect(() => {
    const interval = setInterval(() => {
      // 触发重新渲染以更新过期状态
      setVehicles((prev) => [...prev]);
    }, 1000); // 每秒检查一次

    return () => clearInterval(interval);
  }, []);

  // WebSocket 连接
  useEffect(() => {
    // 创建 WebSocket 连接
    const ws = new WebSocket('ws://localhost:8081/ws');
    wsRef.current = ws;

    ws.onopen = () => {
      console.log('✅ WebSocket 连接成功');
      setWsConnected(true);
    };

    ws.onmessage = (event) => {
      try {
        const telemetryData: TelemetryData = JSON.parse(event.data);
        console.log('📩 收到遥测数据:', telemetryData);

        // 更新匹配 VIN 的车辆数据
        setVehicles((prevVehicles) =>
          prevVehicles.map((vehicle) =>
            vehicle.vin === telemetryData.vehicle_id
              ? {
                  ...vehicle,
                  latitude: telemetryData.latitude,
                  longitude: telemetryData.longitude,
                  speed: telemetryData.speed,
                  lastUpdate: telemetryData.timestamp,
                }
              : vehicle
          )
        );
      } catch (err) {
        console.error('❌ 解析 WebSocket 数据失败:', err);
      }
    };

    ws.onerror = (error) => {
      console.error('❌ WebSocket 错误:', error);
      setWsConnected(false);
    };

    ws.onclose = () => {
      console.log('🔌 WebSocket 连接关闭');
      setWsConnected(false);
    };

    // 清理函数
    return () => {
      if (ws.readyState === WebSocket.OPEN) {
        ws.close();
      }
    };
  }, []);

  // 加载状态
  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen bg-gray-50">
        <div className="text-center">
          <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
          <p className="mt-4 text-gray-600">加载中...</p>
        </div>
      </div>
    );
  }

  // 错误状态
  if (error) {
    return (
      <div className="flex items-center justify-center min-h-screen bg-gray-50">
        <div className="bg-red-50 border border-red-200 rounded-lg p-6 max-w-md">
          <h2 className="text-red-800 text-xl font-semibold mb-2">错误</h2>
          <p className="text-red-600">{error}</p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50 py-8 px-4 sm:px-6 lg:px-8">
      {/* 通知提示 */}
      {notification && (
        <div
          className={`fixed top-4 right-4 z-50 px-6 py-4 rounded-lg shadow-lg animate-fade-in ${
            notification.type === 'success'
              ? 'bg-green-500 text-white'
              : 'bg-red-500 text-white'
          }`}
        >
          <div className="flex items-center gap-2">
            <span className="text-lg">{notification.message}</span>
          </div>
        </div>
      )}

      <div className="max-w-7xl mx-auto">
        {/* 头部 */}
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900">车辆实时监控</h1>
          <div className="mt-2 flex items-center gap-4">
            <p className="text-gray-600">
              总计: <span className="font-semibold">{vehicles.length}</span> 辆车
            </p>
            <div className="flex items-center gap-2">
              <div
                className={`w-3 h-3 rounded-full ${
                  wsConnected ? 'bg-green-500 animate-pulse' : 'bg-red-500'
                }`}
              ></div>
              <span className="text-sm text-gray-600">
                {wsConnected ? 'WebSocket 已连接' : 'WebSocket 未连接'}
              </span>
            </div>
          </div>
        </div>

        {/* 车辆列表 */}
        {vehicles.length === 0 ? (
          <div className="bg-white rounded-lg shadow p-12 text-center">
            <p className="text-gray-500 text-lg">暂无车辆数据</p>
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {vehicles.map((vehicle) => {
              const statusDisplay = getStatusDisplay(vehicle.status);
              
              // 检查数据是否过期（超过 5 秒没有更新）
              const now = Math.floor(Date.now() / 1000);
              const isDataStale = vehicle.lastUpdate && (now - vehicle.lastUpdate) > 5;
              const hasRealtimeData = vehicle.speed !== undefined && !isDataStale;
              
              return (
                <div
                  key={vehicle.id}
                  className="bg-white rounded-lg shadow-md hover:shadow-lg transition-shadow p-6"
                >
                  {/* 车辆基本信息 */}
                  <div className="flex items-start justify-between mb-4">
                    <div>
                      <h3 className="text-lg font-semibold text-gray-900">
                        {vehicle.license_plate}
                      </h3>
                      <p className="text-sm text-gray-500 mt-1">VIN: {vehicle.vin}</p>
                    </div>
                    <span
                      className={`px-3 py-1 rounded-full text-xs font-medium ${statusDisplay.className}`}
                    >
                      {statusDisplay.text}
                    </span>
                  </div>

                  {/* 实时遥测数据 */}
                  {hasRealtimeData ? (
                    <div className="border-t pt-4 space-y-2">
                      <div className="flex justify-between items-center">
                        <span className="text-sm text-gray-600">速度</span>
                        <span className="text-lg font-semibold text-blue-600">
                          {vehicle.speed!.toFixed(1)} km/h
                        </span>
                      </div>
                      <div className="flex justify-between items-center">
                        <span className="text-sm text-gray-600">经度</span>
                        <span className="text-sm font-mono text-gray-900">
                          {vehicle.longitude?.toFixed(4)}
                        </span>
                      </div>
                      <div className="flex justify-between items-center">
                        <span className="text-sm text-gray-600">纬度</span>
                        <span className="text-sm font-mono text-gray-900">
                          {vehicle.latitude?.toFixed(4)}
                        </span>
                      </div>
                      <div className="flex justify-between items-center text-xs mt-2">
                        <span className="text-green-600 flex items-center gap-1">
                          <span className="inline-block w-2 h-2 bg-green-500 rounded-full animate-pulse"></span>
                          实时更新中
                        </span>
                        {vehicle.lastUpdate && (
                          <span className="text-gray-400">
                            {new Date(vehicle.lastUpdate * 1000).toLocaleTimeString()}
                          </span>
                        )}
                      </div>
                    </div>
                  ) : vehicle.speed !== undefined && isDataStale ? (
                    <div className="border-t pt-4 text-center">
                      <p className="text-sm text-orange-600 font-medium">⏸️ 车辆已停止</p>
                      <p className="text-xs text-gray-400 mt-1">
                        最后更新: {vehicle.lastUpdate ? new Date(vehicle.lastUpdate * 1000).toLocaleTimeString() : '未知'}
                      </p>
                    </div>
                  ) : (
                    <div className="border-t pt-4 text-center">
                      <p className="text-sm text-gray-400">暂无实时数据</p>
                    </div>
                  )}

                  {/* 控制按钮 */}
                  <div className="border-t mt-4 pt-4 flex gap-2">
                    <button
                      onClick={() => handleControlVehicle(vehicle.vin, 'STOP')}
                      className="flex-1 text-red-600 hover:text-red-800 font-bold border border-red-200 px-3 py-2 rounded hover:bg-red-50 transition-colors text-sm"
                    >
                      🛑 紧急停车
                    </button>
                    <button
                      onClick={() => handleControlVehicle(vehicle.vin, 'START')}
                      className="flex-1 text-green-600 hover:text-green-800 font-bold border border-green-200 px-3 py-2 rounded hover:bg-green-50 transition-colors text-sm"
                    >
                      ▶️ 恢复
                    </button>
                  </div>
                </div>
              );
            })}
          </div>
        )}
      </div>
    </div>
  );
}

