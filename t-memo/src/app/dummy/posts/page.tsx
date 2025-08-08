import { getCloudflareContext } from "@opennextjs/cloudflare";

export default async function Page() {
  const ctx = await getCloudflareContext({ async: true });
  await ctx.env.POST.put(
    Date.now().valueOf().toString(),
    JSON.stringify({
      title: "foo",
      body: "bar",
    })
  );

  return <div>done!</div>;
}
