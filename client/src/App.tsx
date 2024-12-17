import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import Login from './login';
import './index.css';
import Chat from './chat';

function App() {
  return (
    <Router>
      <div>
        <Routes>
          <Route path="/login" element={<Login />} />
          <Route path="/home" element={<Chat />} />
          <Route path="*" element={<Navigate to="/login" />} />

          {/* <Route path="/app" element={<ChatSidebar />} /> */}
        </Routes>
      </div>
    </Router>
  );
}

export default App;