Dashboard Controller 
====================

Introduction
------------------
Use the Dashboard Controller program to remotely open urls on browser windows.

Our implemented system runs slave code on Raspberry Pis connected to monitors. 
The webserver and master both run on one Mac Mini. 


File Architecture
------------------
 - slave 
  - Receives url from master and loads url in a browser. 
  - Compatible with OS X and Linux operating systems. 
 - master
  - Receives JSON from the webserver and passes JSON to proper slave. 
  - Compatible with OS X. 
 - webserver
  - Receives JSON from user input and passes JSON to master. 
  - Compatible with OS X. 


Getting Started
------------------

If you are not familiar with Go and run into trouble with the following instructions, please visit  ["How to Write Go Code"](https://golang.org/doc/code.html) for more details. 

####Webserver and Master
------------------

You will need [Go](https://golang.org/) installed to run and/or compile the source code. 
Clone this repository to the machine that will run master/webserver. 

To run the code from the cloned repository directory, 

    dashboard-controller$ go run src/master/master.go


####Slave
------------------

Again, you will need [Go](https://golang.org/) installed to run and/or compile the source code. We recommend that you compile the slave.go file to produce a binary file executable on one slave machine, then copy the slave binary file to all your slaves. A slave does not need Go installed to run the binary file, so you can save yourself the 70-100 minutes of installation time for each slave (assuming that you're using Raspberry Pis as slaves). 

From the repository directory, compile slave code with the command, 
 
    dashboard-controller$ go install slave

Then run the executable with,
 
    dashboard-controller$ $GOPATH/bin/slave

####Posting a URL
------------------
Access the website running on your [localhost](https://github.com/prezi/dashboard-controller/blob/master/src/webserver/webserver.go#L126). Fill in the text fields, submit, and see your url post on the indicated slave. :) 
