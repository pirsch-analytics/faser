# Changelog

## 1.7.0

* upgraded to Go version 1.24
* updated dependencies

## 1.6.1

* improved logging
* updated dependencies

## 1.6.0

* refactored project
* added Makefile
* added vendor directory
* updated dependencies
* upgraded to Go 1.22

## 1.5.0

* switched to Chi router
* improved CORS settings
* only allow GET requests
* updated dependencies
* upgraded to Go 1.20

## 1.4.0

* updated dependencies

## 1.3.0

* respond with default icon to "bad" requests (like keywords instead of a valid URL)

## 1.2.1

* use caching headers to serve icons
* opt out of Google's FLoC
* updated dependencies

## 1.2.0

* select a specific default icon (new optional `fallback` parameter)
* updated dependencies

## 1.1.0

* removed deprecated package io/ioutil, the minimum Go version is now 1.16

## 1.0.1

* fixed selecting largest icon
* fixed fallback to /favicon.ico if index cannot be loaded
* fixed grouping sizes

## 1.0.0

Initial release.
