invstick
==========

[![Build Status](https://travis-ci.org/arteev/invstick.svg?branch=master)](https://travis-ci.org/arteev/invstick)

Description
-----------

Tool to generate stickers with QR codes

Installation
------------

This package can be installed with the go get command:

    go get github.com/arteev/invstick
    go install github.com/arteev/invstick
        
License
-------

  MIT

Author
------

Arteev Aleksey


{{range $index, $element := . }}
            <tr style="height: 2.12cm;">
                <td style="width: 3.8cm; background-color: green; ">
                <span>Номер: S -  {{$index}} - {{$element.Num}}</span>
                <img src="{{.QRCode}}"></img>  
                </td>
                
            </tr>                
            {{end}}
            