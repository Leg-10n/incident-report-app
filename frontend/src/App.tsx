import { useState } from "react";
import { useIncidents } from "./hooks/useIncidents";
import { useAuth } from "./hooks/useAuth";
import { IncidentList } from "./components/IncidentList";
import { IncidentForm } from "./components/IncidentForm";
import { AuthForm } from "./components/AuthForm";
import { Incident, IncidentRequest } from "./types/incident";
import { Plus, RefreshCw, AlertTriangle, LogOut } from "lucide-react";

export default function App() {
  const {
    token,
    user,
    error: authError,
    loading: authLoading,
    login,
    register,
    logout,
  } = useAuth();
  const {
    incidents,
    loading,
    error,
    createIncident,
    updateIncident,
    deleteIncident,
    refetch,
  } = useIncidents(token);

  const [showForm, setShowForm] = useState(false);
  const [editingIncident, setEditingIncident] = useState<Incident | null>(null);

  // Auth gate — show login/register screen if no valid session
  if (!token || !user) {
    return (
      <AuthForm
        onLogin={login}
        onRegister={register}
        error={authError}
        loading={authLoading}
      />
    );
  }

  const handleEdit = (incident: Incident) => {
    setEditingIncident(incident);
    setShowForm(true);
  };

  const handleUpdate = async (data: IncidentRequest) => {
    if (!editingIncident) return false;
    return updateIncident(editingIncident.id, data);
  };

  const handleCloseForm = () => {
    setShowForm(false);
    setEditingIncident(null);
  };

  return (
    <div className="app">
      <header className="header">
        <div className="header-inner">
          <div className="header-title">
            <div className="logo-mark">IR</div>
            <div>
              <h1>Incident Reports</h1>
              <p className="subtitle">Track, manage, and resolve incidents</p>
            </div>
          </div>
          <div className="header-actions">
            <span className="username-display">@{user.username}</span>
            <button className="btn btn-ghost" onClick={refetch} title="Refresh">
              <RefreshCw size={16} />
            </button>
            <button
              className="btn btn-primary"
              onClick={() => {
                setEditingIncident(null);
                setShowForm(true);
              }}
            >
              <Plus size={16} />
              New Report
            </button>
            <button className="icon-btn" onClick={logout} title="Logout">
              <LogOut size={16} />
            </button>
          </div>
        </div>
      </header>

      <main className="main">
        {error && (
          <div className="error-banner">
            <AlertTriangle size={16} />
            <span>{error}</span>
          </div>
        )}

        {loading ? (
          <div className="loading">
            <div className="spinner" />
            <p>Loading incidents...</p>
          </div>
        ) : (
          <IncidentList
            incidents={incidents}
            currentUserId={user.id}
            onEdit={handleEdit}
            onDelete={deleteIncident}
          />
        )}
      </main>

      {showForm && (
        <IncidentForm
          incident={editingIncident}
          onSubmit={editingIncident ? handleUpdate : createIncident}
          onClose={handleCloseForm}
        />
      )}
    </div>
  );
}
