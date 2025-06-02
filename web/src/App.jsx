import React from 'react';
import './App.css';
import { TaskProvider } from './context/TaskProvider';
import Home from './pages/Home/Home';
import {
  BrowserRouter as Router,
  Routes,
  Route,
  Navigate,
} from 'react-router-dom';
import OtherPage from './pages/OtherPage/OtherPage';

function App() {
  return (
    <Router>
      <Routes>
        <Route
          path="/"
          element={
            <TaskProvider>
              <Home />
            </TaskProvider>
          }
        ></Route>
        <Route path="/other" element={<OtherPage />}></Route>
      </Routes>
    </Router>
  );
}

export default App;
