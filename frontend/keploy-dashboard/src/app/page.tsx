
// import { Github } from "lucide-react";
// import TimeSeriesChart from "./components/charts/timeSeries";
import GitHubAnalyticsDashboard from "./components/charts/example";
// import { DataProvider } from "@/context/data-context";

export default function Home() {

  const data = [
    { date: '2024-01-01', value: 10 },
    { date: '2024-01-08', value: 15 },
    { date: '2024-01-15', value: 12 },
    { date: '2024-01-22', value: 20 },
    { date: '2024-01-29', value: 18 },
  ];

  
  return (
   <>
   {/* <TimeSeriesChart/> */}
   <GitHubAnalyticsDashboard/>
   </>
  );
}
