import { Routes, Route } from 'react-router-dom';
import AddSet from './pages/AddSet';
import Library from './pages/Library';
import SetDetail from './pages/SetDetail';
import AddCategory from './pages/AddCategory';

function App() {
  return (
    <div>
      <Routes>
        <Route path="/" element={<Library />} />
        <Route path="/add" element={<AddSet/>} />
        <Route path="/sets/:id" element={<SetDetail />} />
        <Route path="/categories" element={<AddCategory />} />
      </Routes>
    </div>
  );
}

export default App;
