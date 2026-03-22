function escapeCsvCell(value: unknown) {
  const text = String(value ?? "");
  const escaped = text.replaceAll('"', '""');
  return /[",\r\n]/.test(text) ? `"${escaped}"` : escaped;
}

export function downloadCsv(filename: string, headers: string[], rows: Array<Array<unknown>>) {
  const content = [headers, ...rows]
    .map((row) => row.map((cell) => escapeCsvCell(cell)).join(","))
    .join("\r\n");

  const blob = new Blob([`\uFEFF${content}`], {
    type: "text/csv;charset=utf-8;",
  });
  const url = URL.createObjectURL(blob);
  const anchor = document.createElement("a");
  anchor.href = url;
  anchor.download = filename;
  document.body.appendChild(anchor);
  anchor.click();
  document.body.removeChild(anchor);
  URL.revokeObjectURL(url);
}
