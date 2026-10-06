import signal,threading
stop=threading.Event()
for s in (signal.SIGTERM,signal.SIGINT):signal.signal(s,lambda *_:stop.set())
print("PIPENV REAL STARTUP",flush=True)
stop.wait()
print("PIPENV REAL SHUTDOWN",flush=True)
