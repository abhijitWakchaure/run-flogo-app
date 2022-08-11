# What is Flogo
Learn more about [**Project Flogo**](https://www.flogo.io/) and [**TIBCO Flogo Enterprise**](https://www.tibco.com/products/tibco-flogo)
## Why I created *run-flogo-app*
Simple, sheer laziness on my part, but this will also save few precious seconds of your life! I hope you find this helpful.

### How to Download
You can download the release from here https://github.com/abhijitWakchaure/run-flogo-app/releases/latest

### How to Install (Linux)
First make the binary executable by running the following command:
```
chmod +x ~/Downloads/run-flogo-app-linux_amd64
``` 
Then install the program with `-install` flag
```
~/Downloads/run-flogo-app-linux_amd64 -install
```

### How to Use
Make sure the binary you downloaded is executable, then you can directly run it as executable
```
./run-flogo-app
```
You can provide `-debug` flag to enable the debug logs in your flogo app, just make sure its the `first` argument to the program like this:
```
./run-flogo-app -debug
```
Also, you can pass command line arguments if your flogo app supports it; like this (the -debug flag is not mandatory for this):
```
./run-flogo-app -debug arg1 arg2 arg3
```

### Flags
You can provide the following flags to the main program like this 
```
./run-flogo-app -flag1 -flag2 -flag3
```

| Flag          | Use                                                   |
| :------------ |:----------------------------------------------------  |
| -debug         | To enable debug logs for your flogo app              |
| -install       | To install **run-flogo-app** (the main program)      |
| -uninstall     | To uninstall run-flogo-app and to remove the config  |
| -version       | To print the version info                            |
| -help          | To print the available flags                         |


### The config file
When the program starts it creates a config file with name `run-flogo-app.config` in your home directory. It is a simple json file which looks like this:
```
{
	"appsDir": "/home/abhijit/Downloads",
	"appPattern": "^.+-linux_amd64.*$",
	"isUpdateAvailable": false,
	"updateURL": "",
	"releaseNotes": ""
}
```
You can override the programs behavior by changing `appsDir` and `appPattern` variables in this file. 
