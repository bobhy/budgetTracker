---
description: Clears the Antigravity application cache to fix blank screen issues. Note: Close Antigravity before running this if possible, or restart immediately after.
---

1. Remove the main Cache directory
// turbo
rm -rf ~/.config/Antigravity/Cache

2. Remove the Code Cache
// turbo
rm -rf ~/.config/Antigravity/Code\ Cache

3. Remove the Service Worker Cache
// turbo
rm -rf ~/.config/Antigravity/Service\ Worker

4. Remove the GPU Cache
// turbo
rm -rf ~/.config/Antigravity/GPUCache

5. Remove WebGPU Caches
// turbo
rm -rf ~/.config/Antigravity/DawnWebGPUCache
// turbo
rm -rf ~/.config/Antigravity/DawnGraphiteCache
