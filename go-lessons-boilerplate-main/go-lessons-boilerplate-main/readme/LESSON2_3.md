# Lesson 2.3: The Project

## Project Overview

You will build a **complete web application** (backend + frontend + database). The application should support user authentication, personalized and public data, pagination, etc.

This is a realistic project that brings together core concepts of web development, RESTful APIs, databases, authorization, business logic, and frontend integration.

---

## Minimum Requirements

### 1. Authentication & Authorization  
- Users must be able to register, log-in, and manage sessions.  
- Users can access their own private data (e.g. posted products, bids, purchases).  

### 2. Persistent Storage  

### 3. Logging  
- Log important actions and errors.  
- Suggestion: Use structured logs (ex. `log/slog` or `logrus`).  

### 4. Pagination  
- When listing large sets of data (e.g. product feed), implement pagination.  
- Example: `GET /products?page=3` returns page 3 of products  
- Show page metadata: `Page 3 / 49`  

### 5. Filtering
- Allow filtering by fields like price, category, etc. 
- Example: `GET /products?category=phones&max_price=500&page=1`  

### 6. Shared / Global Data  
- Some data should be visible to all users (e.g. the public product feed).  

### 7. At Least 3 Domain Tables (excluding users and join tables)  
- Tables should be meaningful and distinct. Examples include:  
  - `products`
  - `bids` or `orders`
  - `reviews`, `categories`, `inventory`, etc.

---

## Example Projects

### Project Idea #1: **Olxa** – Auction-Based Marketplace

**Concept:**  
A marketplace where users can post items for sale. Buyers place bids, and sellers choose the best offer. Items are sold once per posting (like OLX / Craigslist with bidding).

**Key Features:**
- User registration & login
- Product feed & search
- Post item for sale
- Place bids on items
- Seller sees all bids and chooses one to accept
- Sold items are no longer available
- Users view their posted items, received bids, sold items
- (Optional) Reviews for users

---

### Project Idea #2: **Emagino** – Inventory-Based Store

**Concept:**  
A traditional e-commerce site where sellers post items with stock quantities. Buyers add items to cart and complete purchases. Sellers must manage stock levels.

**Key Features:**
- User registration & login
- Product feed & filtering
- Add product to inventory (with quantity)
- Buyers can purchase multiple items
- Stock decrements on purchase
- Prevent overselling
- Seller dashboard: inventory, sold items
- Buyer dashboard: purchases, (Optional) product reviews
- (Optional) Product reviews (per product)

## Final Deliverables

To complete the project successfully, you must submit:

- A **fully working application** (frontend + backend)
- A **GitHub repository** with full source code
- A **SQL schema** or database migration files
- A **database design diagram**  
  - You must include an **Entity-Relationship Diagram (ERD)** showing:
    - Tables
    - Relationships (one-to-many, many-to-many, etc.)
    - Keys (primary/foreign)
    - Types of fields (briefly)

📌 **Note:** The database design is important for showing that you understand the structure and connections in your app. You can draw the diagram using tools like dbdiagram.io, drawSQL, or even hand-drawn scanned images.
