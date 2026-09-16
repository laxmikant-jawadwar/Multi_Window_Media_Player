import React, { useState } from 'react';
import { api } from '../services/api';

export function Playlist({ windowID, playlist, mediaList, onPlaylistChange }) {
  const [selectedMediaId, setSelectedMediaId] = useState('');
  const [error, setError] = useState('');

  const handleAddMedia = async (e) => {
    e.preventDefault();
    if (!selectedMediaId) {
      setError('Select a media item to add');
      return;
    }

    try {
      setError('');
      const nextPosition = playlist.length + 1;
      await api.addPlaylistItem(windowID, selectedMediaId, nextPosition);
      setSelectedMediaId('');
      onPlaylistChange();
    } catch (err) {
      setError(`Failed to add media: ${err.message}`);
    }
  };

  const handleDelete = async (itemId) => {
    try {
      setError('');
      await api.deletePlaylistItem(windowID, itemId);
      onPlaylistChange();
    } catch (err) {
      setError(`Failed to delete item: ${err.message}`);
    }
  };

  const handleMove = async (index, direction) => {
    const targetIndex = index + direction;
    if (targetIndex < 0 || targetIndex >= playlist.length) return;

    const itemToMove = playlist[index];
    const newPosition = targetIndex + 1;

    try {
      setError('');
      await api.updatePlaylistItemPosition(windowID, itemToMove.id, newPosition);
      onPlaylistChange();
    } catch (err) {
      setError(`Failed to move item: ${err.message}`);
    }
  };

  return (
    <div style={{ marginTop: '15px' }}>
      <h4 style={{ margin: '0 0 8px 0' }}>Playlist Items ({playlist.length}):</h4>
      {error && <div style={{ color: 'red', marginBottom: '10px' }}>{error}</div>}

      {playlist.length === 0 ? (
        <p style={{ color: '#666', margin: '5px 0' }}>Playlist is empty.</p>
      ) : (
        <ol style={{ paddingLeft: '20px', margin: '0' }}>
          {playlist.map((item, index) => (
            <li key={item.id || index} style={{ marginBottom: '8px' }}>
              <div className="playlist-item">
                <span>
                  <strong>{item.title || item.media_id}</strong> ({item.media_type || 'media'}, {item.duration_seconds || 0}s)
                </span>
                <span className="playlist-btn-group">
                  <button
                    onClick={() => handleMove(index, -1)}
                    disabled={index === 0}
                  >
                    Up
                  </button>
                  <button
                    onClick={() => handleMove(index, 1)}
                    disabled={index === playlist.length - 1}
                  >
                    Down
                  </button>
                  <button
                    onClick={() => handleDelete(item.id)}
                    style={{ color: 'red' }}
                  >
                    Delete
                  </button>
                </span>
              </div>
            </li>
          ))}
        </ol>
      )}

      <form onSubmit={handleAddMedia} className="add-media-container">
        <strong>Add Media:</strong>
        <div className="add-media-row">
          <select
            value={selectedMediaId}
            onChange={(e) => setSelectedMediaId(e.target.value)}
            className="add-media-select"
          >
            <option value="">-- Select Media --</option>
            {mediaList.map((m) => (
              <option key={m.id} value={m.id}>
                {m.title} ({m.media_type}, {m.duration_seconds}s)
              </option>
            ))}
          </select>
          <button type="submit" style={{ padding: '5px 12px', flexShrink: 0 }}>
            Add
          </button>
        </div>
      </form>
    </div>
  );
}

