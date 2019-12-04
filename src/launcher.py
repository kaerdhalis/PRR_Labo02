#import subprocess
#import json
#import os
#import sys 

with open('src/config/config.json') as json_file:
    data = json.load(json_file)

nmbofProc= data['NumberOfProcesses']
print(nmbofProc)



if(sys.platform.startswith("win32")):
	for x in xrange(0,5):
		os.system("start cmd.exe @cmd /k \"go run src/client/client.go "+str(x)+"\"")

if(sys.platform=="darwin"):
	for x in xrange(0,5):
		subprocess.call(['open', '-W', '-a', 'Terminal.app', 'go', '--args', 'run', 'src/client/client.go',str(x)])

if(sys.platform.startswith("linux")):
	for x in xrange(0,5):
		subprocess.call(['gnome-terminal', '-x', 'go run src/client/client.go '+str(x)])



