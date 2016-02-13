phpsession
===================

PHP session encoder/decoder written in Go  
[![Build Status](https://secure.travis-ci.org/yvasiyarov/phpsession.png?branch=master)](http://travis-ci.org/yvasiyarov/phpsession)

Installation
------------

Install:

- ~~The recommended way~~ to install is using gonuts.io:


    nut get yvasiyarov/phpsession
    for more information, please, go to the http://www.gonuts.io/yvasiyarov/phpsession

- Using default go get tool:


    go get github.com/corestoreio/csfw/util/php/phpsession

Getting started
---------------

Example: load php session data from redis:

    if sessionId, err := req.Cookie("frontend"); err == nil {
        if sessionData, err := redis.Get("PHPREDIS_SESSION:" + sessionId.Value); err == nil {
            decoder := phpsession.NewPhpDecoder(sessionData.String())
            if sessionDataDecoded, err := decoder.Decode(); err == nil {
                //Do something with session data  
            }
        } else {
            //Can not load session - it can be expired
        }
    }

Example: Encode php session data:

    data := make(PhpSession)
    data["make some"] = " changes"
    encoder := NewPhpEncoder(data)
    if result, err := encoder.Encode(); err == nil {
        //Write data to redis/memcached/file/etc
    }

Copyright
----------------------------
2013-2014 Yuriy Vasiyarov   
2014 Yuriy Vasiyarov, Maksim Naumov. 

All rights reserved.
