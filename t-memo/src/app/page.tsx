import { getCloudflareContext } from "@opennextjs/cloudflare";
import PostList from "./components/PostList";
import styles from "./page.module.css";

type Post = {
  id: string;
  title: string;
  body: string;
  updatedAt: string;
};

export default async function Page() {
  const ctx = await getCloudflareContext({ async: true });
  const { keys } = await ctx.env.POST.list();
  const values = (
    await Promise.all(
      keys.map(async (key) => {
        const value = await ctx.env.POST.get(key.name, "json");
        if (value == null) return null;
        return value as Post;
      })
    )
  ).filter((v) => v !== null);

  return (
    <div className={styles.page}>
      <main className={styles.main}>
        <h1>t-memo</h1>
        <PostList posts={values} />
      </main>
    </div>
  );
}
