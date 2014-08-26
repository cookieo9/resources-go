resources-go
============

An assets-loading package for Go.

Applications can use this package to load assets from zip-files (incuding a zip file bundled in the executable),
the filesystem, or other sources through a single interface. Also allows for the building of a search path to access
files sequentially through a set of application defined locations.

This is the development version of the API. It may change at any time, if you want stability use one of the released versions. They are also stored in branches of the git repository.

[![Build Status](https://travis-ci.org/cookieo9/resources-go.svg)](https://travis-ci.org/cookieo9/resources-go)
[![GoDoc](https://godoc.org/github.com/cookieo9/resources-go?status.png)](https://godoc.org/github.com/cookieo9/resources-go)
[![Coverage](http://gocover.io/_badge/github.com/cookieo9/resources-go)](http://gocover.io/github.com/cookieo9/resources-go)

Documentation and Other Versions
--------------------------------

To see the code and documentation for every released version of this package, see: http://gopkg.in/cookieo9/resources-go.v2


Embedding Zip-Files
-------------------

To embed a zip file into your executable do the following:
 - Create executable (eg: go build -o myApp)
 - Create zip file  (eg: zip -r assets.zip assets)
 - Append zip file to executable (eg: cat assets.zip >> myApp)
 - Adjust the offsets in the zip header (optional) (zip -A myApp)

License
-------
http://cookieo9.mit-license.org/2012
