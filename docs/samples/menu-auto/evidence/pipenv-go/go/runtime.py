import os,pathlib,subprocess,json,hashlib,time
base=pathlib.Path('/tmp/mudarro-mud038');root=base/'runtime-final-2';root.mkdir(exist_ok=True)
files={'go.work':'go 1.23\nuse ./app\n','app/go.mod':'module example.test/app\ngo 1.23\n','app/cmd/live/main.go':'''package main
import ("fmt"; "os"; "os/signal"; "syscall")
func main(){fmt.Println("MUD038 Go startup"); ch:=make(chan os.Signal,1); signal.Notify(ch,syscall.SIGTERM,syscall.SIGINT); <-ch; fmt.Println("MUD038 Go shutdown")}
''','app/cmd/live/main_test.go':'package main\nimport "testing"\nfunc TestOffline(t *testing.T){t.Log("MUD038 real go test") }\n','app/cmd/tool/main.go':'//go:build ignore\n\npackage main\nfunc main(){}\n'}
for p,s in files.items(): target=root/p;target.parent.mkdir(parents=True,exist_ok=True);target.write_text(s)
bin=base/'mudarro';env=dict(os.environ,PATH='/home/danielsouza/sdk/go1.27.1/bin:'+str(base)+':'+os.environ['PATH'],GOCACHE='/tmp/mudarro-go-cache',GOPROXY='off',GOTOOLCHAIN='local',GOFLAGS='-buildvcs=false')
results=[]
def run(args,check=True):
 p=subprocess.run([str(a) for a in args],env=env,stdout=subprocess.PIPE,stderr=subprocess.STDOUT,text=True)
 results.append({'argv':[str(a) for a in args],'exit':p.returncode,'output':p.stdout})
 if check and p.returncode: raise RuntimeError(results[-1])
 return p
try:
 scan=run([bin,'scan','--root',root,'--json']);report=json.loads(scan.stdout)
 assert report['warnings'] and report['config']['services'][0]['commands']['start']['args']==['go','run','./cmd/live']
 run([bin,'init','--root',root]);run([bin,'generate','--root',root])
 before={str(p.relative_to(root)):hashlib.sha256(p.read_bytes()).hexdigest() for p in root.rglob('*') if p.is_file()}
 run([bin,'generate','--root',root]);after={str(p.relative_to(root)):hashlib.sha256(p.read_bytes()).hexdigest() for p in root.rglob('*') if p.is_file()};assert before==after
 (base/'runtime-hashes.json').write_text(json.dumps(before,indent=2)+'\n')
 for action in ['test','build','up','status','logs','restart','down','status']:
  p=run(['bash',root/'.mudarro/scripts/app-go'/f'{action}.sh'])
  if action=='up':time.sleep(.15)
  if action=='status' and results[-2]['argv'][-1].endswith('down.sh'):assert 'stopped' in p.stdout
finally:
 run([bin,'run','app-go:down','--root',root],check=False)
 (base/'runtime-final.json').write_text(json.dumps(results,indent=2)+'\n')
print('Offline workspace CLI scan/init/generate2/test/build/up/status/logs/restart/down PASS')
