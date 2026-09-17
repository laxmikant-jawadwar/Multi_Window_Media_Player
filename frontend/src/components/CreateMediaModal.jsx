import React, { useState } from 'react';
import { api } from '../services/api';

export function CreateMediaModal({ windows, onMediaCreated, onClose }) {
  const [title, setTitle] = useState('');
  const [mediaType, setMediaType] = useState('image');
  const [url, setUrl] = useState('');
  const [durationSeconds, setDurationSeconds] = useState(10);
  const [targetWindowId, setTargetWindowId] = useState('');
  const [error, setError] = useState('');
  const [submitting, setSubmitting] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!title.trim()) {
      setError('Title is required');
      return;
    }

    if (mediaType !== 'blank' && !url.trim()) {
      setError('URL is required for image and video media');
      return;
    }

    if (!durationSeconds || durationSeconds <= 0) {
      setError('Duration must be greater than 0');
      return;
    }

    try {
      setError('');
      setSubmitting(true);

      // Create Media Record
      const createdMedia = await api.createMedia(
        title.trim(),
        mediaType,
        mediaType === 'blank' ? '' : url.trim(),
        durationSeconds
      );

      // If target window selected, add to its playlist
      if (targetWindowId && createdMedia && createdMedia.id) {
        const currentPlaylist = await api.getPlaylist(targetWindowId).catch(() => []);
        const nextPosition = Array.isArray(currentPlaylist) ? currentPlaylist.length + 1 : 1;
        await api.addPlaylistItem(targetWindowId, createdMedia.id, nextPosition);
      }

      onMediaCreated();
      onClose();
    } catch (err) {
      setError(`Failed to create media: ${err.message}`);
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div
      style={{
        position: 'fixed',
        top: 0,
        left: 0,
        right: 0,
        bottom: 0,
        background: 'rgba(0,0,0,0.5)',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        zIndex: 1000,
      }}
    >
      <div
        style={{
          background: '#fff',
          padding: '20px',
          borderRadius: '8px',
          maxWidth: '500px',
          width: '90%',
          boxShadow: '0 4px 12px rgba(0,0,0,0.15)',
        }}
      >
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '15px' }}>
          <h3 style={{ margin: 0 }}>Create New Media</h3>
          <button onClick={onClose} style={{ border: 'none', background: 'transparent', fontSize: '18px', cursor: 'pointer' }}>
            ✕
          </button>
        </div>

        {error && <div style={{ color: 'red', marginBottom: '12px', fontSize: '14px' }}>{error}</div>}

        <form onSubmit={handleSubmit} style={{ display: 'flex', flexDirection: 'column', gap: '12px' }}>
          <div>
            <label style={{ display: 'block', fontWeight: 'bold', marginBottom: '4px', fontSize: '14px' }}>Title:</label>
            <input
              type="text"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              placeholder="e.g. Galaxy S26 Product"
              style={{ width: '100%', padding: '8px', boxSizing: 'border-box' }}
            />
          </div>

          <div>
            <label style={{ display: 'block', fontWeight: 'bold', marginBottom: '4px', fontSize: '14px' }}>Media Type:</label>
            <select
              value={mediaType}
              onChange={(e) => setMediaType(e.target.value)}
              style={{ width: '100%', padding: '8px', boxSizing: 'border-box' }}
            >
              <option value="image">Image</option>
              <option value="video">Video</option>
              <option value="blank">Blank</option>
            </select>
          </div>

          {mediaType !== 'blank' && (
            <div>
              <label style={{ display: 'block', fontWeight: 'bold', marginBottom: '4px', fontSize: '14px' }}>
                Media URL ({mediaType === 'image' ? 'Direct Image Link' : 'Video / YouTube Link'}):
              </label>
              <input
                type="text"
                value={url}
                onChange={(e) => setUrl(e.target.value)}
                placeholder={mediaType === 'image' ? 'https://example.com/image.jpg' : 'https://www.youtube.com/watch?v=...'}
                style={{ width: '100%', padding: '8px', boxSizing: 'border-box' }}
              />
            </div>
          )}

          <div>
            <label style={{ display: 'block', fontWeight: 'bold', marginBottom: '4px', fontSize: '14px' }}>Duration (Seconds):</label>
            <input
              type="number"
              min="1"
              value={durationSeconds}
              onChange={(e) => setDurationSeconds(e.target.value)}
              style={{ width: '100%', padding: '8px', boxSizing: 'border-box' }}
            />
          </div>

          <div>
            <label style={{ display: 'block', fontWeight: 'bold', marginBottom: '4px', fontSize: '14px' }}>
              Add directly to window playlist (Optional):
            </label>
            <select
              value={targetWindowId}
              onChange={(e) => setTargetWindowId(e.target.value)}
              style={{ width: '100%', padding: '8px', boxSizing: 'border-box' }}
            >
              <option value="">-- Do Not Add To Playlist Yet --</option>
              {windows.map((w) => (
                <option key={w.id} value={w.id}>
                  {w.name}
                </option>
              ))}
            </select>
          </div>

          <div style={{ display: 'flex', justifyContent: 'flex-end', gap: '10px', marginTop: '10px' }}>
            <button type="button" onClick={onClose} style={{ padding: '8px 16px' }}>
              Cancel
            </button>
            <button type="submit" disabled={submitting} style={{ padding: '8px 16px', background: '#007bff', color: '#fff', border: 'none', borderRadius: '4px' }}>
              {submitting ? 'Creating...' : 'Create Media'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
