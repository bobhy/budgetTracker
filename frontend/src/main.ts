console.log("Main.ts starting...");
import './style.css';
import App from './App.svelte';
import { mount } from 'svelte';

// --- ResizeObserver Error Suppression ---
// Define the counter on the window object
declare global {
    interface Window {
        resizeObserverLoopErrorCount: number;
    }
}
window.resizeObserverLoopErrorCount = 0;

const IGNORED_ERROR = 'ResizeObserver loop completed with undelivered notifications';

// Trap window events (standard mechanism for this error)
window.addEventListener('error', (e) => {
    if (e.message && e.message.includes(IGNORED_ERROR)) {
        window.resizeObserverLoopErrorCount++;
        e.stopImmediatePropagation();
        e.preventDefault();
        // Optional: log to debug if needed, but keeping silent as per request 'suppression'
        // console.debug(`Suppressed ResizeObserver error. Count: ${window.resizeObserverLoopErrorCount}`);
    }
});

// Chain into any existing onerror (or the one we define below)
// In this file, we are defining it, so we can just put logic inside.
window.onerror = function (message, source, lineno, colno, error) {
    if (typeof message === 'string' && message.includes(IGNORED_ERROR)) {
        window.resizeObserverLoopErrorCount++;
        return true; // Suppress the error
    }

    console.error("Global Error:", message, "at", source, lineno, colno);
    alert("Global Error: " + message);
};
// ----------------------------------------

mount(App, { target: document.getElementById('app')! });
