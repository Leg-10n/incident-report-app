export type Category = "Safety" | "Maintenance";
export type Status = "Open" | "In Progress" | "Success";

export interface Incident {
  id: number;
  title: string;
  description: string;
  category: Category;
  status: Status;
  user_id: number;
  owner_username: string;
  created_at: string;
  updated_at: string;
}

export interface IncidentRequest {
  title: string;
  description: string;
  category: Category;
  status: Status;
}
