import { Link } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'

function Navbar() {
  const { token, logout } = useAuth()

  return (
    <nav>
      <Link to="/">Главная</Link>
      {token ? (
        <>
          <Link to="/passes/new">Добавить пропуск</Link>
          <button onClick={logout}>Выйти</button>
        </>
      ) : (
        <>
          <Link to="/login">Вход</Link>
          <Link to="/register">Регистрация</Link>
        </>
      )}
    </nav>
  )
}

export default Navbar