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


Quick start 
-----------

  Generation:
  ```
  invstick -gen -width=45 -height=45 -correction=h -encoding Auto -template templates/stickers-A4-65-38X21.2.gohtml -dir=./out -gen-start=1 -gen-count=65 -barcode=true -mask="%03d" -prefix=#
  
  ```

  ![Result](img/invstick.png)


  From pipe:
  ```
  printf "One\nTwo\n" | invstick -template templates/stickers-A4-65-38X21.2.gohtml -dir=./out
  ```

  From file (use @):
  ```
  invstick -template templates/stickers-A4-65-38X21.2.gohtml -dir=./out -data=@/home/user/list.txt
  ```

  Web:
  ```
    invstick -listen=:8080
  ```

License
-------

  MIT

Author
------

Arteev Aleksey
 