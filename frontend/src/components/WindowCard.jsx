import React, { useState, useEffect, useCallback } from 'react';
import { api } from '../services/api';
import { usePolling } from '../hooks/usePolling';
import { Playlist } from './Playlist';

export function WindowCard({ windowItem, mediaList }) {
  const [currentMedia, setCurrentMedia] = useState(null);
  const [playbackState, setPlaybackState] = useState(null);
  const [playlist, setPlaylist] = useState([]);
  const [error, setError] = useState('');

  const fetchPlaylist = useCallback(async () => {
    try {
      const data = await api.getPlaylist(windowItem.id);
      setPlaylist(Array.isArray(data) ? data : []);
    } catch (err) {
      console.error(`Failed to fetch playlist for ${windowItem.id}:`, err);
    }
  }, [windowItem.id]);

  useEffect(() => {
    fetchPlaylist();
  }, [fetchPlaylist]);

  // Poll current media & playback state every 1 second
  usePolling(async () => {
    try {
      const mediaData = await api.getCurrentMedia(windowItem.id);
      setCurrentMedia(mediaData);
    } catch (err) {
      setCurrentMedia(null);
    }

    try {
      const stateData = await api.getPlaybackState(windowItem.id);
      setPlaybackState(stateData);
    } catch (err) {
      setPlaybackState(null);
    }
  }, true, 1000);

  const handleStartPlayback = async () => {
    try {
      setError('');
      await api.startPlayback(windowItem.id);
      const mediaData = await api.getCurrentMedia(windowItem.id);
      setCurrentMedia(mediaData);
      const stateData = await api.getPlaybackState(windowItem.id);
      setPlaybackState(stateData);
    } catch (err) {
      setError(`Failed to start playback: ${err.message}`);
    }
  };

  const status = currentMedia?.status || playbackState?.status || 'STOPPED / NOT STARTED';

  return (
    <div className="window-card">
      <div>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '10px' }}>
          <h3 style={{ margin: 0 }}>{windowItem.name}</h3>
          <button
            onClick={handleStartPlayback}
            style={{ padding: '6px 12px', cursor: 'pointer' }}
          >
            Start Playback
          </button>
        </div>

        {error && <div style={{ color: 'red', marginBottom: '10px' }}>{error}</div>}

        <div style={{ background: '#f9f9f9', padding: '10px', borderRadius: '4px', marginBottom: '10px' }}>
          <div><strong>Status:</strong> {status}</div>
          {currentMedia && currentMedia.title ? (
            <>
              <div style={{ wordBreak: 'break-word' }}><strong>Current Media:</strong> {currentMedia.title} ({currentMedia.media_type})</div>
              <div><strong>Elapsed:</strong> {currentMedia.elapsed_seconds ?? playbackState?.elapsed_seconds ?? 0}s</div>
              <div><strong>Remaining:</strong> {currentMedia.remaining_seconds ?? playbackState?.remaining_seconds ?? 0}s</div>
              <div style={{ marginTop: '8px' }}>
                {renderMediaDisplay(currentMedia)}
              </div>
            </>
          ) : (
            <div style={{ color: '#666', marginTop: '5px' }}>Playback not active or no media playing.</div>
          )}
        </div>
      </div>

      <Playlist
        windowID={windowItem.id}
        playlist={playlist}
        mediaList={mediaList}
        onPlaylistChange={fetchPlaylist}
      />
    </div>
  );
}

function getYouTubeEmbedUrl(url) {
  if (!url) return null;
  const match = url.match(/(?:youtube\.com\/watch\?v=|youtu\.be\/)([^&]+)/);
  return match ? `https://www.youtube.com/embed/${match[1]}?autoplay=1&mute=1` : null;
}

function renderMediaDisplay(media) {
  if (!media.url) return null;

  const ytEmbedUrl = getYouTubeEmbedUrl(media.url);
  if (ytEmbedUrl) {
    return (
      <div className="media-display">
        <iframe
          src={ytEmbedUrl}
          title={media.title}
          allow="autoplay; encrypted-media; picture-in-picture"
          allowFullScreen
        />
      </div>
    );
  }

  if (media.media_type === 'image') {
    return (
      <div className="media-display">
        <img
          src={media.url}
          alt={media.title}
        />
      </div>
    );
  }

  if (media.media_type === 'video') {
    return (
      <div className="media-display">
        <video
          src={media.url}
          controls
          autoPlay
          muted
        >
          Your browser does not support video playback.
        </video>
      </div>
    );
  }

  return (
    <div className="media-display">
      <code>{media.url}</code>
    </div>
  );
}


