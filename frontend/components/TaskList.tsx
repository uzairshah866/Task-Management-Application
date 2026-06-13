'use client'

import { Task } from './TaskDashboard'
import { TaskCard } from './TaskCard'

interface TaskListProps {
  tasks: Task[]
  onDelete: (id: string) => void
  onEdit: (task: Task) => void
}

export function TaskList({ tasks, onDelete, onEdit }: TaskListProps) {
  return (
    <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
      {tasks?.map((task) => (
        <TaskCard key={task.id} task={task} onDelete={onDelete} onEdit={onEdit} />
      ))}
    </div>
  )
}
