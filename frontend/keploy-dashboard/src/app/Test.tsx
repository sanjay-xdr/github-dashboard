"use client"
import React, { useEffect, useState } from "react";
import { Bar, BarChart, Legend, Tooltip, XAxis, YAxis, ResponsiveContainer, LineChart, Line, PieChart, Pie, Cell } from "recharts";

export const Dashboard = () => {
  const [prData, setPrData] = useState<any[] | null>(null);
  const [repoData, setRepoData] = useState<any[] | null>(null);
  const [workflowData, setWorkflowData] = useState<any[] | null>(null);
  const [filteredPrData, setFilteredPrData] = useState<any[] | null>(null);
  const [filteredRepoData, setFilteredRepoData] = useState<any[] | null>(null);
  const [filteredWorkflowData, setFilteredWorkflowData] = useState<any[] | null>(null);
  const [timeFilter, setTimeFilter] = useState<string>("30days");

  const WORKFLOW_COLORS = {
    success: "#4CAF50", // Green
    failed: "#F44336",  // Red
    cancelled: "#9E9E9E", // Gray
    skipped: "#FF9800",  // Orange
    pending: "#2196F3"   // Blue
  };

  // Generate dummy data for repository stats
  const generateDummyRepoData = () => {
    const data = [];
    const today = new Date();
    let stars = 120;
    let watchers = 45;
    let forks = 30;
    
    for (let i = 30; i >= 0; i--) {
      const date = new Date(today);
      date.setDate(today.getDate() - i);
      
      // Random increment
      stars += Math.floor(Math.random() * 5);
      watchers += Math.floor(Math.random() * 2);
      forks += Math.random() > 0.7 ? 1 : 0;
      
      data.push({
        date: date.toISOString().split('T')[0],
        stars,
        watchers,
        forks
      });
    }
    
    return data;
  };

  // Generate dummy data for workflow status
  const generateDummyWorkflowData = () => {
    const data = [];
    const today = new Date();
    
    for (let i = 30; i >= 0; i--) {
      const date = new Date(today);
      date.setDate(today.getDate() - i);
      
      // Random data for each day
      data.push({
        date: date.toISOString().split('T')[0],
        success: Math.floor(Math.random() * 15) + 10,
        failed: Math.floor(Math.random() * 6),
        cancelled: Math.floor(Math.random() * 3),
        skipped: Math.floor(Math.random() * 4),
        pending: Math.floor(Math.random() * 2)
      });
    }
    
    return data;
  };

  useEffect(() => {
    const fetchData = async () => {
      try {
        // Fetch PR data from API
        const prResponse = await fetch("http://localhost:8080/api/v1/fetch-prs");
        let prData = await prResponse.json();
        prData = prData.map((item: any) => ({
          date: item.date || "",
          openPR: item.openPR ?? 0,
          closedPR: item.closedPR ?? 0,
          mergedPR: item.mergedPR ?? 0,
        }));
        setPrData(prData);

        // Try to fetch repo data, use dummy data if API isn't ready
        try {
          const repoResponse = await fetch("http://localhost:8080/api/v1/fetch-repo-stats");
          let repoData = await repoResponse.json();
          repoData = repoData.map((item: any) => ({
            date: item.date || "",
            stars: item.stars ?? 0,
            watchers: item.watchers ?? 0,
            forks: item.forks ?? 0,
          }));
          setRepoData(repoData);
        } catch (error) {
          console.log("Repository API not ready, using dummy data");
          setRepoData(generateDummyRepoData());
        }

        // Try to fetch workflow data, use dummy data if API isn't ready
        try {
          const workflowResponse = await fetch("http://localhost:8080/api/v1/fetch-workflows");
          let workflowData = await workflowResponse.json();
          workflowData = workflowData.map((item: any) => ({
            date: item.date || "",
            success: item.success ?? 0,
            failed: item.failed ?? 0,
            cancelled: item.cancelled ?? 0,
            skipped: item.skipped ?? 0,
            pending: item.pending ?? 0,
          }));
          setWorkflowData(workflowData);
        } catch (error) {
          console.log("Workflow API not ready, using dummy data");
          setWorkflowData(generateDummyWorkflowData());
        }
      } catch (error) {
        console.error("Error fetching data:", error);
        
        // If PR API fails, generate dummy data for everything 
        if (!prData) {
          console.log("Using dummy data for all charts");
          const dummyPrData = Array.from({ length: 31 }, (_, i) => {
            const date = new Date();
            date.setDate(date.getDate() - (30 - i));
            return {
              date: date.toISOString().split('T')[0],
              openPR: Math.floor(Math.random() * 10) + 5,
              closedPR: Math.floor(Math.random() * 8),
              mergedPR: Math.floor(Math.random() * 12),
            };
          });
          setPrData(dummyPrData);
          setRepoData(generateDummyRepoData());
          setWorkflowData(generateDummyWorkflowData());
        }
      }
    };

    fetchData();
  }, []);

  useEffect(() => {
    if (!prData && !repoData && !workflowData) return;
    
    const filterDataByTime = () => {
      const currentDate = new Date();
      let filterDate = new Date();
      
      switch(timeFilter) {
        case "24hours":
          filterDate.setHours(currentDate.getHours() - 24);
          break;
        case "7days":
          filterDate.setDate(currentDate.getDate() - 7);
          break;
        case "30days":
        default:
          filterDate.setDate(currentDate.getDate() - 30);
          break;
      }
      
      // Filter PR data
      if (prData) {
        const filtered = prData.filter(item => {
          const itemDate = new Date(item.date);
          return itemDate >= filterDate;
        });
        setFilteredPrData(filtered);
      }
      
      // Filter repo data
      if (repoData) {
        const filtered = repoData.filter(item => {
          const itemDate = new Date(item.date);
          return itemDate >= filterDate;
        });
        setFilteredRepoData(filtered);
      }
      
      // Filter workflow data
      if (workflowData) {
        const filtered = workflowData.filter(item => {
          const itemDate = new Date(item.date);
          return itemDate >= filterDate;
        });
        setFilteredWorkflowData(filtered);
      }
    };
    
    filterDataByTime();
  }, [prData, repoData, workflowData, timeFilter]);

  const handleFilterChange = (filter: string) => {
    setTimeFilter(filter);
  };

  // Calculate workflow summary for pie chart
  const getWorkflowSummary = () => {
    if (!filteredWorkflowData || filteredWorkflowData.length === 0) return [];
    
    const summary = {
      success: 0,
      failed: 0,
      cancelled: 0,
      skipped: 0,
      pending: 0
    };
    
    filteredWorkflowData.forEach(item => {
      summary.success += item.success;
      summary.failed += item.failed;
      summary.cancelled += item.cancelled;
      summary.skipped += item.skipped;
      summary.pending += item.pending;
    });
    
    return Object.entries(summary).map(([name, value]) => ({ name, value }));
  };

  if (!prData && !repoData && !workflowData) return <div>Loading...</div>;

  return (
    <div className="bg-gray-800 p-4 text-gray-200">
      <h1 className="text-2xl font-bold mb-4">Repository Dashboard</h1>
      
      <div className="mb-6 flex justify-start space-x-2">
        <button 
          onClick={() => handleFilterChange("24hours")}
          className={`px-4 py-2 rounded ${timeFilter === "24hours" ? "bg-blue-600 text-white" : "bg-gray-700 text-gray-200"}`}
        >
          Last 24 Hours
        </button>
        <button 
          onClick={() => handleFilterChange("7days")}
          className={`px-4 py-2 rounded ${timeFilter === "7days" ? "bg-blue-600 text-white" : "bg-gray-700 text-gray-200"}`}
        >
          Last 7 Days
        </button>
        <button 
          onClick={() => handleFilterChange("30days")}
          className={`px-4 py-2 rounded ${timeFilter === "30days" ? "bg-blue-600 text-white" : "bg-gray-700 text-gray-200"}`}
        >
          Last 30 Days
        </button>
      </div>
      
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        {/* PR Status Chart */}
        <div className="bg-gray-700 p-4 rounded-lg">
          <h2 className="text-xl font-semibold mb-4">Pull Request Status</h2>
          <ResponsiveContainer width="100%" height={300}>
            <BarChart data={filteredPrData || []} margin={{ bottom: 50 }}>
              <XAxis
                dataKey="date"
                stroke="#888"
                interval={0}
                textAnchor="end"
                angle={-45}
                dy={10}
                tickMargin={5}
                tickFormatter={(date) =>
                  new Date(date).toLocaleDateString("en-US", { month: "short", day: "numeric" })
                }
              />
              <YAxis stroke="#888" />
              <Tooltip />
              <Legend />
              <Bar dataKey="openPR" fill="#8884d8" name="Open" />
              <Bar dataKey="closedPR" fill="#82ca9d" name="Closed" />
              <Bar dataKey="mergedPR" fill="#ffc658" name="Merged" />
            </BarChart>
          </ResponsiveContainer>
        </div>
        
        {/* Repository Stats Chart */}
        <div className="bg-gray-700 p-4 rounded-lg">
          <h2 className="text-xl font-semibold mb-4">Repository Stats</h2>
          <ResponsiveContainer width="100%" height={300}>
            <LineChart data={filteredRepoData || []} margin={{ bottom: 50 }}>
              <XAxis
                dataKey="date"
                stroke="#888"
                interval={0}
                textAnchor="end"
                angle={-45}
                dy={10}
                tickMargin={5}
                tickFormatter={(date) =>
                  new Date(date).toLocaleDateString("en-US", { month: "short", day: "numeric" })
                }
              />
              <YAxis stroke="#888" />
              <Tooltip />
              <Legend />
              <Line type="monotone" dataKey="stars" stroke="#FF9800" name="Stars" dot={false} strokeWidth={2} />
              <Line type="monotone" dataKey="watchers" stroke="#03A9F4" name="Watchers" dot={false} strokeWidth={2} />
              <Line type="monotone" dataKey="forks" stroke="#E91E63" name="Forks" dot={false} strokeWidth={2} />
            </LineChart>
          </ResponsiveContainer>
        </div>
        
        {/* Workflow Status Charts */}
        <div className="bg-gray-700 p-4 rounded-lg">
          <h2 className="text-xl font-semibold mb-4">Workflow Runs</h2>
          <ResponsiveContainer width="100%" height={300}>
            <BarChart data={filteredWorkflowData || []} margin={{ bottom: 50 }}>
              <XAxis
                dataKey="date"
                stroke="#888"
                interval={0}
                textAnchor="end"
                angle={-45}
                dy={10}
                tickMargin={5}
                tickFormatter={(date) =>
                  new Date(date).toLocaleDateString("en-US", { month: "short", day: "numeric" })
                }
              />
              <YAxis stroke="#888" />
              <Tooltip />
              <Legend />
              <Bar dataKey="success" fill={WORKFLOW_COLORS.success} name="Success" />
              <Bar dataKey="failed" fill={WORKFLOW_COLORS.failed} name="Failed" />
              <Bar dataKey="cancelled" fill={WORKFLOW_COLORS.cancelled} name="Cancelled" />
              <Bar dataKey="skipped" fill={WORKFLOW_COLORS.skipped} name="Skipped" />
              <Bar dataKey="pending" fill={WORKFLOW_COLORS.pending} name="Pending" />
            </BarChart>
          </ResponsiveContainer>
        </div>
        
        {/* Workflow Summary Pie Chart */}
        <div className="bg-gray-700 p-4 rounded-lg">
          <h2 className="text-xl font-semibold mb-4">Workflow Summary</h2>
          <ResponsiveContainer width="100%" height={300}>
            <PieChart>
              <Pie
                data={getWorkflowSummary()}
                dataKey="value"
                nameKey="name"
                cx="50%"
                cy="50%"
                outerRadius={100}
                label={({ name, percent }) => `${name}: ${(percent * 100).toFixed(0)}%`}
              >
                {getWorkflowSummary().map((entry, index) => (
                  <Cell key={`cell-${index}`} fill={WORKFLOW_COLORS[entry.name as keyof typeof WORKFLOW_COLORS]} />
                ))}
              </Pie>
              <Tooltip formatter={(value) => [`${value} runs`, ""]} />
              <Legend />
            </PieChart>
          </ResponsiveContainer>
        </div>
      </div>
      
      {/* Summary Cards */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mt-6">
        <div className="bg-gray-700 p-4 rounded-lg">
          <h3 className="text-lg font-semibold mb-2">PR Activity</h3>
          <div className="text-3xl font-bold">
            {filteredPrData?.reduce((sum, item) => sum + item.openPR + item.closedPR + item.mergedPR, 0) || 0}
          </div>
          <div className="text-