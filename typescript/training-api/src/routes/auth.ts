import { Hono } from "hono";
import { z } from "zod";
import { zValidator } from "@hono/zod-validator";
async function hashPassword(password: string): Promise<string> {
  const msgUint8 = new TextEncoder().encode(password + "salt");
  const hashBuffer = await crypto.subtle.digest("SHA-256", msgUint8);
  const hashArray = Array.from(new Uint8Array(hashBuffer));
  const hashHex = hashArray.map(b => b.toString(16).padStart(2, '0')).join('');
  return hashHex;
}

async function comparePassword(password: string, hashedPassword: string): Promise<boolean> {
  const hashed = await hashPassword(password);
  return hashed === hashedPassword;
}
import { sign, verify } from "hono/jwt";

type Bindings = {
  DB: D1Database;
  JWT_SECRET: string;
};

const auth = new Hono<{ Bindings: Bindings }>();

const registerSchema = z.object({
  email: z.string().email(),
  password: z.string().min(8),
  name: z.string().min(1),
});

const loginSchema = z.object({
  email: z.string().email(),
  password: z.string(),
});

auth.post(
  "/register",
  zValidator("json", registerSchema),
  async (c) => {
    try {
      const { email, password, name } = c.req.valid("json");
      
      const existingUser = await c.env.DB.prepare(
        "SELECT * FROM users WHERE email = ?"
      ).bind(email).first();
      
      if (existingUser) {
        return c.json({ error: "User with this email already exists" }, 409);
      }
      
      const hashedPassword = await hashPassword(password);
      
      const id = crypto.randomUUID();
      const timestamp = new Date().toISOString();
      
      await c.env.DB.prepare(
        "INSERT INTO users (id, email, password, name, created_at) VALUES (?, ?, ?, ?, ?)"
      ).bind(id, email, hashedPassword, name, timestamp).run();
      
      const token = await sign({ id, email }, c.env.JWT_SECRET || 'dev-jwt-secret');
      
      return c.json({ 
        message: "User registered successfully",
        token,
        user: { id, email, name }
      }, 201);
    } catch (error) {
      console.error("Error registering user:", error);
      return c.json({ error: "Failed to register user" }, 500);
    }
  }
);

auth.post(
  "/login",
  zValidator("json", loginSchema),
  async (c) => {
    try {
      const { email, password } = c.req.valid("json");
      
      const user = await c.env.DB.prepare(
        "SELECT * FROM users WHERE email = ?"
      ).bind(email).first();
      
      if (!user) {
        return c.json({ error: "Invalid credentials" }, 401);
      }
      
      const isMatch = await comparePassword(password, user.password as string);
      
      if (!isMatch) {
        return c.json({ error: "Invalid credentials" }, 401);
      }
      
      const token = await sign({ id: user.id, email: user.email }, c.env.JWT_SECRET || 'dev-jwt-secret');
      
      return c.json({ 
        message: "Login successful",
        token,
        user: { id: user.id, email: user.email, name: user.name }
      });
    } catch (error) {
      console.error("Error logging in:", error);
      return c.json({ error: "Failed to login" }, 500);
    }
  }
);

auth.get("/me", async (c) => {
  try {
    const authHeader = c.req.header("Authorization");
    
    if (!authHeader || !authHeader.startsWith("Bearer ")) {
      return c.json({ error: "Unauthorized" }, 401);
    }
    
    const token = authHeader.split(" ")[1];
    
    try {
      const decoded = await verify(token, c.env.JWT_SECRET || 'dev-jwt-secret');
      
      const user = await c.env.DB.prepare(
        "SELECT id, email, name, created_at FROM users WHERE id = ?"
      ).bind(decoded.id).first();
      
      if (!user) {
        return c.json({ error: "User not found" }, 404);
      }
      
      return c.json({ user });
    } catch (error) {
      return c.json({ error: "Invalid token" }, 401);
    }
  } catch (error) {
    console.error("Error getting user:", error);
    return c.json({ error: "Failed to get user" }, 500);
  }
});

export { auth as authRoutes };
