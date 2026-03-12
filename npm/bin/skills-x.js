#!/usr/bin/env node
'use strict';

const { execFileSync } = require('child_process');
const path = require('path');
const os = require('os');

const ext = os.platform() === 'win32' ? '.exe' : '';
const binary = path.join(__dirname, 'skills-x' + ext);

try {
  execFileSync(binary, process.argv.slice(2), { stdio: 'inherit' });
} catch (e) {
  process.exit(e.status || 1);
}
