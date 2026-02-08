import { BrowserRouter, Routes, Route } from "react-router-dom";
import SheetsPage from "./pages/sheetsPage/SheetsPage";

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<SheetsPage />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
