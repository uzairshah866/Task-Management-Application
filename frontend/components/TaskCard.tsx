'use client'

import { Buttons } from '@/app/ui/Buttons'
import { Task } from './TaskDashboard'

interface TaskCardProps {
  task: Task
  onDelete: (id: string) => void
  onEdit: (task: Task) => void
}

export function TaskCard({ task, onDelete, onEdit }: TaskCardProps) {
  const getStatusColor = (status: string) => {
    switch (status) {
      case 'completed':
        return 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200'
      case 'in_progress':
        return 'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200'
      case 'pending':
        return 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200'
      default:
        return 'bg-gray-100 text-gray-800'
    }
  }

  const getPriorityColor = (priority: string) => {
    switch (priority) {
      case 'high':
        return 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200'
      case 'medium':
        return 'bg-orange-100 text-orange-800 dark:bg-orange-900 dark:text-orange-200'
      case 'low':
        return 'bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-200'
      default:
        return 'bg-gray-100 text-gray-800'
    }
  }

  const formatDate = (dateString: string | null) => {
    if (!dateString) return 'No due date'
    return new Date(dateString).toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      year: 'numeric',
    })
  }

  return (
    <div className="rounded-lg border border-gray-200 bg-white p-6 shadow-sm dark:border-gray-700 dark:bg-gray-800">
      <div className="mb-4 flex items-start justify-between">
        <div className="flex-1">
          <h3 className="text-lg font-semibold text-gray-900 dark:text-white">{task.title}</h3>
          {task.description && (
            <p className="mt-1 text-sm text-gray-600 dark:text-gray-400">{task.description}</p>
          )}
        </div>
      </div>

      <div className="mb-4 flex flex-wrap gap-2">
        <span
          className={`rounded-full px-3 py-1 text-xs font-medium ${getStatusColor(task.status)}`}
        >
          {task.status
            .split('_')
            .map((w) => w.charAt(0).toUpperCase() + w.slice(1))
            .join(' ')}
        </span>
        <span
          className={`rounded-full px-3 py-1 text-xs font-medium ${getPriorityColor(task.priority)}`}
        >
          {task.priority.charAt(0).toUpperCase() + task.priority.slice(1)}
        </span>
      </div>

      {task.due_date && (
        <p className="mb-4 text-sm text-gray-600 dark:text-gray-400">
          Due: {formatDate(task.due_date)}
        </p>
      )}

      <div className="flex gap-2">
        <Buttons variant="primary" className="flex-1" onClick={() => onEdit(task)}>
          Edit
        </Buttons>

        <Buttons variant="danger" className="flex-1" onClick={() => onDelete(task.id)}>
          Delete
        </Buttons>
      </div>
    </div>
  )
}
