import React, { useState,useEffect } from 'react';
import './App.css';

function App() {
const [status, setStatus] = useState({success: false, quote: "Loading..."})
const getNewQuoteFromServer = async()=> {
  try{const response = await fetch("http://localhost:8080/quote");
  const data = await response.json();
  if (data) {
    setStatus({
      success: false,
      quote: data.quote + " - Einstein."
  })
  }} catch(err) {
    setStatus({
      success:true, quote:"Failed to get data. " + err.message
    })
  }
  
}
  const handleClick = ()=> {
  getNewQuoteFromServer()
}
useEffect(()=> {
getNewQuoteFromServer()
}, [])
  return (
    <div className="App">
      <header className="App-header">
        <button onClick={handleClick}>get a new quote</button>
        <h3>{status.quote}</h3>
      </header>
    </div>
  );
}

export default App;
