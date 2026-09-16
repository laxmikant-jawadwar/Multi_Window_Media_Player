import React, { useState, useEffect } from 'react';
import { api } from './services/api';
import { SyncControl } from './components/SyncControl';
import { WindowCard } from './components/WindowCard';

export function App() {
  const [windows, setWindows] = useState([]);
  const [mediaList, setMediaList] = useState([]);
  const [health, setHealth] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  const loadInitialData = async () => {
    setLoading(true);
    setError('');
    try {
      const [healthData, windowsData, mediaData] = await Promise.all([
        api.getHealth().catch(() => ({ status: 'offline' })),
        api.getWindows().catch(() => []),
        api.getMediaList().catch(() => []),
      ]);

      setHealth(healthData);
      setWindows(Array.isArray(windowsData) ? windowsData : []);
      setMediaList(Array.isArray(mediaData) ? mediaData : []);
    } catch (err) {
      setError(`Failed to connect to backend API: ${err.message}`);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadInitialData();
  }, []);

  return (
    <div style={{ maxWidth: '1200px', margin: '0 auto', padding: '20px', fontFamily: 'sans-serif' }}>
      <h1 style={{ borderBottom: '2px solid #333', paddingBottom: '10px' }}>Multi-Window Media Sequencer</h1>

      <div style={{ padding: '8px 12px', background: health?.status === 'ok' ? '#d4edda' : '#f8d7da', marginBottom: '20px', borderRadius: '4px' }}>
        <strong>Backend Health:</strong> {health?.status === 'ok' ? 'Connected (OK)' : 'Offline / Unavailable'}
        {health?.status !== 'ok' && (
          <button onClick={loadInitialData} style={{ marginLeft: '15px', padding: '2px 10px' }}>
            Retry Connection
          </button>
        )}
      </div>

      {error && (
        <div style={{ padding: '10px', background: '#f8d7da', color: '#721c24', marginBottom: '20px', borderRadius: '4px' }}>
          {error}
        </div>
      )}

      {loading ? (
        <p>Loading windows & media list...</p>
      ) : (
        <>
          <SyncControl mediaList={mediaList} />

          <h2 style={{ borderBottom: '1px solid #ccc', paddingBottom: '5px' }}>WINDOWS</h2>
          {windows.length === 0 ? (
            <p>No display windows found. Run database seed to populate demo data.</p>
          ) : (
            <div className="windows-grid">
              {windows.map((w) => (
                <WindowCard key={w.id} windowItem={w} mediaList={mediaList} />
              ))}
            </div>
          )}
        </>
      )}
    </div>
  );

}

export default App;
