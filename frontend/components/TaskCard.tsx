'use client'

import { EditIcon } from '@/app/icons/Edit'
import { Task } from './TaskDashboard'
import { TrashIcon } from '@/app/icons/Trash'
import { CalendarIcon } from '@/app/icons/Calendar'
import { FlameIcon } from '@/app/icons/Flame'
import { CheckCircleIcon } from '@/app/icons/CheckCircle'
import { PulseIcon } from '@/app/icons/Pulse'
import { ClockIcon } from '@/app/icons/Clock'
import { AlertIcon } from '@/app/icons/Alert'
import { LeafIcon } from '@/app/icons/Leaf'

interface TaskCardProps {
  task: Task
  onDelete: (id: string) => void
  onEdit: (task: Task) => void
}

export function TaskCard({ task, onDelete, onEdit }: TaskCardProps) {
  const priorityBarColor =
    {
      high: 'bg-red-500',
      medium: 'bg-amber-400',
      low: 'bg-emerald-500',
    }[task.priority] ?? 'bg-gray-300'

  const statusBadge = {
    completed: {
      className: 'bg-green-100 text-green-800 dark:bg-green-900/50 dark:text-green-300',
      icon: <CheckCircleIcon />,
      label: 'Completed',
    },
    in_progress: {
      className: 'bg-blue-100 text-blue-800 dark:bg-blue-900/50 dark:text-blue-300',
      icon: <PulseIcon />,
      label: 'In Progress',
    },
    pending: {
      className: 'bg-amber-100 text-amber-800 dark:bg-amber-900/50 dark:text-amber-300',
      icon: <ClockIcon />,
      label: 'Pending',
    },
  }[task.status] ?? {
    className: 'bg-gray-100 text-gray-700 dark:bg-gray-700 dark:text-gray-300',
    icon: null,
    label: task.status,
  }

  const priorityBadge = {
    high: {
      className: 'bg-red-100 text-red-800 dark:bg-red-900/50 dark:text-red-300',
      icon: <FlameIcon />,
      label: 'High',
    },
    medium: {
      className: 'bg-amber-100 text-amber-800 dark:bg-amber-900/50 dark:text-amber-300',
      icon: <AlertIcon />,
      label: 'Medium',
    },
    low: {
      className: 'bg-emerald-100 text-emerald-800 dark:bg-emerald-900/50 dark:text-emerald-300',
      icon: <LeafIcon />,
      label: 'Low',
    },
  }[task.priority] ?? {
    className: 'bg-gray-100 text-gray-700 dark:bg-gray-700 dark:text-gray-300',
    icon: null,
    label: task.priority,
  }

  const formatDate = (dateString: string | null) => {
    if (!dateString) return null
    return new Date(dateString).toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      year: 'numeric',
    })
  }

  const isOverdue =
    task.due_date && task.status !== 'completed' && new Date(task.due_date) < new Date()

  const formattedDate = formatDate(task.due_date)

  return (
    <div className="group relative overflow-hidden rounded-xl border border-gray-200 bg-white transition-all duration-200 hover:-translate-y-0.5 hover:border-gray-300 hover:shadow-md dark:border-gray-700 dark:bg-gray-800 dark:hover:border-gray-600">
      {/* Priority accent bar */}
      <div className={`h-[3px] w-full ${priorityBarColor}`} />

      <div className="p-5">
        {/* Title & description */}
        <div className="mb-4">
          <h3 className="text-[15px] leading-snug font-medium text-gray-900 dark:text-white">
            {task.title}
          </h3>
          {task.description && (
            <p className="mt-1.5 text-[13px] leading-relaxed text-gray-500 dark:text-gray-400">
              {task.description}
            </p>
          )}
        </div>

        {/* Badges */}
        <div className="mb-4 flex flex-wrap gap-2">
          <span
            className={`inline-flex items-center gap-1.5 rounded-full px-2.5 py-1 text-[11px] font-medium tracking-wide ${statusBadge.className}`}
          >
            {statusBadge.icon}
            {statusBadge.label}
          </span>
          <span
            className={`inline-flex items-center gap-1.5 rounded-full px-2.5 py-1 text-[11px] font-medium tracking-wide ${priorityBadge.className}`}
          >
            {priorityBadge.icon}
            {priorityBadge.label}
          </span>
        </div>

        {/* Due date */}
        {formattedDate && (
          <div
            className={`mb-4 flex items-center gap-1.5 text-[12px] ${
              isOverdue ? 'text-red-600 dark:text-red-400' : 'text-gray-500 dark:text-gray-400'
            }`}
          >
            <CalendarIcon />
            <span>
              {formattedDate}
              {isOverdue && ' · Overdue'}
            </span>
          </div>
        )}

        {/* Actions */}
        <div className="flex gap-2">
          <button
            onClick={() => onEdit(task)}
            className="flex flex-1 items-center justify-center gap-1.5 rounded-lg border border-gray-200 bg-gray-50 px-3 py-2 text-[13px] font-medium text-gray-700 transition-all duration-150 hover:border-gray-300 hover:bg-gray-100 active:scale-[0.97] dark:border-gray-600 dark:bg-gray-700 dark:text-gray-200 dark:hover:border-gray-500 dark:hover:bg-gray-600"
          >
            <EditIcon />
            Edit
          </button>

          <button
            onClick={() => onDelete(task.id)}
            className="flex flex-1 items-center justify-center gap-1.5 rounded-lg border border-red-100 bg-red-50 px-3 py-2 text-[13px] font-medium text-red-700 transition-all duration-150 hover:border-red-200 hover:bg-red-100 active:scale-[0.97] dark:border-red-900/50 dark:bg-red-900/20 dark:text-red-400 dark:hover:border-red-800 dark:hover:bg-red-900/30"
          >
            <TrashIcon />
            Delete
          </button>
        </div>
      </div>
    </div>
  )
}
