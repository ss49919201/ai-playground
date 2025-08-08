"use server";

import { getCloudflareContext } from "@opennextjs/cloudflare";
import { revalidatePath } from "next/cache";

export async function updatePost(postId: string, title: string, body: string) {
  const ctx = getCloudflareContext();
  
  const post = {
    id: postId,
    title: title,
    body: body,
  };

  await ctx.env.POST.put(postId, JSON.stringify(post));
  revalidatePath("/");
}

export async function createPost(title: string, body: string) {
  const ctx = getCloudflareContext();
  
  const postId = crypto.randomUUID();
  const post = {
    id: postId,
    title: title,
    body: body,
  };

  await ctx.env.POST.put(postId, JSON.stringify(post));
  revalidatePath("/");
}

export async function deletePost(postId: string) {
  const ctx = getCloudflareContext();
  
  await ctx.env.POST.delete(postId);
  revalidatePath("/");
}