# Client-Side Rendering Code Test
## insertAdjacentHTML("beforeend", "markup") vs. Fragment and appendChild()
___
### Test for 100 card insertions with basic string templating

insertAdjacentHTML:    61 ms
Fragment appendChild: 111 ms

insertAdjacentHTML shows a roughly 2x speed-up. Both tests used a function prototype that returns either a HTML string or HTML element respectively. In insertAdjacentHTML, this string was immediately fed to insertAdjacentHTML. In the fragment demo, a document fragement had the returned HTML elements appended to it, and then the fragment was appended to the holder. The values above do not include the time spent appending the fragment to the holder. 

**insertAdjacentHTML should be used for the client-side rendering procedure.**

___
### Test Setup
inxi output

CPU: Quad Core Intel Core i7-8565U (-MT MCP-) 
speed/min/max: 1168/400/4600 MHz
Kernel: 5.4.95-1-MANJARO x86_64
Mem: 2195.3/15834.8 MiB (13.9%)
Storage: 1.94 TiB (16.1% used)
Procs: 255 
