"use server";

import { D1Database } from '@cloudflare/workers-types';
import { getCloudflareContext } from '@opennextjs/cloudflare';

/**
 * データベース接続を取得する
 */
export async function getDB(): Promise<D1Database> {
  // Cloudflare ContextからD1を取得
  const { env } = getCloudflareContext();
  
  if (!env?.DB) {
    throw new Error('D1 database is not available. Make sure D1 binding is properly set up in wrangler.jsonc and open-next.config.ts.');
  }
  
  return env.DB as D1Database;
}

/**
 * ユニークIDを生成する関数
 */
export async function generateId(): Promise<string> {
  return crypto.randomUUID();
}

/**
 * 現在の日時をISO文字列で取得する関数
 */
export async function getCurrentISOString(): Promise<string> {
  return new Date().toISOString();
}