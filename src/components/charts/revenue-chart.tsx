"use client";

import { Area, AreaChart, CartesianGrid, ResponsiveContainer, Tooltip, XAxis, YAxis } from "recharts";

import { formatCurrency } from "@/lib/format";

export function RevenueChart({
  data,
}: {
  data: Array<{ month: string; amount: number }>;
}) {
  return (
    <div className="h-72 w-full">
      <ResponsiveContainer width="100%" height="100%">
        <AreaChart data={data}>
          <defs>
            <linearGradient id="revenueFill" x1="0" x2="0" y1="0" y2="1">
              <stop offset="5%" stopColor="#d97a31" stopOpacity={0.38} />
              <stop offset="95%" stopColor="#d97a31" stopOpacity={0.04} />
            </linearGradient>
          </defs>
          <CartesianGrid stroke="rgba(19,34,56,0.08)" vertical={false} />
          <XAxis dataKey="month" tickLine={false} axisLine={false} tickMargin={10} />
          <YAxis
            tickLine={false}
            axisLine={false}
            tickFormatter={(value) => `¥${Math.round(value / 1000) / 10}k`}
            width={68}
          />
          <Tooltip
            formatter={(value) =>
              formatCurrency(typeof value === "number" ? value : 0)
            }
            contentStyle={{
              borderRadius: 20,
              border: "1px solid rgba(19,34,56,0.1)",
              background: "rgba(255,253,248,0.96)",
            }}
          />
          <Area
            type="monotone"
            dataKey="amount"
            stroke="#d97a31"
            strokeWidth={3}
            fill="url(#revenueFill)"
          />
        </AreaChart>
      </ResponsiveContainer>
    </div>
  );
}
