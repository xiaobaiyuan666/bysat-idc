import "dotenv/config";

import { execSync } from "node:child_process";
import { existsSync, unlinkSync } from "node:fs";
import path from "node:path";

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

function main() {
  const databaseUrl = process.env.DATABASE_URL ?? "file:./dev.db";
  const databasePath = resolveDatabasePath(databaseUrl);

  if (existsSync(databasePath)) {
    unlinkSync(databasePath);
  }

  execSync("npm run db:push", {
    cwd: process.cwd(),
    stdio: "inherit",
  });

  execSync("npm run db:seed", {
    cwd: process.cwd(),
    stdio: "inherit",
  });
}

main();
