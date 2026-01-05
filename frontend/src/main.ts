console.log("Main.ts starting...");
import './style.css';
import App from './App.svelte';
import { mount } from 'svelte';

window.onerror = function (message, source, lineno, colno, error) {
    console.error("Global Error:", message, "at", source, lineno, colno);
    alert("Global Error: " + message);
};

mount(App, { target: document.getElementById('app')! });
