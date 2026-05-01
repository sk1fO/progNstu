import { useState, useEffect } from 'react'
import { useParams, Link, useNavigate } from 'react-router-dom'
import api from '../api'

function PassDetail() {
  const { id } = useParams()
  const navigate = useNavigate()
  const [pass, setPass] = useState(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    const fetchPass = async () => {
      try {
        const res = await api.get(`/passes/${id}`)
        setPass(res.data)
      } catch (err) {
        setError('Пропуск не найден')
      } finally {
        setLoading(false)
      }
    }
    fetchPass()
  }, [id])

  const handleDelete = async () => {
    if (!window.confirm('Удалить пропуск?')) return
    try {
      await api.delete(`/passes/${id}`)
      navigate('/')
    } catch (err) {
      alert('Ошибка удаления')
    }
  }

  if (loading) return <div className="spinner">Загрузка...</div>
  if (error) return <div className="container"><div className="error">{error}</div></div>

  return (
    <div className="container">
      <div className="detail-card">
        <h2>Детали пропуска</h2>
        <p><strong>Сотрудник:</strong> {pass.employeeName}</p>
        <p><strong>Отдел:</strong> {pass.department}</p>
        <p><strong>Действует до:</strong> {pass.validUntil}</p>
        <p><strong>Создатель:</strong> {pass.creatorName}</p>
        <div className="detail-actions">
          <Link to={`/passes/${id}/edit`} className="btn">Редактировать</Link>
          <button onClick={handleDelete} className="btn btn-secondary">Удалить</button>
          <Link to="/" className="btn btn-secondary">Назад к списку</Link>
        </div>
      </div>
    </div>
  )
}

export default PassDetail