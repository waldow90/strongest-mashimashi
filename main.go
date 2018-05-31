package main

import (
	"bufio"
	"errors"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/labstack/gommon/log"
	"golang.org/x/net/context"

	"google.golang.org/appengine"
)

const (
	// inspect word files in advance by "wc -l" command
	nounLen      = 76216
	adjectiveLen = 26664
)

type handler struct {
	noun, adjective []string
}

func (h *handler) Get(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		w.Write([]byte(index))
	case "/api/v1/phrase":
		// adjective 1
		adj1 := rand.Intn(adjectiveLen)
		// adjective 2
		adj2 := rand.Intn(adjectiveLen)
		// noun
		noun3 := rand.Intn(adjectiveLen)
		w.Write([]byte(h.adjective[adj1] + " " + h.adjective[adj2] + " " + h.noun[noun3]))
	}
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	switch r.Method {
	case http.MethodGet:
		h.Get(ctx, w, r)
	default:
		w.Write([]byte("unimplemented method!"))
	}
}

func readWords(filename string, lines int) ([]string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, errors.New("failed to open " + filename + ": " + err.Error())
	}
	defer f.Close()

	words := make([]string, lines)
	var c int
	s := bufio.NewScanner(f)
	for s.Scan() {
		words[c] = s.Text()
		c++
	}
	return words, nil
}

func main() {
	rand.Seed(time.Now().UnixNano())
	var err error
	h := &handler{}

	h.noun, err = readWords("noun.txt", nounLen)
	if err != nil {
		log.Errorf("failed to read noun: %s", err.Error())
		return
	}

	h.adjective, err = readWords("adjective.txt", adjectiveLen)
	if err != nil {
		log.Errorf("failed to read adjective: %s", err.Error())
		return
	}

	http.Handle("/", h)
	appengine.Main()
}

const index = `
<html>
<head>
	<title>Generate a sentence with 3 random words</title>
	<meta name="viewport" content="width=device-width,initial-scale=1">
	<style>
	body {
		font-family: Sans-Serif;
	}
	</style>
	<script>
	window.addEventListener('load', _ => {
		document.getElementById('submit').addEventListener('click', _ => {
			fetch('/api/v1/phrase', {
				method: 'GET'
			}).then(response => {
				if (response.ok) {
					return response.text();
				} else {
					throw new Error();
				}
			}).then(text => {
				document.getElementById('words').textContent = text;
			}).catch(error => {
				console.log(error);
			});
		});
	});

	let copyText = (str) => {
		var tmp = document.createElement('div');
		tmp.appendChild(document.createElement('pre')).textContent = str;

		var s = tmp.style;
		s.position = 'fixed';
		s.left = '-100%';

		document.body.appendChild(tmp);
		document.getSelection().selectAllChildren(tmp);
		document.execCommand('copy');
		document.body.removeChild(tmp);
	}

	let copy = (id) => {
		text = document.getElementById(id).textContent
		if (text == '--- --- ---') {
			return
		}
		copyText(text.replace(/\s+/g, ""));
	}
	</script>
</head>
<body>
<div>Generate a sentence with 3 random words</div>
<br>
<div>Concept</div>
<li>Choice 3 words from English dictionary to generate strong password.</li>
<li>For easy remember, 3 words are choose as "adjective" "adjective" "noun".</li>
<li>Words that are too short (less than 3) or too long (more than 10) are excluded.</li>
<li>Note that since a word may (in dictionary) include space or hyphen, generated sentence may look like more than 3 words.</li>
<br>
<br>
<span><button id="submit">Push to generate</button></span>
<span><button id="copy" onclick="copy('words')">Copy to clipboard (without whitespace)</button></span>
<br>
<br>
<div id="words" style="font-size:x-large">--- --- ---</div>
<input type="hidden" id="wordsforcopy" value="">
<br>
<br>
Contact: @pankona (<a href="https://twitter.com/pankona">twitter</a>, <a href="https://github.com/pankona">github</a>)
</body>
</html>
`
