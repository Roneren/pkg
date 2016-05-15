# CoreStore optimized version

- no empty interfaces
- refactored code
- higher performance
- TODO(CS) update this readme for the code examples which are all ...

A [go](http://www.golang.org) (or 'golang' for search engine friendliness) implementation of [JSON Web Tokens](http://self-issued.info/docs/draft-jones-json-web-token.html)

[![Build Status](https://travis-ci.org/dgrijalva/jwt-go.svg?branch=master)](https://travis-ci.org/dgrijalva/jwt-go)

**NOTICE:** A vulnerability in JWT was [recently published](https://auth0.com/blog/2015/03/31/critical-vulnerabilities-in-json-web-token-libraries/).  As this library doesn't force users to validate the `alg` is what they expected, it's possible your usage is effected.  There will be an update soon to remedy this, and it will likey require backwards-incompatible changes to the API.  In the short term, please make sure your implementation verifies the `alg` is what you expect.


## What the heck is a JWT?

In short, it's a signed JSON object that does something useful (for example, authentication).  It's commonly used for `Bearer` tokens in Oauth 2.  A token is made of three parts, separated by `.`'s.  The first two parts are JSON objects, that have been [base64url](http://tools.ietf.org/html/rfc4648) encoded.  The last part is the signature, encoded the same way.

The first part is called the header.  It contains the necessary information for verifying the last part, the signature.  For example, which encryption method was used for signing and what key was used.

The part in the middle is the interesting bit.  It's called the Claims and contains the actual stuff you care about.  Refer to [the RFC](http://self-issued.info/docs/draft-jones-json-web-token.html) for information about reserved keys and the proper way to add your own.

## What's in the box?

This library supports the parsing and verification as well as the generation and signing of JWTs.  Current supported signing algorithms are HMAC SHA, RSA, RSA-PSS, and ECDSA, though hooks are present for adding your own.

## Parse and Verify

Parsing and verifying tokens is pretty straight forward.  You pass in the token and a function for looking up the key.  This is done as a callback since you may need to parse the token to find out what signing method and key was used.

```go
	token, err := csjwt.Parse(myToken, func(token *csjwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*csjwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return myLookupKey(token.Header["kid"]), nil
	})

	if err == nil && token.Valid {
		deliverGoodness("!")
	} else {
		deliverUtterRejection(":(")
	}
```

## Create a token

```go
	// Create the token
	token := csjwt.NewWithClaims(csjwt.SigningMethodHS256, csjwt.MapClaims{
		"foo": "bar",
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	})
	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString(mySigningKey)
```

## Extensions

This library publishes all the necessary components for adding your own signing methods.  Simply implement the `SigningMethod` interface and register a factory method using `RegisterSigningMethod`.  

Here's an example of an extension that integrates with the Google App Engine signing tools: https://github.com/someone1/gcp-jwt-go

## Project Status & Versioning

This library is considered production ready.  Feedback and feature requests are appreciated.  The API should be considered stable.  There should be very few backwards-incompatible changes outside of major version updates (and only with good reason).

This project uses [Semantic Versioning 2.0.0](http://semver.org).  Accepted pull requests will land on `master`.  Periodically, versions will be tagged from `master`.  You can find all the releases on [the project releases page](https://github.com/dgrijalva/jwt-go/releases).

While we try to make it obvious when we make breaking changes, there isn't a great mechanism for pushing announcements out to users.  You may want to use this alternative package include: `gopkg.in/dgrijalva/jwt-go.v2`.  It will do the right thing WRT semantic versioning.

## Migration Guide from v2 -> v3

Added the ability to supply a typed object for the claims section of the token.

Unfortunately this requires a breaking change. A few new methods were added to support this,
and the old default of `map[string]interface{}` was changed to `csjwt.MapClaims`.

The old example for creating a token looked like this..

```go
	token := csjwt.New(csjwt.SigningMethodHS256)
	token.Claims["foo"] = "bar"
	token.Claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
```

is now directly mapped to...

```go
	token := csjwt.New(csjwt.SigningMethodHS256)
	claims := token.Claims.(csjwt.MapClaims)
	claims["foo"] = "bar"
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
```

However, we added a helper `csjwt.NewWithClaims` which accepts a claims object.

Any type can now be used as the claim object for inside a token so long as it implements the interface `csjwt.Claims`.

So, we added an additional claim type `csjwt.StandardClaims` was added.
This is intended to be used as a base for creating your own types from,
and includes a few helper functions for verifying the claims defined [here](https://tools.ietf.org/html/rfc7519#section-4.1).

```go
	claims := csjwt.StandardClaims{
		Audience: "myapi"
		ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
	}
	token := csjwt.NewWithClaims(csjwt.SigningMethodHS256, claims)
```

On the other end of usage all of the `csjwt.Parse` and friends got a `WithClaims` suffix added to them.

```go
	token, err := csjwt.Parse(token, keyFunc)
	claims := token.Claims.(csjwt.MapClaim)
	//like you used to..
	claims["foo"]
	claims["bar"]
```

New method usage:
```go
	token, err := csjwt.ParseWithClaims(token, keyFunc, &csjwt.StandardClaims{})
	claims := token.Claims.(csjwt.StandardClaims)
	fmt.Println(claims.IssuedAt)
```

## Usage Tips

### Signing vs Encryption

A token is simply a JSON object that is signed by its author. this tells you exactly two things about the data:

* The author of the token was in the possession of the signing secret
* The data has not been modified since it was signed

It's important to know that JWT does not provide encryption, which means anyone who has access to the token can read its contents. If you need to protect (encrypt) the data, there is a companion spec, `JWE`, that provides this functionality. JWE is currently outside the scope of this library.

### Choosing a Signing Method

There are several signing methods available, and you should probably take the time to learn about the various options before choosing one.  The principal design decision is most likely going to be symmetric vs asymmetric.

Symmetric signing methods, such as HSA, use only a single secret. This is probably the simplest signing method to use since any `[]byte` can be used as a valid secret. They are also slightly computationally faster to use, though this rarely is enough to matter. Symmetric signing methods work the best when both producers and consumers of tokens are trusted, or even the same system. Since the same secret is used to both sign and validate tokens, you can't easily distribute the key for validation.

Asymmetric signing methods, such as RSA, use different keys for signing and verifying tokens. This makes it possible to produce tokens with a private key, and allow any consumer to access the public key for verification.

### JWT and OAuth

It's worth mentioning that OAuth and JWT are not the same thing. A JWT token is simply a signed JSON object. It can be used anywhere such a thing is useful. There is some confusion, though, as JWT is the most common type of bearer token used in OAuth2 authentication.

Without going too far down the rabbit hole, here's a description of the interaction of these technologies:

* OAuth is a protocol for allowing an identity provider to be separate from the service a user is logging in to.  For example, whenever you use Facebook to log into a different service (Yelp, Spotify, etc), you are using OAuth.
* OAuth defines several options for passing around authentication data. One popular method is called a "bearer token". A bearer token is simply a string that _should_ only be held by an authenticated user. Thus, simply presenting this token proves your identity. You can probably derive from here why a JWT might make a good bearer token.
* Because bearer tokens are used for authentication, it's important they're kept secret. This is why transactions that use bearer tokens typically happen over SSL.
 
## More

Documentation can be found [on godoc.org](http://godoc.org/github.com/dgrijalva/jwt-go).

## License

The command line utility included in this project (cmd/jwt) provides a straightforward example of token creation and parsing as well as a useful tool for debugging your own integration.  For a more http centric example, see [this gist](https://gist.github.com/cryptix/45c33ecf0ae54828e63b).

Copyright (c) 2012 Dave Grijalva

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

