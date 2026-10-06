if (!globalThis.Bun || !Bun.version) throw new Error("Bun runtime expected");
console.log("MUD033 Bun check PASS " + Bun.version);
