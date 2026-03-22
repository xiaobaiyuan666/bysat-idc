import "dotenv/config";

import { execSync } from "node:child_process";
import path from "node:path";

import Database from "better-sqlite3";

function resolveDatabasePath(url: string) {
  if (!url.startsWith("file:")) {
    throw new Error(`Unsupported DATABASE_URL: ${url}`);
  }

  const rawPath = url.slice("file:".length);
  if (path.isAbsolute(rawPath)) {
    return rawPath;
  }

  return path.resolve(process.cwd(), "prisma", rawPath);
}

function tableExists(db: Database.Database, tableName: string) {
  const row = db
    .prepare(
      "SELECT name FROM sqlite_master WHERE type = 'table' AND name = ? LIMIT 1",
    )
    .get(tableName);

  return Boolean(row);
}

function main() {
  const databaseUrl = process.env.DATABASE_URL ?? "file:./dev.db";
  const databasePath = resolveDatabasePath(databaseUrl);
  const db = new Database(databasePath);
  const hasSchema = tableExists(db, "AdminUser");

  db.close();

  const command = hasSchema
    ? `npx prisma migrate diff --from-url "${databaseUrl}" --to-schema-datamodel prisma/schema.prisma --script`
    : "npx prisma migrate diff --from-empty --to-schema-datamodel prisma/schema.prisma --script";

  const sql = execSync(command, {
    cwd: process.cwd(),
    encoding: "utf8",
    stdio: ["ignore", "pipe", "inherit"],
  });

  const freshDb = new Database(databasePath);
  if (sql.trim()) {
    freshDb.exec(sql);
  }
  freshDb.close();

  console.log(
    `${hasSchema ? "Database schema updated" : "Database schema initialized"} at ${databasePath}`,
  );
}

main();
