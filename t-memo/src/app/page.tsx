import { getCloudflareContext } from "@opennextjs/cloudflare";
import styles from "./page.module.css";

type Post = {
  id: string;
  title: string;
  body: string;
};

export default async function Home() {
  const ctx = getCloudflareContext();
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
        <div className={styles.postList}>
          {values.map((item) => (
            <div key={item.id} className={styles.postItem}>
              <div className={styles.postHeader}>
                <h2 className={`${styles.postTitle}`}>{item.title}</h2>
              </div>
              <p className={styles.postBody}>{item.body}</p>
            </div>
          ))}
        </div>
      </main>
    </div>
  );
}
