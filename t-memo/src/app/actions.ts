"use server";

import { getCloudflareContext } from "@opennextjs/cloudflare";
import { revalidatePath } from "next/cache";

export async function updatePost(postId: string, title: string, body: string) {
  const ctx = await getCloudflareContext({ async: true });

  const post = {
    id: postId,
    title: title,
    body: body,
  };

  await ctx.env.POST.put(postId, JSON.stringify(post));
  revalidatePath("/");
}

export async function createPost(title: string, body: string) {
  const ctx = await getCloudflareContext({ async: true });

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
  const ctx = await getCloudflareContext({ async: true });

  await ctx.env.POST.delete(postId);
  revalidatePath("/");
}
