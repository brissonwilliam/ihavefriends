import logo from './logo.svg';
import './App.css';
import ButtonCounter from './components/button_counter'

function App() {
  return (
    <div className="App">
      <header className="App-header">
        <img src={logo} className="App-logo" alt="logo" />
      </header>
      <div className="App-Body">
        <div calssName="App-Container-ButtonCounter">
          <ButtonCounter />
        </div>
        
      </div>
      
    </div>
  );
}

export default App;
