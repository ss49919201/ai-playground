
declare global {
  interface Env {
    DB: D1Database;
    JWT_SECRET: string;
  }
}

export {};
