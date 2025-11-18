# CheYa Truck Monitoring Platform - Frontend

A modern, enterprise-grade SaaS platform for monitoring autonomous trucks, built with React, TypeScript, and Tailwind CSS.

## Features

- **Authentication**: Login and Registration pages
- **Dashboard**: Real-time truck fleet monitoring with:
  - Active fleet list with status indicators
  - Vehicle instrument panel with live data
  - Action controls for fleet management

## Tech Stack

- **React 18** - UI framework
- **TypeScript** - Type safety
- **Vite** - Build tool and dev server
- **React Router DOM** - Client-side routing
- **Tailwind CSS** - Utility-first CSS framework

## Getting Started

### Installation

```bash
npm install
```

### Development

```bash
npm run dev
```

Visit `http://localhost:5173` to view the application.

### Build

```bash
npm run build
```

### Preview Production Build

```bash
npm run preview
```

## Project Structure

```
front-dashboard/
├── src/
│   ├── pages/
│   │   ├── LoginPage.tsx       # User login page
│   │   ├── RegisterPage.tsx    # User registration page
│   │   └── DashboardPage.tsx   # Main dashboard with fleet monitoring
│   ├── App.tsx                 # Main app with routing
│   ├── main.tsx               # Entry point with BrowserRouter
│   └── index.css              # Global styles with Tailwind
├── index.html
├── package.json
├── tsconfig.json
├── vite.config.ts
└── tailwind.config.js
```

## Routes

- `/` - Redirects to `/login`
- `/login` - Login page
- `/register` - Registration page
- `/dashboard` - Main dashboard (fleet monitoring)

## Status Indicators

- 🟢 **Green** - Online
- ⚫ **Gray** - Offline
- 🟡 **Yellow** - Maintaining

## License

MIT



















