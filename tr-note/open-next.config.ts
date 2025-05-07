import { defineCloudflareConfig } from "@opennextjs/cloudflare";
import kvIncrementalCache from "@opennextjs/cloudflare/kv-cache";

export default defineCloudflareConfig({
  incrementalCache: kvIncrementalCache,
  d1Databases: [
    {
      binding: 'DB',
      databaseName: 'tr_note_db',
    },
  ],
});
