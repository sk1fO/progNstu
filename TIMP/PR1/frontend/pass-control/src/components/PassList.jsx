import { useState, useEffect, useCallback } from 'react';
import { Link } from 'react-router-dom';
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
  PieChart,
  Pie,
  Cell
} from 'recharts';
import api from '../api';
import Spinner from './Spinner';

function PassList() {
  const [passes, setPasses] = useState([]);
  const [loading, setLoading] = useState(true);
  const [refreshing, setRefreshing] = useState(false);
  const [error, setError] = useState('');
  const [chartData, setChartData] = useState([]);
  const [pieData, setPieData] = useState([]);
  const [lastUpdated, setLastUpdated] = useState(new Date());

  const COLORS = ['#2563eb', '#10b981', '#f59e0b', '#ef4444', '#8b5cf6', '#ec489a'];

  const updateCharts = (data) => {
    // Подсчёт пропусков по отделам для гистограммы
    const deptCount = {};
    data.forEach(pass => {
      const dept = pass.department;
      deptCount[dept] = (deptCount[dept] || 0) + 1;
    });
    const barData = Object.entries(deptCount).map(([name, value]) => ({
      name,
      value
    }));
    setChartData(barData);

    // Данные для круговой диаграммы (топ-5 отделов + "Остальные")
    const sorted = [...barData].sort((a, b) => b.value - a.value);
    const top5 = sorted.slice(0, 5);
    const othersValue = sorted.slice(5).reduce((sum, item) => sum + item.value, 0);
    
    let pieChartData = [...top5];
    if (othersValue > 0) {
      pieChartData.push({ name: 'Остальные', value: othersValue });
    }
    setPieData(pieChartData);
  };

  const fetchPasses = useCallback(async (showRefresh = false) => {
    if (showRefresh) {
      setRefreshing(true);
    } else {
      setLoading(true);
    }
    try {
      const res = await api.get('/passes');
      setPasses(res.data);
      updateCharts(res.data);
      setLastUpdated(new Date());
      setError('');
    } catch (err) {
      setError('Не удалось загрузить список пропусков');
    } finally {
      setLoading(false);
      setRefreshing(false);
    }
  }, []);

  // Первоначальная загрузка
  useEffect(() => {
    fetchPasses();
  }, [fetchPasses]);

  // Обработчики действий (обновляют данные сразу после выполнения)
  const handleDelete = async (id) => {
    if (!window.confirm('Вы уверены, что хотите удалить этот пропуск?')) return;
    try {
      await api.delete(`/passes/${id}`);
      // Мгновенное обновление данных
      await fetchPasses(true);
    } catch (err) {
      alert('Ошибка при удалении');
    }
  };

  const handleRefresh = () => {
    fetchPasses(true);
  };

  if (loading) return <Spinner fullPage />;
  if (error) return <div className="container"><div className="error">{error}</div></div>;

  return (
    <div className="container">
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '1.5rem' }}>
        <h2 style={{ margin: 0 }}>Система контроля пропусков</h2>
        <div style={{ display: 'flex', gap: '1rem', alignItems: 'center' }}>
        </div>
      </div>
      
      {/* Индикатор обновления */}
      {refreshing && (
        <div style={{
          position: 'fixed',
          top: '1rem',
          right: '1rem',
          background: 'var(--surface)',
          padding: '0.5rem 1rem',
          borderRadius: '20px',
          boxShadow: 'var(--shadow-md)',
          fontSize: '0.8rem',
          zIndex: 1000,
          display: 'flex',
          alignItems: 'center',
          gap: '0.5rem'
        }}>
          <div className="spinner-small"></div>
          <span>Обновление данных...</span>
        </div>
      )}

      {/* Блок с графиками */}
      <div style={{ 
        display: 'grid', 
        gridTemplateColumns: 'repeat(auto-fit, minmax(400px, 1fr))', 
        gap: '2rem',
        marginBottom: '2rem'
      }}>
        {/* Гистограмма */}
        <div className="chart-card">
          <h3>Количество пропусков по отделам</h3>
          <ResponsiveContainer width="100%" height={300}>
            <BarChart data={chartData}>
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis dataKey="name" />
              <YAxis />
              <Tooltip />
              <Legend />
              <Bar dataKey="value" fill="#2563eb" name="Количество пропусков" />
            </BarChart>
          </ResponsiveContainer>
        </div>

        {/* Круговая диаграмма */}
        <div className="chart-card">
          <h3>Распределение пропусков по отделам</h3>
          <ResponsiveContainer width="100%" height={300}>
            <PieChart>
              <Pie
                data={pieData}
                cx="50%"
                cy="50%"
                labelLine={false}
                label={({ name, percent }) => `${name}: ${(percent * 100).toFixed(0)}%`}
                outerRadius={80}
                fill="#8884d8"
                dataKey="value"
              >
                {pieData.map((entry, index) => (
                  <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                ))}
              </Pie>
              <Tooltip />
              <Legend />
            </PieChart>
          </ResponsiveContainer>
        </div>
      </div>

      {/* Кнопка добавления и список пропусков */}
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '1rem' }}>
        <h2 style={{ margin: 0 }}>Список пропусков</h2>
        <Link to="/passes/new" className="add-link">+ Добавить пропуск</Link>
      </div>
      
      {passes.length === 0 ? (
        <div className="empty-state">Пока нет ни одного пропуска</div>
      ) : (
        <div className="table-container">
          <table>
            <thead>
              <tr>
                <th>Сотрудник</th>
                <th>Отдел</th>
                <th>Действует до</th>
                <th>Создатель</th>
                <th>Действия</th>
              </tr>
            </thead>
            <tbody>
              {passes.map(p => (
                <tr key={p.id}>
                  <td>{p.employeeName}</td>
                  <td>{p.department}</td>
                  <td>{p.validUntil}</td>
                  <td>{p.creatorName}</td>
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