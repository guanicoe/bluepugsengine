

<!-- apt-get install libzmq3-dev

chmod +x bluePugs

go build -o bluePugs src/*.go && ./bluePugs
 -->


<p align="center">
  <img alt="BluePugs" src="https://bluepugs.ovh/static/img/bluepugs_blue.svg" height="100" />

</p>


[![GoDoc](https://godoc.org/github.com/guanicoe/bluepugsengine?status.svg)](https://godoc.org/github.com/guanicoe/bluepugsengine)
[![Build Status](https://travis-ci.com/guanicoe/bluepugsengine.svg?branch=master)](https://travis-ci.com/guanicoe/bluepugsengine)
[![Software License](https://img.shields.io/badge/License-MPL%202.0-brightgreen.svg)](https://github.com/guanicoe/bluepugsengine/blob/master/LICENSE.md)
[![Go Report Card](https://goreportcard.com/badge/github.com/guanicoe/bluepugsengine)](https://goreportcard.com/report/github.com/guanicoe/bluepugsengine)


# BluePugsEngine (BETA)

## Introduction

> After a first try in python using multiprocessing, here's is the new Blue Pugs. Engine. A 100% Go program made to rapidly scan website for emails.The initial idea was to make theHarvester faster, but someone had already used threading. Blue Pugs takes advantage of heavy co-routines it order to send thousands of pugs (workers) on different websites.

<p align="center">
![dirty benchmark](https://github.com/guanicoe/bluepugsengine/blob/master/assets/quick_benchmarks.png)
</p>

## Beta

This program is basically in beta, that is that it basically works,but it needs further testing. If you find any bugs, orif you want a special feature, you can just submit an issue.
<!-- ## Code Samples

> You've gotten their attention in the introduction, now show a few code examples. So they get a visualization and as a bonus, make them copy/paste friendly. -->

## Screenshot

![screenshot](https://github.com/guanicoe/bluepugsengine/blob/master/assets/Screenshot%20from%202020-05-20%2013-51-47.png)

## Installation

The only dependency you need to install is libzmq

```sh
sudo apt-get install libzmq3-dev
```

You can then download the repository with the standard git clone command.

```sh
git clone https://github.com/guanicoe/bluepugsengine
cd bluepugsengine
chmod +x bluepugsengine
```
To run you simply specify the URL, and the scope of the search, that usually is the domain name.

```sh
./bluepugsengine -u http://yourwebsite.com -d yourwebsite
```


If you need to build the program, you should be able, as long as you have go installed, to type

```sh
go build -o bluepugsengine
```
