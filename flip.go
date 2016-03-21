package main

import (
	"fmt"
	"encoding/json"
	"errors"
	"flag"
	"ipernity"
	"io/ioutil"
	"os"
	"strconv"
)

const (
	last_doc_file = "ipernity_last_doc"
)

type Doc struct {
	Docid        string   `json:"docid"`
	OriginalUrl  string   `json:"originalurl"`
	OriginalSize int      `json:"originalsize"`
	Ext          string   `json:"ext"`
	Media        string   `json:"media"`
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	Date         string   `json:"date"`
	Albums       []string `json:"albums"`
	FamilyVisible bool `json:"family"`
	FriendVisible bool `json:"friends"`
	PublicVisible bool `json:"public"`
}

type Docslice []Doc

var (
	Docs Docslice
)

func main() {
	startdocPtr := flag.Int("startdoc", 0, "starting document to download (0 means most recent undownloaded document)")
	numdocsPtr := flag.Int("numdocs", 0, "number of documents to download (0 means all)")

	flag.Parse()

	startdoc := *startdocPtr
	numdocs := *numdocsPtr

	fmt.Println("startdoc: ", startdoc)
	fmt.Println("numdocs: ", numdocs)

	err := ipernity.Login()
	if err != nil {
		fmt.Println("Error logging in to ipernity: " + err.Error())
		os.Exit(1)
	}

	totaldocs, err := GetNumDocs()
	fmt.Println("total number of documents ", totaldocs)
	if err != nil {
		fmt.Println("Error getting number of pictures on ipernity: " + err.Error())
		os.Exit(2)
	}

	if startdoc < 1 {
		startdocascii, err := ioutil.ReadFile(last_doc_file)
		if err != nil {
			startdoc = 1
		} else {
			startdoc, err = strconv.Atoi(string(startdocascii))
			if (err != nil) || (startdoc < 1) {
				fmt.Println("corrupted ipernity last doc file")
				os.Exit(5)
			}
		}
	}

	if numdocs < 1 {
		numdocs = totaldocs - startdoc
	}

	if numdocs < 1 {
		fmt.Println("No documents to download")
		os.Exit(0)
	}

	if startdoc < 1 {
		fmt.Println("Must start at at least document 1")
		os.Exit(6)
	}

	err = EnumeratePictures(startdoc, numdocs, totaldocs)
	if err != nil {
		fmt.Println("Error enumerating pictures on ipernity: " + err.Error())
		os.Exit(3)
	}

	err = SavePictures(startdoc)
	if err != nil {
		fmt.Println("Error saving picture from ipernity: " + err.Error())
		os.Exit(4)
	}
}

func strToBool(s string) bool {
	if s != "0" {
		return true
	}
	return false
}

func GetNumDocs() (int, error) {

	usergetdata, err := ipernity.Call_user_get("")

	if err != nil {
		return 0, err
	}

	numdocs, err := strconv.Atoi(usergetdata.User.Count.Docs)
	if err != nil {
		return 0, err
	}

	return numdocs, nil
}

func EnumeratePictures(startdoc int, numdocs int, totalnumdocs int) error {
	var (
		doclist []ipernity.Docgetlist
	)

	if (startdoc < 1) || (numdocs < 1) || (totalnumdocs < 1) {
		return errors.New("all function arguments must be >= 1")
	}

	fmt.Println("startdoc " + strconv.Itoa(startdoc))
	fmt.Println("numdocs " + strconv.Itoa(numdocs))
	fmt.Println("totalnumdocs " + strconv.Itoa(totalnumdocs))
	numpages := (totalnumdocs / 100)
	if (totalnumdocs % 100) != 0 {
		numpages++
	}
	fmt.Println("numpages " + strconv.Itoa(numpages))

	// XXX: only download the pages we need

	startdocpage := startdoc / 100
	if (startdoc % 100) != 0 {
		startdocpage++
	}

	enddocpage := (startdoc + numdocs) / 100
	if ((startdoc + numdocs) % 100) != 0 {
		enddocpage++
	}

	fmt.Printf("startdocpage %v enddocpage %v\n", startdocpage, enddocpage)

	for page := 1; page <= numpages; page++ {
		fmt.Printf("Downloading photo list data page %v of %v\n",
			page, numpages)
		tmplist, err := ipernity.Call_doc_getList("", page, "")
		if err != nil {
			return err
		}
		doclist = append(doclist, tmplist)
	}

	k := 0
	docpos := 0
	// a list of pages of docs
	for i := len(doclist) - 1; i >= 0; i-- {
		// if we're done, exit the loop over the pages
		if docpos >= (numdocs + startdoc) {
			break
		}
		// one page of docs
		for j := len(doclist[i].Docs.Doc) - 1; j >= 0; j-- {
			// if we're done, exit the loop over the page
			if docpos < startdoc {
				docpos++
				continue;
			}
			// if we shouldn't start yet, skip to the next doc
			if docpos >= (numdocs + startdoc) {
				break
			}
			if (k % 10) == 0 {
				fmt.Printf("Document metadata %v\n", k)
			}
			Docs = append(Docs, Doc{Docid: doclist[i].Docs.Doc[j].Doc_id})

			docinfo, err := ipernity.Call_doc_get(Docs[k].Docid, "")
			if err != nil {
				return err
			}

			Docs[k].OriginalUrl = docinfo.Doc.Original.Url
			Docs[k].OriginalSize, _ = strconv.Atoi(docinfo.Doc.Original.Bytes)
			Docs[k].Title = docinfo.Doc.Title
			Docs[k].FamilyVisible = strToBool(docinfo.Doc.Visibility.Isfamily)
			Docs[k].FriendVisible = strToBool(docinfo.Doc.Visibility.Isfriend)
			Docs[k].PublicVisible = strToBool(docinfo.Doc.Visibility.Ispublic)
			Docs[k].Description = docinfo.Doc.Description
			Docs[k].Date = docinfo.Doc.Dates.Created
			Docs[k].Ext = docinfo.Doc.Original.Ext
			Docs[k].Media = docinfo.Doc.Media

			if docinfo.Doc.Count.Albums != "0" {
				docconts, err := ipernity.Call_doc_getContainers(Docs[k].Docid)
				if err != nil { 
					return err 
				}
				numalbums, err := strconv.Atoi(docconts.Albums.Total)
				if err != nil { 
					return err 
				}
				for l := 0; l < numalbums; l++ {
					Docs[k].Albums = append(Docs[k].Albums, docconts.Albums.Album[l].Title)
				}
			} 

			k++
			docpos++
		}
	}

	return nil
}

func SavePicture(doc Doc, i int, startdoc int) error {
	resp, err := ipernity.HttpClient.Get(doc.OriginalUrl)
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if len(body) != doc.OriginalSize {
		fmt.Printf("downloaded length: %v original length %v\n",
			len(body), doc.OriginalSize)
		return errors.New("length mismatch for URL " + doc.OriginalUrl)
	}

	err = ioutil.WriteFile(doc.Docid+"."+doc.Ext, body, 0644)
	if err != nil {
		return err
	}

	docjson, err := json.Marshal(doc)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(doc.Docid+".json", docjson, 0644)
	if err != nil {
		return err
	}

	istr := strconv.Itoa(i + startdoc + 1)
	err = ioutil.WriteFile(last_doc_file, []byte(istr), 0644)

	return nil
}	

func SavePictures(startdoc int) error {
	for i, doc := range Docs {
		if (i % 10) == 0 {
			fmt.Printf("downloading doc %v\n", i)
		}
		err := SavePicture(doc, i, startdoc)
		if err != nil {
			return err
		}
	}

	return nil
}
