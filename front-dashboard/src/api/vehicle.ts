/**
 * Vehicle API 接口
 * 用于获取车辆列表数据
 */

/** 车辆状态枚举 */
export enum VehicleStatus {
  UNSPECIFIED = 0,
  OFFLINE = 1,
  ONLINE = 2,
}

/** 车辆数据接口 */
export interface Vehicle {
  id: string;
  vin: string;
  license_plate: string;
  status: VehicleStatus;
}

/** 车辆列表响应接口 */
export interface VehicleListResponse {
  code: number;
  data: {
    items: Vehicle[] | null;
    total: number;
  };
}

/**
 * 获取车辆列表
 * @param page - 页码（默认 1）
 * @param pageSize - 每页数量（默认 10）
 * @returns 车辆列表响应
 */
export async function fetchVehicleList(
  page: number = 1,
  pageSize: number = 10
): Promise<VehicleListResponse> {
  try {
    const response = await fetch(
      `http://localhost:8081/api/v1/vehicles?page=${page}&pageSize=${pageSize}`
    );

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const data: VehicleListResponse = await response.json();
    return data;
  } catch (error) {
    console.error('Failed to fetch vehicle list:', error);
    throw error;
  }
}

/**
 * 车辆控制动作类型
 */
export type ControlAction = 'STOP' | 'START';

/**
 * 控制车辆
 * @param vin - 车辆 VIN 码
 * @param action - 控制动作（STOP: 紧急停车, START: 恢复）
 */
export async function controlVehicle(
  vin: string,
  action: ControlAction
): Promise<void> {
  try {
    const response = await fetch(
      `http://localhost:8081/api/v1/vehicles/${vin}/control`,
      {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ action }),
      }
    );

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
  } catch (error) {
    console.error(`Failed to control vehicle ${vin}:`, error);
    throw error;
  }
}

