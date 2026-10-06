#!/usr/bin/env python3
"""Record genuine Mudarro PTY output; injected events are documented separately.
Uses only Python standard library. Runs no preview action and opens no GUI.
Usage: python3 record_cli_preview.py [absolute_mudarro_binary] [output_directory]
"""
import codecs, errno, fcntl, hashlib, json, os, pathlib, pty, select, struct, subprocess, sys, termios, time

BINARY = pathlib.Path(sys.argv[1] if len(sys.argv)>1 else '/tmp/mudarro-ui-bin/mudarro').resolve()
BASE = pathlib.Path(sys.argv[2] if len(sys.argv)>2 else '/tmp/mudarro-ui-recording').resolve()
BASE.mkdir(parents=True,exist_ok=True)
COLS,ROWS=100,32

def record(profile,locale,theme):
    root=BASE/profile
    (root/'scripts').mkdir(parents=True,exist_ok=True)
    source=['#!/usr/bin/env bash','set -euo pipefail','# Aurora: reviewable quality task. Recording never executes this script.']
    source += [f'printf "Aurora quality check {i:02d}: deterministic fixture\\n"' for i in range(1,61)]
    source += ['touch SHOULD_NOT_EXECUTE']
    (root/'scripts/quality.sh').write_text('\n'.join(source)+'\n')
    config={'version':1,'name':'Aurora','ui':{'locale':locale,'theme':theme,'lettering':'auto','preview':{'enabled':True,'mouse':'on'}},'services':[{'id':'aurora','dir':'.','language':'custom','infrastructure':{'kind':'custom'},'commands':{'test':{'args':['bash','scripts/quality.sh'],'group':'qualidade'}}}]}
    (root/'mudarro.json').write_text(json.dumps(config,indent=2)+'\n')
    master,slave=pty.openpty()
    fcntl.ioctl(slave,termios.TIOCSWINSZ,struct.pack('HHHH',ROWS,COLS,0,0))
    env=os.environ.copy();env.pop('NO_COLOR',None);env.update(TERM='xterm-256color',COLUMNS=str(COLS))
    proc=subprocess.Popen([str(BINARY),'menu','--root',str(root)],stdin=slave,stdout=slave,stderr=slave,env=env,close_fds=True)
    os.close(slave)
    events=[]; output=bytearray();started=time.monotonic();decoder=codecs.getincrementaldecoder('utf-8')('replace')
    cast=BASE/(profile+'.cast');raw=BASE/(profile+'.pty.txt')
    fd=cast.open('w');fd.write(json.dumps({'version':2,'width':COLS,'height':ROWS,'timestamp':int(time.time()),'title':'Aurora — genuine Mudarro CLI preview '+profile,'env':{'TERM':'xterm-256color','SHELL':'bash'}})+'\n')
    def pump(duration):
        deadline=time.monotonic()+duration
        while time.monotonic()<deadline:
            ready,_,_=select.select([master],[],[],min(.05,max(0,deadline-time.monotonic())))
            if not ready:continue
            try:data=os.read(master,65536)
            except OSError as e:
                if e.errno==errno.EIO:return
                raise
            if not data:return
            output.extend(data);text=decoder.decode(data)
            if text:fd.write(json.dumps([round(time.monotonic()-started,6),'o',text],ensure_ascii=False)+'\n');fd.flush()
    def inject(label,data,pause=.8):
        events.append({'time':round(time.monotonic()-started,6),'label':label,'hex':data.hex(),'injected':True})
        os.write(master,data);pump(pause)
    try:
        pump(.8)
        inject('select service 1',b'1\n')
        inject('select quality group 1',b'1\n',1.0)
        if b'\x1b[?1049h' not in output:raise RuntimeError('CLI did not enter interactive preview mode')
        inject('Down: expand preview',b'\x1b[B',1.2)
        inject('PageDown: scroll preview',b'\x1b[6~',1.0)
        inject('SGR wheel down (injected mouse event)',b'\x1b[<65;98;15M',1.0)
        inject('SGR scrollbar left press',b'\x1b[<0;100;10M',.5)
        inject('SGR scrollbar drag',b'\x1b[<32;100;24M',1.0)
        inject('SGR scrollbar release',b'\x1b[<0;100;24m',.5)
        inject('Up: collapse preview',b'\x1b[A',1.0)
        inject('q: return from actions without executing',b'q',.6)
        inject('q + Enter: return from groups',b'q\n',.5)
        inject('q + Enter: exit service menu',b'q\n',.5)
        proc.wait(timeout=3);pump(.1)
        if proc.returncode!=0:raise RuntimeError('CLI exit '+str(proc.returncode))
        if (root/'SHOULD_NOT_EXECUTE').exists():raise RuntimeError('Preview executed the fixture action')
        final=decoder.decode(b'',final=True)
        if final:fd.write(json.dumps([round(time.monotonic()-started,6),'o',final])+'\n')
        raw.write_bytes(output)
        meta={'profile':profile,'locale':locale,'theme':theme,'terminal':{'columns':COLS,'rows':ROWS},'binary':str(BINARY),'binary_sha256':hashlib.sha256(BINARY.read_bytes()).hexdigest(),'exit_code':proc.returncode,'duration_seconds':round(time.monotonic()-started,3),'output_bytes':len(output),'events':events,'action_executed':False,'capture':'actual PTY application output; no GUI capture; keyboard and SGR mouse bytes injected by helper'}
        (BASE/(profile+'.metadata.json')).write_text(json.dumps(meta,indent=2)+'\n')
        print(json.dumps({k:meta[k] for k in ['profile','exit_code','duration_seconds','output_bytes','action_executed']}))
    finally:
        fd.close()
        if proc.poll() is None:proc.terminate();proc.wait(timeout=3)
        os.close(master)

record('aurora-dark-en','en','dark')
record('aurora-light-ptbr','pt-BR','light')
