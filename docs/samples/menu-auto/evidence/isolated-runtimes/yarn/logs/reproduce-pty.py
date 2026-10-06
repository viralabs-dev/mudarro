import os,pty,fcntl,termios,struct,subprocess,pathlib,json,time,select,codecs,errno
base=pathlib.Path('/tmp/mudarro-runtimes.K7F4J1');logs=base/'logs/yarn';binary='/tmp/mudarro-timeout-checkpoint/mudarro'
for name in ['classic','modern']:
 root=base/'fixtures/yarn'/name;env=os.environ.copy();env.update(json.loads((logs/(name+'-environment.json')).read_text()));env.pop('NO_COLOR',None);env['TERM']='xterm-256color'
 master,slave=pty.openpty();before=termios.tcgetattr(slave);fcntl.ioctl(slave,termios.TIOCSWINSZ,struct.pack('HHHH',32,100,0,0))
 proc=subprocess.Popen([binary,'menu','--root',str(root)],cwd=root,env=env,stdin=slave,stdout=slave,stderr=slave)
 start=time.monotonic();raw=bytearray();decoder=codecs.getincrementaldecoder('utf8')('replace');events=[]
 cast=(logs/(name+'-menu.cast')).open('w');cast.write(json.dumps({'version':2,'width':100,'height':32,'title':'Mudarro real '+name+' Yarn menu, isolated PTY, injected keys/SGR mouse'})+'\n')
 def pump(seconds):
  end=time.monotonic()+seconds
  while time.monotonic()<end:
   ready,_,_=select.select([master],[],[],min(.02,max(0,end-time.monotonic())))
   if ready:
    try:data=os.read(master,65536)
    except OSError as e:
     if e.errno==errno.EIO:break
     raise
    if not data:break
    raw.extend(data);cast.write(json.dumps([round(time.monotonic()-start,6),'o',decoder.decode(data)],ensure_ascii=False)+'\n');cast.flush()
 try:
  pump(.35)
  for label,data,seconds in [('service',b'1\r',.35),('quality group',b'4\r',.35),('read-only preview down',b'\x1b[B',.5),('preview wheel',b'\x1b[<65;98;20M',.3),('collapse preview',b'\x1b[A',.3),('execute real Yarn test',b'1\r',.8),('return',b'\r',.3),('actions back',b'q',.3),('groups back',b'q',.3),('exit',b'q',.3)]:
   events.append({'label':label,'time':round(time.monotonic()-start,6),'hex':data.hex()});os.write(master,data);pump(seconds)
  proc.wait(timeout=3);assert proc.returncode==0
  assert b'Yarn quality real OK' in raw, 'missing real Yarn action output'
  restored=termios.tcgetattr(slave)==before;assert restored
  (logs/(name+'-menu.pty.txt')).write_bytes(raw);(logs/(name+'-menu-media.json')).write_text(json.dumps({'exit_code':proc.returncode,'termios_restored':restored,'real_yarn_test_output':True,'duration_seconds':round(time.monotonic()-start,3),'events':events,'capture':'genuine application PTY; injected keyboard/SGR mouse, no desktop/browser'},indent=2)+'\n');print(name+' real PTY menu/action PASS')
 finally:
  if proc.poll() is None:proc.terminate();proc.wait(timeout=3)
  cast.close();os.close(slave);os.close(master)
