import { Incident } from '../types/incident'
import { StatusBadge, CategoryBadge } from './StatusBadge'
import { Pencil, Trash2, Clock, User } from 'lucide-react'

interface IncidentCardProps {
  incident:      Incident
  currentUserId: number
  onEdit:        (incident: Incident) => void
  onDelete:      (id: number) => void
}

function formatDate(dateStr: string): string {
  return new Date(dateStr).toLocaleString('en-US', {
    month: 'short', day: 'numeric', year: 'numeric',
    hour: '2-digit', minute: '2-digit',
  })
}

export function IncidentCard({ incident, currentUserId, onEdit, onDelete }: IncidentCardProps) {
  const isOwner = incident.user_id === currentUserId

  const handleDelete = () => {
    if (window.confirm(`Delete "${incident.title}"? This cannot be undone.`)) {
      onDelete(incident.id)
    }
  }

  return (
    <article className="card">
      <div className="card-header">
        <div className="card-badges">
          <CategoryBadge category={incident.category} />
          <StatusBadge status={incident.status} />
        </div>

        {/* Only render action buttons if this user is the owner */}
        {isOwner && (
          <div className="card-actions">
            <button className="icon-btn" onClick={() => onEdit(incident)} title="Edit">
              <Pencil size={15} />
            </button>
            <button className="icon-btn icon-btn-danger" onClick={handleDelete} title="Delete">
              <Trash2 size={15} />
            </button>
          </div>
        )}
      </div>

      <h3 className="card-title">{incident.title}</h3>
      <p className="card-description">{incident.description}</p>

      <div className="card-footer">
        <User size={12} />
        <span>{incident.owner_username}</span>
        <span className="card-footer-sep">·</span>
        <Clock size={12} />
        <span>{formatDate(incident.created_at)}</span>
      </div>
    </article>
  )
}