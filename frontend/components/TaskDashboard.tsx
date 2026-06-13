'use client'

import { useEffect, useState, useCallback, useRef } from 'react'
import { useRouter } from 'next/navigation'
import { useAuthStore } from '@/lib/store'
import { taskApi } from '@/lib/api'
import { TaskList } from './TaskList'
import { TaskForm } from './TaskForm'
import { ThemeToggle } from './ThemeToggle'
import { Buttons } from '@/app/ui/Buttons'

export interface Task {
  id: string
  user_id: string
  title: string
  description: string
  status: 'pending' | 'in_progress' | 'completed'
  priority: 'low' | 'medium' | 'high'
  due_date: string | null
  created_at: string
  updated_at: string
}

interface User {
  id: string
  email: string
}

interface Filters {
  status: string
  priority: string
  search: string
  sortBy: string
  sortOrder: string
  page: number
  pageSize: number
}

export function TaskDashboard({ user }: { user: User }) {
  const router = useRouter()
  const clearAuth = useAuthStore((state) => state.clearAuth)
  const [tasks, setTasks] = useState<Task[]>([])
  const [total, setTotal] = useState(0)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [showForm, setShowForm] = useState(false)
  const [editingTask, setEditingTask] = useState<Task | null>(null)
  const [searchInput, setSearchInput] = useState('')
  const searchTimer = useRef<ReturnType<typeof setTimeout> | null>(null)

  const [filters, setFilters] = useState<Filters>({
    status: '',
    priority: '',
    search: '',
    sortBy: 'created_at',
    sortOrder: 'DESC',
    page: 1,
    pageSize: 10,
  })

  const loadTasks = useCallback(async (f: Filters) => {
    setLoading(true)
    setError('')
    try {
      const response = await taskApi.list(f)
      setTasks(response.data.tasks || [])
      setTotal(response.data.total)
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to load tasks')
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    loadTasks(filters)
  }, [filters, loadTasks])

  const handleSearchChange = (value: string) => {
    setSearchInput(value)
    if (searchTimer.current) clearTimeout(searchTimer.current)
    searchTimer.current = setTimeout(() => {
      setFilters((prev) => ({ ...prev, search: value, page: 1 }))
    }, 300)
  }

  const handleLogout = () => {
    clearAuth()
    router.push('/login')
  }

  const handleSaveTask = async (taskData: any) => {
    if (editingTask) {
      const response = await taskApi.update(editingTask.id, taskData)
      setTasks((prev) => prev.map((t) => (t.id === editingTask.id ? response.data : t)))
      setEditingTask(null)
    } else {
      const response = await taskApi.create(taskData)
      setTasks((prev) => [response.data, ...prev])
      setTotal((prev) => prev + 1)
    }
    setShowForm(false)
  }

  const handleDeleteTask = async (taskId: string) => {
    await taskApi.delete(taskId)
    setTasks((prev) => prev.filter((t) => t.id !== taskId))
    setTotal((prev) => prev - 1)
  }

  const handleEditTask = (task: Task) => {
    setEditingTask(task)
    setShowForm(true)
  }

  const totalPages = Math.ceil(total / filters.pageSize)

  return (
    <div className="min-h-screen pb-8">
      <header className="bg-white shadow-sm dark:bg-gray-900">
        <div className="mx-auto flex max-w-6xl items-center justify-between px-4 py-6 sm:px-6 lg:px-8">
          <div>
            <h1 className="text-3xl font-bold text-gray-900 dark:text-white">Tasks</h1>
            <p className="mt-1 text-sm text-gray-600 dark:text-gray-400">Welcome, {user.email}</p>
          </div>
          <div className="flex items-center gap-4">
            <ThemeToggle />
            <Buttons variant="danger" className="flex-1" onClick={handleLogout}>
              Logout
            </Buttons>
          </div>
        </div>
      </header>

      <main className="mx-auto max-w-6xl px-4 py-8 sm:px-6 lg:px-8">
        <div className="mb-8">
          <Buttons
            variant="primary"
            className="flex-1"
            onClick={() => {
              setEditingTask(null)
              setShowForm(true)
            }}
          >
            + New Task
          </Buttons>
        </div>

        {showForm && (
          <TaskForm
            task={editingTask}
            onSave={handleSaveTask}
            onCancel={() => {
              setShowForm(false)
              setEditingTask(null)
            }}
          />
        )}

        <div className="mb-6 grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">
              Status
            </label>
            <select
              value={filters.status}
              onChange={(e) => setFilters((prev) => ({ ...prev, status: e.target.value, page: 1 }))}
              className="mt-1 w-full rounded-lg border border-gray-300 px-3 py-2 dark:border-gray-600 dark:bg-gray-800 dark:text-white"
            >
              <option value="">All Status</option>
              <option value="pending">Pending</option>
              <option value="in_progress">In Progress</option>
              <option value="completed">Completed</option>
            </select>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">
              Priority
            </label>
            <select
              value={filters.priority}
              onChange={(e) =>
                setFilters((prev) => ({ ...prev, priority: e.target.value, page: 1 }))
              }
              className="mt-1 w-full rounded-lg border border-gray-300 px-3 py-2 dark:border-gray-600 dark:bg-gray-800 dark:text-white"
            >
              <option value="">All Priorities</option>
              <option value="low">Low</option>
              <option value="medium">Medium</option>
              <option value="high">High</option>
            </select>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">
              Sort By
            </label>
            <select
              value={filters.sortBy}
              onChange={(e) => setFilters((prev) => ({ ...prev, sortBy: e.target.value }))}
              className="mt-1 w-full rounded-lg border border-gray-300 px-3 py-2 dark:border-gray-600 dark:bg-gray-800 dark:text-white"
            >
              <option value="created_at">Created Date</option>
              <option value="due_date">Due Date</option>
              <option value="priority">Priority</option>
            </select>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">
              Search
            </label>
            <input
              type="text"
              value={searchInput}
              onChange={(e) => handleSearchChange(e.target.value)}
              placeholder="Search tasks..."
              className="mt-1 w-full rounded-lg border border-gray-300 px-3 py-2 dark:border-gray-600 dark:bg-gray-800 dark:text-white"
            />
          </div>
        </div>

        {loading ? (
          <div className="py-12 text-center">
            <p className="text-gray-500 dark:text-gray-400">Loading tasks...</p>
          </div>
        ) : error ? (
          <div className="py-12 text-center">
            <p className="text-red-500">{error}</p>
          </div>
        ) : tasks.length === 0 ? (
          <div className="py-12 text-center">
            <p className="text-gray-500 dark:text-gray-400">
              No tasks found. Create one to get started!
            </p>
          </div>
        ) : (
          <>
            <TaskList tasks={tasks} onDelete={handleDeleteTask} onEdit={handleEditTask} />

            {totalPages > 1 && (
              <div className="mt-6 flex items-center justify-center gap-2">
                <button
                  onClick={() => setFilters((prev) => ({ ...prev, page: prev.page - 1 }))}
                  disabled={filters.page === 1}
                  className="cursor-pointer rounded-lg border border-gray-300 px-3 py-2 disabled:opacity-50 dark:border-gray-600"
                >
                  Previous
                </button>
                <span className="text-sm text-gray-600 dark:text-gray-400">
                  Page {filters.page} of {totalPages}
                </span>
                <button
                  onClick={() => setFilters((prev) => ({ ...prev, page: prev.page + 1 }))}
                  disabled={filters.page >= totalPages}
                  className="cursor-pointer rounded-lg border border-gray-300 px-3 py-2 disabled:opacity-50 dark:border-gray-600"
                >
                  Next
                </button>
              </div>
            )}
          </>
        )}
      </main>
    </div>
  )
}
