import React, { useState } from 'react';
import { api } from '../services/api';
import { usePolling } from '../hooks/usePolling';

export function SyncControl({ mediaList }) {
  const [selectedMediaId, setSelectedMediaId] = useState('');
  const [duration, setDuration] = useState(30);
  const [syncStatus, setSyncStatus] = useState(null);
  const [error, setError] = useState('');

  // Poll sync status every 1 second
  usePolling(async () => {
    try {
      const data = await api.getSyncStatus();
      setSyncStatus(data);
      setError('');
    } catch (err) {
      if (err.message && (err.message.includes('404') || err.message.includes('not found'))) {
        setSyncStatus({ status: 'INACTIVE' });
      } else {
        // preserve existing if network temporary glitch
      }
    }
  }, true, 1000);

  const handleStartSync = async (e) => {
    e.preventDefault();
    if (!selectedMediaId) {
      setError('Please select media to broadcast sync.');
      return;
    }

    try {
      setError('');
      await api.startSync(selectedMediaId, duration);
      const data = await api.getSyncStatus();
      setSyncStatus(data);
    } catch (err) {
      setError(`Failed to start sync: ${err.message}`);
    }
  };

  const isSyncActive = syncStatus && (syncStatus.status === 'ACTIVE' || syncStatus.active === true);

  return (
    <div style={{ border: '1px solid #999', padding: '15px', marginBottom: '20px', background: '#fafafa' }}>
      <h3 style={{ marginTop: 0 }}>SYNC BROADCAST CONTROL</h3>
      {error && <div style={{ color: 'red', marginBottom: '10px' }}>{error}</div>}

      <form onSubmit={handleStartSync} style={{ display: 'flex', gap: '15px', alignItems: 'center', marginBottom: '15px', flexWrap: 'wrap' }}>
        <label>
          <strong>Sync Media:</strong>{' '}
          <select
            value={selectedMediaId}
            onChange={(e) => setSelectedMediaId(e.target.value)}
            style={{ padding: '5px' }}
          >
            <option value="">-- Select Media --</option>
            {mediaList.map((m) => (
              <option key={m.id} value={m.id}>
                {m.title} ({m.media_type}, {m.duration_seconds}s)
              </option>
            ))}
          </select>
        </label>

        <label>
          <strong>Duration (sec):</strong>{' '}
          <input
            type="number"
            min="1"
            value={duration}
            onChange={(e) => setDuration(e.target.value)}
            style={{ width: '70px', padding: '5px' }}
          />
        </label>

        <button type="submit" style={{ padding: '6px 15px', cursor: 'pointer' }}>
          Start Sync
        </button>
      </form>

      <div style={{ background: '#eee', padding: '10px', borderRadius: '4px' }}>
        <div><strong>Sync Status:</strong> {isSyncActive ? 'ACTIVE' : 'INACTIVE'}</div>
        {isSyncActive && (
          <>
            <div><strong>Media ID:</strong> {syncStatus.media_id || 'N/A'}</div>
            <div><strong>Remaining:</strong> {syncStatus.remaining_seconds ?? syncStatus.remaining ?? 'N/A'} seconds</div>
          </>
        )}
      </div>
    </div>
  );
}
