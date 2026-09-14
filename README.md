# Bootcamp

## Individual Project - FreelanceFi

[![Go](https://img.shields.io/badge/go-1.20+-00ADD8?logo=go\&logoColor=white)](https://golang.org)
[![PostgreSQL](https://img.shields.io/badge/postgresql-14+-316192?logo=postgresql\&logoColor=white)](https://www.postgresql.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

A modern freelancing marketplace built with Go, PostgreSQL, and Templ. **FreelanceFi** connects clients and freelancers through secure authentication, project posting, bidding, and personalized dashboards.

---

## 🌟 Key Features

* **Secure Authentication** with JWT (HTTP‑only cookies)
* **Role‑based Access** for Clients and Freelancers
* **Job Marketplace**: Clients post work, Freelancers bid
* **Personal Dashboards** for active jobs, bids, and earnings
* **Financial Tracking** to monitor payments & history
* **Server‑Side Rendering** via Templ + HTMX + Bootstrap + JS + CSS
* **Type‑Safe Data Layer** using sqlc over PostgreSQL

---

## 🏗️ Architecture & Stack

| Layer       | Technology                             |
| ----------- | -------------------------------------- |
| Language    | Go 1.20+                               |
| Routing     | Chi                                    |
| Database    | PostgreSQL 14+, sqlc‑generated queries |
| Templates   | Templ                                  |
| Frontend UX | Bootstrap 5, HTMX                      |
| Security    | bcrypt (passwords), JWT (sessions)     |

---

## 🚀 Getting Started

### Prerequisites

* Go 1.20+
* PostgreSQL 14+
* `sqlc` and `templ` CLIs

### Clone & Configure

```bash
git clone https://github.com/your-username/freelancefi.git
cd freelancefi

cp config/.env.example .env   # update DB cred & JWT secret
```

### Database Setup

```bash
# Apply migrations
psql -U <user> -d <db> -f migrations/001_init.sql
psql -U <user> -d <db> -f migrations/002_add_gig_schema.sql

# Generate Go types
sqlc generate
```

### Run

```bash
go run main.go
```

Navigate to **[http://localhost:8081](http://localhost:8081)**

---

## 📖 Usage Flow

1. **Register** as a Client or Freelancer
2. **Log in** to receive your session cookie
3. Clients: **Create jobs** and review incoming bids
4. Freelancers: **Browse jobs** and **submit bids**
5. Track progress from **Work History** and **Finance** pages

---

## 📄 Core API Endpoints

| Method | Endpoint                           | Purpose                        |
| ------ | ---------------------------------- | ------------------------------ |
| POST   | `/register`                        | Register new user              |
| POST   | `/login`                           | Authenticate & receive session |
| GET    | `/home`                            | Unified dashboard              |
| GET    | `/jobs_page`                       | Browse available jobs          |
| GET    | `/profile`                         | View personal profile          |
| GET    | `/finance`                         | Financial summary              |
| GET    | `/client/dashboard`                | Client‑specific dashboard      |
| GET    | `/freelancer/freelancer_dashboard` | Freelancer‑specific dashboard  |
| GET    | `/users`                           | List users (admin)             |

---

## Next -> Group Project 🤝 [Availability-Management-System](https://github.com/aletaslauan/Availability-Management-System)

