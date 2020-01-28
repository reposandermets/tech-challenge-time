import React, { useState }  from 'react';
import './App.css';
import { sessionState } from './Requests'

function App() {
  const [sessState, setSessState] = useState('loading...');
  const [init, setInit] = useState(true);

  if(init) {
    setInit(false)
    sessionState()
      .then(r => {
        if(r === 'error') {
          setSessState('Encountered error, please refresh page')
          return
        }

        if(!r.length) {
          setSessState('No ative sessions')
          return
        }

        if(r.length === 1 && !r[0].time_session_partial_end.Valid) {
          setSessState('Active tracking')
          return
        }

        setSessState('Active sessions')
      })
  }

  return (
    <div>
      {sessState}
    </div>
  );
}

export default App;
