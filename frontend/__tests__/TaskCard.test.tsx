import { render, screen, fireEvent } from '@testing-library/react'
import { TaskCard } from '@/components/TaskCard'
import { Task } from '@/components/TaskDashboard'

const mockTask: Task = {
  id: '1',
  title: 'Fix login bug',
  description: 'Users cannot log in with special characters',
  status: 'in_progress',
  priority: 'high',
  due_date: '2026-12-31T00:00:00Z',
  created_at: '2026-06-01T00:00:00Z',
  updated_at: '2026-06-01T00:00:00Z',
}

describe('TaskCard', () => {
  const onDelete = jest.fn()
  const onEdit = jest.fn()

  beforeEach(() => {
    jest.clearAllMocks()
  })

  it('renders task title and description', () => {
    render(<TaskCard task={mockTask} onDelete={onDelete} onEdit={onEdit} />)
    expect(screen.getByText('Fix login bug')).toBeInTheDocument()
    expect(screen.getByText('Users cannot log in with special characters')).toBeInTheDocument()
  })

  it('renders status badge', () => {
    render(<TaskCard task={mockTask} onDelete={onDelete} onEdit={onEdit} />)
    expect(screen.getByText('In Progress')).toBeInTheDocument()
  })

  it('renders priority badge', () => {
    render(<TaskCard task={mockTask} onDelete={onDelete} onEdit={onEdit} />)
    expect(screen.getByText('High')).toBeInTheDocument()
  })

  it('renders due date', () => {
    render(<TaskCard task={mockTask} onDelete={onDelete} onEdit={onEdit} />)
    expect(screen.getByText(/Dec/)).toBeInTheDocument()
  })

  it('does not render due date section when null', () => {
    const taskWithoutDue = { ...mockTask, due_date: null }
    render(<TaskCard task={taskWithoutDue} onDelete={onDelete} onEdit={onEdit} />)
    expect(screen.queryByText(/Due:/)).not.toBeInTheDocument()
  })

  it('does not render description when empty', () => {
    const taskWithoutDesc = { ...mockTask, description: '' }
    render(<TaskCard task={taskWithoutDesc} onDelete={onDelete} onEdit={onEdit} />)
    expect(screen.queryByText('Users cannot log in with special characters')).not.toBeInTheDocument()
  })

  it('calls onEdit with task when Edit is clicked', () => {
    render(<TaskCard task={mockTask} onDelete={onDelete} onEdit={onEdit} />)
    fireEvent.click(screen.getByText('Edit'))
    expect(onEdit).toHaveBeenCalledWith(mockTask)
    expect(onEdit).toHaveBeenCalledTimes(1)
  })

  it('calls onDelete with task id when Delete is clicked', () => {
    render(<TaskCard task={mockTask} onDelete={onDelete} onEdit={onEdit} />)
    fireEvent.click(screen.getByText('Delete'))
    expect(onDelete).toHaveBeenCalledWith('1')
    expect(onDelete).toHaveBeenCalledTimes(1)
  })
})
