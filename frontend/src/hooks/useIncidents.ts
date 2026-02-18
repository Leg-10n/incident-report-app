import { useState, useEffect, useCallback } from "react";
import { Incident, IncidentRequest } from "../types/incident";

const API_BASE = "/api/incidents";

export function useIncidents(token: string | null) {
  const [incidents, setIncidents] = useState<Incident[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Helper so we don't repeat the Authorization header everywhere
  const authHeaders = () => ({
    Authorization: `Bearer ${token}`,
    "Content-Type": "application/json",
  });

  const fetchIncidents = useCallback(async () => {
    if (!token) return;
    try {
      setLoading(true);
      const res = await fetch(API_BASE, {
        headers: { Authorization: `Bearer ${token}` },
      });
      if (!res.ok) throw new Error("Failed to fetch incidents");
      const data: Incident[] = await res.json();
      setIncidents(data);
      setError(null);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Unknown error");
    } finally {
      setLoading(false);
    }
  }, [token]); // re-fetch if token changes (e.g. login/logout)

  useEffect(() => {
    fetchIncidents();
  }, [fetchIncidents]);

  const createIncident = async (req: IncidentRequest): Promise<boolean> => {
    try {
      const res = await fetch(API_BASE, {
        method: "POST",
        headers: authHeaders(),
        body: JSON.stringify(req),
      });
      if (!res.ok) {
        const data = await res.json();
        throw new Error(data.error || "Failed to create incident");
      }
      await fetchIncidents();
      return true;
    } catch (err) {
      setError(err instanceof Error ? err.message : "Unknown error");
      return false;
    }
  };

  const updateIncident = async (
    id: number,
    req: IncidentRequest,
  ): Promise<boolean> => {
    try {
      const res = await fetch(`${API_BASE}/${id}`, {
        method: "PUT",
        headers: authHeaders(),
        body: JSON.stringify(req),
      });
      if (!res.ok) {
        const data = await res.json();
        throw new Error(data.error || "Failed to update incident");
      }
      await fetchIncidents();
      return true;
    } catch (err) {
      setError(err instanceof Error ? err.message : "Unknown error");
      return false;
    }
  };

  const deleteIncident = async (id: number): Promise<boolean> => {
    try {
      const res = await fetch(`${API_BASE}/${id}`, {
        method: "DELETE",
        headers: { Authorization: `Bearer ${token}` },
      });
      if (!res.ok) {
        const data = await res.json();
        throw new Error(data.error || "Failed to delete incident");
      }
      setIncidents((prev) => prev.filter((inc) => inc.id !== id));
      return true;
    } catch (err) {
      setError(err instanceof Error ? err.message : "Unknown error");
      return false;
    }
  };

  return {
    incidents,
    loading,
    error,
    createIncident,
    updateIncident,
    deleteIncident,
    refetch: fetchIncidents,
  };
}
