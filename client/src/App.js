import React, { useState } from 'react';
import { sessionState } from './Requests'
import { ActiveTracking } from './components/activeTracking'
import CssBaseline from '@material-ui/core/CssBaseline';

function App() {
  const [sessState, setSessState] = useState('loading...');
  const [init, setInit] = useState(true);

  if (init) {
    setInit(false)
    sessionState()
      .then(r => {
        if (r === 'error') {
          setSessState('Encountered error, please refresh page')
          return
        }

        if (!r.length) {
          setSessState('NO_ACTIVE')
          return
        }

        if (r.length === 1 && !r[0].time_session_partial_end.Valid) {
          setSessState('ACTIVE_TRACKING')

          return
        }

        setSessState('ACTIVE_SESSION')
      })
  }

  return (
    <React.Fragment>
      <CssBaseline />
      <div style={{padding: "2em"}}>
        <ActiveTracking s={sessState} />
      </div>
      
    </React.Fragment>
  );
}

export default App;
