"use client"
import React, { useMemo } from "react";
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend } from "recharts";
import { format, parseISO } from "date-fns";

export default function TimeSeriesChart() {
  const data = [
    { date: "2024-01-01", value: 10 },
    { date: "2024-01-08", value: 15 },
    { date: "2024-01-15", value: 12 },
    { date: "2024-01-22", value: 20 },
    { date: "2024-01-29", value: 18 },
  ];

  // Ensure tickFormatter runs on client-side only
  const tickFormatter = useMemo(() => (date: string) => format(parseISO(date), "MMM dd"), []);

  return (
    <LineChart width={600} height={300} data={data} margin={{ top: 5, right: 30, left: 20, bottom: 5 }}>
      <CartesianGrid strokeDasharray="3 3" />
      <XAxis dataKey="date" tickFormatter={tickFormatter} />
      <YAxis />
      <Tooltip />
      <Legend />
      <Line type="monotone" dataKey="value" stroke="#8884d8" activeDot={{ r: 8 }} />
    </LineChart>
  );
}
