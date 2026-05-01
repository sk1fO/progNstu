import { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import api from '../api';

function PassList() {
  const [passes, setPasses] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [columns, setColumns] = useState([]); //

  const fetchPasses = async () => {
    setLoading(true);
    try {
      const res = await api.get('/passes');
      setPasses(res.data);
      if (res.data.length > 0) {
        const excludeKeys = ['id', 'userId'];
        const firstItem = res.data[0];
        const dynamicColumns = Object.keys(firstItem)
          .filter(key => !excludeKeys.includes(key))
          .map(key => ({
            key: key,
            label: getColumnLabel(key) // 
          }));
        setColumns(dynamicColumns);
      }
    } catch (err) {
      setError('Не удалось загрузить список пропусков');
    } finally {
      setLoading(false);
    }
  };

  // 
  const getColumnLabel = (key) => {
    const labels = {
      employeeName: 'Сотрудник',
      department: 'Отдел',
      validUntil: 'Действует до',
      creatorName: 'Создатель'
    };
    return labels[key] || key.charAt(0).toUpperCase() + key.slice(1);
  };

  useEffect(() => {
    fetchPasses();
  }, []);

  const handleDelete = async (id) => {
    if (!window.confirm('Вы уверены, что хотите удалить этот пропуск?')) return;
    try {
      await api.delete(`/passes/${id}`);
      setPasses(passes.filter(p => p.id !== id));
    } catch (err) {
      alert('Ошибка при удалении');
    }
  };

  if (loading) return <div className="spinner">Загрузка...</div>;
  if (error) return <div className="container"><div className="error">{error}</div></div>;

  return (
    <div className="container">
      <h2>Список пропусков</h2>
      <Link to="/passes/new" className="add-link">+ Добавить пропуск</Link>
      {passes.length === 0 ? (
        <div className="empty-state">Пока нет ни одного пропуска</div>
      ) : (
        <div className="table-container">
          <table>
            <thead>
              <tr>
                {columns.map(col => (
                  <th key={col.key}>{col.label}</th>
                ))}
                <th>Действия</th>
              </tr>
            </thead>
            <tbody>
              {passes.map(p => (
                <tr key={p.id}>
                  {columns.map(col => (
                    <td key={col.key}>{p[col.key]}</td>
                  ))}
                  <td>
                    <div className="actions">
                      <Link to={`/passes/${p.id}`}>Просмотр</Link>
                      <Link to={`/passes/${p.id}/edit`}>Редактировать</Link>
                      <button onClick={() => handleDelete(p.id)}>Удалить</button>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}

export default PassList;