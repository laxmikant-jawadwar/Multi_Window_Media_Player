const API_BASE_URL = (import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080').replace(/\/+$/, '');

async function request(endpoint, options = {}) {
  const cleanEndpoint = endpoint.startsWith('/') ? endpoint : `/${endpoint}`;
  const url = `${API_BASE_URL}${cleanEndpoint}`;
  const config = {
    headers: {
      'Content-Type': 'application/json',
      ...options.headers,
    },
    ...options,
  };

  try {
    const response = await fetch(url, config);
    if (!response.ok) {
      const errorText = await response.text();
      throw new Error(errorText || `HTTP ${response.status} ${response.statusText}`);
    }
    const contentType = response.headers.get('content-type');
    if (contentType && contentType.includes('application/json')) {
      return await response.json();
    }
    return await response.text();
  } catch (err) {
    console.error(`API Error on ${endpoint}:`, err);
    throw err;
  }
}

export const api = {
  // Health
  getHealth: () => request('/health'),

  // Windows
  getWindows: () => request('/windows'),

  // Media
  getMediaList: () => request('/media'),
  createMedia: (title, mediaType, url, durationSeconds) =>
    request('/media', {
      method: 'POST',
      body: JSON.stringify({
        title,
        media_type: mediaType,
        url: mediaType === 'blank' ? '' : url,
        duration_seconds: Number(durationSeconds),
      }),
    }),

  // Playlist
  getPlaylist: (windowID) => request(`/windows/${windowID}/playlist`),
  addPlaylistItem: (windowID, mediaID, position) =>
    request(`/windows/${windowID}/playlist`, {
      method: 'POST',
      body: JSON.stringify({ media_id: mediaID, position }),
    }),
  deletePlaylistItem: (windowID, itemID) =>
    request(`/windows/${windowID}/playlist/${itemID}`, {
      method: 'DELETE',
    }),
  updatePlaylistItemPosition: (windowID, itemID, position) =>
    request(`/windows/${windowID}/playlist/${itemID}`, {
      method: 'PUT',
      body: JSON.stringify({ position }),
    }),

  // Playback
  startPlayback: (windowID) =>
    request(`/windows/${windowID}/playback/start`, {
      method: 'POST',
    }),
  getPlaybackState: (windowID) => request(`/windows/${windowID}/playback`),
  getCurrentMedia: (windowID) => request(`/windows/${windowID}/playback/current`),

  // Sync
  startSync: (mediaID, durationSeconds) =>
    request('/sync', {
      method: 'POST',
      body: JSON.stringify({ media_id: mediaID, duration_seconds: Number(durationSeconds) }),
    }),
  getSyncStatus: () => request('/sync'),
};
