function escapeCsvCell(value: unknown) {
  const text =
    value === null || value === undefined
      ? ""
      : value instanceof Date
        ? value.toISOString()
        : String(value);

  if (/[",\r\n]/.test(text)) {
    return `"${text.replace(/"/g, '""')}"`;
  }

  return text;
}

export function createCsv(rows: Array<Record<string, unknown>>) {
  if (rows.length === 0) {
    return "\uFEFF";
  }

  const headers = Object.keys(rows[0]);
  const lines = [
    headers.map((header) => escapeCsvCell(header)).join(","),
    ...rows.map((row) =>
      headers.map((header) => escapeCsvCell(row[header])).join(","),
    ),
  ];

  return `\uFEFF${lines.join("\r\n")}`;
}
