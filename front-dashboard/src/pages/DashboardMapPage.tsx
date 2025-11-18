const DashboardMapPage = () => {
  return (
    <div className="h-full min-h-[calc(100vh-200px)] bg-gray-300 dark:bg-gray-700 rounded-lg shadow-md flex items-center justify-center border border-gray-200 dark:border-gray-600">
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
            d="M9 20l-5.447-2.724A1 1 0 013 16.382V5.618a1 1 0 011.447-.894L9 7m0 13l6-3m-6 3V7m6 10l4.553 2.276A1 1 0 0021 18.382V7.618a1 1 0 00-.553-.894L15 4m0 13V4m0 0L9 7" 
          />
        </svg>
        <h2 className="text-3xl font-bold text-gray-700 dark:text-gray-300 mb-2">
          全屏地图占位符
        </h2>
        <p className="text-xl text-gray-600 dark:text-gray-400">
          [Full Screen Map Placeholder]
        </p>
        <p className="text-lg text-gray-500 dark:text-gray-500 mt-4">
          实时车队位置将在此处显示
        </p>
      </div>
    </div>
  )
}

export default DashboardMapPage

