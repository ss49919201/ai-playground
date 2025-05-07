import { describe, it, expect, vi, beforeEach } from 'vitest';
import { Hono } from 'hono';
import { authRoutes } from './auth';

const mockDB = {
  prepare: vi.fn().mockReturnThis(),
  bind: vi.fn().mockReturnThis(),
  first: vi.fn(),
  all: vi.fn(),
  run: vi.fn(),
};

const mockEnv = {
  DB: mockDB,
  JWT_SECRET: 'test-secret',
};

describe('Auth Routes', () => {
  let app: Hono;
  
  beforeEach(() => {
    vi.clearAllMocks();
    app = new Hono();
    app.route('/auth', authRoutes);
  });
  
  describe('POST /auth/register', () => {
    it('should register a new user successfully', async () => {
      mockDB.first.mockResolvedValueOnce(null);
      mockDB.run.mockResolvedValueOnce({});
      
      const res = await app.request('/auth/register', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          email: 'test@example.com',
          password: 'password123',
          name: 'Test User',
        }),
      }, { env: mockEnv });
      
      const data = await res.json();
      
      expect(res.status).toBe(201);
      expect(data.message).toBe('User registered successfully');
      expect(data.token).toBeDefined();
      expect(data.user).toHaveProperty('id');
      expect(data.user.email).toBe('test@example.com');
      expect(data.user.name).toBe('Test User');
      
      expect(mockDB.prepare).toHaveBeenCalledWith(
        expect.stringContaining('SELECT * FROM users WHERE email = ?')
      );
      expect(mockDB.prepare).toHaveBeenCalledWith(
        expect.stringContaining('INSERT INTO users')
      );
    });
    
    it('should return 409 if user already exists', async () => {
      mockDB.first.mockResolvedValueOnce({ id: '123', email: 'test@example.com' });
      
      const res = await app.request('/auth/register', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          email: 'test@example.com',
          password: 'password123',
          name: 'Test User',
        }),
      }, { env: mockEnv });
      
      const data = await res.json();
      
      expect(res.status).toBe(409);
      expect(data.error).toBe('User with this email already exists');
    });
    
    it('should return 400 for invalid input', async () => {
      const res = await app.request('/auth/register', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          email: 'invalid-email',
          password: 'short',
          name: '',
        }),
      }, { env: mockEnv });
      
      expect(res.status).toBe(400);
    });
  });
  
  describe('POST /auth/login', () => {
    it('should login successfully with valid credentials', async () => {
      const hashedPassword = '8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92'; // SHA-256 of 'password123salt'
      mockDB.first.mockResolvedValueOnce({ 
        id: '123', 
        email: 'test@example.com',
        password: hashedPassword,
        name: 'Test User'
      });
      
      const res = await app.request('/auth/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          email: 'test@example.com',
          password: 'password123',
        }),
      }, { env: mockEnv });
      
      const data = await res.json();
      
      expect(res.status).toBe(200);
      expect(data.message).toBe('Login successful');
      expect(data.token).toBeDefined();
      expect(data.user).toHaveProperty('id');
      expect(data.user.email).toBe('test@example.com');
    });
    
    it('should return 401 for invalid credentials', async () => {
      const hashedPassword = 'wrong-password-hash';
      mockDB.first.mockResolvedValueOnce({ 
        id: '123', 
        email: 'test@example.com',
        password: hashedPassword,
        name: 'Test User'
      });
      
      const res = await app.request('/auth/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          email: 'test@example.com',
          password: 'wrong-password',
        }),
      }, { env: mockEnv });
      
      const data = await res.json();
      
      expect(res.status).toBe(401);
      expect(data.error).toBe('Invalid credentials');
    });
    
    it('should return 401 if user does not exist', async () => {
      mockDB.first.mockResolvedValueOnce(null);
      
      const res = await app.request('/auth/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          email: 'nonexistent@example.com',
          password: 'password123',
        }),
      }, { env: mockEnv });
      
      const data = await res.json();
      
      expect(res.status).toBe(401);
      expect(data.error).toBe('Invalid credentials');
    });
  });
  
  describe('GET /auth/me', () => {
    it('should return user information with valid token', async () => {
      mockDB.first.mockResolvedValueOnce({ 
        id: '123', 
        email: 'test@example.com',
        name: 'Test User',
        created_at: '2023-01-01T00:00:00.000Z'
      });
      
      const token = await import('hono/jwt').then(({ sign }) => 
        sign({ id: '123', email: 'test@example.com' }, 'test-secret')
      );
      
      const res = await app.request('/auth/me', {
        method: 'GET',
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      }, { env: mockEnv });
      
      const data = await res.json();
      
      expect(res.status).toBe(200);
      expect(data.user).toHaveProperty('id', '123');
      expect(data.user).toHaveProperty('email', 'test@example.com');
      expect(data.user).toHaveProperty('name', 'Test User');
    });
    
    it('should return 401 with invalid token', async () => {
      const res = await app.request('/auth/me', {
        method: 'GET',
        headers: {
          'Authorization': 'Bearer invalid-token',
        },
      }, { env: mockEnv });
      
      const data = await res.json();
      
      expect(res.status).toBe(401);
      expect(data.error).toBe('Invalid token');
    });
    
    it('should return 401 without authorization header', async () => {
      const res = await app.request('/auth/me', {
        method: 'GET',
      }, { env: mockEnv });
      
      const data = await res.json();
      
      expect(res.status).toBe(401);
      expect(data.error).toBe('Unauthorized');
    });
  });
});
