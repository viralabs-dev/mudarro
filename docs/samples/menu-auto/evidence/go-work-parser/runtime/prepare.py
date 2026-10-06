from pathlib import Path
import json,subprocess,os
base=Path('/tmp/mudarro-mud038-parser/runtime');binary='/tmp/mudarro-mud038-selection/mudarro'
cases={
'valid':{'go.work':'go 1.23\n// local modules\nuse (\n "./app space" // quoted\n ./lib\n)\n','app space/go.mod':'module example.test/app\ngo 1.23\n','app space/main.go':'package main\nimport("fmt";"example.test/lib";"os";"time")\nfunc main(){fmt.Println(lib.Value());if len(os.Args)>1{for{time.Sleep(time.Second)}}}\n','lib/go.mod':'module example.test/lib\ngo 1.23\n','lib/lib.go':'package lib\nfunc Value()string{return "REAL_WORKSPACE_IMPORT_OK"}\n','other/go.mod':'module example.test/other\ngo 1.23\n','other/main.go':'package main\nfunc main(){}\n','app space/nested/go.mod':'module example.test/nested\ngo 1.23\n','app space/nested/lib.go':'package nested\n'},
'missing':{'go.work':'go 1.23\nuse ./missing\n','app/go.mod':'module example.test/app\ngo 1.23\n'},
'invalid':{'go.work':'go 1.23\nuse (\n','app/go.mod':'module example.test/app\ngo 1.23\n'},
'outside':{'go.work':'go 1.23\nuse ../outside-owned\n','app/go.mod':'module example.test/app\ngo 1.23\n'},
'symlink':{'go.work':'go 1.23\nuse ./linked\n','app/go.mod':'module example.test/app\ngo 1.23\n'},
}
for name,files in cases.items():
 root=base/name;root.mkdir(exist_ok=True)
 for rel,body in files.items():p=root/rel;p.parent.mkdir(parents=True,exist_ok=True);p.write_text(body)
 if name=='symlink':(root/'linked').symlink_to(base/'valid/lib',target_is_directory=True)
 p=subprocess.run([binary,'scan','--json','--root',str(root)],capture_output=True,text=True,timeout=15);(base/(name+'-baseline-scan.json')).write_text(p.stdout);(base/(name+'-baseline-scan.stderr')).write_text(p.stderr)
 print(name,p.returncode)
