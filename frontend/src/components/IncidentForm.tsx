import { useState, useEffect } from 'react'
import { Incident, IncidentRequest, Category, Status } from '../types/incident'
import { X } from 'lucide-react'

interface IncidentFormProps {
  incident?: Incident | null
  onSubmit: (data: IncidentRequest) => Promise<boolean>
  onClose: () => void
}

const defaultForm: IncidentRequest = {
  title: '',
  description: '',
  category: 'Safety',
  status: 'Open',
}

export function IncidentForm({ incident, onSubmit, onClose }: IncidentFormProps) {
  const [form, setForm] = useState<IncidentRequest>(defaultForm)
  const [submitting, setSubmitting] = useState(false)
  const [fieldError, setFieldError] = useState<string | null>(null)

  useEffect(() => {
    if (incident) {
      setForm({
        title: incident.title,
        description: incident.description,
        category: incident.category,
        status: incident.status,
      })
    } else {
      setForm(defaultForm)
    }
  }, [incident])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!form.title.trim() || !form.description.trim()) {
      setFieldError('Title and description are required.')
      return
    }
    setFieldError(null)
    setSubmitting(true)
    const success = await onSubmit(form)
    setSubmitting(false)
    if (success) onClose()
  }

  const isEdit = Boolean(incident)

  return (
    <div className="modal-backdrop" onClick={onClose}>
      <div className="modal" onClick={e => e.stopPropagation()}>
        <div className="modal-header">
          <h2>{isEdit ? 'Edit Incident' : 'New Incident Report'}</h2>
          <button className="icon-btn" onClick={onClose}><X size={20} /></button>
        </div>

        <form onSubmit={handleSubmit} className="form">
          {fieldError && <p className="form-error">{fieldError}</p>}

          <div className="form-group">
            <label htmlFor="title">Title</label>
            <input
              id="title"
              type="text"
              value={form.title}
              onChange={e => setForm(f => ({ ...f, title: e.target.value }))}
              placeholder="Brief incident title..."
              maxLength={150}
            />
          </div>

          <div className="form-group">
            <label htmlFor="description">Description</label>
            <textarea
              id="description"
              value={form.description}
              onChange={e => setForm(f => ({ ...f, description: e.target.value }))}
              placeholder="Describe the incident in detail..."
              rows={4}
            />
          </div>

          <div className="form-row">
            <div className="form-group">
              <label htmlFor="category">Category</label>
              <select
                id="category"
                value={form.category}
                onChange={e => setForm(f => ({ ...f, category: e.target.value as Category }))}
              >
                <option value="Safety">Safety</option>
                <option value="Maintenance">Maintenance</option>
              </select>
            </div>

            <div className="form-group">
              <label htmlFor="status">Status</label>
              <select
                id="status"
                value={form.status}
                onChange={e => setForm(f => ({ ...f, status: e.target.value as Status }))}
              >
                <option value="Open">Open</option>
                <option value="In Progress">In Progress</option>
                <option value="Success">Success</option>
              </select>
            </div>
          </div>

          <div className="form-actions">
            <button type="button" className="btn btn-ghost" onClick={onClose}>
              Cancel
            </button>
            <button type="submit" className="btn btn-primary" disabled={submitting}>
              {submitting ? 'Saving...' : isEdit ? 'Save Changes' : 'Create Report'}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}