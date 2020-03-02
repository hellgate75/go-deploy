<p align="center">
<image width="150" height="50" src="images/kube-go.png"></image>&nbsp;
<image width="260" height="410" src="images/golang-logo.png">
&nbsp;<image width="150" height="150" src="images/deploy-logo.png"></image>
</p><br/>
<br/>

# Go Deploy
GoLang deploy manager via command line or service


## Goals

Definition of an automated deploy system, easy to install, easy to update, compliant to innovation programs. No base frameworks, no system library links. Just get it, and us it. Only need is go-lang 1.3 or upper installed. 



## Reference Repositories

Reference is on modules repository:

* [Go Deploy Modules](https://github.com/hellgate75/go-deploy-modules) Modules for go-deploy executions

It has a configuration to use a similar technology TLS server (easy-to-use).

Please take a look at:

* [Go TCP Server](https://github.com/hellgate75/go-tcp-server) Server side TLS secure shell component

* [Go TCP Client](https://github.com/hellgate75/go-tcp-client) Client side TLS secure shell library


## How does it work?

Server starts with one or more input server certificate/key pairs. 

Call ```help``` ```--help``` or ```-h``` from command line to print out the available instructions and  command help.

It reads feeds that contains action instructions, it allows to store, read and use variables and create variables via remote shell command.

It allows to write configuration and variables in following encodings:

* YAML

* XML

* JSON

It executes a main feed that can contain multiple sub-feed, imported in the current one, on selected servers or importing new ones, related to new servers.

We are preparing a site about that features.


## Coming soon

Accordingly to policies we identified following modes for the system:

* Reading from a physical file (single run) - IMPLEMENTED

* Reading from a Web Source - FUTURE

* Reading from a Rest Service - FUTURE

* Reading from a Stream (JMS, IoT, Database flows, etc...) - FUTURE



## Sample code

Source test script is :
```
[test-sample.sh](/test-sample.sh)
```
It accepts optional parameters or the help request.

In order to use this sample please download and install [Go TCP Server](https://github.com/hellgate75/go-tcp-server) and run on same or different machine, then you can change the hosts file [here](/sample/env/hosts.yaml)

Commands included in the main and sub-feeds will give an overview of capabilities provided by the system.


Enjoy the experience.



## License

The library is licensed with [LGPL v. 3.0](/LICENSE) clauses, with prior authorization of author before any production or commercial use. Use of this library or any extension is prohibited due to high risk of damages due to improper use. No warranty is provided for improper or unauthorized use of this library or any implementation.

Any request can be prompted to the author [Fabrizio Torelli](https://www.linkedin.com/in/fabriziotorelli) at the follwoing email address:

[hellgate75@gmail.com](mailto:hellgate75@gmail.com)


