import pathlib,zipfile,email,hashlib,json,urllib.request
root=pathlib.Path('/tmp/mudarro-runtimes.K7F4J1');base=root/'pipenv';rows=[];requirements=[]
for wheel in sorted((base/'wheelhouse').glob('*.whl')):
 with zipfile.ZipFile(wheel) as z:
  metadata=email.message_from_bytes(z.read(next(n for n in z.namelist() if n.endswith('.dist-info/METADATA'))))
 name,version=metadata['Name'],metadata['Version'];digest=hashlib.sha256(wheel.read_bytes()).hexdigest()
 url=f'https://pypi.org/pypi/{name}/{version}/json'
 with urllib.request.urlopen(url,timeout=20) as response:data=json.load(response)
 (base/(name+'-metadata.json')).write_text(json.dumps(data,indent=2)+'\n')
 official=next(u for u in data['urls'] if u['filename']==wheel.name)
 assert official['digests']['sha256']==digest,(wheel,digest)
 if name.lower()=='pipenv':assert version=='2026.8.0' and digest=='8da889c636a8cab59a435c1366f30d252ecee465372341eb8cd4da43b055ad9d'
 rows.append({'name':name,'version':version,'wheel':wheel.name,'sha256':digest,'officialURL':official['url'],'metadataURL':url})
 requirements.append(f'{name}=={version} --hash=sha256:{digest}')
(base/'requirements-hashed.txt').write_text('\n'.join(requirements)+'\n');(base/'verified-wheels.json').write_text(json.dumps(rows,indent=2)+'\n')
print('Verified',len(rows),'official PyPI wheel SHA256 digests including required Pipenv pin')
