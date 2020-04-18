// jablko.js: Main entrypoint for the Jablko Smart Home System
// Cale Overstreet, Corinne Gerhold
// April 17, 2020
// Jablko.js serves as a manager for all Jablko related processes and is responsible for monitoring I/O output for all subprocesses, restarting processes, identifying problems internally.

const fork = require("child_process").fork // Used for forking submodules



