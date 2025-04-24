# AlgoBattle

A platform for hosting stock market trading bots

---

## 👥 Authors & Contributions

| Name             |  Role / Contribution                                      |
|------------------|-----------------------------------------------------------|
| Abhinav Devarakonda|  I built the web application of AlgoBattle, including authentication & authentication pages and a dashboard containing information on all the user's trading bots, including their value over time, their transactions, their api keys, and they cash they have.               |


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


### Assets
- All assets from starter code from Next.js

### References
- [Next.js Docs](https://nextjs.org/docs) – Official documentation
- [Tailwind CSS Docs](https://v2.tailwindcss.com/docs) -Official documentation
- [Firebase Docs](https://firebase.google.com/docs) - Official documentation

---

##  Run locally

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


##  View the current deployment

- [AlgoBattle.app](https://algobattle.vercel.app/) - Click to view the current AlgoBattle web app