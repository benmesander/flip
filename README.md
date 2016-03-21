# flip
flip backs photos stored on the [ipernity photo sharing site](http://www.ipernity.com) up to local storage along with metadata suitable for exporting said photos to another photo sharing service. It is used in conjunction with my flop package to mirror my ipernity photos to flickr.

This program is not currently polished enough to be useful to a general audience; using it requires some programming skills. It is at best a tool for creating your own photo backup solution.

flip is written in the go language, it has only been tested with version 1.5.2 and above.

API Key
=======
To obtain an API key for ipernity, one visits http://www.ipernity.com/apps/key/0 and obtains a noncommercial API key and API secret. These should be exported in your environment in `IPERNITY_API_KEY` and `IPERNITY_API_SECRET` respectively. They can also be placed into src/ipernity.go in the source tree. Ipernity API keys seem to last for a limited number of API calls before being declared invalid. I have many tens of thousands of photos on ipernity, and in order to download them all, I had to re-register several times to get new keys in order to download all my photos.

Usage
=====
flip is intended to be run periodically; perhaps out of `cron`. It has a couple optional command line arguments, but by default, with no arguments, it incrementally downloads all photos uploaded to ipernity since the last time which it was run.

1. `-startdoc #` Allows you to specify the starting document to download from ipernity (1 is the first document in your accouint). The default if not specified is the first document which has not been downloaded before.
2. `-numdocs #` Allows you to specify the number of documents to download. The default is all of the documents which have not been downloaded before.
 
The first time you run flip, and every time you update your API key, you will need to authorize the API key to your account. flip will display a message:

```go to http://some/authorization/url/from/ipernity anmd grant the permissions, then press <ENTER>```

Follow the link in the message with a web browser, grant the API key permissions to your account, press ENTER, and flip will continue. It will store your credentials in a file in the directory where you ran flip, `ipernity_auth_token`. If flip ever fails to authenticate to ipernity for some reason, you can delete this file and the application will go through the authentication procedure above again.

flip will download photos/movies/documents/etc. to the directory in which it was invoked. Each document will consist of itself (ie, `ipernity_document_id.jpg`) as well as a JSON metadata blob (ie, `ipernity_document_id.json`). flip also stores its persistent state in a file named `ipernity_last_doc`. This simply contains the number of the last document which was downloaded. Each document stored in ipernity gets a numeric document id. This is used for the filename and the metadata filename. The metadata stored by flop consists of the following:

- Docid - the ipernity document ID
- OriginalURL - the URL of the "original" (full) size document on ipernity
- OriginalSize - the size, in bytes, of the original document
- Ext - the filename extension of the original document (ie, .jpg, .mov)
- Media - the type of the document (ie, "movie", "photo", "document")
- Title - the title of the document
- Description - the description of the document
- Date - the date which the document was uploaded
- Albums - a list of photo album names which the document was placed in on ipernity
- FamilyVisible - is the document visible to users who are part of your circle of Family members
- FriendVisible - is the document visible to users who are part of your circle of Friends
- PublicVisible - is the document visible to all users

If flip reports an error, the usual reasons for which are network connectivity or some problem with ipernity, it is always safe to re-run it. It is designed to not damage your data.

Example
=======

An example run of flip on an account with 4056 documents, 4040 of which have been downloaded before, so 16 new documents will be downloaded:
```

ben@nederland:~/src/flip$ go run flip.go
startdoc:  0
numdocs:  0
total number of documents  4056
startdoc 4040
numdocs 16
totalnumdocs 4056
numpages 41
startdocpage 41 enddocpage 41
Downloading photo list data page 1 of 41
Downloading photo list data page 2 of 41
Downloading photo list data page 3 of 41
Downloading photo list data page 4 of 41
Downloading photo list data page 5 of 41
Downloading photo list data page 6 of 41
Downloading photo list data page 7 of 41
Downloading photo list data page 8 of 41
Downloading photo list data page 9 of 41
Downloading photo list data page 10 of 41
Downloading photo list data page 11 of 41
Downloading photo list data page 12 of 41
Downloading photo list data page 13 of 41
Downloading photo list data page 14 of 41
Downloading photo list data page 15 of 41
Downloading photo list data page 16 of 41
Downloading photo list data page 17 of 41
Downloading photo list data page 18 of 41
Downloading photo list data page 19 of 41
Downloading photo list data page 20 of 41
Downloading photo list data page 21 of 41
Downloading photo list data page 22 of 41
Downloading photo list data page 23 of 41
Downloading photo list data page 24 of 41
Downloading photo list data page 25 of 41
Downloading photo list data page 26 of 41
Downloading photo list data page 27 of 41
Downloading photo list data page 28 of 41
Downloading photo list data page 29 of 41
Downloading photo list data page 30 of 41
Downloading photo list data page 31 of 41
Downloading photo list data page 32 of 41
Downloading photo list data page 33 of 41
Downloading photo list data page 34 of 41
Downloading photo list data page 35 of 41
Downloading photo list data page 36 of 41
Downloading photo list data page 37 of 41
Downloading photo list data page 38 of 41
Downloading photo list data page 39 of 41
Downloading photo list data page 40 of 41
Downloading photo list data page 41 of 41
Document metadata 0
Document metadata 10
downloading doc 0
downloading doc 10
ben@nederland:~/src/flip$ 
```




