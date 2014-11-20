Dashboard Controller 
====================

![intro_image](https://raw.githubusercontent.com/prezi/dashboard-controller/master/README_images/webpage.png)

Contents
 - [Introduction](https://github.com/prezi/dashboard-controller#introduction)
 - [Architecture](https://github.com/prezi/dashboard-controller#architecture)
 - [Set Up](https://github.com/prezi/dashboard-controller#getting-started)
 - [Posting a URL](https://github.com/prezi/dashboard-controller#posting-a-url)

Introduction
------------------

Use the Dashboard Controller program to remotely open urls on browser windows.

Our implemented system uses [Google Chrome](http://www.google.com/chrome/) web browser. 
Slave runs on Raspberry Pis connected to monitors. 
The webserver and master both run on one Apple Mac Mini. 

Together, they make a delicious Apple-Raspberry-Pi system. 

Architecture
------------------

 - slave 
  - Receives url from master and loads url in a browser. 
  - Compatible with OS X and Linux operating systems with Google Chrome installed. 
 - master
  - Receives JSON from the webserver and passes JSON to proper slave. 
  - Compatible with OS X. 
 - webserver
  - Receives JSON from user input and passes JSON to master. 
  - Compatible with OS X. 


Set Up
------------------

If you are not familiar with Go and run into trouble with the following instructions, please visit  ["How to Write Go Code"](https://golang.org/doc/code.html) for more details. 

####Webserver and Master
------------------

You will need [Go](https://golang.org/) installed to run and/or compile the source code. 
Clone this repository to the machine(s) that will run the master and/or webserver. 

To run the code from the cloned repository directory, 

    dashboard-controller$ go run src/master/master.go

And for the webserver, 

    dashboard-controller$ go run src/webserver/webserver.go 

(There is currently a bug where the webserver must be run from the root directory to ensure the correct path to "form.html".)

####Slave
------------------

Again, you will need [Go](https://golang.org/) installed to run and/or compile the source code. We recommend that you compile the slave.go file to produce a binary file executable on one slave machine, then copy the slave binary file to all your slaves. A slave does not need Go installed to run the binary file, so you can save yourself the Go installation time for each slave. (We used Raspberry Pis as slaves, and each Go installation took 70-100 minutes.)

From the repository directory, compile slave code with the command, 
 
    dashboard-controller$ go install slave

Then run the executable with a speficied slave name and master IP address. 
 
    dashboard-controller$ $GOPATH/bin/slave -slaveName="Main Lobby" -masterIP=10.0.0.195

The slave will begin listening on its default port number: 8080. The master's default port number is 5000.
Optionally, you can specify a port number for the slave and/or master with the -port and -masterPort flags, respectively.
 
    dashboard-controller$ $GOPATH/bin/slave -slaveName="Main Lobby" -port=9999 -masterIP=10.0.0.195 -masterPort=9090 
    
The slave will automatically map itself to the master and periodically emit heartbeats to the master. If the slave's heartbeats do not reach the master for some period of the time, the master will mark the slave as dead and remove it from the map of available slaves. 

Posting a URL
------------------

Access the website running on your [localhost](https://github.com/prezi/dashboard-controller/blob/master/src/webserver/webserver.go#L126). Fill in the text fields, submit, and see your url post on the indicated slave. :) 

Tests
------------------

In order to run the tests, first install [Go](https://golang.org/)'s [testify](https://godoc.org/github.com/stretchr/testify) package with the command

    dashboard-controller$ go get github.com/stretchr/testify

Navigate into the directory where the test file is located, and run the tests with 'go test'. For example, in the master directory that contains master.go and master_test.go, run

    dashboard-controller/src/master$ go test 

The -cover flag will output the test coverage %. 

    dashboard-controller/src/master$ go test -cover

For more testing options, please see the this [Go blog post](http://blog.golang.org/cover). 
