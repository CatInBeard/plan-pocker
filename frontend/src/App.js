import 'bootstrap/dist/css/bootstrap.min.css';
import 'bootstrap-icons/font/bootstrap-icons.css';

import Header from "./components/Header";
import Container from "./components/Container"

import { BrowserRouter as Router, Route, Routes, useParams } from 'react-router-dom';
import CreateGame from './components/CreateGame';
import Game from './components/Game';

function HomePage() {
  
  return  <>
    <Header/>
    <Container>
      <CreateGame/>
    </Container>
    </>
}

function GamePage() {
  const { id } = useParams();

return <Game id={id}/>
}

function App() {

  return (
    <Router>
      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/:id" element={<GamePage />} />
      </Routes>
    </Router>
  );
}

export default App;
