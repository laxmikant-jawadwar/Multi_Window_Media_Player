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

  const status = currentMedia?.status || playbackState?.status || 'STOPPED';
  const isStopped = status === 'STOPPED';
  const hasStarted = Boolean(playbackState?.cycle_started_at || currentMedia?.cycle_started_at);
  const isCycleCompleted = isStopped && hasStarted && (currentMedia?.message === '5-hour playback cycle completed' || (playbackState?.elapsed_seconds != null && playbackState.elapsed_seconds >= 18000));
  const cycleCompletedMessage = isCycleCompleted ? '5-hour playback cycle completed' : null;

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
          <div>
            <strong>Status:</strong> {status}
            {isStopped && cycleCompletedMessage && (
              <span style={{ marginLeft: '8px', color: '#666', fontSize: '0.9em' }}>
                ({cycleCompletedMessage})
              </span>
            )}
          </div>
          {currentMedia && currentMedia.title && !isStopped ? (
            <>
              <div style={{ wordBreak: 'break-word' }}><strong>Current Media:</strong> {currentMedia.title} ({currentMedia.media_type})</div>
              <div><strong>Elapsed:</strong> {currentMedia.elapsed_seconds ?? playbackState?.elapsed_seconds ?? 0}s</div>
              <div><strong>Remaining:</strong> {currentMedia.remaining_seconds ?? playbackState?.remaining_seconds ?? 0}s</div>
              <div style={{ marginTop: '8px' }}>
                <MediaDisplay media={currentMedia} />
              </div>
            </>
          ) : playlist.length === 0 && !hasStarted ? (
            <div style={{ marginTop: '8px' }}>
              <div className="media-display" style={{ flexDirection: 'column', color: '#666', textAlign: 'center', padding: '20px' }}>
                <div style={{ fontWeight: 'bold', fontSize: '1.1em', marginBottom: '4px' }}>Playlist is empty</div>
                <div>Add media to this window to start playback.</div>
              </div>
            </div>
          ) : (
            <div style={{ marginTop: '8px' }}>
              <div className="media-display" style={{ flexDirection: 'column', color: '#666', textAlign: 'center', padding: '20px' }}>
                <div style={{ fontWeight: 'bold', fontSize: '1.1em', marginBottom: '4px' }}>Playback stopped</div>
                <div>{cycleCompletedMessage || 'Click "Start Playback" to begin.'}</div>
              </div>
            </div>
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

function MediaDisplay({ media }) {
  const [hasError, setHasError] = useState(false);

  useEffect(() => {
    setHasError(false);
  }, [media?.media_id, media?.url, media?.title]);

  if (!media) return null;

  if (media.media_type === 'blank') {
    return (
      <div className="media-display">
        <span style={{ color: '#888', fontStyle: 'italic' }}>Blank</span>
      </div>
    );
  }

  if (hasError) {
    return (
      <div className="media-display">
        <span style={{ color: '#888' }}>Media unavailable</span>
      </div>
    );
  }

  if (media.media_type === 'image') {
    if (!media.url) {
      return (
        <div className="media-display">
          <span style={{ color: '#888' }}>Media unavailable</span>
        </div>
      );
    }
    return (
      <div className="media-display">
        <img
          src={media.url}
          alt={media.title || 'Media'}
          onError={() => setHasError(true)}
        />
      </div>
    );
  }

  if (media.media_type === 'video') {
    if (!media.url) {
      return (
        <div className="media-display">
          <span style={{ color: '#888' }}>Media unavailable</span>
        </div>
      );
    }
    const ytEmbedUrl = getYouTubeEmbedUrl(media.url);
    if (ytEmbedUrl) {
      return (
        <div className="media-display">
          <iframe
            src={ytEmbedUrl}
            title={media.title || 'Video'}
            allow="autoplay; encrypted-media; picture-in-picture"
            allowFullScreen
            onError={() => setHasError(true)}
          />
        </div>
      );
    }
    return (
      <div className="media-display">
        <video
          src={media.url}
          controls
          autoPlay
          muted
          onError={() => setHasError(true)}
        >
          Your browser does not support video playback.
        </video>
      </div>
    );
  }

  return (
    <div className="media-display">
      <span style={{ color: '#888' }}>Unsupported media type</span>
    </div>
  );
}


