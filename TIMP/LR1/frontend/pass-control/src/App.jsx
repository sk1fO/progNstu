import { Routes, Route, Navigate } from 'react-router-dom'
import { useAuth } from './context/AuthContext'
import Navbar from './components/Navbar'
import Login from './components/Login'
import Register from './components/Register'
import PassList from './components/PassList'
import PassDetail from './components/PassDetail'
import PassForm from './components/PassForm'

function App() {
  const { token } = useAuth()

  return (
    <>
      <Navbar />
      <div className="container">
        <Routes>
          <Route path="/login" element={!token ? <Login /> : <Navigate to="/" />} />
          <Route path="/register" element={!token ? <Register /> : <Navigate to="/" />} />
          <Route path="/" element={token ? <PassList /> : <Navigate to="/login" />} />
          <Route path="/passes/:id" element={token ? <PassDetail /> : <Navigate to="/login" />} />
          <Route path="/passes/new" element={token ? <PassForm /> : <Navigate to="/login" />} />
          <Route path="/passes/:id/edit" element={token ? <PassForm /> : <Navigate to="/login" />} />
        </Routes>
      </div>
    </>
  )
}

export default App