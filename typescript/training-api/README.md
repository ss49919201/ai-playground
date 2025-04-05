# Training Records API

A RESTful API for managing training records, built with TypeScript, Hono, and Cloudflare Workers.

## Features

- CRUD operations for training records
- User authentication with email and password
- API versioning
- API documentation
- Persistent storage with Cloudflare D1

## Tech Stack

- **Language**: TypeScript
- **Framework**: Hono
- **Infrastructure**: Cloudflare Workers
- **Database**: Cloudflare D1
- **Authentication**: JWT-based authentication

## API Endpoints

### Authentication

- `POST /auth/register` - Register a new user
- `POST /auth/login` - Login with email and password
- `GET /auth/me` - Get current user information

### Training Records (v1)

- `GET /api/v1/records` - Get all training records
- `GET /api/v1/records/:id` - Get a specific training record
- `POST /api/v1/records` - Create a new training record
- `PUT /api/v1/records/:id` - Update a training record
- `DELETE /api/v1/records/:id` - Delete a training record
- `POST /api/v1/records/:id/exercises` - Add an exercise to a training record
- `POST /api/v1/exercises/:id/sets` - Add a set to an exercise

## Database Schema

```sql
-- Users table for authentication
CREATE TABLE IF NOT EXISTS users (
  id TEXT PRIMARY KEY,
  email TEXT UNIQUE NOT NULL,
  password TEXT NOT NULL,
  name TEXT NOT NULL,
  created_at TEXT NOT NULL
);

-- Training records table
CREATE TABLE IF NOT EXISTS training_records (
  id TEXT PRIMARY KEY,
  title TEXT NOT NULL,
  date TEXT NOT NULL,
  description TEXT,
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL,
  user_id TEXT,
  FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Exercises table
CREATE TABLE IF NOT EXISTS exercises (
  id TEXT PRIMARY KEY,
  record_id TEXT NOT NULL,
  name TEXT NOT NULL,
  FOREIGN KEY (record_id) REFERENCES training_records(id) ON DELETE CASCADE
);

-- Sets table
CREATE TABLE IF NOT EXISTS sets (
  id TEXT PRIMARY KEY,
  exercise_id TEXT NOT NULL,
  weight REAL NOT NULL,
  reps INTEGER NOT NULL,
  notes TEXT,
  FOREIGN KEY (exercise_id) REFERENCES exercises(id) ON DELETE CASCADE
);
```

## Development

### Prerequisites

- Node.js and npm
- Wrangler CLI (for Cloudflare Workers)

### Setup

1. Clone the repository
2. Install dependencies:

```bash
npm install
```

3. Create a D1 database:

```bash
wrangler d1 create training_records
```

4. Update the `wrangler.toml` file with your database ID:

```toml
[[d1_databases]]
binding = "DB"
database_name = "training_records"
database_id = "your-database-id"
```

5. Apply the database schema:

```bash
wrangler d1 execute training_records --file=./src/db/schema.sql
```

6. Start the development server:

```bash
npm run dev
```

### Deployment

Deploy to Cloudflare Workers:

```bash
npm run deploy
```

## Environment Variables

- `JWT_SECRET` - Secret key for JWT token generation and verification

## License

MIT
