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


## Write your own modules

In linux system it's possible to write new features using the current client ones or developing new features.

Plugin(s) clients have own interfaces for writing a plugin, available pluggable clients are:

* [Go-TCP Client](https://githib.com/hellgate75/go-tcp-client), TLS custom client

For client(s) plugins (definition of custom clients), you can :

* Develop Proxy Function interface as described and with same name of function [proxy.GetConnectionHandlerFactory](https://github.com/hellgate75/go-deploy-clients/blob/master/proxy/proxy.go)

* Develp Client Wrapper as described in the interface [ConnectionHandler](/net/generic/interfaces.go)

An example of this kind of plugin is available in following repositories:
 
 * [Go-Deploy Client Modules](https://github.com/hellgate75/go-deploy-clients)

For deploy custom command(s) plugins you can:

* Develop Proxy Function interface [GetModulesMap](https://github.com/hellgate75/go-deploy-modules/blob/master/modules/stub.go)

* Develop a Discovery Function and allocating a map of string (unique plugin name) and [ProxyStub](/modules/meta/meta.go) that contains the discovery function, providing the command [Converter](/modules/meta/meta.go) component. Converter interface is used to parse the code from the [Feed](/types/generic.config.go) file and provifing a runnable element implementing [StepRunnable](/types/threadas/pool) interface, filled with parsed data.

An example of this kind of plugin is available in following repositories:
 
 * [Go-Deploy Command Modules](https://github.com/hellgate75/go-deploy-modules)


## Coming soon

Accordingly to policies we identified following modes for the system:

* Reading from a physical file (single run) - IMPLEMENTED

* Reading from a Web Source - FUTURE

* Reading from a Rest Service - FUTURE

* Reading from a Stream (JMS, IoT, Database flows, etc...) - FUTURE


## Official product documentation

Official produict documentation is available at:

* The product [Wiki](https://github.com/hellgate75/go-deploy/wiki) pages, that contain a lot of important information about haw to install and how to use this product.



## Sample code

Source test script is :
```
./test.sh
```
It accepts optional parameters or the help request.

Commands included in the main and sub-feeds will give you an overview of capabilities provided by the framework.

In order to execute the sample you must install [Go! TCP Server](https://github.com/hellate75/go-tcp-server) and execute the binaries in the sample folder 


Enjoy the experience.



## License

The library is licensed with [LGPL v. 3.0](/LICENSE) clauses, with prior authorization of author before any production or commercial use. Use of this library or any extension is prohibited due to high risk of damages due to improper use. No warranty is provided for improper or unauthorized use of this library or any implementation.

Any request can be prompted to the author [Fabrizio Torelli](https://www.linkedin.com/in/fabriziotorelli) at the following email address:

[hellgate75@gmail.com](mailto:hellgate75@gmail.com)


