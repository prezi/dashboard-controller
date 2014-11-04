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

    dashboard-controller$ go run webserver.go

(There is currently a bug where the webserver must be run from the webserver directory to ensure the correct path to "form.html".)

####Slave
------------------

Again, you will need [Go](https://golang.org/) installed to run and/or compile the source code. We recommend that you compile the slave.go file to produce a binary file executable on one slave machine, then copy the slave binary file to all your slaves. A slave does not need Go installed to run the binary file, so you can save yourself the Go installation time for each slave. (We used Raspberry Pis as slaves, and each Go installation took 70-100 minutes.)

From the repository directory, compile slave code with the command, 
 
    dashboard-controller$ go install slave

Then run the executable with,
 
    dashboard-controller$ $GOPATH/bin/slave

The slave will begin listening on a [default port number](https://github.com/prezi/dashboard-controller/blob/master/src/slave/slave.go#L14). Optionally, you can specify a port number with the [-port flag](https://github.com/prezi/dashboard-controller/blob/master/src/slave/slave.go#L50), 
 
    dashboard-controller$ $GOPATH/bin/slave -port=9999
    
    
####Running the tests
------------------

In order to run the tests, first you need to install [Go](https://golang.org/)'s [testify](https://godoc.org/github.com/stretchr/testify) package by running the command,

    dashboard-controller$ go get github.com/stretchr/testify

Then navigate into the folder where the test file is located, and run the tests with the 'go test test\_file\_name\_without\_extension'. Here is an example,

    dashboard-controller/src/master$ go test master
    

Posting a URL
------------------
Access the website running on your [localhost](https://github.com/prezi/dashboard-controller/blob/master/src/webserver/webserver.go#L126). Fill in the text fields, submit, and see your url post on the indicated slave. :) 
