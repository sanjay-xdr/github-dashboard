"use client"
import React from 'react';
import { BarChart, Bar, LineChart, Line, XAxis, YAxis, Tooltip, Legend, PieChart, Pie, Cell, AreaChart, Area } from 'recharts';
import { 
  Activity, 
  GitPullRequest, 
  AlertTriangle, 
  Clock, 
  Users, 
  Star, 
  GitBranch, 
  Code, 
  Check, 
  X 
} from 'lucide-react';
// import { fetchDashboardData } from '@/libs/api';
import { useState,useEffect } from 'react';
import { useData } from '@/context/data-context';


interface RepoOverview{
  name : string,
  value:number
}
// Mock data (you'll replace with actual data fetching)
const pullRequestData = [
  { name: 'Jan', open: 45, closed: 35 },
  { name: 'Feb', open: 50, closed: 40 },
  { name: 'Mar', open: 55, closed: 45 },
];

const issueData = [
  { name: 'Jan', bugs: 20, enhancements: 15, questions: 10 },
  { name: 'Feb', bugs: 22, enhancements: 18, questions: 12 },
  { name: 'Mar', bugs: 25, enhancements: 20, questions: 15 },
];


const contributorData = [
  { name: 'Alice', commits: 250, prs: 45, issues: 30 },
  { name: 'Bob', commits: 200, prs: 35, issues: 25 },
  { name: 'Charlie', commits: 180, prs: 30, issues: 20 },
];

const codeQualityData = [
  { name: 'Jan', coverage: 75, churn: 20 },
  { name: 'Feb', coverage: 80, churn: 15 },
  { name: 'Mar', coverage: 85, churn: 10 },
];

const COLORS = ['#0088FE', '#00C49F', '#FFBB28', '#FF8042'];

const GitHubAnalyticsDashboard = () => {

  // const{data}=useData();
  // console.log(data);
  const [repositoryOverviewData, setRepositoryOverviewData] = useState<any[] | undefined>(undefined);


 

  return (
    <div className="bg-gray-900 text-white min-h-screen p-6">
      <header className="flex justify-between items-center mb-6">
        <h1 className="text-3xl font-bold">GitHub Analytics</h1>
        <div className="flex space-x-4">
          <select className="bg-gray-800 text-white p-2 rounded">
            <option>All Repositories</option>
            <option>Repo 1</option>
            <option>Repo 2</option>
          </select>
          <input 
            type="date" 
            className="bg-gray-800 text-white p-2 rounded"
          />
        </div>
      </header>

      <div className="grid grid-cols-3 gap-6">
        {/* Pull Request Analytics */}
        <div className="bg-gray-800 p-4 rounded-lg">
          <div className="flex items-center mb-4">
            <GitPullRequest className="mr-2" />
            <h2 className="text-xl font-semibold">Pull Request Analytics</h2>
          </div>
          <div className="grid grid-cols-2 gap-2 mb-4">
            <div className="bg-gray-700 p-2 rounded">
              <div className="text-sm text-gray-400">Avg PR Close Time</div>
              <div className="text-lg font-bold">3.2 days</div>
            </div>
            <div className="bg-gray-700 p-2 rounded">
              <div className="text-sm text-gray-400">Merge Ratio</div>
              <div className="text-lg font-bold">85%</div>
            </div>
          </div>
          <LineChart width={350} height={200} data={pullRequestData}>
            <XAxis dataKey="name" stroke="#888" />
            <YAxis stroke="#888" />
            <Tooltip />
            <Legend />
            <Line type="monotone" dataKey="open" stroke="#8884d8" />
            <Line type="monotone" dataKey="closed" stroke="#82ca9d" />
          </LineChart>
        </div>

        {/* Issue Analytics */}
        <div className="bg-gray-800 p-4 rounded-lg">
          <div className="flex items-center mb-4">
            <AlertTriangle className="mr-2" />
            <h2 className="text-xl font-semibold">Issue Analytics</h2>
          </div>
          <div className="grid grid-cols-2 gap-2 mb-4">
            <div className="bg-gray-700 p-2 rounded">
              <div className="text-sm text-gray-400">Avg Issue Close Time</div>
              <div className="text-lg font-bold">4.5 days</div>
            </div>
            <div className="bg-gray-700 p-2 rounded">
              <div className="text-sm text-gray-400">Stale Issues</div>
              <div className="text-lg font-bold">12</div>
            </div>
          </div>
          <BarChart width={350} height={200} data={issueData}>
            <XAxis dataKey="name" stroke="#888" />
            <YAxis stroke="#888" />
            <Tooltip />
            <Legend />
            <Bar dataKey="bugs" fill="#8884d8" />
            <Bar dataKey="enhancements" fill="#82ca9d" />
            <Bar dataKey="questions" fill="#ffc658" />
          </BarChart>
        </div>

        {/* Repository Overview */}
        <div className="bg-gray-800 p-4 rounded-lg">
          <div className="flex items-center mb-4">
            <Star className="mr-2" />
            <h2 className="text-xl font-semibold">Repository Overview</h2>
          </div>
          <PieChart width={350} height={250}>
            <Pie
              data={repositoryOverviewData}
              cx={175}
              cy={125}
              labelLine={false}
              outerRadius={80}
              fill="#8884d8"
              dataKey="value"
            >
              {repositoryOverviewData?.map((entry, index) => (
                <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
              ))}
            </Pie>
            <Tooltip />
            <Legend />
          </PieChart>
          <div className="grid grid-cols-3 gap-2 mt-4">
            {repositoryOverviewData?.map((item, index) => (
              <div key={item.name} className="bg-gray-700 p-2 rounded text-center">
                <div className="text-sm text-gray-400">{item.name}</div>
                <div className="text-lg font-bold">{item.value}</div>
              </div>
            ))}
          </div>
        </div>

        {/* Contributors */}
        <div className="bg-gray-800 p-4 rounded-lg col-span-2">
          <div className="flex items-center mb-4">
            <Users className="mr-2" />
            <h2 className="text-xl font-semibold">Top Contributors</h2>
          </div>
          <table className="w-full">
            <thead className="bg-gray-700">
              <tr>
                <th className="p-2 text-left">Contributor</th>
                <th className="p-2">Commits</th>
                <th className="p-2">Pull Requests</th>
                <th className="p-2">Issues</th>
              </tr>
            </thead>
            <tbody>
              {contributorData.map((contributor) => (
                <tr key={contributor.name} className="border-b border-gray-700">
                  <td className="p-2">{contributor.name}</td>
                  <td className="p-2 text-center">{contributor.commits}</td>
                  <td className="p-2 text-center">{contributor.prs}</td>
                  <td className="p-2 text-center">{contributor.issues}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>

        {/* Code Quality */}
        <div className="bg-gray-800 p-4 rounded-lg">
          <div className="flex items-center mb-4">
            <Code className="mr-2" />
            <h2 className="text-xl font-semibold">Code Quality</h2>
          </div>
          <AreaChart width={350} height={250} data={codeQualityData}>
            <XAxis dataKey="name" stroke="#888" />
            <YAxis stroke="#888" />
            <Tooltip />
            <Area type="monotone" dataKey="coverage" stroke="#8884d8" fill="#8884d8" fillOpacity={0.3} />
            <Area type="monotone" dataKey="churn" stroke="#82ca9d" fill="#82ca9d" fillOpacity={0.3} />
          </AreaChart>
          <div className="grid grid-cols-2 gap-2 mt-4">
            <div className="bg-gray-700 p-2 rounded">
              <div className="text-sm text-gray-400">Test Coverage</div>
              <div className="text-lg font-bold">85%</div>
            </div>
            <div className="bg-gray-700 p-2 rounded">
              <div className="text-sm text-gray-400">Code Churn</div>
              <div className="text-lg font-bold">15%</div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default GitHubAnalyticsDashboard;

