# AlgoBattle

[![Next.js](https://img.shields.io/badge/Next.js-13.0+-000000?style=flat-square&logo=next.js&logoColor=white)](https://nextjs.org/)
[![React](https://img.shields.io/badge/React-18.0+-61DAFB?style=flat-square&logo=react&logoColor=black)](https://react.dev/)
[![Go](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat-square&logo=go&logoColor=white)](https://golang.org/)
[![Firebase](https://img.shields.io/badge/Firebase-9.0+-FFCA28?style=flat-square&logo=firebase&logoColor=black)](https://firebase.google.com/)
[![Tailwind CSS](https://img.shields.io/badge/Tailwind_CSS-3.0+-38B2AC?style=flat-square&logo=tailwind-css&logoColor=white)](https://tailwindcss.com/)

> A professional platform for developing, testing, and competing with algorithmic stock market trading bots

## Table of Contents
- [Overview](#overview)
- [Features](#features)
- [Getting Started](#getting-started)
  - [Option 1: Use Our Demo (Recommended for Beginners)](#option-1-use-our-demo-recommended-for-beginners)
  - [Option 2: Run Locally](#option-2-run-locally)
    - [Web Application Setup](#web-application-setup)
    - [Server Setup](#server-setup)
- [API Documentation](#api-documentation)
- [Authors & Contributions](#authors--contributions)
- [Citations & Attributions](#citations--attributions)
- [Live Deployment](https://github.com/TheScientist101/algobattle/edit/main/README.md#live-deployment)

## Overview

AlgoBattle is a sophisticated platform designed for financial technology enthusiasts, algorithmic traders, and developers interested in stock market automation. The platform enables users to create, deploy, and monitor algorithmic trading bots in a controlled environment with real market data.

The system consists of two main components:
- **Web Dashboard**: A modern, responsive interface for managing bots, viewing performance metrics, and analyzing trading history
- **API Backend**: A robust server-side infrastructure that handles bot interactions, executes trades, and provides market data

## Features

- **Algorithmic Trading**: Develop and deploy custom trading algorithms without risking real capital
- **Real-Time Market Data**: Access to current and historical stock prices from major exchanges
- **Performance Analytics**: Track your bot's performance with comprehensive metrics and visualizations
- **Leaderboard Competition**: Compare your bot's performance against others in the community
- **Secure Authentication**: Firebase-powered user authentication and API key management
- **RESTful API**: Well-documented API for programmatic interaction with the platform

## Getting Started

Choose one of the following options to start using AlgoBattle:

### Option 1: Use Our Demo (Recommended for Beginners)

This option allows you to quickly explore AlgoBattle without setting up the development environment.

#### Step 1: Access the Platform
Visit our live deployment at [https://algobattle.vercel.app](https://algobattle.vercel.app)

#### Step 2: Authentication
- Sign in using your preferred method (Google authentication recommended)
- For a quick demo without registration, use these credentials:
  - Username: `luna@email.com`
  - Password: `imbadatcs`

#### Step 3: Create Your First Bot
- After signing in, click the "Create a bot" button if no bots exist
- Provide a name for your bot and click "Save"
- Refresh the page to view your bot's dashboard
- Note: Historical account value data may take a few minutes to appear

#### Step 4: Connect Your Bot to the API
- **For developers with custom bots**: Use the API key displayed on your dashboard
- **For testing purposes**: Use our [Postman collection](https://www.postman.com/urjithmishra/algobattle-examples/collection/fxkywuo/algobattle?action=share&creator=44405624&active-environment=44405624-e678c307-08d5-4bca-be67-3eb53ee734ee)

#### Setting up the API Key in Postman:
1. Navigate to the "Authorization" tab
2. Hover over "ALGOBATTLE_API_KEY" and click the appearing textbox
3. Enter your API key
4. Select desired API requests from the left panel and click "Send"
5. View responses in the bottom panel

#### Step 5: Monitor Performance
Your dashboard will display transactions and account value updates, which are recalculated daily.

### Option 2: Run Locally

This option is for developers who want to set up the complete AlgoBattle environment on their local machine.

#### Web Application Setup

##### Prerequisites
- Node.js ≥ 18
- npm or yarn package manager

##### Step 1: Clone the Repository
```bash
git clone https://github.com/TheScientist101/algobattle.git
cd algobattle/web
```

##### Step 2: Install Dependencies
```bash
npm install
# or if using yarn
# yarn install
```

##### Step 3: Configure Environment
1. Create a new Firebase project and register a web app
2. Create a `.env` file in the web directory with the following variables:

```bash
# Firebase Configuration
NEXT_PUBLIC_API_KEY="YOUR_FIREBASE_API_KEY"
NEXT_PUBLIC_AUTH_DOMAIN="YOUR_AUTH_DOMAIN"
NEXT_PUBLIC_PROJECT_ID="YOUR_PROJECT_ID"
NEXT_PUBLIC_STORAGE_BUCKET="YOUR_STORAGE_BUCKET"
NEXT_PUBLIC_MESSAGING_SENDER_ID="YOUR_MESSAGING_SENDER_ID"
NEXT_PUBLIC_APP_ID="YOUR_APP_ID"
NEXT_PUBLIC_MEASUREMENT_ID="YOUR_MEASUREMENT_ID"
```

##### Step 4: Start Development Server
```bash
npm run dev
# or if using yarn
# yarn dev
```

The web application will be available at `http://localhost:3000`

##### Deployment (Optional)
For production deployment, we recommend using Vercel:
- Follow the [Vercel deployment guide](https://vercel.com/docs/deployments/)
- Ensure all environment variables are properly configured in your Vercel project settings

#### Server Setup

##### Prerequisites
- Go ≥ 1.23
- Git

##### Step 1: Clone the Repository (if not already done)
```bash
git clone https://github.com/TheScientist101/algobattle.git
cd algobattle/server
```

##### Step 2: Run the Server
```bash
go run urjith.dev/algobattle
```

The API server will be available at `http://localhost:8080`

## API Documentation

AlgoBattle provides a comprehensive RESTful API that allows developers to programmatically interact with the platform. The API enables your trading bots to:

- Retrieve real-time and historical market data
- Execute buy and sell transactions
- Monitor portfolio performance
- Manage watchlists and track specific stocks

For complete technical details, endpoint specifications, and code examples, please refer to our [API Documentation](server/api_documentation.md).

### API Features

| Feature | Description |
|---------|-------------|
| Authentication | Secure API key-based authentication system |
| Portfolio Management | Retrieve account balances, holdings, and transaction history |
| Market Data | Access to current and historical stock prices and indicators |
| Trading | Execute buy and sell orders programmatically |
| Error Handling | Comprehensive error reporting and status codes |

## Authors & Contributions

AlgoBattle was developed by a team of dedicated engineers with expertise in financial technology and web development.

| Contributor         | Role                | Contributions                                             |
|---------------------|---------------------|----------------------------------------------------------|
| Abhinav Devarakonda | Frontend Developer  | • Designed and implemented the web application interface<br>• Built authentication system with Firebase integration<br>• Developed the dashboard for bot monitoring and analytics<br>• Created interactive data visualizations for performance tracking<br>• Implemented the leaderboard system for competitive rankings |
| Urjith Mishra       | Backend Developer   | • Architected and built the server-side API infrastructure<br>• Developed the trading engine and transaction processing system<br>• Implemented market data integration and caching mechanisms<br>• Created portfolio valuation and performance calculation algorithms<br>• Designed the API authentication and security protocols |

## Citations & Attributions

### Frontend Technologies
| Technology | Description | Version |
|------------|-------------|---------|
| [Next.js](https://nextjs.org/) | React framework for server-side rendering | 13.0+ |
| [React](https://react.dev/) | JavaScript library for building user interfaces | 18.0+ |
| [Tailwind CSS](https://tailwindcss.com/) | Utility-first CSS framework | 3.0+ |
| [Firebase](https://firebase.google.com/) | Authentication and database services | 9.0+ |
| [Shadcn UI](https://ui.shadcn.com/) | Component library for Next.js | Latest |
| [UUID](https://www.npmjs.com/package/uuid) | Library for generating unique identifiers | 9.0+ |

### Backend Technologies
| Technology | Description | Version |
|------------|-------------|---------|
| [Go](https://golang.org/) | Programming language for backend services | 1.23+ |
| [Gin Gonic](https://github.com/gin-gonic/gin) | HTTP web framework for Go | Latest |

### Development Tools
| Tool | Purpose |
|------|---------|
| [v0.dev](https://v0.dev/) | AI-powered component generator for Next.js | 
| [ChatGPT](https://openai.com/index/chatgpt/) | AI assistant for documentation and debugging |

### Documentation & References
- [Next.js Documentation](https://nextjs.org/docs) – Official framework guides and API reference
- [Tailwind CSS Documentation](https://v2.tailwindcss.com/docs) – Component styling and utility classes
- [Firebase Documentation](https://firebase.google.com/docs) – Authentication and database integration
- [Go Documentation](https://golang.org/doc/) – Language reference and standard library

## Live Deployment

<div align="center">
  <a href="https://algobattle.vercel.app/" target="_blank">
    <img src="https://img.shields.io/badge/Visit_AlgoBattle-000000?style=for-the-badge&logo=vercel&logoColor=white" alt="Visit AlgoBattle" />
  </a>
  <p>Experience the platform at <a href="https://algobattle.vercel.app/">algobattle.vercel.app</a></p>
</div>
