import { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import api from '../api'

function PassForm() {
  const { id } = useParams()
  const navigate = useNavigate()
  const isEditing = Boolean(id)

  const [formData, setFormData] = useState({
    employeeName: '',
    department: '',
    validUntil: ''
  })
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [fetchLoading, setFetchLoading] = useState(isEditing)

  useEffect(() => {
    if (isEditing) {
      const fetchPass = async () => {
        try {
          const res = await api.get(`/passes/${id}`)
          setFormData({
            employeeName: res.data.employeeName,
            department: res.data.department,
            validUntil: res.data.validUntil
          })
        } catch (err) {
          setError('Не удалось загрузить данные пропуска')
        } finally {
          setFetchLoading(false)
        }
      }
      fetchPass()
    }
  }, [id, isEditing])

  const handleChange = (e) => {
    setFormData({ ...formData, [e.target.name]: e.target.value })
  }

  const validate = () => {
    if (!formData.employeeName.trim()) return 'Имя сотрудника обязательно'
    if (!formData.department.trim()) return 'Отдел обязателен'
    if (!formData.validUntil) return 'Дата окончания обязательна'
    return null
  }

  const handleSubmit = async (e) => {
    e.preventDefault()
    const validationError = validate()
    if (validationError) {
      setError(validationError)
      return
    }

    setLoading(true)
    setError('')
    try {
      if (isEditing) {
        await api.put(`/passes/${id}`, formData)
      } else {
        await api.post('/passes', formData)
      }
      navigate('/')
    } catch (err) {
      setError(err.response?.data?.error || 'Ошибка сохранения')
    } finally {
      setLoading(false)
    }
  }

  if (fetchLoading) return <div className="spinner">Загрузка данных...</div>

  return (
    <div className="container">
      <form onSubmit={handleSubmit}>
        <h2>{isEditing ? 'Редактирование пропуска' : 'Добавление пропуска'}</h2>
        {error && <div className="error">{error}</div>}
        <div>
          <label htmlFor="employeeName">ФИО сотрудника</label>
          <input
            id="employeeName"
            type="text"
            name="employeeName"
            value={formData.employeeName}
            onChange={handleChange}
            disabled={loading}
          />
        </div>
        <div>
          <label htmlFor="department">Отдел</label>
          <input
            id="department"
            type="text"
            name="department"
            value={formData.department}
            onChange={handleChange}
            disabled={loading}
          />
        </div>
        <div>
          <label htmlFor="validUntil">Действует до</label>
          <input
            id="validUntil"
            type="date"
            name="validUntil"
            value={formData.validUntil}
            onChange={handleChange}
            disabled={loading}
          />
        </div>
        <div style={{ display: 'flex', gap: '1rem' }}>
          <button type="submit" disabled={loading}>
            {loading ? 'Сохранение...' : 'Сохранить'}
          </button>
          <button type="button" onClick={() => navigate('/')} className="btn-secondary">
            Отмена
          </button>
        </div>
      </form>
    </div>
  )
}

export default PassForm