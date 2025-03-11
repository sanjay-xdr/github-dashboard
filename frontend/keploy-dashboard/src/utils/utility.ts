export const fetchDashboardData = async () => {
    const response = await fetch('http://localhost:8080/api/v1/repodata');
    if (!response.ok) {
      throw new Error('Failed to fetch dashboard data');
    }
    return response.json();
  };