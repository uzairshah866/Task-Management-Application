import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import { TaskForm } from '@/components/TaskForm'
import { Task } from '@/components/TaskDashboard'

const mockTask: Task = {
  id: '1',
  title: 'Existing task',
  description: 'Existing description',
  status: 'in_progress',
  priority: 'high',
  due_date: '2026-12-31T00:00:00Z',
  created_at: '2026-06-01T00:00:00Z',
  updated_at: '2026-06-01T00:00:00Z',
}

describe('TaskForm', () => {
  const onSave = jest.fn()
  const onCancel = jest.fn()

  beforeEach(() => {
    jest.clearAllMocks()
  })

  it('renders "Create New Task" heading when no task prop', () => {
    render(<TaskForm onSave={onSave} onCancel={onCancel} />)
    expect(screen.getByText('Create New Task')).toBeInTheDocument()
  })

  it('renders "Edit Task" heading when task prop provided', () => {
    render(<TaskForm task={mockTask} onSave={onSave} onCancel={onCancel} />)
    expect(screen.getByText('Edit Task')).toBeInTheDocument()
  })

  it('renders all form fields', () => {
    render(<TaskForm onSave={onSave} onCancel={onCancel} />)
    expect(screen.getByPlaceholderText('Enter task title')).toBeInTheDocument()
    expect(screen.getByPlaceholderText('Enter task description')).toBeInTheDocument()
    expect(screen.getByText('Save Task')).toBeInTheDocument()
    expect(screen.getByText('Cancel')).toBeInTheDocument()
  })

  it('pre-fills fields when editing an existing task', () => {
    render(<TaskForm task={mockTask} onSave={onSave} onCancel={onCancel} />)
    expect(screen.getByDisplayValue('Existing task')).toBeInTheDocument()
    expect(screen.getByDisplayValue('Existing description')).toBeInTheDocument()
    expect(screen.getByDisplayValue('In Progress')).toBeInTheDocument()
    expect(screen.getByDisplayValue('High')).toBeInTheDocument()
  })

  it('defaults to pending status and medium priority for new task', () => {
    render(<TaskForm onSave={onSave} onCancel={onCancel} />)
    expect(screen.getByDisplayValue('Pending')).toBeInTheDocument()
    expect(screen.getByDisplayValue('Medium')).toBeInTheDocument()
  })

  it('shows validation error when submitting with empty title', async () => {
    const { container } = render(<TaskForm onSave={onSave} onCancel={onCancel} />)
    fireEvent.submit(container.querySelector('form')!)
    await waitFor(() => {
      expect(screen.getByText('Title is required')).toBeInTheDocument()
    })
    expect(onSave).not.toHaveBeenCalled()
  })

  it('shows validation error when submitting whitespace-only title', async () => {
    const { container } = render(<TaskForm onSave={onSave} onCancel={onCancel} />)
    fireEvent.change(screen.getByPlaceholderText('Enter task title'), {
      target: { value: '   ' },
    })
    fireEvent.submit(container.querySelector('form')!)
    await waitFor(() => {
      expect(screen.getByText('Title is required')).toBeInTheDocument()
    })
    expect(onSave).not.toHaveBeenCalled()
  })

  it('calls onSave with trimmed form data when valid', async () => {
    onSave.mockResolvedValueOnce(undefined)
    const { container } = render(<TaskForm onSave={onSave} onCancel={onCancel} />)
    fireEvent.change(screen.getByPlaceholderText('Enter task title'), {
      target: { value: '  New task  ' },
    })
    fireEvent.submit(container.querySelector('form')!)
    await waitFor(() => {
      expect(onSave).toHaveBeenCalledWith(
        expect.objectContaining({ title: 'New task' })
      )
    })
  })

  it('calls onCancel when Cancel is clicked', () => {
    render(<TaskForm onSave={onSave} onCancel={onCancel} />)
    fireEvent.click(screen.getByText('Cancel'))
    expect(onCancel).toHaveBeenCalledTimes(1)
  })

  it('shows error message when onSave rejects', async () => {
    onSave.mockRejectedValueOnce({ message: 'Server error' })
    render(<TaskForm onSave={onSave} onCancel={onCancel} />)
    fireEvent.change(screen.getByPlaceholderText('Enter task title'), {
      target: { value: 'Valid title' },
    })
    fireEvent.click(screen.getByText('Save Task'))
    await waitFor(() => {
      expect(screen.getByText('Server error')).toBeInTheDocument()
    })
  })
})
