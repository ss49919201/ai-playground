"use server";

import { D1Database } from '@cloudflare/workers-types';
import { User, Session } from '../types';
import { getDB, generateId, getCurrentISOString } from './index';
import bcrypt from 'bcrypt';

const SALT_ROUNDS = 10;

/**
 * パスワードをハッシュ化する
 */
export async function hashPassword(password: string): Promise<string> {
  return bcrypt.hash(password, SALT_ROUNDS);
}

/**
 * パスワードを検証する
 */
export async function verifyPassword(password: string, hash: string): Promise<boolean> {
  return bcrypt.compare(password, hash);
}

/**
 * ユーザーを作成する
 */
export async function createUser(
  email: string, 
  password: string, 
  name?: string
): Promise<User> {
  const db = await getDB();
  const id = await generateId();
  const timestamp = await getCurrentISOString();
  const passwordHash = await hashPassword(password);
  
  await db.prepare(
    `INSERT INTO users (id, email, password_hash, name, created_at, updated_at)
     VALUES (?, ?, ?, ?, ?, ?)`
  ).bind(
    id,
    email.toLowerCase(),
    passwordHash,
    name || null,
    timestamp,
    timestamp
  ).run();
  
  return {
    id,
    email: email.toLowerCase(),
    name,
    createdAt: timestamp,
    updatedAt: timestamp
  };
}

/**
 * メールアドレスからユーザーを検索する
 */
export async function getUserByEmail(email: string): Promise<User | null> {
  const db = await getDB();
  const user = await db.prepare(
    `SELECT id, email, name, created_at as createdAt, updated_at as updatedAt
     FROM users
     WHERE email = ?`
  ).bind(email.toLowerCase()).first<User>();
  
  return user || null;
}

/**
 * IDからユーザーを検索する
 */
export async function getUserById(id: string): Promise<User | null> {
  const db = await getDB();
  const user = await db.prepare(
    `SELECT id, email, name, created_at as createdAt, updated_at as updatedAt
     FROM users
     WHERE id = ?`
  ).bind(id).first<User>();
  
  return user || null;
}

/**
 * ユーザー情報を更新する
 */
export async function updateUser(
  id: string,
  updates: { name?: string; email?: string; password?: string }
): Promise<User | null> {
  const db = await getDB();
  const timestamp = await getCurrentISOString();
  
  // 現在のユーザー情報を取得
  const existingUser = await getUserById(id);
  if (!existingUser) {
    return null;
  }
  
  // 更新するフィールドを準備
  const updateFields: string[] = [];
  const values: any[] = [];
  
  if (updates.name !== undefined) {
    updateFields.push("name = ?");
    values.push(updates.name || null);
  }
  
  if (updates.email !== undefined) {
    updateFields.push("email = ?");
    values.push(updates.email.toLowerCase());
  }
  
  if (updates.password !== undefined) {
    updateFields.push("password_hash = ?");
    values.push(await hashPassword(updates.password));
  }
  
  // 更新日時は常に更新
  updateFields.push("updated_at = ?");
  values.push(timestamp);
  
  // IDは最後に追加
  values.push(id);
  
  if (updateFields.length > 0) {
    await db.prepare(
      `UPDATE users
       SET ${updateFields.join(", ")}
       WHERE id = ?`
    ).bind(...values).run();
  }
  
  return await getUserById(id);
}

/**
 * ユーザーのパスワードハッシュを取得する
 */
export async function getUserPasswordHash(id: string): Promise<string | null> {
  const db = await getDB();
  const result = await db.prepare(
    `SELECT password_hash as passwordHash
     FROM users
     WHERE id = ?`
  ).bind(id).first<{ passwordHash: string }>();
  
  return result?.passwordHash || null;
}

/**
 * 新しいセッションを作成する
 */
export async function createSession(userId: string, expiresInHours: number = 24): Promise<Session> {
  const db = await getDB();
  const id = await generateId();
  const createdAt = await getCurrentISOString();
  
  // 有効期限を計算
  const expiresAt = new Date();
  expiresAt.setHours(expiresAt.getHours() + expiresInHours);
  const expiresAtString = expiresAt.toISOString();
  
  await db.prepare(
    `INSERT INTO sessions (id, user_id, expires_at, created_at)
     VALUES (?, ?, ?, ?)`
  ).bind(id, userId, expiresAtString, createdAt).run();
  
  return {
    id,
    userId,
    expiresAt: expiresAtString,
    createdAt
  };
}

/**
 * セッションを取得する
 */
export async function getSessionById(id: string): Promise<Session | null> {
  const db = await getDB();
  const session = await db.prepare(
    `SELECT id, user_id as userId, expires_at as expiresAt, created_at as createdAt
     FROM sessions
     WHERE id = ?`
  ).bind(id).first<Session>();
  
  if (!session) {
    return null;
  }
  
  // 有効期限が切れていないか確認
  if (new Date(session.expiresAt) < new Date()) {
    await deleteSession(id);
    return null;
  }
  
  return session;
}

/**
 * ユーザーIDに関連する全てのセッションを取得する
 */
export async function getSessionsByUserId(userId: string): Promise<Session[]> {
  const db = await getDB();
  
  const sessions = await db.prepare(
    `SELECT id, user_id AS userId, expires_at AS expiresAt, created_at AS createdAt
     FROM sessions
     WHERE user_id = ?`
  ).bind(userId).all<Session>();
  
  return sessions.results;
}

/**
 * セッションを削除する
 */
export async function deleteSession(id: string): Promise<void> {
  const db = await getDB();
  await db.prepare(
    `DELETE FROM sessions
     WHERE id = ?`
  ).bind(id).run();
}

/**
 * ユーザーの全てのセッションを削除する
 */
export async function deleteAllUserSessions(userId: string): Promise<void> {
  const db = await getDB();
  
  await db.prepare(
    `DELETE FROM sessions
     WHERE user_id = ?`
  ).bind(userId).run();
}

/**
 * 期限切れのセッションをクリーンアップする
 */
export async function cleanupExpiredSessions(): Promise<number> {
  const db = await getDB();
  const now = await getCurrentISOString();
  
  const result = await db.prepare(
    `DELETE FROM sessions
     WHERE expires_at < ?`
  ).bind(now).run();
  
  return result.meta.changes;
}