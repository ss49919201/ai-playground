import { Hono } from "hono";
import { jwt } from "hono/jwt";

interface Env {
  DB: D1Database;
  JWT_SECRET: string;
}

const api = new Hono<{ Bindings: Env }>();

api.use("*", async (c, next) => {
  if (process.env.NODE_ENV === "development") {
    return next();
  }
  
  return jwt({
    secret: c.env.JWT_SECRET,
  })(c, next);
});

api.get("/records", async (c) => {
  try {
    const { results } = await c.env.DB.prepare(
      "SELECT * FROM training_records ORDER BY created_at DESC"
    ).all();
    
    return c.json({ records: results });
  } catch (error) {
    console.error("Error fetching records:", error);
    return c.json({ error: "Failed to fetch training records" }, 500);
  }
});

api.get("/records/:id", async (c) => {
  const id = c.req.param("id");
  
  try {
    const record = await c.env.DB.prepare(
      "SELECT * FROM training_records WHERE id = ?"
    ).bind(id).first();
    
    if (!record) {
      return c.json({ error: "Training record not found" }, 404);
    }
    
    const { results: exercises } = await c.env.DB.prepare(
      "SELECT * FROM exercises WHERE record_id = ?"
    ).bind(id).all();
    
    for (const exercise of exercises) {
      const { results: sets } = await c.env.DB.prepare(
        "SELECT * FROM sets WHERE exercise_id = ?"
      ).bind(exercise.id).all();
      
      exercise.sets = sets;
    }
    
    record.exercises = exercises;
    
    return c.json({ record });
  } catch (error) {
    console.error("Error fetching record:", error);
    return c.json({ error: "Failed to fetch training record" }, 500);
  }
});

api.post("/records", async (c) => {
  try {
    const data = await c.req.json();
    const { title, date, description, exercises = [] } = data;
    
    if (!title || !date) {
      return c.json({ error: "Title and date are required" }, 400);
    }
    
    const timestamp = new Date().toISOString();
    const id = crypto.randomUUID();
    
    await c.env.DB.prepare(
      "INSERT INTO training_records (id, title, date, description, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)"
    ).bind(id, title, date, description || "", timestamp, timestamp).run();
    
    for (const exercise of exercises) {
      const exerciseId = crypto.randomUUID();
      
      await c.env.DB.prepare(
        "INSERT INTO exercises (id, record_id, name) VALUES (?, ?, ?)"
      ).bind(exerciseId, id, exercise.name).run();
      
      for (const set of exercise.sets || []) {
        const setId = crypto.randomUUID();
        
        await c.env.DB.prepare(
          "INSERT INTO sets (id, exercise_id, weight, reps, notes) VALUES (?, ?, ?, ?, ?)"
        ).bind(setId, exerciseId, set.weight, set.reps, set.notes || "").run();
      }
    }
    
    return c.json({ id, message: "Training record created successfully" }, 201);
  } catch (error) {
    console.error("Error creating record:", error);
    return c.json({ error: "Failed to create training record" }, 500);
  }
});

api.put("/records/:id", async (c) => {
  const id = c.req.param("id");
  
  try {
    const data = await c.req.json();
    const { title, date, description } = data;
    
    if (!title || !date) {
      return c.json({ error: "Title and date are required" }, 400);
    }
    
    const record = await c.env.DB.prepare(
      "SELECT * FROM training_records WHERE id = ?"
    ).bind(id).first();
    
    if (!record) {
      return c.json({ error: "Training record not found" }, 404);
    }
    
    const timestamp = new Date().toISOString();
    
    await c.env.DB.prepare(
      "UPDATE training_records SET title = ?, date = ?, description = ?, updated_at = ? WHERE id = ?"
    ).bind(title, date, description || "", timestamp, id).run();
    
    return c.json({ message: "Training record updated successfully" });
  } catch (error) {
    console.error("Error updating record:", error);
    return c.json({ error: "Failed to update training record" }, 500);
  }
});

api.delete("/records/:id", async (c) => {
  const id = c.req.param("id");
  
  try {
    const record = await c.env.DB.prepare(
      "SELECT * FROM training_records WHERE id = ?"
    ).bind(id).first();
    
    if (!record) {
      return c.json({ error: "Training record not found" }, 404);
    }
    
    const { results: exercises } = await c.env.DB.prepare(
      "SELECT id FROM exercises WHERE record_id = ?"
    ).bind(id).all();
    
    for (const exercise of exercises) {
      await c.env.DB.prepare(
        "DELETE FROM sets WHERE exercise_id = ?"
      ).bind(exercise.id).run();
    }
    
    await c.env.DB.prepare(
      "DELETE FROM exercises WHERE record_id = ?"
    ).bind(id).run();
    
    await c.env.DB.prepare(
      "DELETE FROM training_records WHERE id = ?"
    ).bind(id).run();
    
    return c.json({ message: "Training record deleted successfully" });
  } catch (error) {
    console.error("Error deleting record:", error);
    return c.json({ error: "Failed to delete training record" }, 500);
  }
});

api.post("/records/:id/exercises", async (c) => {
  const recordId = c.req.param("id");
  
  try {
    const data = await c.req.json();
    const { name, sets = [] } = data;
    
    if (!name) {
      return c.json({ error: "Exercise name is required" }, 400);
    }
    
    const record = await c.env.DB.prepare(
      "SELECT * FROM training_records WHERE id = ?"
    ).bind(recordId).first();
    
    if (!record) {
      return c.json({ error: "Training record not found" }, 404);
    }
    
    const exerciseId = crypto.randomUUID();
    
    await c.env.DB.prepare(
      "INSERT INTO exercises (id, record_id, name) VALUES (?, ?, ?)"
    ).bind(exerciseId, recordId, name).run();
    
    for (const set of sets) {
      const setId = crypto.randomUUID();
      
      await c.env.DB.prepare(
        "INSERT INTO sets (id, exercise_id, weight, reps, notes) VALUES (?, ?, ?, ?, ?)"
      ).bind(setId, exerciseId, set.weight, set.reps, set.notes || "").run();
    }
    
    const timestamp = new Date().toISOString();
    await c.env.DB.prepare(
      "UPDATE training_records SET updated_at = ? WHERE id = ?"
    ).bind(timestamp, recordId).run();
    
    return c.json({ id: exerciseId, message: "Exercise added successfully" }, 201);
  } catch (error) {
    console.error("Error adding exercise:", error);
    return c.json({ error: "Failed to add exercise" }, 500);
  }
});

api.post("/exercises/:id/sets", async (c) => {
  const exerciseId = c.req.param("id");
  
  try {
    const data = await c.req.json();
    const { weight, reps, notes } = data;
    
    if (weight === undefined || reps === undefined) {
      return c.json({ error: "Weight and reps are required" }, 400);
    }
    
    const exercise = await c.env.DB.prepare(
      "SELECT * FROM exercises WHERE id = ?"
    ).bind(exerciseId).first();
    
    if (!exercise) {
      return c.json({ error: "Exercise not found" }, 404);
    }
    
    const setId = crypto.randomUUID();
    
    await c.env.DB.prepare(
      "INSERT INTO sets (id, exercise_id, weight, reps, notes) VALUES (?, ?, ?, ?, ?)"
    ).bind(setId, exerciseId, weight, reps, notes || "").run();
    
    const timestamp = new Date().toISOString();
    await c.env.DB.prepare(
      "UPDATE training_records SET updated_at = ? WHERE id = (SELECT record_id FROM exercises WHERE id = ?)"
    ).bind(timestamp, exerciseId).run();
    
    return c.json({ id: setId, message: "Set added successfully" }, 201);
  } catch (error) {
    console.error("Error adding set:", error);
    return c.json({ error: "Failed to add set" }, 500);
  }
});

export { api as apiRoutes };
