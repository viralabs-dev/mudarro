const fs=require('node:fs');fs.writeFileSync('probe-marker.json',JSON.stringify({cwd:process.cwd(),args:process.argv.slice(2)}));console.log('Yarn quality real OK');
