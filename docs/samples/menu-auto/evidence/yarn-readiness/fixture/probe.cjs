const fs = require('node:fs'); console.log('Yarn fixture probe executed'); fs.writeFileSync('runtime-marker.txt', process.cwd());
