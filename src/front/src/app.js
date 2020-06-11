import React from 'react';
import logo from './logo.svg';
import './app.css';
import SearchAppBar from './navbar';
import FullWidthGrid from './main';

function App() {
  return (
    <div className="App">
      <SearchAppBar></SearchAppBar>
      <FullWidthGrid></FullWidthGrid>
    </div>
    
  );
}

export default App;
