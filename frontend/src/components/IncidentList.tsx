import { useState } from "react";
import { Incident, Status, Category } from "../types/incident";
import { IncidentCard } from "./IncidentCard";
import { Filter } from "lucide-react";

interface IncidentListProps {
  incidents: Incident[];
  currentUserId: number;
  onEdit: (incident: Incident) => void;
  onDelete: (id: number) => void;
}

export function IncidentList({
  incidents,
  currentUserId,
  onEdit,
  onDelete,
}: IncidentListProps) {
  const [filterStatus, setFilterStatus] = useState<Status | "All">("All");
  const [filterCategory, setFilterCategory] = useState<Category | "All">("All");

  const filtered = incidents.filter(
    (inc) =>
      (filterStatus === "All" || inc.status === filterStatus) &&
      (filterCategory === "All" || inc.category === filterCategory),
  );

  if (incidents.length === 0) {
    return (
      <div className="empty-state">
        <p>No incidents reported yet.</p>
        <p className="empty-sub">
          Click <strong>New Report</strong> to get started.
        </p>
      </div>
    );
  }

  return (
    <div>
      <div className="filter-bar">
        <Filter size={14} />
        <span>Filter:</span>

        <select
          value={filterStatus}
          onChange={(e) => setFilterStatus(e.target.value as Status | "All")}
        >
          <option value="All">All Statuses</option>
          <option value="Open">Open</option>
          <option value="In Progress">In Progress</option>
          <option value="Success">Success</option>
        </select>

        <select
          value={filterCategory}
          onChange={(e) =>
            setFilterCategory(e.target.value as Category | "All")
          }
        >
          <option value="All">All Categories</option>
          <option value="Safety">Safety</option>
          <option value="Maintenance">Maintenance</option>
        </select>

        <span className="filter-count">
          {filtered.length} of {incidents.length}
        </span>
      </div>

      {filtered.length === 0 ? (
        <div className="empty-state">
          <p>No incidents match the current filters.</p>
        </div>
      ) : (
        <div className="grid">
          {filtered.map((inc) => (
            <IncidentCard
              key={inc.id}
              incident={inc}
              currentUserId={currentUserId}
              onEdit={onEdit}
              onDelete={onDelete}
            />
          ))}
        </div>
      )}
    </div>
  );
}
