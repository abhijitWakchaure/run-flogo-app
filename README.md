# Run Flogo App

Run the downloaded flogo app binaries with single command

## What is Flogo

Learn more about [**Project Flogo**](https://www.flogo.io/) and [**TIBCO Flogo Enterprise**](https://www.tibco.com/products/tibco-flogo)

## Why *run-flogo-app*

Simple, sheer laziness on my part, but this will also save few precious seconds of your life! I hope you find this helpful.

## How to

### How to Download

You can download the release from here https://github.com/abhijitWakchaure/run-flogo-app/releases/latest

### How to Install (Linux)

First make the binary executable by running the following command:

```bash
chmod +x ~/Downloads/run-flogo-app-linux_amd64
```

Then install the program with `install` command

```bash
$ ~/Downloads/run-flogo-app-linux_amd64 install
#> Installing run-flogo-app...done
#> You can now directly execute run-flogo-app
```

### How to Use

After installing you can directly run it as any other command

```bash
$ run-flogo-app
#> Do you want to execute this app "/home/abhijit/Downloads/hello-world-linux_amd64" [Y/n]: y
#> Making app executable...

#> Executing: /home/abhijit/Downloads/hello-world-linux_amd64

```

You can also use this program as a stand alone binary, simple open a command promt and run the program.
But make sure the program is executable by running following command:

```bash
chmod +x ~/Downloads/run-flogo-app-linux-amd64
```

After this just run the executable with command:

```bash
~/Downloads/run-flogo-app-linux-amd64
```

You can provide `-d` (or `--debug`) flag to enable the debug logs in your flogo app, just make sure its the `first` argument to the program like this:

```bash
run-flogo-app -d
```

Also, you can pass command line arguments if your flogo app supports it; like this (the -d flag is not mandatory):

```bash
run-flogo-app -d arg1 arg2 arg3
```

## Commands and flags

### run-flogo-app

Run the most recent flogo app from your apps dir

#### Synopsis

Run the most recent flogo app from your configured apps dir. If the apps dir is not configured, the default will be used

```bash
run-flogo-app [flags]
```

#### Options

```text
  -d, --debug         Enable debug logs
  -h, --help          help for run-flogo-app
  -l, --list          List last 5 apps and choose a number to run
  -n, --name string   Run app with given (partial) name
```

#### SEE ALSO

* [run-flogo-app config](docs/run-flogo-app_config.md) - Print current config file
* [run-flogo-app delete](docs/run-flogo-app_delete.md) - Delete all the flogo apps in apps dir
* [run-flogo-app install](docs/run-flogo-app_install.md) - Install the program
* [run-flogo-app uninstall](docs/run-flogo-app_uninstall.md) - Uninstall the program
* [run-flogo-app version](docs/run-flogo-app_version.md) - Print the version info of the program

## Config file

When the program starts it creates a config file with name `.run-flogo-app` in your home directory. It is a simple json file which looks like this:

```json
{
  "appsDir": "/home/abhijit/Downloads",
  "appPattern": "^.+-linux_amd64.*$",
  "isUpdateAvailable": false,
  "updateURL": "",
  "releaseNotes": ""
}
```

You can override the programs' behavior by changing `appsDir` and `appPattern` variables in this file.
