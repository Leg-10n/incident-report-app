import { Status, Category } from '../types/incident'

interface StatusBadgeProps {
  status: Status
}

interface CategoryBadgeProps {
  category: Category
}

export function StatusBadge({ status }: StatusBadgeProps) {
  const styles: Record<Status, string> = {
    'Open': 'badge badge-open',
    'In Progress': 'badge badge-progress',
    'Success': 'badge badge-success',
  }
  return <span className={styles[status]}>{status}</span>
}

export function CategoryBadge({ category }: CategoryBadgeProps) {
  const styles: Record<Category, string> = {
    'Safety': 'badge badge-safety',
    'Maintenance': 'badge badge-maintenance',
  }
  return <span className={styles[category]}>{category}</span>
}