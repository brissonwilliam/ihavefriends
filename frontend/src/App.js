import './app.css';
import { BrowserRouter as Router, Routes, Route, Navigate } from "react-router-dom";
import Dashboard from './components/dashboard';
import Login, {IsJWTValid} from './components/login';
import { useCookies } from 'react-cookie';
import { Outlet } from 'react-router';

function App() {
  const [cookies] = useCookies(['jwt', 'jwtExpiration']);

  function RequireAuth() {
    if (!IsJWTValid(cookies.jwt, cookies.jwtExpiration)) {
      return (<Navigate replace to="/login"/>);
    }
    return <Outlet />;
  };

  return (
    <div className="App">
      <header className="App-header">
        
      </header>

      <Router>
        <Routes>
            <Route path="/login" element={<Login/>}/>

            <Route element={<RequireAuth/>}>
              <Route path="/dashboard" element={<Dashboard/>}/>
              <Route path="/" element={<Navigate replace to="/dashboard"/>}/>
            </Route>
        </Routes>
      </Router>
    </div>
  );
}


export default App;
