import { useEffect, useRef } from 'react';

export function usePolling(callback, delayEnabled = true, delayMs = 1000) {
  const savedCallback = useRef(callback);

  useEffect(() => {
    savedCallback.current = callback;
  }, [callback]);

  useEffect(() => {
    if (!delayEnabled) return;

    // Run callback immediately
    savedCallback.current();

    const intervalId = setInterval(() => {
      savedCallback.current();
    }, delayMs);

    return () => clearInterval(intervalId);
  }, [delayEnabled, delayMs]);
}
