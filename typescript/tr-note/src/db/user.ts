"use server";

import { D1Database } from '@cloudflare/workers-types';
import { User, Session } from '../types';
import { getDB, generateId, getCurrentISOString } from './index';

/**
 * メールアドレスからユーザーを検索する
 */
export async function getUserByEmail(email: string): Promise<User | null> {
  const db = await getDB();
  const user = await db.prepare(
    `SELECT id, email, name, created_at as createdAt, updated_at as updatedAt
     FROM users
     WHERE email = ?`
  ).bind(email).first<User>();
  
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