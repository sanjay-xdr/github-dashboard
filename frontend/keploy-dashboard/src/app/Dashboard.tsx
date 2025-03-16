"use client"
import React, { useEffect, useState } from "react";
import { Bar, BarChart, Legend, Tooltip, XAxis, YAxis, ResponsiveContainer } from "recharts";

export const Dashboard = () => {
  const [prData, setPrData] = useState<any[] | null>(null);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await fetch("http://localhost:8080/api/v1/fetch-prs");
        let data = await response.json();

        data = data.map((item: any) => ({
          date: item.date || "",
          openPR: item.openPR ?? 0,
          closedPR: item.closedPR ?? 0,
          mergedPR: item.mergedPR ?? 0,
        }));

        setPrData(data);
      } catch (error) {
        console.error("Error fetching data:", error);
      }
    };

    fetchData();
  }, []);

  if (!prData) return <div>Loading...</div>;

  return (
    <div className="bg-gray-800 p-4 ">
      <ResponsiveContainer width="50%" height={300}>
        <BarChart data={prData} margin={{ bottom: 50 }}>
          <XAxis
            dataKey="date"
            stroke="#888"
            interval={0}
            textAnchor="end"
            dy={10}
            tickMargin={-10}
            tickFormatter={(date) =>
              new Date(date).toLocaleDateString("en-US", { month: "short", day: "numeric" })
            }
          />
          <YAxis stroke="#888" />
          <Tooltip />
          <Legend />
          <Bar dataKey="openPR" fill="#8884d8" />
          <Bar dataKey="closedPR" fill="#82ca9d" />
          <Bar dataKey="mergedPR" fill="#ffc658" />
        </BarChart>
      </ResponsiveContainer>
    </div>
  );
};
