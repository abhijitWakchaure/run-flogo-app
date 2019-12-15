# What is Flogo
Learn more about [**Project Flogo**](https://www.flogo.io/) and [**TIBCO Flogo Enterprise**](https://www.tibco.com/products/tibco-flogo)
# run-flogo-app
This program will execute the latest TIBCO Flogo Enterprise app in the directory specified by you

### How to Download
You can download the release from here https://github.com/abhijitWakchaure/run-flogo-app/releases/latest

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


### The config file
When the program starts it creates a config file with name `.run-flogo-app` in the same directory. It is a simple json file which looks like this:
```
{
	"rfAppDir": "/home/abhijit/Downloads",
	"rfAppPattern": "^.+-linux_amd64.*$"
}
```
You can override the programs behavior by changing this file. 
