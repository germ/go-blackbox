# Go-Blackbox
## A library for automagic paintings
[![GoDoc](https://godoc.org/github.com/germ/go-blackbox?status.svg)](https://godoc.org/github.com/germ/go-blackbox)

This library is a simple wrapper for the API provided by [The BlackBox](http://theblackbox.tk), API keys can be obtained from http://theblackbox.tk/api or by using the included developer credentials.

Note: The functionality allowing inspection/wrapping of existing io.Readers is experimental, tread carefully.
In addition, payment processing via Stripe is in test mode. Use card # 4242 4242 4242 4242 to test charges

 General usage:
	Using your ID make a call to create or info, this will return a 
	struct with the session ID/all sessions.

	The structure has a Upload method that will drain reader r and
	store it on my server. Max upload size is 10MB, split larger files
	or host elsewhere and provide a link. Please don't upload more then
	500MB of data, to my server thank you.

	Lastly when all assets have been submitted, call the Finalize method
	to submit it for approval. If I decline the work, an email will be sent
	to the registered account. Otherwise a work will show up in the future!

	The Session structure has a flag "DevMode", with this set true all calls 
	will not generate side effects.
