import { Suspense, lazy } from 'react';
import { Routes, Route, Navigate } from 'react-router-dom';
import { useAuth } from './context/AuthContext';
import Navbar from './components/Navbar';
import Spinner from './components/Spinner';

// Временно для тестирования спиннера
const delay = (ms) => new Promise(resolve => setTimeout(resolve, ms));

const Login = lazy(() => delay(2000).then(() => import('./components/Login')));
const Register = lazy(() => delay(2000).then(() => import('./components/Register')));
const PassList = lazy(() => delay(2000).then(() => import('./components/PassList')));
const PassDetail = lazy(() => delay(2000).then(() => import('./components/PassDetail')));
const PassForm = lazy(() => delay(2000).then(() => import('./components/PassForm')));

function App() {
  const { token } = useAuth();

  return (
    <>
      <Navbar />
      <div className="container">
        <Suspense fallback={<Spinner fullPage />}>
          <Routes>
            <Route path="/login" element={!token ? <Login /> : <Navigate to="/" />} />
            <Route path="/register" element={!token ? <Register /> : <Navigate to="/" />} />
            <Route path="/" element={token ? <PassList /> : <Navigate to="/login" />} />
            <Route path="/passes/:id" element={token ? <PassDetail /> : <Navigate to="/login" />} />
            <Route path="/passes/new" element={token ? <PassForm /> : <Navigate to="/login" />} />
            <Route path="/passes/:id/edit" element={token ? <PassForm /> : <Navigate to="/login" />} />
          </Routes>
        </Suspense>
      </div>
    </>
  );
}

export default App;