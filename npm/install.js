#!/usr/bin/env node

const fs = require('fs');
const path = require('path');

const BINARY_NAME = 'skills-x';

// Determine platform and architecture
function getPlatformArch() {
  const platform = process.platform;
  const arch = process.arch;

  let os;
  if (platform === 'darwin') os = 'darwin';
  else if (platform === 'linux') os = 'linux';
  else if (platform === 'win32') os = 'windows';
  else throw new Error(`Unsupported platform: ${platform}`);

  let cpu;
  if (arch === 'x64') cpu = 'amd64';
  else if (arch === 'arm64') cpu = 'arm64';
  else throw new Error(`Unsupported architecture: ${arch}`);

  return { os, cpu };
}

function main() {
  try {
    const { os, cpu } = getPlatformArch();
    console.log(`Platform: ${os}-${cpu}`);

    const binDir = path.join(__dirname, 'bin');
    const ext = os === 'windows' ? '.exe' : '';
    
    // Source binary (platform-specific)
    const srcBinary = path.join(binDir, `${BINARY_NAME}-${os}-${cpu}${ext}`);
    // Target binary (what skills-x.js expects)
    const dstBinary = path.join(binDir, `${BINARY_NAME}${ext}`);

    if (!fs.existsSync(srcBinary)) {
      throw new Error(`Binary not found: ${srcBinary}`);
    }

    // Copy/rename to the expected name
    fs.copyFileSync(srcBinary, dstBinary);
    
    // Make executable on Unix
    if (os !== 'windows') {
      fs.chmodSync(dstBinary, 0o755);
    }

    console.log(`✓ skills-x installed successfully!`);
  } catch (err) {
    console.error(`\n⚠ Installation failed: ${err.message}`);
    console.error('\nYou can install manually:');
    console.error('  go install github.com/anthropics/skills-x/cmd/skills-x@latest');
    process.exit(1);
  }
}

main();
