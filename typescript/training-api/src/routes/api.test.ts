import { describe, it, expect, vi, beforeEach } from 'vitest';
import { Hono } from 'hono';
import { apiRoutes } from './api';

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

describe('API Routes', () => {
  let app: Hono;
  
  beforeEach(() => {
    vi.clearAllMocks();
    app = new Hono();
    app.route('/api/v1', apiRoutes);
  });
  
  describe('GET /api/v1/records', () => {
    it('should return all training records', async () => {
      const mockRecords = [
        { id: '1', title: 'Workout 1', date: '2023-01-01' },
        { id: '2', title: 'Workout 2', date: '2023-01-02' },
      ];
      
      mockDB.all.mockResolvedValueOnce({ results: mockRecords });
      
      const res = await app.request('/api/v1/records', {
        method: 'GET',
      }, { env: mockEnv });
      
      const data = await res.json();
      
      expect(res.status).toBe(200);
      expect(data.records).toEqual(mockRecords);
      expect(mockDB.prepare).toHaveBeenCalledWith(
        expect.stringContaining('SELECT * FROM training_records')
      );
    });
    
    it('should handle database errors', async () => {
      mockDB.all.mockRejectedValueOnce(new Error('Database error'));
      
      const res = await app.request('/api/v1/records', {
        method: 'GET',
      }, { env: mockEnv });
      
      const data = await res.json();
      
      expect(res.status).toBe(500);
      expect(data.error).toBe('Failed to fetch training records');
    });
  });
  
  describe('GET /api/v1/records/:id', () => {
    it('should return a specific training record with exercises and sets', async () => {
      const mockRecord = { id: '1', title: 'Workout 1', date: '2023-01-01' };
      
      const mockExercises = [
        { id: 'e1', name: 'Bench Press', record_id: '1' },
        { id: 'e2', name: 'Squat', record_id: '1' },
      ];
      
      const mockSets1 = [
        { id: 's1', weight: 100, reps: 10, exercise_id: 'e1' },
        { id: 's2', weight: 110, reps: 8, exercise_id: 'e1' },
      ];
      
      const mockSets2 = [
        { id: 's3', weight: 150, reps: 8, exercise_id: 'e2' },
        { id: 's4', weight: 160, reps: 6, exercise_id: 'e2' },
      ];
      
      mockDB.first.mockResolvedValueOnce(mockRecord);
      mockDB.all.mockResolvedValueOnce({ results: mockExercises });
      mockDB.all.mockResolvedValueOnce({ results: mockSets1 });
      mockDB.all.mockResolvedValueOnce({ results: mockSets2 });
      
      const res = await app.request('/api/v1/records/1', {
        method: 'GET',
      }, { env: mockEnv });
      
      const data = await res.json();
      
      expect(res.status).toBe(200);
      expect(data.record).toEqual({
        ...mockRecord,
        exercises: [
          { ...mockExercises[0], sets: mockSets1 },
          { ...mockExercises[1], sets: mockSets2 },
        ],
      });
    });
    
    it('should return 404 if record not found', async () => {
      mockDB.first.mockResolvedValueOnce(null);
      
      const res = await app.request('/api/v1/records/999', {
        method: 'GET',
      }, { env: mockEnv });
      
      const data = await res.json();
      
      expect(res.status).toBe(404);
      expect(data.error).toBe('Training record not found');
    });
  });
  
  describe('POST /api/v1/records', () => {
    it('should create a new training record', async () => {
      const originalRandomUUID = crypto.randomUUID;
      crypto.randomUUID = vi.fn()
        .mockReturnValueOnce('record-id')
        .mockReturnValueOnce('exercise-id')
        .mockReturnValueOnce('set-id');
      
      mockDB.run.mockResolvedValue({});
      
      const res = await app.request('/api/v1/records', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          title: 'New Workout',
          date: '2023-01-03',
          description: 'Leg day',
          exercises: [
            {
              name: 'Squat',
              sets: [
                { weight: 150, reps: 8, notes: 'Felt good' },
              ],
            },
          ],
        }),
      }, { env: mockEnv });
      
      const data = await res.json();
      
      expect(res.status).toBe(201);
      expect(data.id).toBe('record-id');
      expect(data.message).toBe('Training record created successfully');
      
      crypto.randomUUID = originalRandomUUID;
    });
    
    it('should return 400 for invalid input', async () => {
      const res = await app.request('/api/v1/records', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          description: 'Invalid workout',
        }),
      }, { env: mockEnv });
      
      const data = await res.json();
      
      expect(res.status).toBe(400);
      expect(data.error).toBe('Title and date are required');
    });
  });
  
  describe('PUT /api/v1/records/:id', () => {
    it('should update a training record', async () => {
      mockDB.first.mockResolvedValueOnce({ id: '1', title: 'Old Title' });
      mockDB.run.mockResolvedValueOnce({});
      
      const res = await app.request('/api/v1/records/1', {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          title: 'Updated Workout',
          date: '2023-01-03',
          description: 'Updated description',
        }),
      }, { env: mockEnv });
      
      const data = await res.json();
      
      expect(res.status).toBe(200);
      expect(data.message).toBe('Training record updated successfully');
      expect(mockDB.prepare).toHaveBeenCalledWith(
        expect.stringContaining('UPDATE training_records SET')
      );
    });
    
    it('should return 404 if record not found', async () => {
      mockDB.first.mockResolvedValueOnce(null);
      
      const res = await app.request('/api/v1/records/999', {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          title: 'Updated Workout',
          date: '2023-01-03',
        }),
      }, { env: mockEnv });
      
      const data = await res.json();
      
      expect(res.status).toBe(404);
      expect(data.error).toBe('Training record not found');
    });
  });
  
  describe('DELETE /api/v1/records/:id', () => {
    it('should delete a training record and its related data', async () => {
      mockDB.first.mockResolvedValueOnce({ id: '1' });
      mockDB.all.mockResolvedValueOnce({ results: [{ id: 'e1' }, { id: 'e2' }] });
      mockDB.run.mockResolvedValue({});
      
      const res = await app.request('/api/v1/records/1', {
        method: 'DELETE',
      }, { env: mockEnv });
      
      const data = await res.json();
      
      expect(res.status).toBe(200);
      expect(data.message).toBe('Training record deleted successfully');
      
      expect(mockDB.prepare).toHaveBeenCalledWith(
        expect.stringContaining('DELETE FROM sets WHERE exercise_id = ?')
      );
      expect(mockDB.prepare).toHaveBeenCalledWith(
        expect.stringContaining('DELETE FROM exercises WHERE record_id = ?')
      );
      expect(mockDB.prepare).toHaveBeenCalledWith(
        expect.stringContaining('DELETE FROM training_records WHERE id = ?')
      );
    });
    
    it('should return 404 if record not found', async () => {
      mockDB.first.mockResolvedValueOnce(null);
      
      const res = await app.request('/api/v1/records/999', {
        method: 'DELETE',
      }, { env: mockEnv });
      
      const data = await res.json();
      
      expect(res.status).toBe(404);
      expect(data.error).toBe('Training record not found');
    });
  });
  
  describe('POST /api/v1/records/:id/exercises', () => {
    it('should add an exercise to a training record', async () => {
      const originalRandomUUID = crypto.randomUUID;
      crypto.randomUUID = vi.fn()
        .mockReturnValueOnce('exercise-id')
        .mockReturnValueOnce('set-id');
      
      mockDB.first.mockResolvedValueOnce({ id: '1' });
      mockDB.run.mockResolvedValue({});
      
      const res = await app.request('/api/v1/records/1/exercises', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          name: 'Deadlift',
          sets: [
            { weight: 200, reps: 5, notes: 'PR attempt' },
          ],
        }),
      }, { env: mockEnv });
      
      const data = await res.json();
      
      expect(res.status).toBe(201);
      expect(data.id).toBe('exercise-id');
      expect(data.message).toBe('Exercise added successfully');
      
      crypto.randomUUID = originalRandomUUID;
    });
    
    it('should return 400 for invalid input', async () => {
      const res = await app.request('/api/v1/records/1/exercises', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          sets: [],
        }),
      }, { env: mockEnv });
      
      const data = await res.json();
      
      expect(res.status).toBe(400);
      expect(data.error).toBe('Exercise name is required');
    });
  });
  
  describe('POST /api/v1/exercises/:id/sets', () => {
    it('should add a set to an exercise', async () => {
      const originalRandomUUID = crypto.randomUUID;
      crypto.randomUUID = vi.fn().mockReturnValueOnce('set-id');
      
      mockDB.first.mockResolvedValueOnce({ id: 'e1', record_id: '1' });
      mockDB.run.mockResolvedValue({});
      
      const res = await app.request('/api/v1/exercises/e1/sets', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          weight: 100,
          reps: 10,
          notes: 'Last set',
        }),
      }, { env: mockEnv });
      
      const data = await res.json();
      
      expect(res.status).toBe(201);
      expect(data.id).toBe('set-id');
      expect(data.message).toBe('Set added successfully');
      
      crypto.randomUUID = originalRandomUUID;
    });
    
    it('should return 400 for invalid input', async () => {
      const res = await app.request('/api/v1/exercises/e1/sets', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          reps: 10,
        }),
      }, { env: mockEnv });
      
      const data = await res.json();
      
      expect(res.status).toBe(400);
      expect(data.error).toBe('Weight and reps are required');
    });
  });
});
