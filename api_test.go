package blackbox

import (
	"bytes"
	"io/ioutil"
	"testing"
)

/* General usage:
	Using your ID make a call to create or info, this will return a 
	structure with the session ID/all sessions. Save these for future reference

	The structure has a Upload method that will drain reader r and
	store it on my server. Max upload size is 10MB, split larger files
	or host elsewhere and provide a link. Please don't upload more then
	500MB of data, to my server thank you.

	Lastly when all assets have been submitted, call the Finalize method
	to submit it for approval. If I decline the work, an email will be sent
	to the registered account. Otherwise a work will show up in the future!

	The Session structure has a flag "DevMode", with this set true all calls 
	will not generate side effects.
*/

var TestCreds = DevCreds

var TestData = `you say it isn't logical
if it's not black or white
it's either positive or negative
either day or night
can't be 6 of one 
half dozen of the other
you know what I mean
know what I'm sayin brother

make up your mind
just give me the truth
don't wrap me in a cord
in a telephone booth
is it “A” or “B”
it's gotta be part of a set
I work with truths
before I place my bet

binary numbers that intersect
ands or nots or or's it can be
part of the superset
the limbs of the tree
true or false
you just gotta decide
algebraic notation
proves if you lied

could you be wrong
could there be areas of gray
in matters of love
it's not just what you say
sometimes it's what's missing
that matters the most
no salty or sweet
like a piece of dry toast    

is science perfect
how the hell would I know
can only go by
the factors that show
but I got this feeling
it's more than neurologic
in matters of the heart
it takes more than boolean logic

Gomer Lepoet`

func TestCreate(t *testing.T) {
	t.Log("Create Test")
	s, err := Create(TestCreds.UID)
	t.Log(s, err)
	
	if err == nil {
		TestCreds = s
	}
}

func TestUpload(t *testing.T) {
	t.Log("Upload Test")
	buf := bytes.NewBufferString(TestData)
	err := TestCreds.Upload(buf)
	t.Log(err)
}

func TestInfo(t *testing.T) {
	t.Log("Info Test")
	s, err := Info(TestCreds.UID)
	t.Log(s, err)
}

func TestAttach(t *testing.T) {
	readone := TestCreds.Attach(bytes.NewBufferString(TestData[:30]))

	data, _ := ioutil.ReadAll(readone)
	t.Log("One: ", string(data))
}

func TestFinal(t *testing.T) {
	t.Log("Final Test")
	err := TestCreds.Finalize(10 * 100)
	t.Log(err)
}

