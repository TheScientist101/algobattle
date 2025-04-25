# AlgoBattle

A platform for hosting stock market trading bots


## Use our demo (RECOMMENDED FOR BEGINNERS)

Go to: https://algobattle.vercel.app

When prompted to login, choose either of the sign-up options, we recommend Signing In With Google. If you are not interested in making an account, sign in using `luna@email.com` for the username and `imbadatcs` for the password.

When you sign in, if there is no example bot created, click on the "Create a bot" button. Choose a name for the bot and click on the "Save" button. Reload the page to show information about your newly created bot. You may notice that there is no historical account value, and this is expected. It may take a couple minutes for the data to appear.

If you have made a bot that is compatible with our API, go ahead and use your API key to connect your bot to our servers. If not, use [this](https://urjithmishra.postman.co/workspace/Urjith-Mishra's-Workspace~43538fc0-c30d-40e3-8045-90c077511b1d/collection/44405624-4dc9bb5a-f0dd-4c42-91c9-679ae5220c1d?action=share&creator=44405624) Postman link that has examples.

You will notice that transactions appear on the dashboard, and your account will be revalued daily. We are working on having more frequently updated stock data so you will have more accurate account estimates.

##  Hosting (web)

These instructions will help you set up the project locally for development and testing.

### ✅ Prerequisites

- Node.js ≥ 18
- npm or yarn

### 📦 Installation

```bash
git clone https://github.com/TheScientist101/algobattle.git
cd algobattle
cd web
npm i
```
### ⚙ Set up enviornment variables

1. **Create a new firebase project, create a web app and add your enviornment variables to a .env file in the root directory of the web folder**

```bash
//.env contents
NEXT_PUBLIC_API_KEY = "YOUR_KEY"
NEXT_PUBLIC_AUTH_DOMAIN = "YOUR_AUTH_DOMAIN"
NEXT_PUBLIC_PROJECT_ID = "YOUR_PROJECT_ID"
NEXT_PUBLIC_STORAGE_BUCKET = "YOUR_STORAGE_BUCKET"
NEXT_PUBLIC_MESSAGING_SENDER_ID = "YOUR_MESSAGING_SENDER_ID"
NEXT_PUBLIC_APP_ID = "YOUR_APP_ID"
NEXT_PUBLIC_MEASUREMENT_ID = "YOUR_MEASUREMENT_ID"
```
##  Run the website

### 🚀 Run the website locally

```bash
npm run dev
```
##  Deploy the website

These instructions will help you deploy the AlgoBattle web app. Remember to properly define your enviornment variables.
- [Deploy on Vercel](https://vercel.com/docs/deployments/) - Click to view the documentation

## Hosting (server)

### ✅ Prerequisites

- Go ≥ 1.23

### 🚀 Run the server locally

```bash
git clone https://github.com/TheScientist101/algobattle.git
cd algobattle
cd server
go run urjith.dev/algobattle
```

## 👥 Authors & Contributions

| Name             |  Role / Contribution                                      |
|------------------|-----------------------------------------------------------|
| Abhinav Devarakonda|  I built the web application of AlgoBattle, including authentication pages and a dashboard containing information on all the user's trading bots, including their value over time, their transactions, their api keys, their account value, and the cash they have. Furthermore, I created a learboard page that shows which bots have made the most money.               |
| Urjith Mishra| I built the server-side code that any bot you write will interface with. This code also generates the value histories that are viewed on the web dashboard page, executes trades, calculates account values, provides stock pricing data, and validates API keys.|


---

## 📚 Citations & Attributions

### Code Libraries & Tools
- [Next.js](https://nextjs.org/) – React Framework for SSR and SSG
- [React](https://react.dev/) – JavaScript library for building UI
- [Tailwind CSS](https://tailwindcss.com/) – Utility-first CSS framework
- [Firebase](https://firebase.google.com/) - A suite of tools for web development includuding a cloud database and authentication manager
- [Shadcn](https://ui.shadcn.com/) – Nextjs Component Library
- [v0.dev](https://v0.dev/) - AI Nextjs component generator/editor
- [ChatGPT](https://openai.com/index/chatgpt/) - AI chatbot used to create this Readme and to diagnose logical missteps and help with stying
- [uuid](https://www.npmjs.com/package/uuid) - A library for creating unique ids
- [gin gonic](https://github.com/gin-gonic/gin) - HTTP web framework for Go


### Assets
- All assets from starter code from Next.js

### References
- [Next.js Docs](https://nextjs.org/docs) – Official documentation
- [Tailwind CSS Docs](https://v2.tailwindcss.com/docs) -Official documentation
- [Firebase Docs](https://firebase.google.com/docs) - Official documentation

##  View the current deployment

- [AlgoBattle.app](https://algobattle.vercel.app/) - Click to view the current AlgoBattle web app