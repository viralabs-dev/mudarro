import os,time,signal
print("UV startup",os.getpid(),flush=True)
def stop(sig,frame):
 print("UV shutdown",sig,flush=True);raise SystemExit(0)
signal.signal(signal.SIGTERM,stop)
while True:time.sleep(.1)
