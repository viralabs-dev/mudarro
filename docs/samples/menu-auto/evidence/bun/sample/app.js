console.log("MUD033 Bun startup " + process.pid);
const timer = setInterval(() => {}, 1000);
for (const signal of ["SIGTERM", "SIGINT"]) process.on(signal, () => { console.log("MUD033 Bun shutdown " + signal); clearInterval(timer); process.exit(0); });
