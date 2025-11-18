const DashboardListPage = () => {
  return (
    <div className="h-full min-h-[calc(100vh-200px)] bg-white dark:bg-gray-800 rounded-lg shadow-md flex items-center justify-center p-8 border border-gray-200 dark:border-gray-700">
      <div className="text-center">
        <svg 
          className="w-24 h-24 mx-auto mb-6 text-gray-600 dark:text-gray-400" 
          fill="none" 
          stroke="currentColor" 
          viewBox="0 0 24 24"
        >
          <path 
            strokeLinecap="round" 
            strokeLinejoin="round" 
            strokeWidth={2} 
            d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01" 
          />
        </svg>
        <h2 className="text-3xl font-bold text-gray-800 dark:text-gray-200 mb-2">
          专用车辆列表
        </h2>
        <p className="text-xl text-gray-600 dark:text-gray-400">
          [Dedicated Truck List Placeholder]
        </p>
        <p className="text-lg text-gray-500 dark:text-gray-500 mt-4">
          详细的车辆列表将在此处显示
        </p>
      </div>
    </div>
  )
}

export default DashboardListPage

